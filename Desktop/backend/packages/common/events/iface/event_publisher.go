package iface

type EventPublisher interface {
	PublishConnectionCreatedEvent(userId string, service string) error
}
