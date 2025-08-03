package private_register

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"project/internal/dal"
	model "project/internal/model"
	"project/mqtt_private"

	"github.com/sirupsen/logrus"
)

type HTTPRes struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func GatewayRegister(devID string) error {
	addr := mqtt_private.Webaddr + ":9999"
	// 构造HTTP GET请求URL
	url := fmt.Sprintf("http://%s/api/v1/device/gateway-register", addr)

	// 构造请求数据
	reqData := model.GatewayRegisterReq{
		GatewayId: devID, // 将设备 ID 放入请求数据中
		TenantId:  model.DefaultTenantId,
		Model:     model.DefaultGatewayCfgName,
		Name:      mqtt_private.HostName,
		Version:   "V1.0.0",
	}

	// 将请求数据编码为 JSON
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		logrus.Println("Failed to marshal request data:", err)
		return err
	}
	logrus.Println("gatewayRegister Request URL:", url, " data:", string(reqBody))
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
func SubDevRegister(regsiterReq model.DeviceRegisterReq) error {
	addr := mqtt_private.Webaddr + ":9999"
	// 构造HTTP GET请求URL
	url := fmt.Sprintf("http://%s/api/v1/device/gateway-sub-register", addr)

	// 将请求数据编码为 JSON
	reqBody, err := json.Marshal(regsiterReq)
	if err != nil {
		logrus.Println("Failed to marshal request data:", err)
		return err
	}
	logrus.Println("gatewayRegister Request URL:", url, " data:", string(reqBody))
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
	var cfgID = model.DefaultGatewayCfgID
	req := model.GetDeviceListByPageReq{
		DeviceConfigId: &cfgID,
	}
	total, list, err := dal.GetDeviceListByPage(&req, model.DefaultTenantId)
	if err != nil {
		logrus.Error("Failed to get device list by cfgID:", model.DefaultGatewayCfgID)
		return "", err
	}
	if total == 1 {
		GatewayRegister(list[0].ID)
	} else {
		logrus.Error("dev Number error:", model.DefaultGatewayCfgID)
		return "", fmt.Errorf("dev Number error")
	}
	return list[0].ID, nil
}

func RegisterSubDev(gatewayID string) error {
	// 注册网关设备
	// 1. 查找到气象站设备的ID
	var cfgID = model.DefaultESP32CfgID
	deviceListReq := model.GetDeviceListByPageReq{
		PageReq:        model.PageReq{Page: 1, PageSize: 1000},
		DeviceConfigId: &cfgID,
	}
	total, list, err := dal.GetDeviceListByPage(&deviceListReq, model.DefaultTenantId)
	if err != nil {
		logrus.Error("Failed to get device list by cfgID:", model.DefaultGatewayCfgID)
		return err
	}
	var regsiterReq model.DeviceRegisterReq
	regsiterReq.DeviceId = gatewayID
	regsiterReq.Type = ""
	items := make([]model.DeviceSubItem, 0, total)
	for _, dev := range list {
		item := model.DeviceSubItem{
			SubAddr: dev.ID,
			Model:   dev.DeviceConfigName,
			Name:    dev.Name,
			Version: dev.CurrentVersion,
		}
		items = append(items, item)
	}
	regsiterReq.Registers = items
	SubDevRegister(regsiterReq)
	return nil
}

func PrivateRegisterInit() {
	gatewayID, err := RegisterGateway()
	if err != nil {
		logrus.Error("Failed to get device list by cfgID:", model.DefaultGatewayCfgID)
		return
	}
	RegisterSubDev(gatewayID)
}
