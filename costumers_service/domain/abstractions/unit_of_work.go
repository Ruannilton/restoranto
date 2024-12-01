package abstractions

import "database/sql"

type IUnitOfWork interface {
	GetTransaction() *sql.Tx
	Commit() error
	Rollback() error
}
