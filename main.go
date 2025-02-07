package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"project/initialize"
	"project/internal/app"
	"project/internal/query"
	"project/mqtt"
	"project/mqtt/device"
	"project/mqtt/publish"
	"project/mqtt/subscribe"
	"time"

	router "project/router"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	initialize.ViperInit("./configs/conf.yml")
	//initialize.ViperInit("./configs/conf-localdev.yml")
	initialize.RsaDecryptInit("./configs/rsa_key/private_key.pem")
	initialize.LogInIt()
	db, err := initialize.PgInit()
	if err != nil {
		logrus.Fatal(err)
	}
	_, err = initialize.RedisInit()
	if err != nil {
		logrus.Fatal(err)
	}

	query.SetDefault(db)

	err = mqtt.MqttInit()
	if err != nil {
		logrus.Fatal(err)
	}
	go device.InitDeviceStatus()
	err = subscribe.SubscribeInit()
	if err != nil {
		logrus.Fatal(err)
	}
	publish.PublishInit()
	//定时任务
	//croninit.CronInit()
}

// @title           ThingsPanel API
// @version         1.0
// @description     ThingsPanel API.
// @schemes         http
// @host      localhost:9999
// @BasePath
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        x-token
func main() {
	// 初始化服务管理器 (订阅redis过期事件，然后发送对应的MQTT消息)
	manager := app.NewManager()
	if err := manager.Start(); err != nil {
		logrus.Fatalf("Failed to start services: %v", err)
	}
	defer manager.Stop()
	// gin.SetMode(gin.ReleaseMode)

	// TODO: 替换gin默认日志，默认日志不支持日志级别设置
	host, port := loadConfig()
	router := router.RouterInit()
	srv := initServer(host, port, router)

	// 启动服务
	go startServer(srv, host, port)

	// 优雅关闭
	gracefulShutdown(srv)

}

func loadConfig() (host, port string) {
	host = viper.GetString("service.http.host")
	if host == "" {
		host = "localhost"
		logrus.Println("Using default host:", host)
	}

	port = viper.GetString("service.http.port")
	if port == "" {
		port = "9999"
		logrus.Println("Using default port:", port)
	}

	return host, port
}

func initServer(host, port string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         net.JoinHostPort(host, port),
		Handler:      handler,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
}

func startServer(srv *http.Server, host, port string) {
	logrus.Println("Listening and serving HTTP on", host, ":", port)
	successInfo()
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Fatalf("listen: %s\n", err)
	}
}

func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logrus.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server Shutdown:", err)
	}
	logrus.Println("Server exiting")
}

func successInfo() {
	// 获取当前时间
	startTime := time.Now().Format("2006-01-02 15:04:05")

	// 打印启动成功消息
	fmt.Println("----------------------------------------")
	fmt.Println("        TingsPanel 启动成功!")
	fmt.Println("----------------------------------------")
	fmt.Printf("启动时间: %s\n", startTime)
	fmt.Println("版本: v1.1.4社区版")
	fmt.Println("----------------------------------------")
	fmt.Println("欢迎使用 TingsPanel！")
	fmt.Println("如需帮助，请访问: http://thingspanel.io")
	fmt.Println("----------------------------------------")
}
