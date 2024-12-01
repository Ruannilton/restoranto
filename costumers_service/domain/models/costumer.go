package models

import (
	"time"
)

type Costumer struct {
	Id             int
	Name           string
	Cpf            string
	Birthdate      time.Time
	Deleted        bool
	PhoneValidated bool
	EmailValidated bool
	phone          string
	email          string
}

func NewCostumer(Id int, Name string, Cpf string, Birthdate time.Time, Deleted bool, PhoneValidated bool, EmailValidated bool, Phone string, Email string) Costumer {
	return Costumer{
		Id:             Id,
		phone:          Phone,
		email:          Email,
		Name:           Name,
		Cpf:            Cpf,
		Birthdate:      Birthdate,
		Deleted:        Deleted,
		PhoneValidated: PhoneValidated,
		EmailValidated: EmailValidated,
	}
}

func (costumer *Costumer) Phone() string {
	return costumer.phone
}

func (costumer *Costumer) SetPhone(phone string) {
	costumer.phone = phone
	costumer.PhoneValidated = false
}

func (costumer *Costumer) Email() string {
	return costumer.email
}

func (costumer *Costumer) SetEmail(email string) {
	costumer.email = email
	costumer.EmailValidated = false
}
