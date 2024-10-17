package request

type FileUploadRequest struct {
	ChatId   string `json:"chatId"`
	FileId   string `json:"fileId"`
	Filename string `json:"filename"`
}
