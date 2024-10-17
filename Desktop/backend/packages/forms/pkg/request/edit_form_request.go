package request

type EditForm struct {
	ID int `json:"id" validate:"required"`
	// Title must be at least 4 characters
	Title string `json:"title" validate:"required,min=1"`
	// description is not required
	Description string `json:"description"`
}
