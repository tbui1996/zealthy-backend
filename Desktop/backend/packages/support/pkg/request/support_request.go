package request

type SupportRequestSend struct {
	Action  string  `json:"action"`
	Payload Payload `json:"payload"`
}
type Payload struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type SupportRequestReceive struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
