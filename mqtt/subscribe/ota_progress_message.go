package subscribe

import (
	"encoding/json"
	initialize "project/initialize"
	"project/internal/query"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type DeviceProgressMsg struct {
	UpgradeProgress interface{} `json:"step,omitempty" alias:"进度"`
	StatusDetail    string      `json:"desc" alias:"描述"`
	//Module          string      `json:"module,omitempty" alias:"模块"`
	//UpgradeStatus    string      `json:"upgrade_status,omitempty"`
	//StatusUpdateTime string      `json:"status_update_time" alias:"升级更新时间"`
}

// 接收OTA升级进度消息
func OtaUpgrade(payload []byte, topic string) {
	/*
		消息规范：topic:ota/devices/progress/315d9d82-5c76-3197-4eab-8c0a641ccdc9/112233445566
				 payload是json格式的消息
				 {"step":"100","desc":"升级进度100%"}
				 {"step":"-1","desc":"OTA升级失败，请求不到升级包信息。"}
	*/
	// 验证消息有效性
	datas := strings.Split(string(topic), "/")
	if len(datas) != 5 {
		logrus.Error("ota msg topic length error")
		return
	}

	cfgID, devID := datas[3], datas[4]
	logrus.Debugln("cfgID:", cfgID, " mdevID:", devID)

	// 处理消息
	device, err := initialize.GetDeviceCacheById(devID)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	var progressMsg DeviceProgressMsg
	err = json.Unmarshal(payload, &progressMsg)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	switch progressMsg.UpgradeProgress.(type) {
	case float64: // json转的数值 100%为 float64
		// 数值转文本
		progressMsg.UpgradeProgress = strconv.FormatInt(int64(progressMsg.UpgradeProgress.(float64)), 10) // 抛弃小数

	case string:
		// 直接越过
	default:
		logrus.Error("不支持的数据类型")
		return
	}
	logrus.Debugf("UpgradeProgress:%s StatusDetail:%s", progressMsg.UpgradeProgress, progressMsg.StatusDetail)
	// 查询对应设备升级信息
	otaTaskDetail, err := query.OtaUpgradeTaskDetail.
		Where(query.OtaUpgradeTaskDetail.DeviceID.Eq(device.ID),
			query.OtaUpgradeTaskDetail.Status.In(2, 3),
		).First()
	if err != nil && otaTaskDetail != nil {
		logrus.Errorf("未找到对应升级任务")
		return
	}

	intProgress, err := strconv.Atoi(progressMsg.UpgradeProgress.(string))
	if err != nil {
		desc := progressMsg.UpgradeProgress.(string) + " " + progressMsg.StatusDetail
		otaTaskDetail.StatusDescription = &desc
	}
	var transferValue int16
	switch {
	case intProgress == -1:
		desc := "错误码-1,升级失败 " + progressMsg.StatusDetail
		otaTaskDetail.Status = 5
		otaTaskDetail.StatusDescription = &desc
	case intProgress == -2:
		desc := "错误码-2,下载失败 " + progressMsg.StatusDetail
		otaTaskDetail.Status = 5
		otaTaskDetail.StatusDescription = &desc
	case intProgress == -3:
		desc := "错误码-3,校验失败 " + progressMsg.StatusDetail
		otaTaskDetail.Status = 5
		otaTaskDetail.StatusDescription = &desc
	case intProgress == -4:
		desc := "错误码-4,烧写失败 " + progressMsg.StatusDetail
		otaTaskDetail.Status = 5
		otaTaskDetail.StatusDescription = &desc
	case intProgress >= 1 && intProgress < 100:
		otaTaskDetail.Status = 3
		transferValue = int16(intProgress)
		otaTaskDetail.Step = &transferValue
		otaTaskDetail.StatusDescription = &progressMsg.StatusDetail
	case intProgress == 100:
		otaTaskDetail.Status = 4
		transferValue = int16(intProgress)
		otaTaskDetail.Step = &transferValue
		otaTaskDetail.StatusDescription = &progressMsg.StatusDetail
	default:
		logrus.Error("数据格式有问题")
		return
	}
	_, err = query.OtaUpgradeTaskDetail.Where(query.OtaUpgradeTaskDetail.ID.Eq(otaTaskDetail.ID)).Updates(otaTaskDetail)
	if err != nil {
		logrus.Error(err)
		return
	}
}
