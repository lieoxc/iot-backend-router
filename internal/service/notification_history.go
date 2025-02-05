package service

import (
	dal "project/internal/dal"
	model "project/internal/model"
	"project/pkg/errcode"
)

type NotificationHisory struct{}

// NotificationHistory orm define:
// type NotificationHistory struct {
// 	ID               string    `gorm:"column:id;primaryKey" json:"id"`
// 	SendTime         time.Time `gorm:"column:send_time;not null" json:"send_time"`
// 	SendContent      *string   `gorm:"column:send_content" json:"send_content"`
// 	SendTarget       string    `gorm:"column:send_target;not null" json:"send_target"`
// 	SendResult       *string   `gorm:"column:send_result" json:"send_result"`
// 	NotificationType string    `gorm:"column:notification_type;not null" json:"notification_type"`
// 	TenantID         string    `gorm:"column:tenant_id;not null" json:"tenant_id"`
// 	Remark           *string   `gorm:"column:remark" json:"remark"`
// }

func (*NotificationHisory) GetNotificationHistoryListByPage(pageParam *model.GetNotificationHistoryListByPageReq) (map[string]interface{}, error) {
	total, list, err := dal.GetNotificationHisoryListByPage(pageParam)
	if err != nil {
		return nil, errcode.WithData(errcode.CodeDBError, map[string]interface{}{
			"sql_error": err.Error(),
		})
	}
	notificationListRsp := make(map[string]interface{})
	notificationListRsp["total"] = total
	notificationListRsp["list"] = list

	return notificationListRsp, err
}

func (*NotificationHisory) SaveNotificationHistory(req *model.NotificationHistory) error {
	err := dal.CreateNotificationHistory(req)
	if err != nil {
		return err
	}
	return nil

}
