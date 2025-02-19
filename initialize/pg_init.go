package initialize

import (
	"context"
	"fmt"
	"os"
	"time"

	global "project/pkg/global"
	utils "project/pkg/utils"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 数据库配置
type DbConfig struct {
	Host          string
	Port          int
	DbName        string
	Username      string
	Password      string
	TimeZone      string
	LogLevel      int
	SlowThreshold int
	IdleConns     int
	OpenConns     int
}

func PgInit(cfgPath string) (*gorm.DB, error) {
	// 初始化配置
	config, err := LoadDbConfig()
	if err != nil {
		logrus.Errorf("加载数据库配置失败: %v", err)
		return nil, err
	}

	// 初始化数据库
	db, err := PgConnect(config)
	if err != nil {
		logrus.Error("连接数据库失败:", err)
		return nil, err
	}
	global.DB = db

	// casbin 初始化
	CasbinInit(cfgPath)

	// 检查版本
	err = CheckVersion(db, cfgPath)
	if err != nil {
		fmt.Println(err)
	}

	return db, nil
}

// LoadDbConfig 从配置文件加载数据库配置
func LoadDbConfig() (*DbConfig, error) {
	config := &DbConfig{
		Host:          viper.GetString("db.psql.host"),
		Port:          viper.GetInt("db.psql.port"),
		DbName:        viper.GetString("db.psql.dbname"),
		Username:      viper.GetString("db.psql.username"),
		Password:      viper.GetString("db.psql.password"),
		TimeZone:      viper.GetString("db.psql.time_zone"),
		LogLevel:      viper.GetInt("db.psql.log_level"),
		SlowThreshold: viper.GetInt("db.psql.slow_threshold"),
		IdleConns:     viper.GetInt("db.psql.idle_conns"),
		OpenConns:     viper.GetInt("db.psql.open_conns"),
	}

	// 设置默认值
	if config.Host == "" {
		config.Host = "localhost"
	}
	if config.Port == 0 {
		config.Port = 5432
	}
	if config.TimeZone == "" {
		config.TimeZone = "Asia/Shanghai"
	}
	if config.LogLevel == 0 {
		config.LogLevel = 1
	}
	if config.SlowThreshold == 0 {
		config.SlowThreshold = 200
	}
	if config.IdleConns == 0 {
		config.IdleConns = 10
	}
	if config.OpenConns == 0 {
		config.OpenConns = 50
	}

	// 检查必要的配置
	if config.DbName == "" || config.Username == "" || config.Password == "" {
		return nil, fmt.Errorf("database configuration is incomplete")
	}

	return config, nil
}

// GormLogger 自定义 GORM 日志器
type GormLogger struct {
	LogLevel logger.LogLevel
}

// LogMode 设置日志级别
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l
}

// Info 输出 Info 级别日志
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		logrus.WithContext(ctx).Infof(msg, data...)
	}
}

// Warn 输出 Warn 级别日志
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		logrus.WithContext(ctx).Warnf(msg, data...)
	}
}

// Error 输出 Error 级别日志
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		logrus.WithContext(ctx).Errorf(msg, data...)
	}
}

// Trace 输出 Trace 级别日志
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	logEntry := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"elapsed": elapsed,
		"rows":    rows,
	})

	switch {
	case err != nil && l.LogLevel >= logger.Error:
		logEntry.Errorf("SQL: %s, Error: %v", sql, err)
	case l.LogLevel >= logger.Info:
		logEntry.Infof("SQL: %s", sql)
	}
}

// PgInit 初始化数据库连接
func PgConnect(config *DbConfig) (*gorm.DB, error) {
	dataSource := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable TimeZone=%s",
		config.Host, config.Port, config.DbName, config.Username, config.Password, config.TimeZone)

	var err error
	db, err := gorm.Open(postgres.Open(dataSource), &gorm.Config{
		Logger:                 &GormLogger{LogLevel: logger.Info},
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false, // use singular table name, table for `User` would be `user` with this option enabled
		},
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取原生数据库连接失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(config.IdleConns)
	sqlDB.SetMaxOpenConns(config.OpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logrus.Println("连接数据库完成...")

	return db, nil
}

/*
注意 sql中不要有sys_version表
1. 检查版本表是否存在: 检查数据库版本，如果没有sys_version表，创建sys_version表，插入版本序号0，版本号0.0.0
2. 程序版本低于数据版本: 提示升级
3. 数据版本低于程序版本: 执行sql文件，更新版本号
*/
// 检查版本，在表sys_version中的version字段
func CheckVersion(db *gorm.DB, cfgPath string) error {
	version := global.VERSION
	versionNumber := global.VERSION_NUMBER // 当前程序版本号
	var dataVersionNumber int              // 数据库版本号

	// 判断有没有sys_version的表
	var exists bool
	result := db.Raw("SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema='public' AND table_name='sys_version')").Scan(&exists)
	if result.Error != nil {
		return result.Error
	}
	// 创建事务
	logrus.Info("----", exists)
	if !exists { // 如果不存在sys_version表，创建sys_version表
		logrus.Info("创建sys_version表")
		dataVersionNumber = 0
		t := db.Exec("CREATE TABLE sys_version (version_number INT NOT NULL DEFAULT 0, version varchar(255) NOT NULL, PRIMARY KEY (version_number))")
		if t.Error != nil {
			return t.Error
		}

	}
	tx := db.Begin()
	// 查询版本号
	result = db.Table("sys_version").Select("version_number").Scan(&dataVersionNumber)
	if result.Error != nil {
		return result.Error
	}
	// 如果版本号为空，插入版本号
	if dataVersionNumber == 0 {
		t := tx.Exec("INSERT INTO sys_version (version_number, version) VALUES (?, ?)", 0, "0.0.0")
		if t.Error != nil {
			// 回滚
			tx.Rollback()
			return t.Error
		}
	}
	if dataVersionNumber > global.VERSION_NUMBER {
		// 回滚
		tx.Rollback()
		return fmt.Errorf("当前数据版本高于程序版本，请升级程序")
	} else if dataVersionNumber < global.VERSION_NUMBER {
		logrus.Println("数据版本：", dataVersionNumber)
		logrus.Println("程序版本：", global.VERSION_NUMBER)
		logrus.Println("开始升级...")
		// sql文件名为：版本编号.sql，执行所大于当前数据版本小于等于程序版本的sql文件
		for i := dataVersionNumber + 1; i <= global.VERSION_NUMBER; i++ {
			fileName := fmt.Sprintf(cfgPath+"/sql/%d.sql", i)
			// 检查文件是否存在
			if !utils.FileExist(fileName) {
				// 回滚
				tx.Rollback()
				return fmt.Errorf("sql文件不存在,可能需要手动升级：%s", fileName)
			}
			logrus.Println("执行sql文件：", fileName)
			// 读取 SQL 脚本文件
			sqlFile, err := os.ReadFile(fileName)
			if err != nil {
				panic(err)
			}
			fmt.Println("执行sql脚本...")
			// 执行 SQL 脚本
			t := tx.Exec(string(sqlFile))
			if t.Error != nil {
				// 回滚
				tx.Rollback()
				return t.Error
			}
		}
		// 更新版本号
		t := tx.Exec("UPDATE sys_version SET version_number = ?, version = ?", versionNumber, version)
		if t.Error != nil {
			// 回滚
			tx.Rollback()
			return t.Error
		}
		logrus.Println("升级成功")
	}
	return tx.Commit().Error
}

func ExecuteSQLFile(db *gorm.DB, fileName string) error {
	// 读取 SQL 脚本文件
	sqlFile, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	// 执行 SQL 脚本
	t := db.Exec(string(sqlFile))
	if t.Error != nil {
		return t.Error
	}

	return nil
}
