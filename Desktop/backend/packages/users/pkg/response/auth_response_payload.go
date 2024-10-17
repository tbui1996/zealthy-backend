package response

type AuthResponsePayload struct {
	ID           string `json:"id"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken,omitempty"`
}
