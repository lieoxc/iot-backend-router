package private_register

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"project/internal/dal"
	model "project/internal/model"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	DefaultTenantId       = "d616bcbb"
	DefaultGatewayCfgName = "气象站"
	DefaultESP32CfgName   = "ESP32"
)

var (
	DefaultGatewayCfgID = "964d6220-ecbf-a043-1960-85b1a2758cea"
	DefaultESP32CfgID   = "315d9d82-5c76-3197-4eab-8c0a641ccdc9"
)

type HTTPRes struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func gatewayRegister(devID string) error {
	addr := viper.GetString("web.addr")
	if addr == "" {
		addr = "localhost:9999"
		logrus.Println("Using default broker:", addr)
	}
	// 构造HTTP GET请求URL
	url := fmt.Sprintf("http://%s/api/v1/device/gateway-register", addr)
	logrus.Println("gatewayRegister Request URL:", url)

	// 构造请求数据
	reqData := model.GatewayRegisterReq{
		GatewayId: devID, // 将设备 ID 放入请求数据中
		TenantId:  DefaultTenantId,
		Model:     DefaultGatewayCfgName,
	}

	// 将请求数据编码为 JSON
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		logrus.Println("Failed to marshal request data:", err)
		return err
	}

	// 发送 HTTP POST 请求
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		logrus.Println("Failed to send POST request:", err)
		return err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Println("Failed to read response body:", err)
		return err
	}
	logrus.Debug(string(body))
	var respData HTTPRes
	err = json.Unmarshal(body, &respData)
	if err != nil {
		logrus.Println("Failed to json.Unmarshal:", err)
		return err
	}
	logrus.Debugf("MqttClientId:%v", respData.Data)
	return nil
}
func subDevRegister(regsiterReq model.DeviceRegisterReq) error {
	addr := viper.GetString("web.addr")
	if addr == "" {
		addr = "localhost:9999"
		logrus.Println("Using default broker:", addr)
	}
	// 构造HTTP GET请求URL
	url := fmt.Sprintf("http://%s/api/v1/device/gateway-sub-register", addr)
	logrus.Println("sub dev Register Request URL:", url)

	// 将请求数据编码为 JSON
	reqBody, err := json.Marshal(regsiterReq)
	if err != nil {
		logrus.Println("Failed to marshal request data:", err)
		return err
	}

	// 发送 HTTP POST 请求
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		logrus.Println("Failed to send POST request:", err)
		return err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Println("Failed to read response body:", err)
		return err
	}
	logrus.Debug(string(body))
	var respData model.DeviceRegisterRes
	err = json.Unmarshal(body, &respData)
	if err != nil {
		logrus.Println("Failed to json.Unmarshal:", err)
		return err
	}
	logrus.Debugf("subDevRegister Status:%v", respData.Status)
	return nil
}
func RegisterGateway() (string, error) {
	// 注册网关设备
	// 1. 查找到气象站设备的ID
	req := model.GetDeviceListByPageReq{
		DeviceConfigId: &DefaultGatewayCfgID,
	}
	total, list, err := dal.GetDeviceListByPage(&req, DefaultTenantId)
	if err != nil {
		logrus.Error("Failed to get device list by cfgID:", DefaultGatewayCfgID)
		return "", err
	}
	if total == 1 {
		gatewayRegister(list[0].ID)
	} else {
		logrus.Error("dev Number error:", DefaultGatewayCfgID)
		return "", fmt.Errorf("dev Number error")
	}
	return list[0].ID, nil
}

func RegisterSubDev(gatewayID string) error {
	// 注册网关设备
	// 1. 查找到气象站设备的ID
	deviceListReq := model.GetDeviceListByPageReq{
		PageReq:        model.PageReq{Page: 1, PageSize: 1000},
		DeviceConfigId: &DefaultESP32CfgID,
	}
	total, list, err := dal.GetDeviceListByPage(&deviceListReq, DefaultTenantId)
	if err != nil {
		logrus.Error("Failed to get device list by cfgID:", DefaultGatewayCfgID)
		return err
	}
	var regsiterReq model.DeviceRegisterReq
	regsiterReq.DeviceId = gatewayID
	regsiterReq.Type = ""
	items := make([]model.DeviceSubItem, 0, total)
	for _, dev := range list {
		item := model.DeviceSubItem{
			SubAddr: dev.ID,
			Model:   DefaultESP32CfgName,
		}
		items = append(items, item)
	}
	regsiterReq.Registers = items
	subDevRegister(regsiterReq)
	return nil
}

func PrivateRegisterInit() {
	gatewayID, err := RegisterGateway()
	if err != nil {
		logrus.Error("Failed to get device list by cfgID:", DefaultGatewayCfgID)
		return
	}
	RegisterSubDev(gatewayID)
}
