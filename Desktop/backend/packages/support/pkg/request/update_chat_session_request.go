package request

type UpdateChatSessionRequest struct {
	ID            string `json:"id"`
	Open          bool   `json:"open"`
	RideScheduled bool   `json:"rideScheduled"`
}
