package model

type SessionDescriptor struct {
	SessionId int        `json:"sessionId"`
	Name      Descriptor `json:"name"`
	Value     string     `json:"value"`
}

type Descriptor string

const (
	TOPIC          Descriptor = "TOPIC"
	STARRED        Descriptor = "STARRED"
	RIDE_SCHEDULED Descriptor = "RIDE_SCHEDULED"
)
