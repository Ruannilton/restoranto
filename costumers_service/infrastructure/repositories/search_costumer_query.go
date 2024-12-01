package repositories

import (
	"context"
	"costumers-api/domain/dependencies"
	"costumers-api/domain/domain_errors"
	"costumers-api/domain/models"
	"costumers-api/infrastructure/databases/costumers_db"
	"database/sql"
)

type SearchCostumerQuery struct {
	queries *costumers_db.Queries
	ctx     context.Context
}

func CreateSearchCostumerQuery(db *sql.DB, ctx context.Context) SearchCostumerQuery {

	queries := costumers_db.New(db)

	return SearchCostumerQuery{
		queries: queries,
		ctx:     ctx,
	}
}

func (query SearchCostumerQuery) SearchCostumers(params dependencies.SearchCostumerQueryParams) ([]models.Costumer, error) {

	queryParams := costumers_db.SearchCostumersParams{
		Name:      params.Name,
		Cpf:       params.Cpf,
		Email:     params.Email,
		Phone:     params.Phone,
		Page:      int32(params.Page),
		Pagecount: int32(params.PageCount),
	}

	dbCostumers, dbErr := query.queries.SearchCostumers(query.ctx, queryParams)

	if dbErr != nil {
		return nil, domain_errors.NewInfrastructureError(CostumersDbServiceName, GetErrorCode(dbErr), "failed to search costumers", dbErr)
	}

	mapedCostumers := make([]models.Costumer, len(dbCostumers))

	for i, a := range dbCostumers {
		mapedCostumers[i] = models.NewCostumer(int(a.ID), a.Name, a.Cpf, a.Birthdate, a.Deleted, a.PhoneValidated, a.EmailValidated, a.Phone, a.Email)
	}

	return mapedCostumers, nil
}
