package usecases

import (
	"costumers-api/domain/abstractions"
	"costumers-api/domain/dependencies"
	"costumers-api/domain/validators"
)

type ValidateCostumerPhoneParams struct {
	Id   int
	Code string
}

type ValidateCostumerPhoneResponse struct {
	Valid bool
}

type ValidateCostumerPhoneUseCase struct {
	costumerRepository dependencies.ICostumerRepository
	distributedCache   dependencies.IDistributedCacheRepository
	unitOfWork         abstractions.IUnitOfWork
}

func (params ValidateCostumerPhoneParams) Validate() error {
	return validators.IsIdValid(params.Id)
}

func NewValidateCostumerPhoneUseCase(
	costumerRepository dependencies.ICostumerRepository,
	distributedCache dependencies.IDistributedCacheRepository,
	unitOfWork abstractions.IUnitOfWork,
) ValidateCostumerPhoneUseCase {
	return ValidateCostumerPhoneUseCase{
		costumerRepository: costumerRepository,
		distributedCache:   distributedCache,
		unitOfWork:         unitOfWork,
	}
}

func (service ValidateCostumerPhoneUseCase) Execute(params ValidateCostumerPhoneParams) (ValidateCostumerPhoneResponse, error) {

	costumer, selecError := service.costumerRepository.SelectCostumer(params.Id)

	if selecError != nil {
		return ValidateCostumerPhoneResponse{}, nil
	}

	code, getCodeError := service.distributedCache.Get(costumer.Phone())

	if getCodeError != nil {
		return ValidateCostumerPhoneResponse{Valid: false}, nil
	}

	if code != params.Code {
		return ValidateCostumerPhoneResponse{Valid: false}, nil
	}

	costumer.PhoneValidated = true
	transaction := service.unitOfWork.GetTransaction()
	updateErr := service.costumerRepository.UpdateCostumer(transaction, costumer)

	if updateErr != nil {
		// TODO: notify error to try again
		// TODO: create apropriate error code
		return ValidateCostumerPhoneResponse{}, updateErr
	}

	return ValidateCostumerPhoneResponse{Valid: true}, nil
}
