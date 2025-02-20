package subscribe

import (
	"encoding/json"
	initialize "project/initialize"
	"project/internal/model"
	"project/pkg/utils"

	service "project/internal/service"

	"github.com/sirupsen/logrus"
)

var AccessWay = "A"
var TenantID = "d616bcbb"

type RegisterSubSt struct {
	CfgID string `json:"cfgID"`
	Mac   string `json:"mac"`
	Name  string `json:"name,omitempty"`
}

func RegisterMessages(payload []byte, topic string) {
	logrus.Debugln(string(payload))

	//byte转map
	var regMsg RegisterSubSt
	err := json.Unmarshal(payload, &regMsg)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	device, err := initialize.GetDeviceCacheById(regMsg.Mac)
	if device != nil && device.ID != "" {
		logrus.Warnf("deviceID:%s is exist,do not create.", regMsg.Mac)
		return
	}
	var createDevReq model.CreateDeviceReq
	createDevReq.AccessWay = &AccessWay
	if regMsg.Name != "" {
		createDevReq.Name = &regMsg.Name
	} else {
		createDevReq.Name = &regMsg.Mac
	}

	createDevReq.DeviceConfigId = &regMsg.CfgID
	createDevReq.DeviceNumber = &regMsg.Mac
	createDevReq.Label = &regMsg.Mac

	var claims utils.UserClaims
	claims.TenantID = TenantID
	data, err := service.GroupApp.Device.CreateDevice(createDevReq, &claims)
	if err != nil {
		logrus.Debugln("Create Device Error:")
		return
	}
	logrus.Debugln("create Successed:", data.DeviceNumber)
}
