package usecases

import (
	"costumers-api/domain/abstractions"
	"costumers-api/domain/dependencies"
	"costumers-api/domain/validators"
)

type DeleteCostumerParams struct {
	Id int
}

type DeleteCostumerResponse struct {
	Deleted bool
}

type DeleteCostumerUseCase struct {
	repository dependencies.ICostumerRepository
	unitOfWork abstractions.IUnitOfWork
}

func (params DeleteCostumerParams) Validate() error {
	return validators.IsIdValid(params.Id)
}

func NewDeleteCostumerUseCase(repository dependencies.ICostumerRepository, unitOfWork abstractions.IUnitOfWork) DeleteCostumerUseCase {
	return DeleteCostumerUseCase{
		repository: repository,
		unitOfWork: unitOfWork,
	}
}

func (service DeleteCostumerUseCase) Execute(params DeleteCostumerParams) (DeleteCostumerResponse, error) {

	transaction := service.unitOfWork.GetTransaction()

	selectErr := service.repository.DeleteCostumer(transaction, params.Id)

	if selectErr != nil {
		return DeleteCostumerResponse{}, selectErr
	}

	return DeleteCostumerResponse{Deleted: true}, nil
}
