package repositories

import (
	"context"
	"costumers-api/domain/models"
	"costumers-api/infrastructure/databases/costumers_db"
	"database/sql"
	"encoding/json"

	"github.com/Ruannilton/go-msg-contracts/pkg/events"
)

type EventRepository struct {
	queries *costumers_db.Queries
	ctx     context.Context
}

func NewEventRepository(db *sql.DB, ctx context.Context) EventRepository {
	queries := costumers_db.New(db)
	return EventRepository{
		queries: queries,
		ctx:     ctx,
	}
}

func (repo EventRepository) InsertEvent(transaction *sql.Tx, event events.IEvent) error {
	bytes, err := json.Marshal(event)

	if err != nil {
		//TODO: handle error
		return err
	}

	params := costumers_db.AddEventParams{
		EventName: event.GetEventName(),
		Message:   bytes,
	}

	query := repo.queries.WithTx(transaction)

	err = query.AddEvent(repo.ctx, params)

	if err != nil {
		//TODO: handle error
		return err
	}

	return nil
}

func (repo EventRepository) GetEvents(count int) ([]models.OutboxMessage, error) {
	events, err := repo.queries.GetEvents(repo.ctx, int32(count))

	if err != nil {
		//TODO: handle error
		return []models.OutboxMessage{}, err
	}

	mappedEvents := make([]models.OutboxMessage, len(events))

	for i, e := range events {
		mappedEvents[i] = models.OutboxMessage{
			Id:        int(e.ID),
			Name:      e.EventName,
			Message:   e.Message,
			CreatedAt: e.CreatedAt,
		}
	}

	return mappedEvents, nil
}

func (repo EventRepository) ProcessEvent(event models.OutboxMessage) error {
	err := repo.queries.ProcessEvent(repo.ctx, int64(event.Id))

	if err != nil {
		//TODO: handle error
		return err
	}

	return nil
}
