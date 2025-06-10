package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"project/initialize"
	"project/internal/dal"
	"project/internal/model"
	"project/pkg/common"
	global "project/pkg/global"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-basic/uuid"
	pkgerrors "github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type Automate struct {
	device          *model.Device
	formExt         AutomateFromExt
	mu              sync.Mutex
	automateLimiter *initialize.AutomateLimiter
}

// 1. 定义一个自定义类型（例如 int）
type PolicyType int

// 2. 使用 const + iota 定义枚举值
const (
	Normal              PolicyType = iota // 0 普通策略
	WindProtection                        // 1 防风策略
	WindProtectionLeave                   // 2 解除防风策略
)

type RunFlag int

const (
	ProtectionLeaveWithAck    RunFlag = iota // 0 解除防风 标记，停止下发解除防风命令
	Protection                               // 1 防风下发
	ProtectionLeaveWithoutAck                // 2 解除防风，还未收到反馈；
)

var conditionAfterDecoration = []ConditionAfterFunc{
	ConditionAfterAlarm,
}

var actionAfterDecoration = []ActionAfterFunc{
	ActionAfterAlarm,
}

type ConditionAfterFunc = func(ok bool, conditions initialize.DTConditions, deviceId string, contents []string) error
type ActionAfterFunc = func(actions []model.ActionInfo, err error) error

type AutomateFromExt struct {
	TriggerParamType string
	TriggerParam     []string
	TriggerValues    map[string]interface{}
}

func (a *Automate) conditionAfterDecorationRun(ok bool, conditions initialize.DTConditions, deviceId string, contents []string) {
	defer a.ErrorRecover()
	for _, fc := range conditionAfterDecoration {
		err := fc(ok, conditions, deviceId, contents)
		if err != nil {
			logrus.Error(err)
		}
	}
}

func (a *Automate) actionAfterDecorationRun(actions []model.ActionInfo, err error) {
	defer a.ErrorRecover()
	for _, fc := range actionAfterDecoration {
		err := fc(actions, err)
		if err != nil {
			logrus.Error(err)
		}
	}
}

func (*Automate) ErrorRecover() func() {
	return func() {
		if r := recover(); r != nil {
			// 获取当前的调用堆栈
			stack := string(debug.Stack())
			// 打印堆栈信息
			logrus.Error("自动化 执行异常:\n", r, "\nStack trace:\n", stack)
		}
	}
}

// Execute
// @description 遥测设置上报执行自动化（读取缓存信息 缓存无信息数据库查询保存缓存信息）
// @params deviceInfo *model.Device
// @return error
func (a *Automate) Execute(deviceInfo *model.Device, fromExt AutomateFromExt) error {
	defer a.ErrorRecover()
	a.device = deviceInfo
	a.formExt = fromExt
	//
	logrus.Debugf("Automate Run, devName:%s DeviceConfigID:%v TriggerParamType:%s",
		*deviceInfo.Name, *deviceInfo.DeviceConfigID, fromExt.TriggerParamType)
	//单类设备t
	if deviceInfo.DeviceConfigID != nil {
		deviceConfigId := *deviceInfo.DeviceConfigID
		err := a.telExecute(deviceInfo.ID, deviceConfigId, fromExt)
		if err != nil {
			logrus.Error("自动化执行失败", err)
		}
	}
	logrus.Debug("run telExecute without DeviceConfigID")
	return a.telExecute(deviceInfo.ID, "", fromExt)

}

func (a *Automate) telExecute(deviceId, deviceConfigId string, fromExt AutomateFromExt) error {
	info, resultInt, err := initialize.NewAutomateCache().GetCacheByDeviceId(deviceId, deviceConfigId)
	logrus.Debugf("Automate run: info:%#v, resultInt:%d", info, resultInt)
	if err != nil {
		return pkgerrors.Wrap(err, "查询缓存信息失败")
	}
	//当前设备没自动化任务
	if resultInt == initialize.AUTOMATE_CACHE_RESULT_NOT_TASK {
		return nil
	}
	//缓存未查询到数据 数据查询存入缓存
	if resultInt == initialize.AUTOMATE_CACHE_RESULT_NOT_FOUND {
		info, resultInt, err = a.QueryAutomateInfoAndSetCache(deviceId, deviceConfigId)
		if err != nil {
			return pkgerrors.Wrap(err, "查询设置 设置缓存失败")
		}
		//当前设备没自动化任务
		if resultInt == initialize.AUTOMATE_CACHE_RESULT_NOT_TASK {
			return nil
		}
	}
	logrus.Debugf("Automate run before Filter: info:%#v, resultInt:%d", info, resultInt)
	//过滤自动化触发条件
	info = a.AutomateFilter(info, fromExt)

	logrus.Debugf("Automate run after Filter: info:%#v, resultInt:%v", info, fromExt)
	//执行自动化
	return a.ExecuteRun(info)
}

func (a *Automate) AutomateFilter(info initialize.AutomateExecteParams, fromExt AutomateFromExt) initialize.AutomateExecteParams {
	var sceneInfo []initialize.AutomateExecteSceneInfo
	for _, scene := range info.AutomateExecteSceeInfos {
		var isExists bool
		for _, cond := range scene.GroupsCondition {
			if cond.TriggerParamType == nil || cond.TriggerParam == nil {
				continue
			}
			condTriggerParamType := strings.ToUpper(*cond.TriggerParamType)
			switch fromExt.TriggerParamType {
			case model.TRIGGER_PARAM_TYPE_TEL: //遥测触发
				if condTriggerParamType == model.TRIGGER_PARAM_TYPE_TEL || condTriggerParamType == model.TRIGGER_PARAM_TYPE_TELEMETRY {
					if a.containString(fromExt.TriggerParam, *cond.TriggerParam) {
						isExists = true
					}
				}
			case model.TRIGGER_PARAM_TYPE_STATUS:
				if condTriggerParamType == model.TRIGGER_PARAM_TYPE_STATUS {
					isExists = true
				}
			case model.TRIGGER_PARAM_TYPE_EVT:
				if (condTriggerParamType == model.TRIGGER_PARAM_TYPE_EVT || condTriggerParamType == model.TRIGGER_PARAM_TYPE_EVENT) && a.containString(fromExt.TriggerParam, *cond.TriggerParam) {
					isExists = true
				}
			case model.TRIGGER_PARAM_TYPE_ATTR:
				if condTriggerParamType == model.TRIGGER_PARAM_TYPE_ATTR && a.containString(fromExt.TriggerParam, *cond.TriggerParam) {
					isExists = true
				}
			}
		}
		if isExists {
			sceneInfo = append(sceneInfo, scene)
		}
	}
	info.AutomateExecteSceeInfos = sceneInfo
	return info
}

func (*Automate) containString(slice []string, str string) bool {
	for _, v := range slice {
		logrus.Info(v, " : ", str)
		if v == str {
			return true
		}
	}
	return false
}

// 限流实现 1秒一次 安场景实现
func (a *Automate) LimiterAllow(id string) bool {
	if a.automateLimiter == nil {
		a.automateLimiter = initialize.NewAutomateLimiter(initialize.AutomateRateLimitConfig)
	}
	return a.automateLimiter.GetLimiter(fmt.Sprintf("SceneAutomationId:%s", id)).Allow()
}

var runProtectionLeaveCount int // 执行解除防风的次数
const MaxRunTimes = 15

// ExecuteRun
// @description  自动化场景联动执行
// @params info initialize.AutomateExecteParams
// @return error
func (a *Automate) ExecuteRun(info initialize.AutomateExecteParams) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	logrus.Debug("len of AutomateExecteSceeInfos:", len(info.AutomateExecteSceeInfos))
	newSceneInfo := make([]initialize.AutomateExecteSceneInfo, 2, 2)
	for _, v := range info.AutomateExecteSceeInfos {
		if a.CheckSceneAutomationWindProtectionLeave(v.SceneAutomationId) { // 解除防风策略
			newSceneInfo[0] = v
		} else if a.CheckSceneAutomationWindProtection(v.SceneAutomationId) { // 防风策略
			newSceneInfo[1] = v
		}
	}
	info.AutomateExecteSceeInfos = newSceneInfo

	sceneFlag := false

	for _, v := range info.AutomateExecteSceeInfos {
		//场景频率限制(根据场景id)
		// if !a.LimiterAllow(v.SceneAutomationId) {
		// 	logrus.Debug("ExecuteRun Limiter Not Allow")
		// 	continue
		// }
		if sceneFlag {
			logrus.Debug("WindProtectionLeave Run, Now Break This Policy Check")
			continue
		}
		logrus.Debugf("查询自动化是否关闭1: info:%#v,", v.SceneAutomationId)
		//查询自动化是否关闭
		if a.CheckSceneAutomationHasClose(v.SceneAutomationId) {
			continue
		}
		logrus.Debugf("判断条件是否成立: info:%#v,", info)

		pType := Normal
		if a.CheckSceneAutomationWindProtectionLeave(v.SceneAutomationId) {
			pType = WindProtectionLeave // 解除防风策略
		} else if a.CheckSceneAutomationWindProtection(v.SceneAutomationId) {
			pType = WindProtection //防风策略
		}
		//条件判断
		if !a.AutomateConditionCheck(v.GroupsCondition, info.DeviceId) {
			// 条件判断不成立
			if pType != WindProtection {
				continue
			}
			// 防风策略特殊处理
			flag, err := getRunFlag()
			if err != nil {
				logrus.Error("getRunFlag err:", err)
				continue
			}
			// 当前不处于防风模式时跳过
			if RunFlag(flag) != Protection {
				continue
			}

			logrus.Debug("WindProtection handler,auto run WindProtection.")
		}
		var runID int
		//动作执行前判断
		if pType == WindProtectionLeave { // 解除防风策略
			sceneFlag = true
			if runProtectionLeaveCount >= MaxRunTimes {
				logrus.Debug("WindProtectionLeave Policy Check, Device Report Policy Cmd Run Max 5,Now Do Not Run This Policy")
				continue
			}
			runProtectionLeaveCount++
			//未执行成功，继续下发 退出防风

			// 1. 先获取 policyRunID
			runID, err := global.REDIS.Get(context.Background(), "policyRunID").Int()
			if err != nil {
				if errors.Is(err, redis.Nil) { // 没有这个标记表示之前一直未执行过 防风策略
					runID = 1 //表示设备第一次执行防风策略
				} else {
					logrus.Error("WindProtectionLeave Policy check RunFlag err", err)
				}
			}
			logrus.Debug("now ready run WindProtectionLeave Policy,runID:", runID)
			remarkValue := fmt.Sprintf("runID:%d", runID)
			// 把policyRunID 保存到 Remark 中，在后续的命令下发中可以解析出来
			for i := range v.Actions {
				v.Actions[i].Remark = &remarkValue
			}
		} else if pType == WindProtection { // 防风策略
			runProtectionLeaveCount = 0
			var err error
			// 1. 先获取 policyRunID
			runID, err = global.REDIS.Get(context.Background(), "policyRunID").Int()
			if err != nil {
				if errors.Is(err, redis.Nil) { // 没有这个标记表示之前一直未执行过 防风策略
					runID = 1 //表示设备第一次执行防风策略
				} else {
					logrus.Error("WindProtection Policy check RunFlag err", err)
				}
			}
			// 2. 根据RunFlag的值，判断 policyRunID 是否需要更新
			resultFlag, err := global.REDIS.Get(context.Background(), "RunFlag").Int()
			if err != nil {
				logrus.Debug("WindProtectionLeave Policy check RunFlag err", err)
			} else if RunFlag(resultFlag) == ProtectionLeaveWithAck { // 上一次 解除防风 已经正确收到响应， 增加policyRunID
				// policyRunID + 1
				runID += 1
			}
			logrus.Debug("now ready run WindProtection Policy,runID:", runID)
			remarkValue := fmt.Sprintf("runID:%d", runID)
			// 把policyRunID 保存到 Remark 中，在后续的命令下发中可以解析出来
			for i := range v.Actions {
				v.Actions[i].Remark = &remarkValue
			}
		}
		// 场景联动 动作执行
		err := a.SceneAutomateExecute(v.SceneAutomationId, []string{info.DeviceId}, v.Actions)
		// 执行后，添加redis 标记
		if err == nil {
			// 防风策略的特殊处理
			if pType == WindProtection {
				// 1. 开始处理 policyRunID 的更新问题
				resultFlag, err := global.REDIS.Get(context.Background(), "RunFlag").Int()
				if err != nil {
					if errors.Is(err, redis.Nil) { // 没有这个标记表示之前一直未执行过 防风策略，增加policyRunID
						logrus.Debug("WindProtection Set policyRunID Value: 1")
						global.REDIS.Set(context.Background(), "policyRunID", fmt.Sprintf("%d", 1), 0) // key 用不过期
					}
				} else if RunFlag(resultFlag) == ProtectionLeaveWithAck { // 上一次 解除防风 已经正确收到响应， 增加policyRunID
					// policyRunID + 1
					logrus.Debug("WindProtection Set policyRunID Value:", runID)
					global.REDIS.Set(context.Background(), "policyRunID", fmt.Sprintf("%d", runID), 0) // key 用不过期
				}
				// 更新状态标记
				global.REDIS.Set(context.Background(), "RunFlag", fmt.Sprintf("%d", Protection), 0) // key 用不过期
			} else if pType == WindProtectionLeave { // 解除防风策略
				logrus.Debug("WindProtection Set RunFlag to 2")
				global.REDIS.Set(context.Background(), "RunFlag", fmt.Sprintf("%d", ProtectionLeaveWithoutAck), 0) // key 用不过期
			}
		}
		// 场景动作之后装饰
		a.actionAfterDecorationRun(v.Actions, err)
	}

	return nil
}

func getRunFlag() (int, error) {
	resultFlag, err := global.REDIS.Get(context.Background(), "RunFlag").Int()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return -1, nil
		}
		return -1, err
	}
	return resultFlag, nil
}

// CheckSceneAutomationHasClose
// @description 查询是否关闭了自动化
func (*Automate) CheckSceneAutomationHasClose(sceneAutomationId string) bool {
	ok := dal.CheckSceneAutomationHasClose(sceneAutomationId)
	//删除缓存
	if ok {
		_ = initialize.NewAutomateCache().DeleteCacheBySceneAutomationId(sceneAutomationId)
	}
	return ok
}

func (*Automate) CheckSceneAutomationWindProtection(sceneAutomationId string) bool {
	return dal.CheckWindProtection(sceneAutomationId)
}

func (*Automate) CheckSceneAutomationWindProtectionLeave(sceneAutomationId string) bool {
	return dal.CheckWindProtectionLeave(sceneAutomationId)
}

// SceneAutomateExecute
// @description 场景联动 动作执行
// @params info initialize.AutomateExecteParams
// @return error
func (a *Automate) SceneAutomateExecute(sceneAutomationId string, deviceIds []string, actions []model.ActionInfo) error {
	tenantID := dal.GetSceneAutomationTenantID(context.Background(), sceneAutomationId)
	//执行动作
	details, err := a.AutomateActionExecute(sceneAutomationId, deviceIds, actions, tenantID)

	_ = a.sceneExecuteLogSave(sceneAutomationId, details, err)

	return err
}

// ActiveSceneExecute
// @description 场景激活
// @params info initialize.AutomateExecteParams
// @return error
func (a *Automate) ActiveSceneExecute(scene_id, tenantID, remark string) error {

	actions, err := dal.GetActionInfoListBySceneId([]string{scene_id})
	if err != nil {
		return nil
	}
	var (
		deviceIds      []string
		deviceConfigId []string
	)
	for i, v := range actions {
		if v.ActionType == model.AUTOMATE_ACTION_TYPE_MULTIPLE && v.ActionTarget != nil {
			deviceConfigId = append(deviceConfigId, *v.ActionTarget)
		}
		if remark != "" {
			actions[i].Remark = &remark
		}
	}
	if len(deviceConfigId) > 0 {
		deviceIds, err = dal.GetDeviceIdsByDeviceConfigId(deviceConfigId)
		if err != nil {
			return err
		}
	}
	details, err := a.AutomateActionExecute(scene_id, deviceIds, actions, tenantID)
	var exeResult string
	if err == nil {
		exeResult = "S"
	} else {
		exeResult = "F"
	}
	logrus.Debug(details)
	return dal.SceneLogInsert(&model.SceneLog{
		ID:              uuid.New(),
		SceneID:         scene_id,
		ExecutedAt:      time.Now().UTC(),
		Detail:          details,
		ExecutionResult: exeResult,
		TenantID:        tenantID,
	})
}

// @description sceneExecuteLogSave 自动化场景联动执行
// @params info initialize.AutomateExecteParams
// @return error
func (*Automate) sceneExecuteLogSave(scene_id, details string, err error) error {
	var exeResult string
	if err == nil {
		exeResult = "S"
	} else {
		exeResult = "F"
	}
	logrus.Debug(details)
	return dal.SceneAutomationLogInsert(&model.SceneAutomationLog{
		SceneAutomationID: scene_id,
		ExecutedAt:        time.Now().UTC(),
		Detail:            details,
		ExecutionResult:   exeResult,
		TenantID:          dal.GetSceneAutomationTenantID(context.Background(), scene_id),
	})
}

// AutomateConditionCheck
// @description  自动化条件判断 复合其中一组条件就返回true
// @params conditions []initialize.DTConditions
// @return bool true 表示可以执行动作
func (a *Automate) AutomateConditionCheck(conditions initialize.DTConditions, deviceId string) bool {
	logrus.Debug("条件判断开始...")
	//key是groupId val是条件列表
	conditionsByGroupId := make(map[string]initialize.DTConditions)
	for _, v := range conditions {
		conditionsByGroupId[v.GroupID] = append(conditionsByGroupId[v.GroupID], v)
	}
	var result bool
	for _, val := range conditionsByGroupId {
		ok, _ := a.AutomateConditionCheckWithGroup(val, deviceId)
		if ok {
			result = true
		}
		//组条件执行完成装饰
		//a.conditionAfterDecorationRun(ok, val, deviceId, contents)
	}
	return result
}

// AutomateConditionCheckWithGroup
// @description  一组条件比较 一个为假结果就为假
// @params conditions initialize.DTConditions
// @return bool
func (a *Automate) AutomateConditionCheckWithGroup(conditions initialize.DTConditions, deviceId string) (bool, []string) {
	var (
		result   []string
		resultOk bool = true
	)
	for _, val := range conditions {
		ok, content := a.AutomateConditionCheckWithGroupOne(val, deviceId)
		result = append(result, content)
		if !ok {
			resultOk = false
			break
		}
	}

	return resultOk, result
}

// @description AutomateConditionCheckWithGroupOne 单个条件验证
// @params cond model.DeviceTriggerCondition
// @return bool
func (a *Automate) AutomateConditionCheckWithGroupOne(cond model.DeviceTriggerCondition, deviceId string) (bool, string) {
	logrus.Debug("条件type:", cond.TriggerConditionType)
	switch cond.TriggerConditionType {
	case model.DEVICE_TRIGGER_CONDITION_TYPE_TIME:
		return a.automateConditionCheckWithTime(cond), ""
	case model.DEVICE_TRIGGER_CONDITION_TYPE_ONE, model.DEVICE_TRIGGER_CONDITION_TYPE_MULTIPLE: //单类或者单个设备
		return a.automateConditionCheckWithDevice(cond, deviceId)
	default:
		return true, ""
	}
}

// @description automateConditionCheckWithTime 单个条件时间范围验证
// @params cond model.DeviceTriggerCondition
// @return bool
func (*Automate) automateConditionCheckWithTime(cond model.DeviceTriggerCondition) bool {
	logrus.Debug("时间范围对比开始... 条件:", cond.TriggerValue)
	nowTime := time.Now().UTC()
	if cond.TriggerValue == "" {
		return false
	}
	valParts := strings.Split(cond.TriggerValue, "|")
	if len(valParts) < 3 {
		return false
	}
	var ok bool
	//获取当前星期
	weekDay := common.GetWeekDay(nowTime)
	//判断当前时间和条件星期
	for _, char := range valParts[0] {
		num, _ := strconv.Atoi(string(char))
		if weekDay == num {
			ok = true
			continue
		}
	}
	//没有在当前指定的星期中
	if !ok {
		return false
	}
	nowTimeNotDay, _ := time.Parse("15:04:05-07:00", nowTime.Format("15:04:05-07:00"))
	startTime, err := time.Parse("15:04:05-07:00", valParts[1])
	if err != nil {
		logrus.Error("时间格式不正确, 字符串", cond.TriggerValue)
		return false
	}
	if startTime.After(nowTimeNotDay) {
		return false
	}

	endTime, err := time.Parse("15:04:05-07:00", valParts[2])
	if err != nil {
		logrus.Error("时间格式不正确, 字符串", cond.TriggerValue)
		return false
	}
	if endTime.Before(nowTimeNotDay) {
		return false
	}
	logrus.Debug("时间范围对比结束。OK")
	return true
}

func (a *Automate) getActualValue(deviceId string, key string, triggerParamType string) (interface{}, error) {
	for k, v := range a.formExt.TriggerValues {
		if key == k {
			return v, nil
		}
	}
	switch triggerParamType {
	case model.TRIGGER_PARAM_TYPE_TEL:
		return dal.GetCurrentTelemetryDataOneKeys(deviceId, key)
	case model.TRIGGER_PARAM_TYPE_ATTR:
		return dal.GetAttributeOneKeys(deviceId, key)
	case model.TRIGGER_PARAM_TYPE_EVT:
		return dal.GetDeviceEventOneKeys(deviceId, key)
	case model.TRIGGER_PARAM_TYPE_STATUS:
		return dal.GetDeviceCurrentStatus(deviceId)
	}

	return nil, nil
}
func (a *Automate) automateConditionCheckWithDevice(cond model.DeviceTriggerCondition, deviceId string) (bool, string) {
	logrus.Debug("设备条件验证开始...")
	//设备id不存在 返回假
	if cond.TriggerSource == nil {
		return false, ""
	}
	//单类设置 获取上报的设置 单个设置 使用设置的设备id
	if cond.TriggerConditionType == model.DEVICE_TRIGGER_CONDITION_TYPE_ONE {
		deviceId = *cond.TriggerSource
	}

	//条件查询
	var (
		actualValue     interface{}
		trigger         string
		triggerValue    string
		triggerOperator string
		triggerKey      string
		result          string
		deviceName      string
	)
	if a.device.Name != nil {
		deviceName = *a.device.Name
	}
	if cond.TriggerOperator == nil {
		triggerOperator = "="
	} else {
		triggerOperator = *cond.TriggerOperator
	}

	logrus.Debug("设备条件验证开始...", strings.ToUpper(*cond.TriggerParamType))
	switch strings.ToUpper(*cond.TriggerParamType) {
	case model.TRIGGER_PARAM_TYPE_TEL, model.TRIGGER_PARAM_TYPE_TELEMETRY: //遥测
		trigger = "遥测"
		//actualValue, _ = dal.GetCurrentTelemetryDataOneKeys(deviceId, *cond.TriggerParam)
		actualValue, _ = a.getActualValue(deviceId, *cond.TriggerParam, model.TRIGGER_PARAM_TYPE_TEL)
		triggerValue = cond.TriggerValue
		triggerKey = *cond.TriggerParam
		logrus.Debugf("GetCurrentTelemetryDataOneKeys:triggerOperator:%s, TriggerParam:%s, triggerValue:%v, actualValue:%v", triggerOperator, *cond.TriggerParam, triggerValue, actualValue)
		dataValue := a.getTriggerParamsValue(triggerKey, dal.GetIdentifierNameTelemetry())
		result = fmt.Sprintf("设备(%s)%s [%s]: %v %s %v", deviceName, trigger, dataValue, actualValue, triggerOperator, triggerValue)
	case model.TRIGGER_PARAM_TYPE_ATTR: //属性
		trigger = "属性"
		actualValue, _ = a.getActualValue(deviceId, *cond.TriggerParam, model.TRIGGER_PARAM_TYPE_ATTR)
		triggerValue = cond.TriggerValue
		triggerKey = *cond.TriggerParam
		dataValue := a.getTriggerParamsValue(triggerKey, dal.GetIdentifierNameAttribute())
		result = fmt.Sprintf("设备(%s)%s [%s]: %v %s %v", deviceName, trigger, dataValue, actualValue, triggerOperator, triggerValue)
	case model.TRIGGER_PARAM_TYPE_EVT, model.TRIGGER_PARAM_TYPE_EVENT: //事件
		trigger = "事件"
		actualValue, _ = a.getActualValue(deviceId, *cond.TriggerParam, model.TRIGGER_PARAM_TYPE_EVT)
		triggerValue = cond.TriggerValue
		triggerKey = *cond.TriggerParam
		logrus.Debugf("事件...actualValue:%#v, triggerValue:%#v", actualValue, triggerValue)
		dataValue := a.getTriggerParamsValue(triggerKey, dal.GetIdentifierNameEvent())
		result = fmt.Sprintf("设备(%s)%s [%s]: %v %s %v", deviceName, trigger, dataValue, actualValue, triggerOperator, triggerValue)
	case model.TRIGGER_PARAM_TYPE_STATUS: //状态
		trigger = "下线"
		actualValue, _ = a.getActualValue(deviceId, "login", model.TRIGGER_PARAM_TYPE_STATUS)
		triggerValue = *cond.TriggerParam
		if strings.ToUpper(actualValue.(string)) == "ON-LINE" {
			trigger = "上线"
		}
		result = fmt.Sprintf("设备(%s)已%s", deviceName, trigger)
		triggerOperator = "="
		if strings.ToUpper(triggerValue) == "ALL" {
			return true, result
		}
	}
	logrus.Debugf("automateConditionCheckByOperator:设备条件验证参数 triggerOperator:%v triggerValue:%v actualValue:%v",
		triggerOperator, triggerValue, actualValue)
	ok := a.automateConditionCheckByOperator(triggerOperator, triggerValue, actualValue)
	logrus.Debugf("比较结果:%t", ok)
	return ok, result
}

type DataIdentifierName func(device_template_id, identifier string) string

func (*Automate) getTriggerParamsValue(triggerKey string, fc DataIdentifierName) string {
	tempId, _ := dal.GetDeviceTemplateIdByDeviceId(triggerKey)
	if tempId == "" {
		return triggerKey
	}

	return fc(tempId, triggerKey)
}

// automateConditionCheckByOperator
// @description  运算符处理
// @params cond model.DeviceTriggerCondition
// @return bool
func (a *Automate) automateConditionCheckByOperator(operator string, condValue string, actualValue interface{}) bool {
	//logrus.Warningf("比较:operator:%s, condValue:%s, actualValue: %s, result:%d", operator, condValue, actualValue, strings.Compare(actualValue, condValue))
	switch value := actualValue.(type) {
	case string:
		return a.automateConditionCheckByOperatorWithString(operator, condValue, value)
	case float64:
		return a.automateConditionCheckByOperatorWithFloat(operator, condValue, value)
	case bool:
		return a.automateConditionCheckByOperatorWithString(operator, condValue, fmt.Sprintf("%t", value))
	}
	return false
}

func float64Equal(a, b float64) bool {
	const threshold = 1e-9
	return math.Abs(a-b) < threshold
}

// automateConditionCheckByOperatorWithString
// @description  运算符处理
// @params cond model.DeviceTriggerCondition
// @return bool
func (*Automate) automateConditionCheckByOperatorWithFloat(operator string, condValue string, actualValue float64) bool {
	//logrus.Warningf("比较:operator:%s, condValue:%s, actualValue: %s, result:%d", operator, condValue, actualValue, strings.Compare(actualValue, condValue))

	switch operator {
	case model.CONDITION_TRIGGER_OPERATOR_EQ:
		condValueFloat, err := strconv.ParseFloat(condValue, 64)
		if err != nil {
			logrus.Error(err)
			return false
		}
		return float64Equal(condValueFloat, actualValue)
	case model.CONDITION_TRIGGER_OPERATOR_NEQ:
		condValueFloat, err := strconv.ParseFloat(condValue, 64)
		if err != nil {
			logrus.Error(err)
			return false
		}
		logrus.Debugf("condValueFloat:%f, actualValue:%f, 结果：%t", condValueFloat, actualValue, float64Equal(condValueFloat, actualValue))
		return !float64Equal(condValueFloat, actualValue)
	case model.CONDITION_TRIGGER_OPERATOR_GT:
		condValueFloat, err := strconv.ParseFloat(condValue, 64)
		if err != nil {
			logrus.Error(err)
			return false
		}
		return actualValue > condValueFloat
	case model.CONDITION_TRIGGER_OPERATOR_LT:
		condValueFloat, err := strconv.ParseFloat(condValue, 64)
		if err != nil {
			logrus.Error(err)
			return false
		}
		return actualValue < condValueFloat
	case model.CONDITION_TRIGGER_OPERATOR_GTE:
		condValueFloat, err := strconv.ParseFloat(condValue, 64)
		if err != nil {
			logrus.Error(err)
			return false
		}
		return actualValue >= condValueFloat
	case model.CONDITION_TRIGGER_OPERATOR_LTE:
		condValueFloat, err := strconv.ParseFloat(condValue, 64)
		if err != nil {
			logrus.Error(err)
			return false
		}
		return actualValue <= condValueFloat
	case model.CONDITION_TRIGGER_OPERATOR_BETWEEN:
		valParts := strings.Split(condValue, "-")
		if len(valParts) != 2 {
			return false
		}

		valParts0Float64, err := strconv.ParseFloat(valParts[0], 64)
		if err != nil {
			return false
		}
		valParts1Float64, err := strconv.ParseFloat(valParts[1], 64)
		if err != nil {
			return false
		}
		return actualValue >= valParts0Float64 && actualValue <= valParts1Float64
	case model.CONDITION_TRIGGER_OPERATOR_IN:
		valParts := strings.Split(condValue, ",")
		for _, v := range valParts {
			vFloat, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return false
			}
			if float64Equal(vFloat, actualValue) {
				return true
			}
		}
	}
	return false
}

// automateConditionCheckByOperatorWithString
// @description  运算符处理
// @params cond model.DeviceTriggerCondition
// @return bool
func (*Automate) automateConditionCheckByOperatorWithString(operator string, condValue string, actualValue string) bool {
	logrus.Warningf("比较:operator:%s, condValue:%s, actualValue: %s, result:%d", operator, condValue, actualValue, strings.Compare(actualValue, condValue))
	switch operator {
	case model.CONDITION_TRIGGER_OPERATOR_EQ:
		return strings.EqualFold(strings.ToUpper(actualValue), strings.ToUpper(condValue))
	case model.CONDITION_TRIGGER_OPERATOR_NEQ:
		return strings.Compare(actualValue, condValue) != 0
	case model.CONDITION_TRIGGER_OPERATOR_GT:
		actualValueFloat64, err := strconv.ParseFloat(actualValue, 64)
		if err != nil {
			return strings.Compare(actualValue, condValue) > 0
		}
		condValueFloat64, err := strconv.ParseFloat(condValue, 64)
		if err != nil {
			return strings.Compare(actualValue, condValue) > 0
		}
		return actualValueFloat64 > condValueFloat64
	case model.CONDITION_TRIGGER_OPERATOR_LT:
		actualValueFloat64, err := strconv.ParseFloat(actualValue, 64)
		if err != nil {
			return strings.Compare(actualValue, condValue) < 0
		}
		condValueFloat64, err := strconv.ParseFloat(condValue, 64)
		if err != nil {
			return strings.Compare(actualValue, condValue) < 0
		}
		return actualValueFloat64 < condValueFloat64
	case model.CONDITION_TRIGGER_OPERATOR_GTE:
		actualValueFloat64, err := strconv.ParseFloat(actualValue, 64)
		if err != nil {
			return strings.Compare(actualValue, condValue) >= 0
		}
		condValueFloat64, err := strconv.ParseFloat(condValue, 64)
		if err != nil {
			return strings.Compare(actualValue, condValue) >= 0
		}
		return actualValueFloat64 >= condValueFloat64
	case model.CONDITION_TRIGGER_OPERATOR_LTE:
		actualValueFloat64, err := strconv.ParseFloat(actualValue, 64)
		if err != nil {
			return strings.Compare(actualValue, condValue) <= 0
		}
		condValueFloat64, err := strconv.ParseFloat(condValue, 64)
		if err != nil {
			return strings.Compare(actualValue, condValue) <= 0
		}
		return actualValueFloat64 <= condValueFloat64
	case model.CONDITION_TRIGGER_OPERATOR_BETWEEN:
		valParts := strings.Split(condValue, "-")
		if len(valParts) != 2 {
			return false
		}
		actualValueFloat64, err := strconv.ParseFloat(actualValue, 64)
		if err != nil {
			return actualValue >= valParts[0] && actualValue <= valParts[1]
		}
		valParts0Float64, err := strconv.ParseFloat(valParts[0], 64)
		if err != nil {
			return actualValue >= valParts[0] && actualValue <= valParts[1]
		}
		valParts1Float64, err := strconv.ParseFloat(valParts[1], 64)
		if err != nil {
			return actualValue >= valParts[0] && actualValue <= valParts[1]
		}
		return actualValueFloat64 >= valParts0Float64 && actualValueFloat64 <= valParts1Float64
	case model.CONDITION_TRIGGER_OPERATOR_IN:
		valParts := strings.Split(condValue, ",")
		for _, v := range valParts {
			if v == actualValue {
				return true
			}
		}
	}
	return false
}

// AutomateActionExecute
// @description  自动化动作执行
// @params deviceId string
// @params actions []model.ActionInf
// @return void
func (a *Automate) AutomateActionExecute(sceneAutomationId string, deviceIds []string, actions []model.ActionInfo, tenantID string) (string, error) {
	logrus.Debug("动作开始执行:")
	var (
		result    string
		resultErr error
	)
	if len(actions) == 0 {
		return "未找到执行动作", errors.New("未找到执行动作")
	}

	for _, action := range actions {
		var actionService AutomateTelemetryAction
		logrus.Debug("actionType:", action.ActionType)
		switch action.ActionType {
		case model.AUTOMATE_ACTION_TYPE_ONE: //单个设置
			actionService = &AutomateTelemetryActionOne{TenantID: tenantID}
		case model.AUTOMATE_ACTION_TYPE_ALARM: //告警触发
			actionService = &AutomateTelemetryActionAlarm{}
		case model.AUTOMATE_ACTION_TYPE_MULTIPLE: //单类设置
			actionService = &AutomateTelemetryActionMultiple{DeviceIds: deviceIds, TenantID: tenantID}
		case model.AUTOMATE_ACTION_TYPE_SCENE: //激活场景
			actionService = &AutomateTelemetryActionScene{TenantID: tenantID}
		case model.AUTOMATE_ACTION_TYPE_SERVICE: //服务
			actionService = &AutomateTelemetryActionService{}
		}
		if actionService == nil {
			logrus.Error("暂不支持的动作类型")
			return "暂不支持的动作类型", errors.New("暂不支持的动作类型")
		}
		// type commandInfo struct {
		// 	Method string `json:"method"`
		// 	Params string `json:"params"`
		// }
		// var pType PolicyType
		// if a.CheckSceneAutomationWindProtection(sceneAutomationId) { // 防风策略的特殊处理
		// 	var info commandInfo
		// 	err := json.Unmarshal([]byte(*action.ActionValue), &info)
		// 	if err != nil {
		// 		logrus.Error("防风策略，原始命令解析失败 ActionValue:", *action.ActionValue)
		// 		return "暂不支持的动作类型", errors.New("暂不支持的动作类型")
		// 	}
		// 	pType = WindProtection
		// 	// 标记防风执行一次
		// 	runid, err := global.REDIS.Incr(context.Background(), policyRunID).Result()
		// 	if err != nil {
		// 		logrus.Error("set policyRunID err:", err)
		// 		return "内部错误", errors.New("内部错误")
		// 	}
		// 	if info.Params != "" {
		// 		paramValue := make(map[string]interface{}, 0)
		// 		if err := json.Unmarshal([]byte(info.Params), &paramValue); err != nil {
		// 			logrus.Error("防风策略，原始命令解析失败")
		// 			return "内部错误", errors.New("内部错误")
		// 		}
		// 		paramValue["runid"] = runid
		// 		newcmdStr, err := json.Marshal(paramValue)
		// 		if err != nil {
		// 			logrus.Error("json Marshal err:", err)
		// 			return "内部错误", errors.New("内部错误")
		// 		}
		// 		info.Params = string(newcmdStr)
		// 		newValue, _ := json.Marshal(info)
		// 		*action.ActionValue = string(newValue)
		// 	}

		// } else if a.CheckSceneAutomationWindProtectionLeave(sceneAutomationId) { // 解除防风策略的特殊处理
		// 	var info commandInfo
		// 	err := json.Unmarshal([]byte(*action.ActionValue), &info)
		// 	if err != nil {
		// 		logrus.Error("解除防风策略，原始命令解析失败")
		// 		return "暂不支持的动作类型", errors.New("暂不支持的动作类型")
		// 	}
		// 	pType = WindProtectionLeave
		// }
		actionMessage, err := actionService.AutomateActionRun(action)
		if err != nil && resultErr == nil {
			resultErr = err
		}
		if err != nil {
			result += fmt.Sprintf("%s 执行失败;", actionMessage)
		} else {
			result += fmt.Sprintf("%s 下发成功;", actionMessage)
		}
	}
	logrus.Debug("result:", result)
	return result, resultErr
}

// QueryAutomateInfoAndSetCache
// @description  查询设备自动化信息并缓存
// @params deviceId string
// @return initialize.AutomateExecteParams, int, error
func (*Automate) QueryAutomateInfoAndSetCache(deviceId, deviceConfigId string) (initialize.AutomateExecteParams, int, error) {
	automateExecuteParams := initialize.AutomateExecteParams{
		DeviceId:       deviceId,
		DeviceConfigId: deviceConfigId,
	}
	var (
		groups []model.DeviceTriggerCondition
		err    error
	)
	//deviceConfigId 存在 表示单类设置
	if deviceConfigId != "" {
		groups, err = dal.GetDeviceTriggerConditionByDeviceId(deviceConfigId, model.DEVICE_TRIGGER_CONDITION_TYPE_MULTIPLE)
	} else {
		groups, err = dal.GetDeviceTriggerConditionByDeviceId(deviceId, model.DEVICE_TRIGGER_CONDITION_TYPE_ONE)
	}
	logrus.Debugf("设备配置id: %s, 查询条件：%v", deviceConfigId, groups)
	if err != nil {
		return automateExecuteParams, 0, pkgerrors.Wrap(err, "根据设备id查询自动化条件失败")
	}
	//没有查询到该设备自动化信息
	if len(groups) == 0 {
		err := initialize.NewAutomateCache().SetCacheByDeviceIdWithNoTask(deviceId, deviceConfigId)
		if err != nil {
			return automateExecuteParams, 0, pkgerrors.Wrap(err, "设置设备无自动化信息缓存失败")
		}
		return automateExecuteParams, initialize.AUTOMATE_CACHE_RESULT_NOT_TASK, nil
	}
	sceneAutomateGroups := make(map[string]bool)
	var (
		sceneAutomateIds []string
		groupIds         []string
	)

	for _, groupInfo := range groups {
		if _, ok := sceneAutomateGroups[groupInfo.SceneAutomationID]; !ok {
			sceneAutomateIds = append(sceneAutomateIds, groupInfo.SceneAutomationID)
			sceneAutomateGroups[groupInfo.SceneAutomationID] = true
		}
		groupIds = append(groupIds, groupInfo.GroupID)
	}
	//查询场景所有的group条件
	groups, err = dal.GetDeviceTriggerConditionByGroupIds(groupIds)
	if err != nil {
		return automateExecuteParams, 0, pkgerrors.Wrap(err, "查询自动化条件失败")
	}
	//查询场景执行动作
	actionInfos, err := dal.GetActionInfoListBySceneAutomationId(sceneAutomateIds)
	if err != nil {
		return automateExecuteParams, 0, pkgerrors.Wrap(err, "查询自动化执行失败")
	}
	logrus.Debugf("设备配置id2: %s, 查询条件：%v, 动作: %v", deviceConfigId, groups, actionInfos)
	//设置自动化缓存
	err = initialize.NewAutomateCache().SetCacheByDeviceId(deviceId, deviceConfigId, groups, actionInfos)
	if err != nil {
		return automateExecuteParams, 0, pkgerrors.Wrap(err, "设置自动化缓存失败")
	}

	return initialize.NewAutomateCache().GetCacheByDeviceId(deviceId, deviceConfigId)
}
