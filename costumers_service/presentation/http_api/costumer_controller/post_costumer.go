package costumercontroller

import (
	"costumers-api/domain/abstractions"
	"costumers-api/domain/usecases"
	httpmodels "costumers-api/presentation/http_models"
	httputils "costumers-api/presentation/http_utils"
	"net/http"
)

func (controller *CostumerController) postCostumer(w http.ResponseWriter, r *http.Request) {

	type PostCostumerResponse struct {
		Id int
	}

	body, parseErr := httputils.ReadJson[httpmodels.PostCostumerRequest](r)

	if parseErr != nil {
		httputils.ErrorResponse(w, parseErr)
		return
	}

	unityOfWork, uowErr := controller.unitOfWorkFactory.NewUnitOfWork(r.Context())

	if uowErr != nil {
		httputils.ErrorResponse(w, uowErr)
		return
	}

	useCase := usecases.NewCreateCostumerUseCase(controller.costumerRepository, controller.eventRepository, unityOfWork)
	useCaseDecorator := abstractions.NewStatefullUseCase(useCase, unityOfWork)

	useCaseParams := usecases.CreateCostumerParams{
		Name:      body.Name,
		Cpf:       body.Cpf,
		Birthdate: body.Birthdate,
		Phone:     body.Phone,
		Email:     body.Email,
	}

	useCaseResponse, useCaseErr := useCaseDecorator.Execute(useCaseParams)

	if useCaseErr != nil {
		httputils.ErrorResponse(w, useCaseErr)
		return
	}

	bodyResponse := PostCostumerResponse{
		Id: useCaseResponse.Id,
	}

	httputils.JsonOk(w, bodyResponse)
}
