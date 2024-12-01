package usecases

import (
	"costumers-api/domain/abstractions"
	"costumers-api/domain/dependencies"
	"costumers-api/domain/domain_errors"
	"costumers-api/domain/models"
	"costumers-api/domain/validators"
	"errors"
	"strings"
	"time"

	"github.com/Ruannilton/go-msg-contracts/pkg/events"
)

type CreateCostumerParams struct {
	Name      string
	Cpf       string
	Birthdate time.Time
	Phone     string
	Email     string
}

type CreateCostumerResponse struct {
	Id int
}

type CreateCostumerUseCase struct {
	costumerRepository dependencies.ICostumerRepository
	eventRepository    dependencies.IEventRepository
	unityOfWork        abstractions.IUnitOfWork
}

func (param CreateCostumerParams) Validate() error {
	nameErr := validators.IsNameValid(param.Name)
	cpfErr := validators.IsCpfValid(param.Cpf)
	birthErr := validators.IsBirthdateValid(param.Birthdate)
	phoneErr := validators.IsPhoneValid(param.Phone)
	emailErr := validators.IsEmailValid(param.Email)

	err := errors.Join(nameErr, cpfErr, birthErr, phoneErr, emailErr)

	if err != nil {
		valErr := domain_errors.NewValidationError(err)

		if nameErr != nil {
			valErr.AddField("name", nameErr)
		}

		if cpfErr != nil {
			valErr.AddField("cpf", cpfErr)
		}

		if birthErr != nil {
			valErr.AddField("birthdate", birthErr)
		}

		if phoneErr != nil {
			valErr.AddField("phone", phoneErr)
		}

		if emailErr != nil {
			valErr.AddField("email", emailErr)
		}

		return valErr
	}

	return nil
}

func NewCreateCostumerUseCase(repository dependencies.ICostumerRepository, eventRepository dependencies.IEventRepository, unitOfWork abstractions.IUnitOfWork) CreateCostumerUseCase {
	return CreateCostumerUseCase{
		costumerRepository: repository,
		eventRepository:    eventRepository,
		unityOfWork:        unitOfWork,
	}
}

func (service CreateCostumerUseCase) Execute(param CreateCostumerParams) (CreateCostumerResponse, error) {

	param.Cpf = strings.ReplaceAll(param.Cpf, ".", "")
	param.Cpf = strings.ReplaceAll(param.Cpf, "-", "")

	costumer := models.Costumer{
		Name:      param.Name,
		Cpf:       param.Cpf,
		Birthdate: param.Birthdate,
	}

	costumer.SetEmail(param.Email)
	costumer.SetPhone(param.Phone)

	transaction := service.unityOfWork.GetTransaction()
	id, insertCostumerError := service.costumerRepository.InsertCostumer(transaction, costumer)

	if insertCostumerError != nil {
		return CreateCostumerResponse{}, insertCostumerError
	}

	createdCostumerEvent := events.CostumerCreatedEvent{
		Id:        id,
		Name:      costumer.Name,
		Cpf:       costumer.Cpf,
		Birthdate: costumer.Birthdate,
		Phone:     costumer.Phone(),
		Email:     costumer.Email(),
	}

	insertEventError := service.eventRepository.InsertEvent(transaction, createdCostumerEvent)

	if insertEventError != nil {
		return CreateCostumerResponse{}, insertEventError
	}

	return CreateCostumerResponse{Id: id}, nil
}
