package request

type DeleteFile struct {
	ID int `json:"id" validate:"required"`
}
