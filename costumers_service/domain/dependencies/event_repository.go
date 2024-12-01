package dependencies

import (
	"costumers-api/domain/models"
	"database/sql"

	"github.com/Ruannilton/go-msg-contracts/pkg/events"
)

type IEventRepository interface {
	InsertEvent(transaction *sql.Tx, event events.IEvent) error
	GetEvents(count int) ([]models.OutboxMessage, error)
	ProcessEvent(event models.OutboxMessage) error
}
