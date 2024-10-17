package iface

type OfflineMessageNotificationRepo interface {
	Create(userID string) (bool, error)
	Remove(userID string) error
	UpdateStatus(userId string, status string) (bool, error)
}
