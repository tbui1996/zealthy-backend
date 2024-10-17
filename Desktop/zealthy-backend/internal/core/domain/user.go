package domain

type Address struct {
	Street *string `json:"street,omitempty"`
	City   *string `json:"city,omitempty"`
	State  *string `json:"state,omitempty"`
	Zip    *string `json:"zip,omitempty"`
}

type User struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	Password  string   `json:"password,omitempty"`
	AboutMe   *string  `json:"aboutMe,omitempty"`
	Address   *Address `json:"address,omitempty"`
	Birthdate *string  `json:"birthdate,omitempty"`
}
