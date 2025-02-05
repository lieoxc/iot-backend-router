package apps

import (
	"project/internal/api"

	"github.com/gin-gonic/gin"
)

type OTA struct{}

func (*OTA) InitOTA(Router *gin.RouterGroup) {
	otaapi := Router.Group("ota")
	{
		upgradePackage := otaapi.Group("package")
		{
			upgradePackage.POST("", api.Controllers.OTAApi.CreateOTAUpgradePackage)
			upgradePackage.DELETE(":id", api.Controllers.OTAApi.DeleteOTAUpgradePackage)
			upgradePackage.PUT("", api.Controllers.OTAApi.UpdateOTAUpgradePackage)
			upgradePackage.GET("", api.Controllers.OTAApi.HandleOTAUpgradePackageByPage)
		}

		task := otaapi.Group("task")
		{
			task.POST("", api.Controllers.OTAApi.CreateOTAUpgradeTask)

			task.DELETE(":id", api.Controllers.OTAApi.DeleteOTAUpgradeTask)

			task.GET("", api.Controllers.OTAApi.HandleOTAUpgradeTaskByPage)

			task.GET("detail", api.Controllers.OTAApi.HandleOTAUpgradeTaskDetailByPage)

			task.PUT("detail", api.Controllers.OTAApi.UpdateOTAUpgradeTaskStatus)
		}
	}
}
