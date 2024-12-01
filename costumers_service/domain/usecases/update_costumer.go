package usecases

import (
	"costumers-api/domain/abstractions"
	"costumers-api/domain/dependencies"
	"costumers-api/domain/domain_errors"
	"costumers-api/domain/validators"
	"errors"
)

type UpdateCostumersParams struct {
	Id    int
	Name  string
	Cpf   string
	Email string
	Phone string
}

type UpdateCostumerUseCase struct {
	repository dependencies.ICostumerRepository
	unitOfWork abstractions.IUnitOfWork
}

func (params UpdateCostumersParams) Validate() error {
	var (
		nameErr  error = nil
		cpfErr   error = nil
		emailErr error = nil
		phoneErr error = nil
	)

	idErr := validators.IsIdValid(params.Id)

	if len(params.Name) > 0 {
		nameErr = validators.IsNameValid(params.Name)
	}

	if len(params.Cpf) > 0 {
		cpfErr = validators.IsCpfValid(params.Cpf)
	}

	if len(params.Email) > 0 {
		emailErr = validators.IsEmailValid(params.Email)
	}

	if len(params.Phone) > 0 {
		phoneErr = validators.IsPhoneValid(params.Phone)
	}

	return errors.Join(idErr, nameErr, cpfErr, emailErr, phoneErr)
}

func NewUpdateCostumersUseCase(repository dependencies.ICostumerRepository, unitOfWork abstractions.IUnitOfWork) UpdateCostumerUseCase {
	return UpdateCostumerUseCase{
		repository: repository,
		unitOfWork: unitOfWork,
	}
}

func (service UpdateCostumerUseCase) Execute(params UpdateCostumersParams) (interface{}, error) {

	costumer, selectErr := service.repository.SelectCostumer(params.Id)

	if selectErr != nil {
		return nil, domain_errors.ErrCostumerNotFound
	}

	if len(params.Name) > 0 {
		costumer.Name = params.Name
	}

	if len(params.Cpf) > 0 {
		costumer.Cpf = params.Cpf
	}

	if len(params.Email) > 0 {
		costumer.SetEmail(params.Email)
	}

	if len(params.Phone) > 0 {
		costumer.SetPhone(params.Phone)
	}

	transaction := service.unitOfWork.GetTransaction()
	uptErr := service.repository.UpdateCostumer(transaction, costumer)

	return nil, uptErr
}
