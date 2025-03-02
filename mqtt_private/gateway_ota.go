package mqtt_private

import (
	"encoding/json"
	"project/initialize"
	"project/internal/model"
	"project/mqtt/publish"
	"strings"

	"github.com/sirupsen/logrus"
)

// topic: ota/gateway/infrom/cfgID/devID
func GatewayDeviceOta(payload []byte, topic string) error {
	logrus.Debugln("gateway Device Ota", string(payload))
	// 验证消息有效性
	datas := strings.Split(string(topic), "/")
	if len(datas) < 5 {
		logrus.Error("Invalid topic:", topic)
		return nil
	}

	cfgID, gatewayID := datas[3], datas[4]
	logrus.Debugln("gateway cfgID:", cfgID, " gatewayID:", gatewayID)
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
		payload, err := json.Marshal(&v)
		if err != nil {
			logrus.Error(err.Error())
			return err
		}
		err = publish.PublishOtaAdress(model.DefaultESP32CfgID, devID, payload)
		if err != nil {
			logrus.Error("gateway PublishCommand err", err)
		}
	}
	return nil
}
