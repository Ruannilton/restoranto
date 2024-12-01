package usecases

import (
	"costumers-api/domain/dependencies"
	"costumers-api/domain/domain_errors"
	"costumers-api/domain/models"
)

type SearchCostumersParams struct {
	Name      string
	Cpf       string
	Email     string
	Phone     string
	Page      int
	PageCount int
}

type SearchCostumersResponse struct {
	Costumer []models.Costumer
	Page     int
	Count    int
}

type SearcCostumersUseCase struct {
	searchQuery dependencies.ISearchCostumerQuery
}

func (params SearchCostumersParams) Validate() error {
	if params.Page < 1 || params.PageCount <= 0 {
		return domain_errors.ErrInvalidParameters
	}

	return nil
}

func NewSearchCostumersUseCase(searchQuery dependencies.ISearchCostumerQuery) SearcCostumersUseCase {
	return SearcCostumersUseCase{
		searchQuery: searchQuery,
	}
}

func (service SearcCostumersUseCase) Execute(params SearchCostumersParams) (SearchCostumersResponse, error) {

	query := dependencies.SearchCostumerQueryParams{
		Name:      params.Name,
		Cpf:       params.Cpf,
		Email:     params.Email,
		Phone:     params.Phone,
		Page:      params.Page,
		PageCount: params.PageCount,
	}

	costumers, queryErr := service.searchQuery.SearchCostumers(query)

	if queryErr != nil {
		return SearchCostumersResponse{}, queryErr
	}

	return SearchCostumersResponse{
		Costumer: costumers,
		Page:     params.Page,
		Count:    len(costumers),
	}, nil
}
