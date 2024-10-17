package request

type FeedbackRequest struct {
	Email         *string `json:"email"`
	UserName      string  `json:"userName"`
	Activity      string  `json:"activity"`
	ActivityNotes string  `json:"activityNotes"`
	Suggestion    string  `json:"suggestion"`
}
