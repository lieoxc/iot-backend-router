package service

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"project/initialize"
	config "project/mqtt"
	"project/mqtt/publish"
	simulationpublish "project/mqtt/simulation_publish"
	"project/pkg/constant"
	"project/pkg/errcode"
	"project/pkg/utils"
	"strconv"
	"strings"
	"time"

	"github.com/go-basic/uuid"
	"github.com/mintance/go-uniqid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"

	dal "project/internal/dal"
	model "project/internal/model"
)

type TelemetryData struct{}

func (*TelemetryData) GetCurrentTelemetrData(device_id string) (interface{}, error) {
	// d, err := dal.GetCurrentTelemetrData(device_id)
	// 数据源替换
	d, err := dal.GetCurrentTelemetryDataEvolution(device_id)
	if err != nil {
		return nil, err
	}

	// 查询设备信息
	deviceInfo, err := dal.GetDeviceByID(device_id)
	if err != nil {
		return nil, err
	}
	var telemetryModelMap = make(map[string]*model.DeviceModelTelemetry)
	var telemetryModelUintMap = make(map[string]interface{})
	var telemetryModelRWMap = make(map[string]interface{})
	// 是否有设备配置
	if deviceInfo.DeviceConfigID != nil {
		// 查询设备配置
		deviceConfig, err := dal.GetDeviceConfigByID(*deviceInfo.DeviceConfigID)
		if err != nil {
			return nil, err
		}
		// 是否有设备模板
		if deviceConfig.DeviceTemplateID != nil {
			// 查询遥测模型
			telemetryModel, err := dal.GetDeviceModelTelemetryDataList(*deviceConfig.DeviceTemplateID)
			if err != nil {
				return nil, err
			}
			if len(telemetryModel) > 0 {
				// 遍历并转换为map
				for _, v := range telemetryModel {
					telemetryModelMap[v.DataIdentifier] = v
					telemetryModelUintMap[v.DataIdentifier] = v.Unit
					telemetryModelRWMap[v.DataIdentifier] = v.ReadWriteFlag
				}
			}
		}
	}
	// 格式化返回值
	data := make([]map[string]interface{}, 0)
	if len(d) > 0 {
		for _, v := range d {
			tmp := make(map[string]interface{})
			tmp["device_id"] = v.DeviceID
			tmp["key"] = v.Key
			tmp["ts"] = v.T
			tmp["tenant_id"] = v.TenantID
			if v.BoolV != nil {
				tmp["value"] = v.BoolV
			}
			if v.NumberV != nil {
				tmp["value"] = v.NumberV
			}
			if v.StringV != nil {
				tmp["value"] = v.StringV
			}
			// 是否有设备模型
			if len(telemetryModelMap) > 0 {
				telemetryModel, ok := telemetryModelMap[v.Key]
				if ok {
					tmp["label"] = telemetryModel.DataName
					tmp["unit"] = telemetryModelUintMap[v.Key]
					tmp["read_write_flag"] = telemetryModelRWMap[v.Key]
					tmp["data_type"] = telemetryModel.DataType
					if telemetryModel.DataType != nil && *telemetryModel.DataType == "Enum" {
						var enumItems []model.EnumItem
						json.Unmarshal([]byte(*telemetryModel.AdditionalInfo), &enumItems)
						tmp["enum"] = enumItems
					}
				}
			}
			data = append(data, tmp)
		}
	}

	return data, err
}

// 根据设备ID和key获取当前遥测数据
func (*TelemetryData) GetCurrentTelemetrDataKeys(req *model.GetTelemetryCurrentDataKeysReq) (interface{}, error) {
	// d, err := dal.GetCurrentTelemetrData(device_id)
	// 数据源替换
	d, err := dal.GetCurrentTelemetryDataEvolutionByKeys(req.DeviceID, req.Keys)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	// 查询设备信息
	deviceInfo, err := dal.GetDeviceByID(req.DeviceID)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	var telemetryModelMap = make(map[string]*model.DeviceModelTelemetry)
	var telemetryModelUintMap = make(map[string]interface{})
	// 是否有设备配置
	if deviceInfo.DeviceConfigID != nil {
		// 查询设备配置
		deviceConfig, err := dal.GetDeviceConfigByID(*deviceInfo.DeviceConfigID)
		if err != nil {
			return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
				"sql_error": err.Error(),
			})
		}
		// 是否有设备模板
		if deviceConfig.DeviceTemplateID != nil {
			// 查询遥测模型
			telemetryModel, err := dal.GetDeviceModelTelemetryDataList(*deviceConfig.DeviceTemplateID)
			if err != nil {
				return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
					"sql_error": err.Error(),
				})
			}
			if len(telemetryModel) > 0 {
				// 遍历并转换为map
				for _, v := range telemetryModel {
					telemetryModelMap[v.DataIdentifier] = v
					telemetryModelUintMap[v.DataIdentifier] = v.Unit
				}
			}
		}
	}
	// 格式化返回值
	data := make([]map[string]interface{}, 0)
	if len(d) > 0 {
		for _, v := range d {
			tmp := make(map[string]interface{})

			tmp["device_id"] = v.DeviceID
			tmp["key"] = v.Key
			tmp["ts"] = v.T
			tmp["tenant_id"] = v.TenantID
			if v.BoolV != nil {
				tmp["value"] = v.BoolV
			}
			if v.NumberV != nil {
				tmp["value"] = v.NumberV
			}
			if v.StringV != nil {
				tmp["value"] = v.StringV
			}
			// 是否有设备模型
			if len(telemetryModelMap) > 0 {
				telemetryModel, ok := telemetryModelMap[v.Key]
				if ok {
					tmp["label"] = telemetryModel.DataName
					tmp["unit"] = telemetryModelUintMap[v.Key]
					tmp["data_type"] = telemetryModel.DataType
					if telemetryModel.DataType != nil && *telemetryModel.DataType == "Enum" {
						var enumItems []model.EnumItem
						json.Unmarshal([]byte(*telemetryModel.AdditionalInfo), &enumItems)
						tmp["enum"] = enumItems
					}
				}
			}
			data = append(data, tmp)
		}
	}

	return data, err
}

// 返回数据格式{"key":value,"key1":value1}
func (*TelemetryData) GetCurrentTelemetrDataForWs(device_id string) (interface{}, error) {
	// d, err := dal.GetCurrentTelemetrData(device_id)

	// 数据源替换
	d, err := dal.GetCurrentTelemetryDataEvolution(device_id)
	if err != nil {
		return nil, err
	}

	// 格式化返回值
	data := make(map[string]interface{})
	if len(d) > 0 {
		for _, v := range d {
			if v.BoolV != nil {
				data[v.Key] = v.BoolV
			}
			if v.NumberV != nil {
				data[v.Key] = v.NumberV
			}
			if v.StringV != nil {
				data[v.Key] = v.StringV
			}
		}
	}
	return data, err
}

// 返回数据格式{"key":value,"key1":value1}
func (*TelemetryData) GetCurrentTelemetrDataKeysForWs(device_id string, keys []string) (interface{}, error) {
	// d, err := dal.GetCurrentTelemetrData(device_id)

	// 数据源替换
	d, err := dal.GetCurrentTelemetryDataEvolutionByKeys(device_id, keys)
	if err != nil {
		return nil, err
	}

	// 格式化返回值
	data := make(map[string]interface{})
	if len(d) > 0 {
		for _, v := range d {
			if v.BoolV != nil {
				data[v.Key] = v.BoolV
			}
			if v.NumberV != nil {
				data[v.Key] = v.NumberV
			}
			if v.StringV != nil {
				data[v.Key] = v.StringV
			}
		}
	}
	return data, err
}

func (*TelemetryData) GetTelemetrHistoryData(req *model.GetTelemetryHistoryDataReq) (interface{}, error) {
	// 时间戳转换
	sT := req.StartTime * 1000
	eT := req.EndTime * 1000

	d, err := dal.GetHistoryTelemetrData(req.DeviceID, req.Key, sT, eT)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	// 格式化返回值
	data := make([]map[string]interface{}, 0)
	if len(d) > 0 {
		for _, v := range d {
			tmp := make(map[string]interface{})

			tmp["device_id"] = v.DeviceID
			tmp["key"] = v.Key
			tmp["ts"] = v.T
			tmp["tenant_id"] = v.TenantID
			if v.BoolV != nil {
				tmp["value"] = v.BoolV
			}
			if v.NumberV != nil {
				tmp["value"] = v.NumberV
			}
			if v.StringV != nil {
				tmp["value"] = v.StringV
			}
			data = append(data, tmp)
		}
	}

	return data, nil
}

func (*TelemetryData) DeleteTelemetrData(req *model.DeleteTelemetryDataReq) error {
	err := dal.DeleteTelemetrData(req.DeviceID, req.Key)
	if err != nil {
		return err
	}
	// 删除当前值
	err = dal.DeleteCurrentTelemetryData(req.DeviceID, req.Key)
	return err
}

func (*TelemetryData) GetCurrentTelemetrDetailData(device_id string) (interface{}, error) {
	data, err := dal.GetCurrentTelemetrDetailData(device_id)

	if err != nil {
		return nil, err
	}

	dataMap := make(map[string]interface{})

	dataMap["device_id"] = data.DeviceID
	dataMap["key"] = data.Key
	dataMap["ts"] = data.T
	dataMap["tenant_id"] = data.TenantID

	if data.BoolV != nil {
		dataMap["value"] = data.BoolV
	}

	if data.NumberV != nil {
		dataMap["value"] = data.NumberV
	}

	if data.StringV != nil {
		dataMap["value"] = data.StringV
	}

	return dataMap, err
}

func (*TelemetryData) GetTelemetrHistoryDataByPage(req *model.GetTelemetryHistoryDataByPageReq) (interface{}, error) {

	if *req.ExportExcel {
		var addr string
		f := excelize.NewFile()
		f.SetCellValue("Sheet1", "A1", "时间")
		f.SetCellValue("Sheet1", "B1", "数值")

		batchSize := 100000
		offset := 0
		rowNumber := 2

		for {
			datas, err := dal.GetHistoryTelemetrDataByExport(req, offset, batchSize)
			if err != nil {
				return addr, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
					"sql_error": err.Error(),
				})
			}
			if len(datas) == 0 {
				break
			}
			for _, data := range datas {
				t := time.Unix(0, data.T*int64(time.Millisecond))
				f.SetCellValue("Sheet1", fmt.Sprintf("A%d", rowNumber), t.Format("2006-01-02 15:04:05"))
				f.SetCellValue("Sheet1", fmt.Sprintf("B%d", rowNumber), *data.NumberV)
				rowNumber++
			}
			offset += batchSize
		}

		// 创建保存目录
		uploadDir := "./files/excel/"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return "", errcode.WithVars(errcode.CodeFilePathGenError, map[string]interface{}{
				"error": err.Error(),
			})
		}

		// 生成文件名
		uniqidStr := uniqid.New(uniqid.Params{
			Prefix:      "excel",
			MoreEntropy: true,
		})
		addr = "files/excel/数据列表" + uniqidStr + ".xlsx"

		// 保存文件
		if err := f.SaveAs(addr); err != nil {
			return "", errcode.WithVars(errcode.CodeFileSaveError, map[string]interface{}{
				"error": err.Error(),
			})
		}

		return addr, nil
	}

	//  暂时忽略总数
	_, data, err := dal.GetHistoryTelemetrDataByPage(req)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	// 格式化
	var easyData []map[string]interface{}
	for _, v := range data {
		d := make(map[string]interface{})
		d["ts"] = v.T
		d["key"] = v.Key
		if v.StringV != nil {
			d["value"] = v.StringV
		} else if v.NumberV != nil {
			d["value"] = v.NumberV
		} else if v.BoolV != nil {
			d["value"] = v.BoolV
		} else {
			d["value"] = ""
		}

		easyData = append(easyData, d)
	}
	return easyData, nil
}

func (*TelemetryData) GetTelemetrHistoryDataByPageV2(req *model.GetTelemetryHistoryDataByPageReq) (interface{}, error) {

	if req.ExportExcel != nil && *req.ExportExcel {
		var addr string
		f := excelize.NewFile()
		f.SetCellValue("Sheet1", "A1", "时间")
		f.SetCellValue("Sheet1", "B1", "数值")

		batchSize := 100000
		offset := 0
		rowNumber := 2

		for {
			datas, err := dal.GetHistoryTelemetrDataByExport(req, offset, batchSize)
			if err != nil {
				return addr, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
					"sql_error": err.Error(),
				})
			}
			if len(datas) == 0 {
				break
			}
			for _, data := range datas {
				t := time.Unix(0, data.T*int64(time.Millisecond))
				f.SetCellValue("Sheet1", fmt.Sprintf("A%d", rowNumber), t.Format("2006-01-02 15:04:05"))
				f.SetCellValue("Sheet1", fmt.Sprintf("B%d", rowNumber), *data.NumberV)
				rowNumber++
			}
			offset += batchSize
		}

		// 创建保存目录
		uploadDir := "./files/excel/"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return "", errcode.WithVars(errcode.CodeFilePathGenError, map[string]interface{}{
				"error": err.Error(),
			})
		}

		// 生成文件名
		uniqidStr := uniqid.New(uniqid.Params{
			Prefix:      "excel",
			MoreEntropy: true,
		})
		addr = "files/excel/数据列表" + uniqidStr + ".xlsx"

		// 保存文件
		if err := f.SaveAs(addr); err != nil {
			return "", errcode.WithVars(errcode.CodeFileSaveError, map[string]interface{}{
				"error": err.Error(),
			})
		}

		return addr, nil
	}

	//  暂时忽略总数
	total, data, err := dal.GetHistoryTelemetrDataByPage(req)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	// 格式化
	var easyData []map[string]interface{}
	for _, v := range data {
		d := make(map[string]interface{})
		d["ts"] = v.T
		d["key"] = v.Key
		if v.StringV != nil {
			d["value"] = v.StringV
		} else if v.NumberV != nil {
			d["value"] = v.NumberV
		} else if v.BoolV != nil {
			d["value"] = v.BoolV
		} else {
			d["value"] = ""
		}

		easyData = append(easyData, d)
	}
	dataRsp := make(map[string]interface{})
	dataRsp["total"] = total
	dataRsp["list"] = easyData
	return dataRsp, nil
}

// 获取模拟设备发送遥测数据的回显数据
func (*TelemetryData) ServeEchoData(req *model.ServeEchoDataReq) (interface{}, error) {
	// 获取设备信息
	deviceInfo, err := dal.GetDeviceByID(req.DeviceId)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	voucher := deviceInfo.Voucher
	// 校验voucher是否json
	if !IsJSON(voucher) {
		return nil, errcode.NewWithMessage(errcode.CodeParamError, "voucher is not json")
	}
	var voucherMap map[string]interface{}
	err = json.Unmarshal([]byte(voucher), &voucherMap)
	if err != nil {
		return nil, err
	}
	// 判断是否有username字段
	var username, password, host, post, payload, clientID string
	if _, ok := voucherMap["username"]; !ok {
		return nil, errcode.NewWithMessage(errcode.CodeParamError, "username is not exist")
	}
	username = voucherMap["username"].(string)
	// 判断是否有password字段
	if _, ok := voucherMap["password"]; !ok {
		password = ""
	} else {
		password = voucherMap["password"].(string)
	}

	accessAddress := viper.GetString("mqtt.access_address")
	if accessAddress == "" {
		return nil, errcode.NewWithMessage(errcode.CodeParamError, "mqtt access address is not exist")
	}
	accessAddressList := strings.Split(accessAddress, ":")
	host = accessAddressList[0]
	post = accessAddressList[1]
	topic := config.MqttConfig.Telemetry.SubscribeTopic
	clientID = "mqtt_" + uuid.New()[0:12] //代表随机生成
	payload = `{\"test_data1\":25.5,\"test_data2\":60}`
	// 拼接命令
	command := utils.BuildMosquittoPubCommand(host, post, username, password, topic, payload, clientID)
	return command, nil

}

// 模拟设备发送遥测数据
func (*TelemetryData) TelemetryPub(mosquittoCommand string) (interface{}, error) {
	// 解析mosquitto_pub命令
	params, err := utils.ParseMosquittoPubCommand(mosquittoCommand)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"error": err.Error(),
		})
	}
	// 根据凭证信息查询设备信息
	// 组装凭证信息
	var voucher string
	if params.Password == "" {
		voucher = fmt.Sprintf("{\"username\":\"%s\"}", params.Username)
	} else {
		voucher = fmt.Sprintf("{\"username\":\"%s\",\"password\":\"%s\"}", params.Username, params.Password)
	}
	// 查询设备信息
	deviceInfo, err := dal.GetDeviceByVoucher(voucher)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	var isOnline int
	if deviceInfo.IsOnline == int16(1) {
		isOnline = 1
	}

	// 发送mqtt消息
	logrus.Debug("params:", params)
	err = simulationpublish.PublishMessage(params.Host, params.Port, params.Topic, params.Payload, params.Username, params.Password, params.ClientId)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeSystemError, map[string]interface{}{
			"error": err.Error(),
		})
	}
	go func() {
		time.Sleep(3 * time.Second)
		// 更新设备状态
		if isOnline == 1 {
			dal.UpdateDeviceOnlineStatus(deviceInfo.ID, int16(isOnline))
			// 发送上线消息
			// 发送mqtt消息
			err = publish.PublishOnlineMessage(deviceInfo.ID, []byte("1"))
			if err != nil {
				logrus.Error("publish online message failed:", err)
			}
		}
	}()
	return nil, nil
}

func (*TelemetryData) GetTelemetrSetLogsDataListByPage(req *model.GetTelemetrySetLogsListByPageReq) (interface{}, error) {
	count, data, err := dal.GetTelemetrySetLogsListByPage(req)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	dataMap := make(map[string]interface{})
	dataMap["count"] = count
	dataMap["list"] = data
	return dataMap, nil

}

/*
 1. 部分参数说明：
    aggregate_window [聚合间隔]
    - no_aggregate 不聚合
    - "30s","1m","2m","5m","10m","30m","1h","3h","6h","1d","7d","1mo"
    time_range
    - 时间范围，后端支持的参数有：custom，last_5m，last_15m，last_30m，last_1h，last_3h 当选择自定义时，后端会根据开始和结束时间来判断是否超过3小时，如过超过3小时，则不能选择“不聚合”
    aggregate_function [聚合方法]
    - avg 平均数
    - max 最大值
 2. 前端筛选联动规则：
    - 页面初始化：最近1小时 - 不聚合 - 默认不展示计算方式，当选择了间隔后 展示两种计算方式（平均值/最大值）
    - 最近5分钟 - 展示全部聚合间隔
    - 最近15分钟 - 展示全部聚合间隔
    - 最近30分钟 - 展示全部聚合间隔
    - 最近1小时 - 展示全部聚合间隔
    - 最近3小时 - 间隔默认选择“30秒”（不聚合不可选） - 计算方式默认为 “平均值”
    - 最近6小时 - 间隔默认选择“1分钟”（不聚合，小于等于30秒的不可选） - 计算方式默认为 “平均值”
    - 最近12小时 - 间隔默认选择“2分钟”（不聚合，小于等于1分钟的不可选） - 计算方式默认为 “平均值”
    - 最近24小时 - 间隔默认选择“5分钟”（不聚合，小于等于2分钟的不可选） - 计算方式默认为 “平均值”
    - 最近3天 - 间隔默认选择“10分钟”（不聚合，小于等于5分钟的不可选） - 计算方式默认为 “平均值”
    - 最近7天 - 间隔默认选择“30分钟”（不聚合，小于等于10分钟的不可选） - 计算方式默认为 “平均值”
    - 最近15天 - 间隔默认选择“1小时”（不聚合，小于等于30分钟的不可选） - 计算方式默认为 “平均值”
    - 最近30天 - 间隔默认选择“1小时”（不聚合，小于等于30分钟的不可选） - 计算方式默认为 “平均值”
    - 最近60天 - 间隔默认选择“3小时”（不聚合，小于等于1小时的不可选） - 计算方式默认为 “平均值”
    - 最近90天 - 间隔默认选择“6小时”（不聚合，小于等于3小时的不可选） - 计算方式默认为 “平均值”
    - 最近6个月 - 间隔默认选择“6小时”（不聚合，小于等于3小时的不可选） - 计算方式默认为 “平均值”
    - 最近1年 - 间隔默认选择“1个月”（不聚合，小于等于7天的不可选） - 计算方式默认为 “平均值”
    - 今天 - 间隔默认选择“5分钟”（不聚合，小于等于2分钟的不可选） - 计算方式默认为 “平均值”
    - 昨天 - 间隔默认选择“5分钟”（不聚合，小于等于2分钟的不可选） - 计算方式默认为 “平均值”
    - 前天 - 间隔默认选择“5分钟”（不聚合，小于等于2分钟的不可选） - 计算方式默认为 “平均值”
    - 上周今日 - 间隔默认选择“5分钟”（不聚合，小于等于2分钟的不可选） - 计算方式默认为 “平均值”
    - 本周 - 间隔默认选择“30分钟”（不聚合，小于等于10分钟的不可选） - 计算方式默认为 “平均值”
    - 上周 - 间隔默认选择“30分钟”（不聚合，小于等于10分钟的不可选） - 计算方式默认为 “平均值”
    - 本月 - 间隔默认选择“1小时”（不聚合，小于等于30分钟的不可选） - 计算方式默认为 “平均值”
    - 上个月 - 间隔默认选择“1小时”（不聚合，小于等于30分钟的不可选） - 计算方式默认为 “平均值”
    - 今年 - 间隔默认选择“1个月”（不聚合，小于等于7天的不可选） - 计算方式默认为 “平均值”
    - 去年 - 间隔默认选择“1个月”（不聚合，小于等于7天的不可选） - 计算方式默认为 “平均值”

请求参数示例，前端可以直接用这个开发：
```

	{
	    "device_id": "4a5b326c-ba99-9ea2-34a9-1c484d69a1ab",
	    "key": "temperature",
	    "start_time": 1691048558615446,
	    "end_time": 1691048693603021,
	    "aggregate_window": "no_aggregate",
	    "time_range": "custom"
	}

```
30秒最大值
```

	{
	    "device_id": "4a5b326c-ba99-9ea2-34a9-1c484d69a1ab",
	    "key": "temperature",
	    "start_time": 1691048558615446,
	    "end_time": 1691048693603021,
	    "aggregate_window": "30s",
	    "aggregate_function":"max"
	}

```
*/
func (*TelemetryData) GetTelemetrServeStatisticData(req *model.GetTelemetryStatisticReq) (any, error) {
	// 处理时间范围
	if err := processTimeRange(req); err != nil {
		return nil, errcode.WithData(errcode.CodeParamError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	// 获取数据
	rspData, err := fetchTelemetryData(req)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	// 如果不需要导出且无数据，返回空切片
	if !req.IsExport {
		if len(rspData) == 0 {
			return []map[string]interface{}{}, nil
		}
		return rspData, nil
	}

	// 处理导出逻辑
	data, err := exportToCSV(req, rspData)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return data, nil
}

// 处理时间范围
func processTimeRange(req *model.GetTelemetryStatisticReq) error {
	if req.TimeRange == "custom" {
		if req.StartTime == 0 || req.EndTime == 0 || req.StartTime > req.EndTime {
			return fmt.Errorf("time range is invalid")
		}
		return nil
	}

	timeRanges := map[string]time.Duration{
		"last_5m":  5 * time.Minute,
		"last_15m": 15 * time.Minute,
		"last_30m": 30 * time.Minute,
		"last_1h":  time.Hour,
		"last_3h":  3 * time.Hour,
		"last_6h":  6 * time.Hour,
		"last_12h": 12 * time.Hour,
		"last_24h": 24 * time.Hour,
		"last_3d":  72 * time.Hour,
		"last_7d":  7 * 24 * time.Hour,
		"last_15d": 15 * 24 * time.Hour,
		"last_30d": 30 * 24 * time.Hour,
		"last_60d": 60 * 24 * time.Hour,
		"last_90d": 90 * 24 * time.Hour,
		"last_6m":  180 * 24 * time.Hour,
		"last_1y":  365 * 24 * time.Hour,
	}

	duration, ok := timeRanges[req.TimeRange]
	if !ok {
		return fmt.Errorf("unknown time range: %s", req.TimeRange)
	}

	now := time.Now()
	req.EndTime = now.UnixNano() / 1e6
	req.StartTime = now.Add(-duration).UnixNano() / 1e6
	return nil
}

// 获取遥测数据
func fetchTelemetryData(req *model.GetTelemetryStatisticReq) ([]map[string]interface{}, error) {
	if req.AggregateWindow == "no_aggregate" {
		return dal.GetTelemetrStatisticData(req.DeviceId, req.Key, req.StartTime, req.EndTime)
	}

	if err := validateAggregateWindow(req.StartTime, req.EndTime, req.AggregateWindow); err != nil {
		return nil, err
	}

	if req.AggregateFunction == "" {
		req.AggregateFunction = "avg"
	}

	return dal.GetTelemetrStatisticaAgregationData(
		req.DeviceId,
		req.Key,
		req.StartTime,
		req.EndTime,
		dal.StatisticAggregateWindowMillisecond[req.AggregateWindow],
		req.AggregateFunction,
	)
}

func exportToCSV(req *model.GetTelemetryStatisticReq, data []map[string]interface{}) (map[string]interface{}, error) {
	// 检查数据是否为空
	if len(data) == 0 {
		return nil, errcode.New(202100) // 导出数据不能为空
	}

	// 创建导出目录
	exportDir := "./files/excel/telemetry/"
	if err := os.MkdirAll(exportDir, os.ModePerm); err != nil {
		return nil, errcode.WithVars(202101, map[string]interface{}{
			"error": err.Error(),
		})
	}

	// 生成文件名和路径
	fileName := fmt.Sprintf("%s_%s_%d_%d.csv", req.DeviceId, req.Key, req.StartTime, req.EndTime)
	filePath := filepath.Join(exportDir, fileName)

	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		return nil, errcode.WithVars(202102, map[string]interface{}{
			"error": err.Error(),
		})
	}

	// 确保文件最终会被关闭和同步
	defer func() {
		syncErr := file.Sync()
		closeErr := file.Close()
		if err == nil {
			err = syncErr
		}
		if err == nil {
			err = closeErr
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	if err := writer.Write([]string{"时间戳", "数值"}); err != nil {
		return nil, errcode.WithVars(202103, map[string]interface{}{
			"error": err.Error(),
		})
	}

	// 写入数据
	for _, row := range data {
		timestamp, ok := row["x"].(int64)
		if !ok {
			return nil, errcode.New(202105) // 无效的时间戳格式
		}

		value, ok := row["y"].(float64)
		if !ok {
			return nil, errcode.New(202106) // 无效的数值格式
		}

		// 格式化时间
		t := time.Unix(0, timestamp*int64(time.Millisecond))
		formattedTime := t.Format("2006-01-02 15:04:05.000")

		if err := writer.Write([]string{formattedTime, fmt.Sprintf("%.3f", value)}); err != nil {
			return nil, errcode.WithVars(202104, map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	logrus.Info("CSV文件已创建:", filePath)

	return map[string]interface{}{
		"file_name": fileName,
		"file_path": filePath,
	}, nil
}

// AggregateRule 定义聚合规则结构
type AggregateRule struct {
	Days         int    // 天数
	MinInterval  string // 最小允许的聚合间隔
	FriendlyDesc string // 友好描述
}

// validateAggregateWindow 校验聚合窗口设置
func validateAggregateWindow(startTime, endTime int64, aggregateWindow string) error {
	// 计算天数
	days := int((endTime - startTime) / (24 * 60 * 60 * 1000))

	// 定义规则（从大到小排序）
	rules := []AggregateRule{
		{365, "7d", "一年"},
		{180, "1d", "6个月"},
		{90, "6h", "90天"},
		{60, "3h", "60天"},
		{30, "1h", "30天"},
		{15, "30m", "15天"},
		{7, "10m", "7天"},
		{3, "5m", "3天"},
		{1, "2m", "1天"},
	}

	// 检查规则
	for _, rule := range rules {
		if days > rule.Days && !isValidInterval(aggregateWindow, rule.MinInterval) {
			return fmt.Errorf(
				"查询时间范围超过%s，聚合间隔不能小于%s\n\n"+
					"当前配置:\n"+
					"- 时间范围：%s 至 %s（%d天）\n"+
					"- 聚合间隔：%s\n\n"+
					"建议：\n"+
					"1. 使用更大的聚合间隔（>= %s）\n"+
					"2. 或缩短查询时间范围（<= %d天）",
				rule.FriendlyDesc, rule.MinInterval,
				formatTime(startTime), formatTime(endTime), days,
				aggregateWindow,
				rule.MinInterval, rule.Days,
			)
		}
	}

	return nil
}

// isValidInterval 检查聚合间隔是否符合最小要求
func isValidInterval(current, minInterval string) bool {
	// 定义间隔的排序权重
	weights := map[string]int{
		"30s": 1,
		"1m":  2,
		"2m":  3,
		"5m":  4,
		"10m": 5,
		"30m": 6,
		"1h":  7,
		"3h":  8,
		"6h":  9,
		"1d":  10,
		"7d":  11,
		"1mo": 12,
	}

	currentWeight, exists := weights[current]
	if !exists {
		return false
	}

	minWeight, exists := weights[minInterval]
	if !exists {
		return false
	}

	return currentWeight >= minWeight
}

// formatTime 格式化时间戳为可读字符串
func formatTime(timestamp int64) string {
	return time.Unix(timestamp/1000, 0).Format("2006-01-02 15:04:05")
}

func (*TelemetryData) TelemetryPutMessage(ctx context.Context, userID string, param *model.PutMessage, operationType string) error {
	var (
		log = dal.TelemetrySetLogsQuery{}

		errorMessage string
	)
	// 校验param.Value必须是json
	if !json.Valid([]byte(param.Value)) {
		errorMessage = "value must be json"
	}

	deviceInfo, err := initialize.GetDeviceCacheById(param.DeviceID)
	if err != nil {
		logrus.Error(ctx, "[TelemetryPutMessage][GetDeviceCacheById]failed:", err)
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"error": err.Error(),
		})
	}
	// 获取设备配置
	var protocolType string
	var deviceConfig *model.DeviceConfig
	var deviceType string

	if deviceInfo.DeviceConfigID != nil {
		deviceConfig, err = dal.GetDeviceConfigByID(*deviceInfo.DeviceConfigID)
		if err != nil {
			logrus.Error(ctx, "[TelemetryPutMessage][GetDeviceConfigByID]failed:", err)
			return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
				"error": err.Error(),
			})
		}
		deviceType = deviceConfig.DeviceType

		if deviceConfig.ProtocolType != nil {
			protocolType = *deviceConfig.ProtocolType
		} else {
			return errcode.WithData(errcode.CodeParamError, map[string]interface{}{
				"error": "protocolType is nil",
			})
		}
	} else {
		protocolType = "MQTT"
		deviceType = "1"

	}
	var topic string
	if protocolType == "MQTT" {
		// 网关和子设备需要特殊处理
		//messageID := common.GetMessageID()
		topic, err = getTopicByDevice(deviceInfo, deviceType, param)
		if err != nil {
			logrus.Error(ctx, "failed to get topic", err)
			return errcode.WithData(errcode.CodeParamError, map[string]interface{}{
				"error": err.Error(),
			})
		}
	} else {
		// 获取主题前缀
		subTopicPrefix, err := dal.GetServicePluginSubTopicPrefixByDeviceConfigID(*deviceInfo.DeviceConfigID)
		if err != nil {
			logrus.Error(ctx, "failed to get sub topic prefix", err)
			return errcode.WithData(errcode.CodeParamError, map[string]interface{}{
				"error": err.Error(),
			})
		}
		topic = fmt.Sprintf("%s%s%s", subTopicPrefix, config.MqttConfig.Telemetry.PublishTopic, deviceInfo.ID)

	}
	err = publish.PublishTelemetryMessage(topic, deviceInfo, param)
	if err != nil {
		logrus.Error(ctx, "下发失败", err)
		errorMessage = err.Error()
	}
	//operationType := strconv.Itoa(constant.Manual)

	description := "下发遥测日志记录"
	logInfo := &model.TelemetrySetLog{
		ID:            uuid.New(),
		DeviceID:      param.DeviceID,
		OperationType: &operationType,
		Datum:         &(param.Value),
		Status:        nil,
		ErrorMessage:  &errorMessage,
		CreatedAt:     time.Now().UTC(),
		Description:   &description,
		UserID:        &userID,
	}
	// 系统自动发送
	if userID == "" {
		logInfo.UserID = nil
	}
	if err != nil {
		logInfo.ErrorMessage = &errorMessage
		status := strconv.Itoa(constant.StatusFailed)
		logInfo.Status = &status
	} else {
		status := strconv.Itoa(constant.StatusOK)
		logInfo.Status = &status
	}
	_, err = log.Create(ctx, logInfo)
	if err != nil {
		logrus.Error(ctx, "failed to create telemetry set log", err)
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"error": err.Error(),
		})
	}
	return err
}

// getTopicByDevice 根据设备信息获取要发送的控制主题（内置MQTT协议）
func getTopicByDevice(deviceInfo *model.Device, deviceType string, param *model.PutMessage) (string, error) {
	switch deviceType {
	case "1":
		// 处理独立设备
		return fmt.Sprintf("%s%s", config.MqttConfig.Telemetry.PublishTopic, deviceInfo.DeviceNumber), nil
	case "2", "3":
		// 处理网关设备和子设备
		gatewayID := deviceInfo.ID
		if deviceType == "3" {
			if deviceInfo.ParentID == nil {
				return "", fmt.Errorf("parentID 为空")
			}
			gatewayID = *deviceInfo.ParentID
		}

		gatewayInfo, err := initialize.GetDeviceCacheById(gatewayID)
		if err != nil {
			return "", fmt.Errorf("获取网关信息失败: %v", err)
		}

		// 修改payload
		var inputData map[string]interface{}
		if err := json.Unmarshal([]byte(param.Value), &inputData); err != nil {
			return "", fmt.Errorf("解析输入 JSON 失败: %v", err)
		}

		var outputData map[string]interface{}
		if deviceType == "3" {
			if deviceInfo.SubDeviceAddr == nil {
				return "", fmt.Errorf("subDeviceAddr 为空")
			}
			outputData = map[string]interface{}{
				"sub_device_data": map[string]interface{}{
					*deviceInfo.SubDeviceAddr: inputData,
				},
			}
		} else {
			outputData = map[string]interface{}{
				"gateway_data": inputData,
			}
		}

		output, err := json.Marshal(outputData)
		if err != nil {
			return "", fmt.Errorf("生成输出 JSON 失败: %v", err)
		}
		param.Value = string(output)

		return fmt.Sprintf(config.MqttConfig.Telemetry.GatewayPublishTopic, gatewayInfo.DeviceNumber), nil
	default:
		return "", fmt.Errorf("未知的设备类型")
	}
}

func (*TelemetryData) ServeMsgCountByTenantId(tenantId string) (int64, error) {
	cnt, err := dal.GetTelemetryDataCountByTenantId(tenantId)
	if err != nil {
		return 0, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	return cnt, err
}
