package mqtt_private

import (
	"encoding/json"
	"project/internal/model"

	"github.com/sirupsen/logrus"
)

func makeGatewayPubPayload(cfgID, devID string, payload []byte) ([]byte, error) {
	var dataMap = make(map[string]interface{})
	err := json.Unmarshal(payload, &dataMap)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	gatewayData := make(map[string]interface{})
	subDeviceData := make(map[string]map[string]interface{})

	if cfgID == model.DefaultGatewayCfgID { // 网关自身
		gatewayData = dataMap
	} else { // 子设备
		subDeviceData[devID] = dataMap
	}
	// 创建结构体实例
	publish := model.GatewayPublish{
		GatewayData:   &gatewayData,
		SubDeviceData: &subDeviceData,
	}
	// 序列化为 JSON
	jsonData, err := json.MarshalIndent(publish, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func makeGatewayEventPayload(cfgID, devID string, payload []byte) ([]byte, error) {
	var dataMap model.EventInfo
	err := json.Unmarshal(payload, &dataMap)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	var gatewayData model.EventInfo
	subDeviceData := make(map[string]model.EventInfo)

	if cfgID == model.DefaultGatewayCfgID { // 网关自身
		gatewayData = dataMap
	} else { // 子设备
		subDeviceData[devID] = dataMap
	}
	// 创建结构体实例
	publish := model.GatewayCommandPulish{
		GatewayData:   &gatewayData,
		SubDeviceData: &subDeviceData,
	}
	// 序列化为 JSON
	jsonData, err := json.MarshalIndent(publish, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

// 上报telemetry消息
func ForwardTelemetryMessage(cfgID, devID string, payload []byte) error {
	if PrivateMqttClient == nil || !PrivateMqttClient.IsConnected() || gatewayID == "" {
		logrus.Debug("privateMqttClient is not connected")
		return nil
	}
	qos := byte(PrivateMqttConfig.Telemetry.QoS)
	pubTelemetryTopic := PrivateMqttConfig.Telemetry.GatewayPublishTopic + "/" + model.DefaultGatewayCfgID + "/" + gatewayID
	logrus.Debug("topic:", pubTelemetryTopic, "value:", string(payload))
	// 转发消息
	jsonData, err := makeGatewayPubPayload(cfgID, devID, payload)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	logrus.Debugf("privateMqttClient pub topic:%v payload%v", pubTelemetryTopic, string(jsonData))
	token := PrivateMqttClient.Publish(pubTelemetryTopic, qos, false, jsonData)
	if token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
	}
	return token.Error()
}

// 上报attributes消息
func ForwardAttributesMessage(cfgID, devID string, payload []byte) error {
	if PrivateMqttClient == nil || !PrivateMqttClient.IsConnected() || gatewayID == "" {
		logrus.Debug("privateMqttClient is not connected")
		return nil
	}
	qos := byte(PrivateMqttConfig.Telemetry.QoS)
	pubAttributesTopic := PrivateMqttConfig.Attributes.GatewayPublishTopic + "/" + model.DefaultGatewayCfgID + "/" + gatewayID
	logrus.Debug("topic:", pubAttributesTopic, "value:", string(payload))
	// 发布消息
	jsonData, err := makeGatewayPubPayload(cfgID, devID, payload)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	logrus.Debugf("privateMqttClient pub topic:%v payload%v", pubAttributesTopic, string(jsonData))

	token := PrivateMqttClient.Publish(pubAttributesTopic, qos, false, jsonData)
	if token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
	}
	return token.Error()
}
func ForwardCommandRespMessage(cfgID, devID, msgID string, payload []byte) error {
	if PrivateMqttClient == nil || !PrivateMqttClient.IsConnected() || gatewayID == "" {
		logrus.Debug("privateMqttClient is not connected")
		return nil
	}
	qos := byte(PrivateMqttConfig.Telemetry.QoS)
	pubCmdTopic := PrivateMqttConfig.Commands.GatewayPublishTopic + "/" + cfgID + "/" + devID + "/" + msgID
	logrus.Debug("topic:", pubCmdTopic, "  value:", string(payload))

	token := PrivateMqttClient.Publish(pubCmdTopic, qos, false, payload)
	if token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
	}
	return token.Error()
}

// 上报tevents消息
func ForwardEventsMessage(cfgID, devID string, payload []byte) error {
	if PrivateMqttClient == nil || !PrivateMqttClient.IsConnected() || gatewayID == "" {
		logrus.Debug("privateMqttClient is not connected")
		return nil
	}
	qos := byte(PrivateMqttConfig.Telemetry.QoS)
	pubEventsTopic := PrivateMqttConfig.Events.GatewayPublishTopic + "/" + model.DefaultGatewayCfgID + "/" + gatewayID
	logrus.Debug("topic:", pubEventsTopic, "value:", string(payload))
	// 发布消息
	jsonData, err := makeGatewayEventPayload(cfgID, devID, payload)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	logrus.Debugf("privateMqttClient pub topic:%v payload%v", pubEventsTopic, string(jsonData))
	token := PrivateMqttClient.Publish(pubEventsTopic, qos, false, jsonData)
	if token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
	}
	return token.Error()
}

// 上报ota消息
func ForwardOtaMessage(cfgID, devID string, payload []byte) error {
	if PrivateMqttClient == nil || !PrivateMqttClient.IsConnected() || gatewayID == "" {
		logrus.Debug("privateMqttClient is not connected")
		return nil
	}
	qos := byte(PrivateMqttConfig.OTA.QoS)
	pubOtaTopic := PrivateMqttConfig.OTA.GatewayPublishTopic + "/" + model.DefaultGatewayCfgID + "/" + gatewayID
	// 发布消息
	jsonData, err := makeGatewayPubPayload(cfgID, devID, payload)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	logrus.Debugf("privateMqttClient pub topic:%v payload%v", pubOtaTopic, string(jsonData))
	token := PrivateMqttClient.Publish(pubOtaTopic, qos, false, jsonData)
	if token.Wait() && token.Error() != nil {
		logrus.Error(token.Error())
	}
	return token.Error()
}
