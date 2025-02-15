package router

import (
	"net/http"
	middleware "project/internal/middleware"
	"project/internal/middleware/response"
	"project/pkg/global"
	"project/router/apps"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	// gin-swagger middleware

	api "project/internal/api"
)

// swagger embed files

func RouterInit() *gin.Engine {
	//gin.SetMode(gin.ReleaseMode) //开启生产模式
	gin.DefaultWriter = logrus.StandardLogger().Out
	gin.DefaultErrorWriter = logrus.StandardLogger().Out
	router := gin.Default()
	// 提供 Vue 编译后的静态文件
	router.Static("/assets", "./dist/assets") // 将 static 目录中的文件暴露为静态资源
	router.LoadHTMLFiles("./dist/index.html")
	router.StaticFile("/favicon.ico", "./dist/favicon.ico")             // ico
	router.StaticFile("/EasyWasmPlayer.js", "./dist/EasyWasmPlayer.js") // ico
	router.StaticFile("/libDecoder.wasm ", "./dist/libDecoder.wasm")    // ico
	// 默认路由：返回 index.html
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil) // 返回编译后的 index.html
	})

	// 处理文件访问请求
	router.GET("/files/*filepath", func(c *gin.Context) {
		filepath := c.Param("filepath")
		c.File("./files" + filepath)
	})

	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(middleware.Cors())
	// 初始化响应处理器
	handler, err := response.NewHandler("configs/messages.yaml")
	if err != nil {
		logrus.Fatalf("初始化响应处理器失败: %v", err)
	}
	// 全局使用
	global.ResponseHandler = handler
	// 使用中间件
	router.Use(handler.Middleware())

	controllers := new(api.Controller)
	// 健康检查
	router.GET("/health", controllers.SystemApi.HealthCheck)

	api := router.Group("api")
	{
		// 无需权限校验
		v1 := api.Group("v1")
		{
			v1.POST("plugin/heartbeat", controllers.Heartbeat)
			v1.POST("plugin/device/config", controllers.HandleDeviceConfigForProtocolPlugin)
			v1.POST("plugin/service/access/list", controllers.HandlePluginServiceAccessList)
			v1.POST("plugin/service/access", controllers.HandlePluginServiceAccess)
			v1.POST("login", controllers.Login)
			v1.GET("verification/code", controllers.HandleVerificationCode)
			v1.POST("reset/password", controllers.ResetPassword)
			v1.GET("logo", controllers.HandleLogoList)
			// 设备遥测（ws）
			v1.GET("telemetry/datas/current/ws", controllers.TelemetryDataApi.ServeCurrentDataByWS)
			// 设备在线离线状态（ws）
			v1.GET("device/online/status/ws", controllers.TelemetryDataApi.ServeDeviceStatusByWS)
			// 设备遥测keys（ws）
			v1.GET("telemetry/datas/current/keys/ws", controllers.TelemetryDataApi.ServeCurrentDataByKey)
			v1.GET("ota/download/files/upgradePackage/:path/:file", controllers.OTAApi.DownloadOTAUpgradePackage)
			// 获取系统时间
			v1.GET("systime", controllers.SystemApi.HandleSystime)
			// 查询系统功能设置
			v1.GET("sys_function", controllers.SysFunctionApi.HandleSysFcuntion)
			// 租户邮箱注册
			v1.POST("/tenant/email/register", controllers.UserApi.EmailRegister)
			// 网关自动注册
			v1.POST("/device/gateway-register", controllers.DeviceApi.GatewayRegister)
			// 网关子设备注册
			v1.POST("/device/gateway-sub-register", controllers.DeviceApi.GatewaySubRegister)
		}
		// 需要权限校验
		v1.Use(middleware.JWTAuth())

		// 需要权限校验
		v1.Use(middleware.CasbinRBAC())
		// SSE服务
		SSERouter(v1)
		// 记录操作日志
		v1.Use(middleware.OperationLogs())
		{
			apps.Model.User.InitUser(v1) // 用户模块

			apps.Model.Role.Init(v1) // 角色管理

			apps.Model.Casbin.Init(v1) // 权限管理

			apps.Model.Dict.InitDict(v1) // 字典模块

			apps.Model.OTA.InitOTA(v1) // OTA模块

			apps.Model.UpLoad.Init(v1) // 文件上传

			apps.Model.ProtocolPlugin.InitProtocolPlugin(v1) // 协议插件模块

			apps.Model.Device.InitDevice(v1) // 设备

			apps.Model.UiElements.Init(v1) // UI元素控制

			apps.Model.Board.InitBoard(v1) // 首页

			apps.Model.EventData.InitEventData(v1) // 事件数据

			apps.Model.TelemetryData.InitTelemetryData(v1) // 遥测数据

			apps.Model.AttributeData.InitAttributeData(v1) // 属性数据

			apps.Model.CommandData.InitCommandData(v1) //命令数据

			apps.Model.OperationLog.Init(v1) // 操作日志

			apps.Model.Logo.Init(v1) // logo

			apps.Model.DataPolicy.Init(v1) // 数据清理

			apps.Model.DeviceConfig.Init(v1) // 设备配置

			apps.Model.DataScript.Init(v1) // 数据处理脚本

			apps.Model.NotificationGroup.InitNotificationGroup(v1) // 通知组

			apps.Model.NotificationHistoryGroup.InitNotificationHistory(v1) // 通知组

			apps.Model.NotificationServicesConfig.Init(v1) // 通知服务配置

			apps.Model.Alarm.Init(v1) // 告警模块

			apps.Model.Scene.Init(v1) // 场景

			apps.Model.SceneAutomations.Init(v1) // 场景联动

			apps.Model.SysFunction.Init(v1) //功能设置

			apps.Model.ServicePlugin.Init(v1) // 插件管理

			apps.Model.ExpectedData.InitExpectedData(v1)
		}
	}

	return router
}
