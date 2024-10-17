package request

// InputSubmissionRequest TODO:  Form Sent ID in this request
type InputSubmissionRequest struct {
	FormSentId int         `json:"formSentId" validate:"required"`
	SubmitData []InputData `json:"submitData" validate:"required"`
}

type InputData struct {
	ID       int    `json:"id" validate:"required"`
	Response string `json:"response" validate:"required"`
}
