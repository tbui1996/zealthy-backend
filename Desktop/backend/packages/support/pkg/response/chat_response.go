package response

type InternalChatResponse struct {
	Sender    string `json:"sender"`
	Session   string `json:"session"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
	File      string `json:"file"`
}
