package repositories

import (
	"context"
	"costumers-api/domain/domain_errors"
	"costumers-api/domain/models"
	"costumers-api/infrastructure/databases/costumers_db"
	"database/sql"
	"errors"
)

type CostumerRepository struct {
	queries *costumers_db.Queries
	ctx     context.Context
}

func CreateCostumerRepository(db *sql.DB, ctx context.Context) CostumerRepository {

	queries := costumers_db.New(db)

	return CostumerRepository{
		queries: queries,
		ctx:     ctx,
	}
}

func (repo CostumerRepository) SelectCostumer(id int) (models.Costumer, error) {

	dbCostumer, dbErr := repo.queries.GetCostumer(repo.ctx, int32(id))

	if dbErr != nil {
		if errors.Is(dbErr, sql.ErrNoRows) {
			return models.Costumer{}, domain_errors.ErrCostumerNotFound
		}

		err := domain_errors.NewInfrastructureError(CostumersDbServiceName, GetErrorCode(dbErr), "failed to get costumer", dbErr)

		return models.Costumer{}, err
	}

	costumer := models.NewCostumer(
		int(dbCostumer.ID),
		dbCostumer.Name,
		dbCostumer.Cpf,
		dbCostumer.Birthdate,
		dbCostumer.Deleted,
		dbCostumer.PhoneValidated,
		dbCostumer.EmailValidated,
		dbCostumer.Phone,
		dbCostumer.Email,
	)

	return costumer, nil
}

func (repo CostumerRepository) InsertCostumer(transaction *sql.Tx, costumer models.Costumer) (int, error) {
	addParameters := costumers_db.AddCostumerParams{
		Name:      costumer.Name,
		Cpf:       costumer.Cpf,
		Birthdate: costumer.Birthdate,
		Phone:     costumer.Phone(),
		Email:     costumer.Email(),
	}

	query := repo.queries.WithTx(transaction)

	dbCostumer, dbErr := query.AddCostumer(repo.ctx, addParameters)

	if dbErr != nil {
		err := domain_errors.NewInfrastructureError(CostumersDbServiceName, GetErrorCode(dbErr), "failed to persist costumer", dbErr)
		return 0, err
	}

	return int(dbCostumer.ID), nil
}

func (repo CostumerRepository) UpdateCostumer(transaction *sql.Tx, costumer models.Costumer) error {

	updtParameters := costumers_db.UpdateCostumerParams{
		ID:             int32(costumer.Id),
		Name:           costumer.Name,
		Cpf:            costumer.Cpf,
		Birthdate:      costumer.Birthdate,
		Email:          costumer.Email(),
		Phone:          costumer.Phone(),
		PhoneValidated: costumer.PhoneValidated,
		EmailValidated: costumer.EmailValidated,
	}

	query := repo.queries.WithTx(transaction)

	dbErr := query.UpdateCostumer(repo.ctx, updtParameters)

	if dbErr != nil {
		err := domain_errors.NewInfrastructureError(CostumersDbServiceName, GetErrorCode(dbErr), "failed to update costumer", dbErr)
		return err
	}

	return nil
}

func (repo CostumerRepository) DeleteCostumer(transaction *sql.Tx, id int) error {

	query := repo.queries.WithTx(transaction)

	dbErr := query.DeleteCostumer(repo.ctx, int32(id))

	if dbErr != nil {
		return domain_errors.NewInfrastructureError(CostumersDbServiceName, GetErrorCode(dbErr), "failed to delete costumer", dbErr)
	}

	return nil
}
