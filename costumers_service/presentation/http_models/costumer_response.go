package httpmodels

import (
	"costumers-api/domain/models"
	"time"
)

type CostumerResponse struct {
	Id             int       `json:"id"`
	Name           string    `json:"name"`
	Cpf            string    `json:"cpf"`
	Birthdate      time.Time `json:"birthdate"`
	PhoneValidated bool      `json:"phoneValidated"`
	EmailValidated bool      `json:"emailValidated"`
	Phone          string    `json:"phone"`
	Email          string    `json:"email"`
}

func NewCostumerResponse(costumer models.Costumer) CostumerResponse {
	return CostumerResponse{
		Id:             costumer.Id,
		Name:           costumer.Name,
		Cpf:            costumer.Cpf,
		Birthdate:      costumer.Birthdate,
		PhoneValidated: costumer.PhoneValidated,
		EmailValidated: costumer.EmailValidated,
		Phone:          costumer.Phone(),
		Email:          costumer.Email(),
	}
}
