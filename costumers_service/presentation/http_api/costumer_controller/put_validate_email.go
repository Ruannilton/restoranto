package costumercontroller

import (
	"costumers-api/domain/abstractions"
	"costumers-api/domain/usecases"
	httputils "costumers-api/presentation/http_utils"
	"net/http"
	"strconv"
)

func (controller *CostumerController) validateEmail(w http.ResponseWriter, r *http.Request) {
	pathId := httputils.ReadPathVariable(r, "id")
	code := httputils.ReadPathVariable(r, "code")
	costumerId, errId := strconv.Atoi(pathId)

	if errId != nil {
		httputils.ErrorResponse(w, errId)
		return
	}

	useCaseParams := usecases.ValidateCostumerEmailParams{
		Id:   costumerId,
		Code: code,
	}

	unityOfWork, uowErr := controller.unitOfWorkFactory.NewUnitOfWork(r.Context())

	if uowErr != nil {
		httputils.ErrorResponse(w, uowErr)
		return
	}

	useCase := usecases.NewValidateCostumerEmailUseCase(controller.costumerRepository, controller.distributedCache, unityOfWork)
	useCaseDecorator := abstractions.NewStatefullUseCase(useCase, unityOfWork)

	useCaseResponse, useCaseErr := useCaseDecorator.Execute(useCaseParams)

	if useCaseErr != nil {
		httputils.ErrorResponse(w, useCaseErr)
		return
	}

	httputils.JsonOk(w, useCaseResponse)
}
