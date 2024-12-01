package costumercontroller

import (
	"costumers-api/domain/abstractions"
	"costumers-api/domain/usecases"
	httputils "costumers-api/presentation/http_utils"
	"net/http"
	"strconv"
)

func (controller *CostumerController) deleteCostumers(w http.ResponseWriter, r *http.Request) {
	pathId := httputils.ReadPathVariable(r, "id")

	costumerId, errId := strconv.Atoi(pathId)

	if errId != nil {
		httputils.ErrorResponse(w, errId)
		return
	}

	useCaseParams := usecases.DeleteCostumerParams{
		Id: costumerId,
	}

	unityOfWork, uowErr := controller.unitOfWorkFactory.NewUnitOfWork(r.Context())

	if uowErr != nil {
		httputils.ErrorResponse(w, uowErr)
		return
	}

	useCase := usecases.NewDeleteCostumerUseCase(controller.costumerRepository, unityOfWork)
	useCaseDecorator := abstractions.NewStatefullUseCase(useCase, unityOfWork)

	_, useCaseErr := useCaseDecorator.Execute(useCaseParams)

	if useCaseErr != nil {
		httputils.ErrorResponse(w, useCaseErr)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
