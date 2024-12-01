package usecases

import (
	"costumers-api/domain/dependencies"
	"costumers-api/domain/models"
	"costumers-api/domain/validators"
	"errors"

	"github.com/google/uuid"
)

type CreateEstablishmentParams struct {
	Name       string
	Cnpj       string
	CostumerId int
}

type CreateEstablishmentResponse struct {
	Id  int
	Key string
}

type CreateEstablishmentUseCase struct {
	establishmentRepository dependencies.IEstablishmentRepository
	costumersRepository     dependencies.ICostumerRepository
}

func (params CreateEstablishmentParams) Validate() error {
	nameErr := validators.IsNameValid(params.Name)
	cnpjErr := validators.IsCnpjValid(params.Cnpj)

	return errors.Join(nameErr, cnpjErr)
}

func (service CreateEstablishmentUseCase) Execute(params CreateEstablishmentParams) (CreateEstablishmentResponse, error) {

	if validationErr := params.Validate(); validationErr != nil {
		return CreateEstablishmentResponse{}, validationErr
	}

	costumer, findCostumerErr := service.costumersRepository.SelectCostumer(params.CostumerId)

	if findCostumerErr != nil {
		return CreateEstablishmentResponse{}, findCostumerErr
	}

	establishment := models.Establishment{
		CompanyName:      params.Name,
		Cnpj:             params.Cnpj,
		Costumer:         costumer,
		EstablishmentKey: uuid.NewString(),
	}

	id, err := service.establishmentRepository.Insert(establishment)

	if err != nil {
		return CreateEstablishmentResponse{}, err
	}

	return CreateEstablishmentResponse{
		Id:  id,
		Key: establishment.EstablishmentKey,
	}, nil
}
