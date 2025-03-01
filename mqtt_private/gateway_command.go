package mqtt_private

import (
	"project/initialize"
	"project/internal/model"
	"project/mqtt/publish"
	"strings"

	"encoding/json"

	"github.com/sirupsen/logrus"
)

// topic: gateway/command/cfgID/devID/{message_id
func GatewayDeviceCommand(payload []byte, topic string) error {
	logrus.Debugln("gateway DeviceCommand", string(payload))
	// 验证消息有效性
	datas := strings.Split(string(topic), "/")
	if len(datas) < 5 {
		logrus.Error("Invalid topic:", topic)
		return nil
	}

	messageID, gatewayID := datas[4], datas[3]
	logrus.Debugln("gateway dev:", gatewayID, " messageID:", messageID)
	//校验一下网关是否存在
	_, err := initialize.GetDeviceCacheById(gatewayID)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	// 解析消息
	var command model.GatewayPublish
	err = json.Unmarshal(payload, &command)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	if command.SubDeviceData == nil {
		logrus.Warn("SubDeviceData is nil, has not command data !!!!!")
		return nil
	}
	// 子设备命令
	for k, v := range *command.SubDeviceData {
		logrus.Debugln("sub devID:", k)
		devID := k
		// 子设备命令处理
		subCommandTopic := PrivateMqttConfig.Commands.PublishTopic + "/" + model.DefaultESP32CfgID + "/" + devID + "/" + messageID
		payload, err := json.Marshal(&v)
		if err != nil {
			logrus.Error(err.Error())
			return err
		}
		err = publish.PublishCommandMessage(subCommandTopic, payload)
		if err != nil {
			logrus.Error("gateway PublishCommand err", err)
		}

	}
	return nil
}
