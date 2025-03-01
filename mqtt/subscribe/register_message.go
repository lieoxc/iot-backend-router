package subscribe

import (
	"encoding/json"
	initialize "project/initialize"
	dal "project/internal/dal"
	"project/internal/model"
	"project/pkg/utils"

	"project/internal/private_register"
	service "project/internal/service"

	"github.com/sirupsen/logrus"
)

var AccessWay = "A"
var GlobalGatewayID = ""

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
	claims.TenantID = model.DefaultTenantId
	data, err := service.GroupApp.Device.CreateDevice(createDevReq, &claims)
	if err != nil {
		logrus.Debugln("Create Device Error:")
		return
	}
	logrus.Debugln("create Successed:", data.DeviceNumber)

	//本地注册完后，需要去内网监控站注册
	//TODO 因该还需要检查页面是否有配置私有服务器地址
	registerToPrivateGateway(createDevReq)
	//TODO 还需要注册到公网服务器
}
func registerToPrivateGateway(device model.CreateDeviceReq) {
	if *device.DeviceConfigId == model.DefaultGatewayCfgID {
		// 网关注册
		err := private_register.GatewayRegister(*device.DeviceNumber)
		if err != nil {
			logrus.Error(err.Error())
		}
	}
	// 子设备注册
	gatewayID := GetGatewayID()
	if gatewayID == "" {
		logrus.Error("Failed to get gatewayID")
		return
	}
	// 构造字设备注册结构体
	var regsiterReq model.DeviceRegisterReq
	regsiterReq.DeviceId = gatewayID
	regsiterReq.Type = ""
	items := make([]model.DeviceSubItem, 0, 1)
	item := model.DeviceSubItem{
		SubAddr: *device.DeviceNumber,
		Model:   model.DefaultESP32CfgName,
	}
	items = append(items, item)

	regsiterReq.Registers = items
	private_register.SubDevRegister(regsiterReq)
}

// 获取气象站网关ID，用于注册字设备到私有服务器
func GetGatewayID() string {
	if GlobalGatewayID != "" {
		return GlobalGatewayID
	}
	var cfgID = model.DefaultGatewayCfgID
	req := model.GetDeviceListByPageReq{
		DeviceConfigId: &cfgID,
	}
	total, list, err := dal.GetDeviceListByPage(&req, model.DefaultTenantId)
	if err != nil {
		logrus.Error("Failed to get device list by cfgID:", model.DefaultGatewayCfgID)
		return ""
	}
	if total == 1 {
		GlobalGatewayID = list[0].ID
		return GlobalGatewayID
	}
	logrus.Error("dev Number error:", model.DefaultGatewayCfgID)
	return ""
}
