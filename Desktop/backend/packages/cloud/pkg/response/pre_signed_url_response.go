package response

type PreSignedUrlResponse struct {
	URL string `json:"url"`
	Key string `json:"key"`
}
