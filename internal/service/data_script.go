package service

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	initialize "project/initialize"
	dal "project/internal/dal"
	model "project/internal/model"
	"project/pkg/errcode"
	global "project/pkg/global"
	"project/pkg/metrics"
	utils "project/pkg/utils"

	"github.com/go-basic/uuid"
	"github.com/sirupsen/logrus"
)

type DataScript struct{}

// DelDataScriptCache 根据脚本删除数据脚本缓存
func DelDataScriptCache(data_script *model.DataScript) error {
	deviceIDs, err := dal.GetDeviceIDsByDataScriptID(data_script.ID)
	if err != nil {
		logrus.Error(err)
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	for _, deviceID := range deviceIDs {
		_ = global.REDIS.Del(context.Background(), deviceID+"_"+data_script.ScriptType+"_script").Err()
	}
	return nil
}

func (*DataScript) CreateDataScript(req *model.CreateDataScriptReq) (data_script model.DataScript, err error) {

	data_script.ID = uuid.New()
	data_script.Name = req.Name
	data_script.Description = req.Description
	data_script.DeviceConfigID = req.DeviceConfigId
	data_script.EnableFlag = "N"
	data_script.Content = req.Content
	data_script.ScriptType = req.ScriptType
	data_script.LastAnalogInput = req.LastAnalogInput

	t := time.Now().UTC()
	data_script.CreatedAt = &t
	data_script.UpdatedAt = &t

	data_script.Remark = req.Remark
	err = dal.CreateDataScript(&data_script)

	if err != nil {
		logrus.Error(err)
		return data_script, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	return data_script, err
}

func (*DataScript) UpdateDataScript(UpdateDataScriptReq *model.UpdateDataScriptReq) error {

	err := dal.UpdateDataScript(UpdateDataScriptReq)
	if err != nil {
		logrus.Error(err)
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	new_script, _ := dal.GetDataScriptById(UpdateDataScriptReq.Id)
	err = DelDataScriptCache(new_script)
	if err != nil {
		logrus.Error(err)
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	return err
}

func (*DataScript) DeleteDataScript(id string) error {
	new_script, err := dal.GetDataScriptById(id)
	if err != nil {
		logrus.Error(err)
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	err = dal.DeleteDataScript(id)
	if err != nil {
		logrus.Error(err)
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	if new_script.EnableFlag == "Y" {
		_ = DelDataScriptCache(new_script)
	}
	return err
}

func (*DataScript) GetDataScriptListByPage(Params *model.GetDataScriptListByPageReq) (map[string]interface{}, error) {

	total, list, err := dal.GetDataScriptListByPage(Params)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	data_scriptListRsp := make(map[string]interface{})
	data_scriptListRsp["total"] = total
	data_scriptListRsp["list"] = list

	return data_scriptListRsp, nil
}

func (*DataScript) QuizDataScript(req *model.QuizDataScriptReq) (string, error) {
	if strings.HasPrefix(req.AnalogInput, "0x") {
		msg, err := hex.DecodeString(strings.ReplaceAll(req.AnalogInput, "0x", ""))
		if err != nil {
			return "", errcode.WithVars(100002, map[string]interface{}{
				"error": "hex decode error",
				"input": req.AnalogInput,
			})
		}
		data, error := utils.ScriptDeal(req.Content, msg, req.Topic)
		if error != nil {
			return data, errcode.WithVars(200052, map[string]interface{}{
				"error": error.Error(),
			})
		}
		return data, nil
	}

	data, err := utils.ScriptDeal(req.Content, []byte(req.AnalogInput), req.Topic)
	if err != nil {
		return data, errcode.WithVars(200052, map[string]interface{}{
			"error": err.Error(),
		})
	}
	return data, nil

}

func (*DataScript) EnableDataScript(req *model.EnableDataScriptReq) error {

	if req.EnableFlag == "Y" {
		if ok, err := dal.OnlyOneScriptTypeEnabled(req.Id); !ok {
			return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
				"sql_error": err.Error(),
			})
		}
	}

	var data_script model.DataScript
	data_script.ID = req.Id
	data_script.EnableFlag = req.EnableFlag

	err := dal.EnableDataScript(&data_script)
	if err != nil {
		logrus.Error(err)
		return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}

	if req.EnableFlag == "N" {
		new_script, _ := dal.GetDataScriptById(req.Id)
		err = DelDataScriptCache(new_script)
		if err != nil {
			logrus.Error(err)
			return errcode.WithData(errcode.CodeDBError, map[string]interface{}{
				"sql_error": err.Error(),
			})
		}
	}

	return err
}

func (*DataScript) Exec(device *model.Device, scriptType string, msg []byte, topic string) ([]byte, error) {
	var err error

	script, err := initialize.GetScriptByDeviceAndScriptType(device, scriptType)
	if err != nil {
		return msg, err
	}
	if script == nil {
		return msg, nil
	}
	newMsg, err := utils.ScriptDeal(*script.Content, msg, topic)
	if err != nil {
		return msg, err
	}
	return []byte(newMsg), nil
}

func (*DataScript) RunScript() {
	ins := metrics.NewInstance()
	ins.Instan()
	ins.Count = dal.GetDevicesCount()
	fmt.Println("设备数量:", ins.Count, "个", ins)
	ins.SendSignedRequest()
}
