package router

type Router interface {
	Send(input *RouterSendInput) error
}
