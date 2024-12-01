package infrastructure

import (
	"database/sql"
	"fmt"
)

type UnitOfWork struct {
	transaction *sql.Tx
}

func (uw UnitOfWork) GetTransaction() *sql.Tx {
	return uw.transaction
}

func (uw UnitOfWork) Commit() error {
	err := uw.transaction.Commit()

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Transaction commited")
	}
	return err
}

func (uw UnitOfWork) Rollback() error {
	return uw.transaction.Rollback()
}
