package costumercontroller

import (
	"costumers-api/domain/abstractions"
	"costumers-api/domain/usecases"
	httpmodels "costumers-api/presentation/http_models"
	httputils "costumers-api/presentation/http_utils"
	"net/http"
	"strconv"
)

func (controller *CostumerController) getCostumer(w http.ResponseWriter, r *http.Request) {

	pathId := httputils.ReadPathVariable(r, "id")

	costumerId, errId := strconv.Atoi(pathId)

	if errId != nil {
		httputils.ErrorResponse(w, errId)
		return
	}

	useCaseParams := usecases.FindCostumerParams{
		Id: costumerId,
	}

	useCase := usecases.NewFindCostumerUseCase(controller.costumerRepository)
	useCaseDecorator := abstractions.NewStatelessUseCase(useCase)

	useCaseResponse, useCaseErr := useCaseDecorator.Execute(useCaseParams)

	if useCaseErr != nil {
		httputils.ErrorResponse(w, useCaseErr)
		return
	}

	bodyResponse := httpmodels.NewCostumerResponse(useCaseResponse.Costumer)

	httputils.JsonOk(w, bodyResponse)
}
