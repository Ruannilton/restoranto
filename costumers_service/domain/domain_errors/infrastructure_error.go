package domain_errors

import "fmt"

type InfrastructureError struct {
	Service  string
	Code     string
	Message  string
	InnerErr error
}

func NewInfrastructureError(service string, code string, message string, err error) InfrastructureError {
	return InfrastructureError{
		Service:  service,
		Code:     code,
		Message:  message,
		InnerErr: err,
	}
}

func (e InfrastructureError) Error() string {
	return fmt.Sprintf("%s: code: %s\nmessage: %s\nerror: (%v)", e.Service, e.Code, e.Message, e.InnerErr)
}
