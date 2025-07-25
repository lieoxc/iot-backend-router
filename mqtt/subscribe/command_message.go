package subscribe

import (
	"encoding/json"
	"project/internal/model"
	config "project/mqtt"
	"strings"

	"github.com/sirupsen/logrus"
)

// 接收设备命令的响应消息
func DeviceCommand(payload []byte, topic string) (string, error) {
	/*
		消息规范：topic:devices/command/response/+/+/+
				 +是cfgID/devID/msgID
				 payload是json格式的命令消息
	*/
	// 验证消息有效性
	// TODO处理消息
	datas := strings.Split(string(topic), "/")
	if len(datas) != 6 {
		logrus.Error("commamd response msg topic length error")
		return "", nil
	}

	cfgID, devID, messageId := datas[3], datas[4], datas[5]
	logrus.Debugln("cfgID:", cfgID, " mdevID:", devID)

	logrus.Debug("command response message:", string(payload))
	var responseMsg model.MqttResponse
	err := json.Unmarshal(payload, &responseMsg)
	if err != nil {
		logrus.Error(err.Error())
		return "", err
	}
	if ch, ok := config.MqttDirectResponseFuncMap[messageId]; ok {
		ch <- responseMsg
	}
	return messageId, err
}
