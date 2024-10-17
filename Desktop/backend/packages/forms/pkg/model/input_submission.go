package model

type InputSubmission struct {
	ID               int    `json:"id"`
	FormSubmissionId int    `json:"formSubmissionId"`
	InputId          int    `json:"inputId"`
	Response         string `json:"response"`
}
