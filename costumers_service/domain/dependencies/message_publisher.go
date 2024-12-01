package dependencies

import (
	"github.com/Ruannilton/go-msg-contracts/pkg/events"
)

type IMessagePublisher interface {
	Publish(message events.IEvent) error
}
