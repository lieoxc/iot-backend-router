package apps

import (
	"project/internal/api"

	"github.com/gin-gonic/gin"
)

type NotificationServicesConfig struct{}

func (*NotificationServicesConfig) Init(Router *gin.RouterGroup) {
	url := Router.Group("notification/services/config")
	{
		// 创建/修改
		url.POST("", api.Controllers.NotificationServicesConfigApi.SaveNotificationServicesConfig)

		// 查询
		url.GET(":type", api.Controllers.NotificationServicesConfigApi.HandleNotificationServicesConfig)

		// 调试
		url.POST("e-mail/test", api.Controllers.NotificationServicesConfigApi.SendTestEmail)
	}
}
