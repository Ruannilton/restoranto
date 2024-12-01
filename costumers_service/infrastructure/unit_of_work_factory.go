package infrastructure

import (
	"context"
	"costumers-api/domain/abstractions"
	"database/sql"
)

type UnitOfWorkFactory struct {
	db *sql.DB
}

func NewUnitOfWorkFactory(db *sql.DB) UnitOfWorkFactory {
	return UnitOfWorkFactory{
		db: db,
	}
}

func (factory UnitOfWorkFactory) NewUnitOfWork(ctx context.Context) (abstractions.IUnitOfWork, error) {
	tx, err := factory.db.BeginTx(ctx, nil)

	if err != nil {
		return UnitOfWork{}, err
	}

	return UnitOfWork{
		transaction: tx,
	}, nil
}
