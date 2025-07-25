package subscribe

import (
	"path"
	"strings"
	"time"

	config "project/mqtt"
	"project/mqtt/publish"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-basic/uuid"
	"github.com/panjf2000/ants"
	"github.com/sirupsen/logrus"
)

var SubscribeMqttClient mqtt.Client
var TelemetryMessagesChan chan map[string]interface{}

func GenTopic(topic string) string {
	topic = path.Join("$share/mygroup", topic)
	return topic
}

func SubscribeInit() error {

	//实例限流客户端
	//initialize.NewAutomateLimiter()
	// 创建mqtt客户端
	subscribeMqttClient()
	// 创建消息队列
	telemetryMessagesChan()

	//消息订阅
	err := subscribe()
	return err
}

func subscribe() error {
	// 订阅attribute消息
	err := SubscribeAttribute()
	if err != nil {
		logrus.Error(err)
		return err
	}
	// 订阅event消息
	err = SubscribeEvent()
	if err != nil {
		logrus.Error(err)
		return err
	}
	//订阅telemetry消息
	err = SubscribeTelemetry()
	if err != nil {
		logrus.Error(err)
		return err
	}
	//订阅设备注册消息
	err = SubscribeRegister()
	if err != nil {
		logrus.Error(err)
		return err
	}

	//订阅一些非标准的注册消息
	err = SubscribeCustomer()
	if err != nil {
		logrus.Error(err)
		return err
	}

	// 订阅设备命令消息
	err = SubscribeCommand()
	if err != nil {
		logrus.Error(err)
		return err
	}

	// 订阅OTA命令消息
	err = SubscribeOtaUpprogress()
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func subscribeMqttClient() {
	// 初始化配置
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.MqttConfig.Broker)
	opts.SetUsername(config.MqttConfig.User)
	opts.SetPassword(config.MqttConfig.Pass)
	id := "thingspanel-go-sub-" + uuid.New()[0:8]
	opts.SetClientID(id)
	logrus.Info("clientid: ", id)

	// 干净会话
	opts.SetCleanSession(true)
	// 恢复客户端订阅，需要broker支持
	opts.SetResumeSubs(true)
	// 自动重连
	opts.SetAutoReconnect(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetMaxReconnectInterval(200 * time.Second)
	// 消息顺序
	opts.SetOrderMatters(false)
	opts.SetOnConnectHandler(func(_ mqtt.Client) {
		logrus.Println("mqtt connect success")
	})
	// 断线重连
	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		logrus.Println("mqtt connect  lost: ", err)
		SubscribeMqttClient.Disconnect(250)
		for {
			if token := SubscribeMqttClient.Connect(); token.Wait() && token.Error() != nil {
				logrus.Error("MQTT Broker 1 连接失败:", token.Error())
				time.Sleep(5 * time.Second)
				continue
			}
			subscribe()
			break
		}
	})

	SubscribeMqttClient = mqtt.NewClient(opts)
	// 等待连接成功，失败重新连接
	for {
		if token := SubscribeMqttClient.Connect(); token.Wait() && token.Error() != nil {
			logrus.Error("MQTT Broker 1 连接失败:", token.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

}

// 创建消息队列
func telemetryMessagesChan() {
	TelemetryMessagesChan = make(chan map[string]interface{}, config.MqttConfig.ChannelBufferSize)
	writeWorkers := config.MqttConfig.WriteWorkers
	for i := 0; i < writeWorkers; i++ {
		go MessagesChanHandler(TelemetryMessagesChan)
	}
}

// 订阅telemetry消息
func SubscribeTelemetry() error {
	p, err := ants.NewPool(config.MqttConfig.Telemetry.PoolSize)
	if err != nil {
		return err
	}
	deviceTelemetryMessageHandler := func(_ mqtt.Client, d mqtt.Message) {
		err = p.Submit(func() {
			// 处理消息
			TelemetryMessages(d.Payload(), d.Topic())
		})
		if err != nil {
			logrus.Error(err)
		}
	}

	topic := config.MqttConfig.Telemetry.SubscribeTopic
	logrus.Debug("subscribe topic:", topic)
	qos := byte(config.MqttConfig.Telemetry.QoS)
	if token := SubscribeMqttClient.Subscribe(topic, qos, deviceTelemetryMessageHandler); token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
		return err
	}
	return nil
}

// 订阅设备Register消息，暂不需要线程池，不需要消息队列
func SubscribeRegister() error {
	// 订阅command消息
	deviceRegisterHandler := func(_ mqtt.Client, d mqtt.Message) {
		// 处理消息
		logrus.Debugf("[MQTT] \n Topic:%s \n payload:%s", d.Topic(), string(d.Payload()))
		RegisterMessages(d.Payload(), d.Topic())
	}
	topic := config.MqttConfig.Register.SubscribeTopic
	logrus.Debug("subscribe topic:", topic)
	qos := byte(config.MqttConfig.Commands.QoS)
	if token := SubscribeMqttClient.Subscribe(topic, qos, deviceRegisterHandler); token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
		return token.Error()
	}
	return nil
}

// 订阅设备Register消息，暂不需要线程池，不需要消息队列
func SubscribeCustomer() error {
	// 订阅attribute消息
	deviceCustomerHandler := func(_ mqtt.Client, d mqtt.Message) {
		// 处理消息
		logrus.Debugf("[MQTT]\n Topic:%s \n payload:%s", d.Topic(), string(d.Payload()))
		var devID, cfgID string
		topicList := strings.Split(d.Topic(), "/")
		if len(topicList) < 5 {
			devID = ""
		} else {
			cfgID = topicList[3]
			devID = topicList[4]
		}

		logrus.Debug("NTP 请求, cfgId:", cfgID, "devID:", devID)
		if cfgID != "" && devID != "" {
			// 响应设备NTP请求
			publish.PublishNtpResponseMessage(cfgID, devID, nil)
		}
	}
	topic := config.MqttConfig.Customer.SubscribeTopic
	//topic = GenTopic(topic)
	logrus.Info("subscribe topic:", topic)
	qos := byte(config.MqttConfig.Customer.QoS)
	if token := SubscribeMqttClient.Subscribe(topic, qos, deviceCustomerHandler); token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
		return token.Error()
	}
	return nil
}

// 订阅attribute消息，暂不需要线程池，不需要消息队列
func SubscribeAttribute() error {
	// 订阅attribute消息
	deviceAttributeHandler := func(_ mqtt.Client, d mqtt.Message) {
		// 处理消息
		logrus.Debugf("[MQTT]\n Topic:%s \n payload:%s", d.Topic(), string(d.Payload()))
		_, err := DeviceAttributeReport(d.Payload(), d.Topic())
		//logrus.Debug("响应设备属性上报", deviceNumber, err)
		if err != nil {
			logrus.Error(err)
		}
	}
	topic := config.MqttConfig.Attributes.SubscribeTopic
	logrus.Info("subscribe topic:", topic)
	qos := byte(config.MqttConfig.Attributes.QoS)
	if token := SubscribeMqttClient.Subscribe(topic, qos, deviceAttributeHandler); token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
		return token.Error()
	}
	return nil
}

// 订阅command消息，暂不需要线程池，不需要消息队列
func SubscribeCommand() error {
	// 订阅command消息
	deviceCommandHandler := func(_ mqtt.Client, d mqtt.Message) {
		// 处理消息
		logrus.Debugf("[MQTT]\n Topic:%s \n payload:%s", d.Topic(), string(d.Payload()))
		messageID, err := DeviceCommand(d.Payload(), d.Topic())
		if err != nil || messageID == "" {
			logrus.Debug("设备命令响应上报失败", messageID, err)
			logrus.Error(err)
		}
	}
	topic := config.MqttConfig.Commands.SubscribeTopic
	topic = GenTopic(topic)
	logrus.Info("subscribe topic:", topic)
	qos := byte(config.MqttConfig.Commands.QoS)
	if token := SubscribeMqttClient.Subscribe(topic, qos, deviceCommandHandler); token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
		return token.Error()
	}
	return nil
}

// 订阅event消息，暂不需要线程池，不需要消息队列
func SubscribeEvent() error {
	// 订阅event消息
	deviceEventHandler := func(_ mqtt.Client, d mqtt.Message) {
		// 处理消息
		logrus.Debugf("[MQTT]\n Topic:%s \n payload:%s", d.Topic(), string(d.Payload()))
		deviceNumber, method, err := DeviceEvent(d.Payload(), d.Topic())
		logrus.Debug("响应设备属性上报", deviceNumber, method, err)
		if err != nil {
			logrus.Error(err)
		}
	}
	topic := config.MqttConfig.Events.SubscribeTopic
	qos := byte(config.MqttConfig.Events.QoS)
	if token := SubscribeMqttClient.Subscribe(topic, qos, deviceEventHandler); token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
		return token.Error()
	}
	return nil
}

func SubscribeOtaUpprogress() error {
	// 订阅ota升级消息
	otaUpgradeHandler := func(_ mqtt.Client, d mqtt.Message) {
		// 处理消息
		logrus.Debugf("[MQTT]\n Topic:%s \n payload:%s", d.Topic(), string(d.Payload()))
		OtaUpgrade(d.Payload(), d.Topic())
	}
	topic := config.MqttConfig.OTA.SubscribeTopic
	qos := byte(config.MqttConfig.OTA.QoS)
	if token := SubscribeMqttClient.Subscribe(topic, qos, otaUpgradeHandler); token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
		return token.Error()
	}
	return nil
}
