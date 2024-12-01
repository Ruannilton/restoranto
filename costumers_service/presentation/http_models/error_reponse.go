package httpmodels

import (
	"costumers-api/domain/domain_errors"
	"errors"
	"net/http"
)

type ValidationError struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

type ValidationProblemDetails struct {
	Type   string            `json:"type"`
	Title  string            `json:"title"`
	Errors []ValidationError `json:"errors"`
}

type ProblemDetails struct {
	Type     string `json:"type"`
	Status   int    `json:"status"`
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

func NewProblemDetails(err error) ProblemDetails {
	var infErr domain_errors.InfrastructureError

	if errors.As(err, &infErr) {
		return ProblemDetails{
			Type:   "about:blank",
			Title:  infErr.Service + ": " + infErr.Code,
			Detail: infErr.Message,
			Status: http.StatusBadRequest,
		}
	}

	return ProblemDetails{
		Type:   "about:blank",
		Title:  "unknow",
		Detail: err.Error(),
		Status: http.StatusBadRequest,
	}
}

func NewValidationProblemDetails(err domain_errors.ValidationError) ValidationProblemDetails {

	errors := make([]ValidationError, len(err.Fields))

	i := 0
	for key, value := range err.Fields {
		e := ValidationError{
			Name:   key,
			Reason: value.Error(),
		}

		errors[i] = e
		i++
	}

	return ValidationProblemDetails{
		Type:   "about:blank",
		Title:  "validation error",
		Errors: errors,
	}
}
