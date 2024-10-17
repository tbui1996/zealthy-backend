package dynamo

type UnconfirmedConnectionItem struct {
	ConnectionId string
	Email        string
}

type ConnectionItem struct {
	ConnectionId string
	UserID       string
}
