package dependencies

type IMessageInput interface {
	GetContent() []byte
	Ack() error
	Nack(requeue bool) error
}

type IMessageListener interface {
	Listen(queueName string) (chan IMessageInput, error)
	Close() error
}
