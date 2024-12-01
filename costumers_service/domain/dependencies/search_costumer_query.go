package dependencies

import "costumers-api/domain/models"

type SearchCostumerQueryParams struct {
	Name      string
	Cpf       string
	Email     string
	Phone     string
	Page      int
	PageCount int
}

type ISearchCostumerQuery interface {
	SearchCostumers(params SearchCostumerQueryParams) ([]models.Costumer, error)
}
