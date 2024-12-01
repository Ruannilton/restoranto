package usecases

import (
	"costumers-api/domain/abstractions"
	"costumers-api/domain/dependencies"
	"costumers-api/domain/validators"
)

type ValidateCostumerEmailParams struct {
	Id   int
	Code string
}

type ValidateCostumerEmailResponse struct {
	Valid bool
}

type ValidateCostumerEmailUseCase struct {
	costumerRepository dependencies.ICostumerRepository
	distributedCache   dependencies.IDistributedCacheRepository
	unitOfWork         abstractions.IUnitOfWork
}

func (params ValidateCostumerEmailParams) Validate() error {

	return validators.IsIdValid(params.Id)
}

func NewValidateCostumerEmailUseCase(
	costumerRepository dependencies.ICostumerRepository,
	distributedCache dependencies.IDistributedCacheRepository,
	unitOfWork abstractions.IUnitOfWork,
) ValidateCostumerEmailUseCase {
	return ValidateCostumerEmailUseCase{
		costumerRepository: costumerRepository,
		distributedCache:   distributedCache,
		unitOfWork:         unitOfWork,
	}
}

func (service ValidateCostumerEmailUseCase) Execute(params ValidateCostumerEmailParams) (ValidateCostumerEmailResponse, error) {

	costumer, selecError := service.costumerRepository.SelectCostumer(params.Id)

	if selecError != nil {
		return ValidateCostumerEmailResponse{}, nil
	}

	code, getCodeError := service.distributedCache.Get(costumer.Email())

	if getCodeError != nil {
		return ValidateCostumerEmailResponse{Valid: false}, nil
	}

	if code != params.Code {
		return ValidateCostumerEmailResponse{Valid: false}, nil
	}

	costumer.EmailValidated = true

	transaction := service.unitOfWork.GetTransaction()
	updateErr := service.costumerRepository.UpdateCostumer(transaction, costumer)

	if updateErr != nil {
		// TODO: notify error to try again
		// TODO: create apropriate error code
		return ValidateCostumerEmailResponse{}, updateErr
	}

	return ValidateCostumerEmailResponse{Valid: true}, nil
}
