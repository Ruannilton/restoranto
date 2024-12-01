package dependencies

type IMessagePublisher interface {
	Publish(message any, queueName string) error
}
