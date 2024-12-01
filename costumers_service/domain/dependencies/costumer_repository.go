package dependencies

import (
	"costumers-api/domain/models"
	"database/sql"
)

type ICostumerRepository interface {
	SelectCostumer(id int) (models.Costumer, error)
	InsertCostumer(transaction *sql.Tx, costumer models.Costumer) (int, error)
	UpdateCostumer(transaction *sql.Tx, costumer models.Costumer) error
	DeleteCostumer(transaction *sql.Tx, id int) error
}
