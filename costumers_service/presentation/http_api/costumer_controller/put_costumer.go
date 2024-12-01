package costumercontroller

import (
	"costumers-api/domain/abstractions"
	"costumers-api/domain/usecases"
	httputils "costumers-api/presentation/http_utils"
	"net/http"
	"strconv"
)

func (controller *CostumerController) putCostumer(w http.ResponseWriter, r *http.Request) {
	type PutCostumersRequest struct {
		Name  string `json:"name"`
		Cpf   string `json:"cpf"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	pathId := httputils.ReadPathVariable(r, "id")

	costumerId, errId := strconv.Atoi(pathId)

	if errId != nil {
		httputils.ErrorResponse(w, errId)
		return
	}

	body, parseErr := httputils.ReadJson[PutCostumersRequest](r)

	if parseErr != nil {
		httputils.ErrorResponse(w, parseErr)
		return
	}

	useCaseParams := usecases.UpdateCostumersParams{
		Name:  body.Name,
		Cpf:   body.Cpf,
		Phone: body.Phone,
		Email: body.Email,
		Id:    costumerId,
	}

	unityOfWork, uowErr := controller.unitOfWorkFactory.NewUnitOfWork(r.Context())

	if uowErr != nil {
		httputils.ErrorResponse(w, uowErr)
		return
	}

	useCase := usecases.NewUpdateCostumersUseCase(controller.costumerRepository, unityOfWork)
	useCaseDecorator := abstractions.NewStatefullUseCase(useCase, unityOfWork)

	_, useCaseErr := useCaseDecorator.Execute(useCaseParams)

	if useCaseErr != nil {
		httputils.ErrorResponse(w, useCaseErr)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
