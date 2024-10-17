package chatHelper

import (
	"encoding/json"
	"fmt"

	"github.com/circulohealth/sonar-backend/packages/support/pkg/request"
	"github.com/circulohealth/sonar-backend/packages/support/pkg/session/iface"
)

var HandleReadReceipt = func(payload string, time int64, repo iface.ChatSessionRepository) error {
	var readReceiptRequest request.ReadReceiptRequest
	if err := json.Unmarshal([]byte(payload), &readReceiptRequest); err != nil {
		return fmt.Errorf("failed to unmarshal read receipt request payload")
	}

	sess, _, err := repo.GetEntityWithStatus(readReceiptRequest.SessionID)

	if err != nil {
		return fmt.Errorf("failed to get chat session container to determine if this session is pending or not pending (%s)", err.Error())
	}

	return repo.HandleReadReceipt(sess, time, readReceiptRequest)
}
