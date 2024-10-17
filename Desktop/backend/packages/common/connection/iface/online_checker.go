package iface

import "github.com/circulohealth/sonar-backend/packages/common/connection/dto"

type OnlineChecker interface {
	IsUserOnline(userID string) (dto.UserOnlineStatus, error)
}
