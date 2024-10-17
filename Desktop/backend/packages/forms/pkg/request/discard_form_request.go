package request

type DiscardFormRequest struct {
	FormSentId int `json:"formSentId" validate:"required"`
}
