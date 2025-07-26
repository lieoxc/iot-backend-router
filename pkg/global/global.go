package global

import (
	"project/internal/middleware/response"

	"github.com/casbin/casbin/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var VERSION = "0.0.7"
var VERSION_NUMBER = 7
var DB *gorm.DB
var REDIS *redis.Client
var STATUS_REDIS *redis.Client
var CasbinEnforcer *casbin.Enforcer
var OtaAddress string
var TPSSEManager *SSEManager
var ResponseHandler *response.Handler

type EventData struct {
	Name    string
	Message string
}

// 事件通道
var EventChan chan EventData
