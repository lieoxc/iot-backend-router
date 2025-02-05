package apps

import (
	"project/internal/api"

	"github.com/gin-gonic/gin"
)

type ProtocolPlugin struct {
}

func (*ProtocolPlugin) InitProtocolPlugin(Router *gin.RouterGroup) {
	protocolPluginApi := Router.Group("protocol_plugin")
	{
		// // 增
		protocolPluginApi.POST("", api.Controllers.ProtocolPluginApi.CreateProtocolPlugin)

		// // 删
		protocolPluginApi.DELETE(":id", api.Controllers.ProtocolPluginApi.DeleteProtocolPlugin)

		// // 改
		protocolPluginApi.PUT("", api.Controllers.ProtocolPluginApi.UpdateProtocolPlugin)

		// // 查
		protocolPluginApi.GET("", api.Controllers.ProtocolPluginApi.HandleProtocolPluginListByPage)

		// 根据设备id获取设备配置表单
		protocolPluginApi.GET("device_config_form", api.Controllers.ProtocolPluginApi.HandleProtocolPluginForm)

		// 根据协议类型和设备类型获取设备配置的配置表单
		protocolPluginApi.GET("config_form", api.Controllers.ProtocolPluginApi.HandleProtocolPluginFormByProtocolType)
	}
}
