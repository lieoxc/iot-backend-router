package mqtt_private

import (
	"fmt"
	"project/internal/model"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

var privateMqttClient mqtt.Client

func createMqttClient() {
	// 初始化配置
	opts := mqtt.NewClientOptions()
	opts.AddBroker(PrivateMqttConfig.Broker)
	opts.SetUsername(PrivateMqttConfig.User)
	opts.SetPassword(PrivateMqttConfig.Pass)
	opts.SetClientID("weather-Station")
	// 干净会话
	opts.SetCleanSession(true)
	// 恢复客户端订阅，需要broker支持
	opts.SetResumeSubs(true)
	// 自动重连
	opts.SetAutoReconnect(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetMaxReconnectInterval(20 * time.Second)
	// 消息顺序
	opts.SetOrderMatters(false)
	opts.SetOnConnectHandler(func(_ mqtt.Client) {
		logrus.Debug("mqtt connect success")
	})
	// 断线重连
	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		logrus.Error("mqtt connect  lost: ", err)
		privateMqttClient.Disconnect(250)
		// 等待连接成功，失败重新连接
		for {
			token := privateMqttClient.Connect()
			if token.Wait() && token.Error() == nil {
				fmt.Println("Reconnected to MQTT broker")
				break
			}
			fmt.Printf("Reconnect failed: %v\n", token.Error())
			time.Sleep(5 * time.Second)
		}
	})

	privateMqttClient = mqtt.NewClient(opts)
	for {
		if token := privateMqttClient.Connect(); token.Wait() && token.Error() != nil {
			logrus.Error("MQTT Broker 1 连接失败:", token.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
}
func subscribe() error {
	err := SubscribeCommand()
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

// 订阅内网监控站下发的 命令消息，通过网关进行转发
// TODO 目前只做了中转，没有做命令的响应处理
func SubscribeCommand() error {
	// 订阅command消息
	commandHandler := func(_ mqtt.Client, d mqtt.Message) {
		logrus.Debug("command message received")
		err := GatewayDeviceCommand(d.Payload(), d.Topic())
		if err != nil {
			logrus.Error("private forward comman err:", err)
		}
	}
	// topic: /gateway/command/cfgID/gatewayID/+
	topic := PrivateMqttConfig.Commands.GatewaySubscribeTopic + "/" + model.DefaultGatewayCfgID + "/" + gatewayID + "/+"
	logrus.Debug("mqtt private subscribe topic:", topic)
	qos := byte(PrivateMqttConfig.Commands.QoS)
	if token := privateMqttClient.Subscribe(topic, qos, commandHandler); token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
		return token.Error()
	}
	return nil
}
