// /initialize/croninit/cron.go
package croninit

import (
	"project/internal/service"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

var (
	c = cron.New()
)

// 定义任务初始化
func CronInit() {

	// 初始化设备统计定时任务
	//InitDeviceStatsCron(c)

	// 每天凌晨2点执行数据清理
	c.AddFunc("0 2 * * *", func() {
		logrus.Debug("系统数据清理任务开始：")
		service.GroupApp.CleanSystemDataByCron()
	})

	c.Start()
}
