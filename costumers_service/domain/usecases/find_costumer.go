package usecases

import (
	"costumers-api/domain/dependencies"
	"costumers-api/domain/models"
	"costumers-api/domain/validators"
)

type FindCostumerParams struct {
	Id int
}

type FindCostumerResponse struct {
	Costumer models.Costumer
}

type FindCostumerUseCase struct {
	repository dependencies.ICostumerRepository
}

func (params FindCostumerParams) Validate() error {
	return validators.IsIdValid(params.Id)
}

func NewFindCostumerUseCase(repository dependencies.ICostumerRepository) FindCostumerUseCase {
	return FindCostumerUseCase{
		repository: repository,
	}
}

func (service FindCostumerUseCase) Execute(params FindCostumerParams) (FindCostumerResponse, error) {

	costumer, selectError := service.repository.SelectCostumer(params.Id)

	if selectError != nil {
		return FindCostumerResponse{}, selectError
	}

	return FindCostumerResponse{Costumer: costumer}, nil
}
