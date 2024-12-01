package costumercontroller

import (
	"costumers-api/domain/abstractions"
	"costumers-api/domain/usecases"
	httpmodels "costumers-api/presentation/http_models"
	httputils "costumers-api/presentation/http_utils"
	"net/http"
)

func (controller *CostumerController) getCostumers(w http.ResponseWriter, r *http.Request) {
	type GetCostumersRequest struct {
		Name      string `json:"name"`
		Cpf       string `json:"cpf"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		Page      int    `json:"page"`
		PageCount int    `json:"pageCount"`
	}

	type GetCostumersResponse struct {
		Costumers []httpmodels.CostumerResponse `json:"costumers"`
		Page      int                           `json:"page"`
		PageCount int                           `json:"pageCount"`
		Total     int                           `json:"total"`
	}

	body, parseErr := httputils.ReadJson[GetCostumersRequest](r)

	if parseErr != nil {
		httputils.ErrorResponse(w, parseErr)
		return
	}

	usecase := usecases.NewSearchCostumersUseCase(controller.searchCostumerQuery)
	usecaseDecorator := abstractions.NewStatelessUseCase(usecase)

	useCaseParams := usecases.SearchCostumersParams{
		Name:      body.Name,
		Cpf:       body.Cpf,
		Email:     body.Email,
		Phone:     body.Phone,
		Page:      body.Page,
		PageCount: body.PageCount,
	}

	useCaseResponse, searchErr := usecaseDecorator.Execute(useCaseParams)

	if searchErr != nil {
		httputils.ErrorResponse(w, searchErr)
		return
	}

	bodyResponse := GetCostumersResponse{
		Page:      useCaseParams.Page,
		PageCount: useCaseParams.PageCount,
		Total:     0,
		Costumers: func() []httpmodels.CostumerResponse {
			maped := make([]httpmodels.CostumerResponse, len(useCaseResponse.Costumer))

			for i, v := range useCaseResponse.Costumer {
				maped[i] = httpmodels.NewCostumerResponse(v)
			}

			return maped
		}(),
	}

	httputils.JsonOk(w, bodyResponse)
}
