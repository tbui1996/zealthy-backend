package request

type Chat struct {
	Session string `json:"session"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
	File    string `json:"file"`
}
