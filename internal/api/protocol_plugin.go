package api

import (
	model "project/internal/model"
	service "project/internal/service"

	"github.com/gin-gonic/gin"
)

type ProtocolPluginApi struct{}

// GetProtocolPluginForm 根据协议类型获取设备配置表单
// @Router   /api/v1/protocol_plugin/config_form [get]
func (*ProtocolPluginApi) HandleProtocolPluginFormByProtocolType(c *gin.Context) {
	var req model.GetProtocolPluginFormByProtocolType
	if !BindAndValidate(c, &req) {
		return
	}

	data, err := service.GroupApp.ServicePlugin.GetProtocolPluginFormByProtocolType(req.ProtocolType, req.DeviceType)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}
