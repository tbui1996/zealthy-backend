package request

type RouterTypeRequest struct {
	Message string `json:"message,omitempty"`
	Type    string `json:"type"`
}
