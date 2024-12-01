package httputils

import (
	"costumers-api/domain/domain_errors"
	httpmodels "costumers-api/presentation/http_models"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

func JsonOk[T any](w http.ResponseWriter, object T) {
	JsonResponse(w, object, http.StatusOK)
}

func JsonResponse[T any](w http.ResponseWriter, object T, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(object); err != nil {
		http.Error(w, `{"error": "failed to encode JSON"}`, http.StatusInternalServerError)
	}
}

func ErrorResponse(w http.ResponseWriter, err error) {

	var valErr domain_errors.ValidationError

	if errors.As(err, &valErr) {
		validationDetails := httpmodels.NewValidationProblemDetails(valErr)
		JsonResponse(w, validationDetails, http.StatusBadRequest)
		return
	}

	problemDetails := httpmodels.NewProblemDetails(err)
	JsonResponse(w, problemDetails, problemDetails.Status)
}

func ReadJson[T any](r *http.Request) (T, error) {
	var object T
	decodeErr := json.NewDecoder(r.Body).Decode(&object)

	return object, decodeErr
}

func ReadPathVariable(r *http.Request, name string) string {
	vars := mux.Vars(r)
	return vars[name]
}
