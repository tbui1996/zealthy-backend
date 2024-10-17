package model

type Chat struct {
	Session   string `json:"session"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
}

var ChatTypes = map[ChatType]string{
	CIRCULATOR: "internals_program_manager",
	GENERAL:    "internals_general_support",
}

var SupportTypes = map[string]ChatType{
	"internals_program_manager":   CIRCULATOR,
	"internals_general_support":   GENERAL,
	"internals_development_admin": CIRCULATOR, // this is a temporary solution for convinience so that admins can still do testing of program manager flows
}
