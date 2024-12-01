package domain_errors

import (
	"errors"
)

var (
	ErrInvalidId         = errors.New("invalid id")
	ErrInvalidCnpj       = errors.New("invalid cnpj")
	ErrInvalidCpf        = errors.New("invalid cpf")
	ErrInvalidName       = errors.New("invalid name")
	ErrInvalidPhone      = errors.New("invalid phone")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrUnderage          = errors.New("underage")
	ErrInvalidBirthdate  = errors.New("invalid birthdate")
	ErrInvalidParameters = errors.New("invalid parameter(s)")
	ErrCostumerNotFound  = errors.New("costumer not found")
)

const (
	CodeTimeOut          = "TIME_OUT"
	CodeConnectionErr    = "CONNECTION_ERR"
	CodeInvalidOperation = "INVALID_OPERATION"
	CodeEntryNotFound    = "ENTRY_NOT_FOUND"
)
