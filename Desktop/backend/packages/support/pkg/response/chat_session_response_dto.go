package response

import "time"

// Holds common properties that can be returned to the client
type ChatSessionResponseDTO struct {
	ID                   string          `json:"ID"`
	InternalUserID       string          `json:"internalUserID"`
	Email                string          `json:"email"`
	UserID               string          `json:"userID"`
	Topic                string          `json:"topic"`
	Patient              PatientResponse `json:"patient"`
	ChatType             string          `json:"chatType"`
	Notes                *string         `json:"notes"`
	ChatOpen             bool            `json:"chatOpen"`
	CreatedTimestamp     int64           `json:"createdTimestamp"`
	UserLastRead         int64           `json:"userLastRead"`
	InternalUserLastRead int64           `json:"internalUserLastRead"`
	IsPending            bool            `json:"isPending"`
	LastMessageTimestamp int64           `json:"lastMessageTimestamp"`
	LastMessagePreview   string          `json:"lastMessagePreview"`
	LastMessageSenderID  string          `json:"lastMessageSenderID"`
	Starred              bool            `json:"starred"`
}

type PatientResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	LastName    string    `json:"lastName"`
	Address     string    `json:"address"`
	InsuranceID string    `json:"insuranceID"`
	Birthday    time.Time `json:"birthday"`
}
