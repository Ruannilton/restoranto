package costumercontroller

import (
	"costumers-api/domain/abstractions"
	"costumers-api/domain/dependencies"

	"github.com/gorilla/mux"
)

type CostumerController struct {
	costumerRepository  dependencies.ICostumerRepository
	searchCostumerQuery dependencies.ISearchCostumerQuery
	distributedCache    dependencies.IDistributedCacheRepository
	eventRepository     dependencies.IEventRepository
	unitOfWorkFactory   abstractions.IUnitOfWorkFactory
}

func NewCostumerController(
	costumerRepository dependencies.ICostumerRepository,
	searchCostumerQuery dependencies.ISearchCostumerQuery,
	distributedCache dependencies.IDistributedCacheRepository,
	unitOfWorkFactory abstractions.IUnitOfWorkFactory,
	eventRepository dependencies.IEventRepository,
) CostumerController {
	controller := CostumerController{
		costumerRepository:  costumerRepository,
		searchCostumerQuery: searchCostumerQuery,
		distributedCache:    distributedCache,
		unitOfWorkFactory:   unitOfWorkFactory,
		eventRepository:     eventRepository,
	}

	return controller
}

func (controller *CostumerController) Inject(mux *mux.Router) {
	mux.HandleFunc("/costumer", controller.postCostumer).Methods("POST")
	mux.HandleFunc("/costumer/{id}", controller.getCostumer).Methods("GET")
	mux.HandleFunc("/costumer", controller.getCostumers).Methods("GET")
	mux.HandleFunc("/costumer/{id}", controller.putCostumer).Methods("PUT")
	mux.HandleFunc("/costumer/{id}", controller.deleteCostumers).Methods("DELETE")
	mux.HandleFunc("/costumer/{id}/validate-phone/{code}", controller.validatePhone).Methods("POST")
	mux.HandleFunc("/costumer/{id}/validate-email/{code}", controller.validateEmail).Methods("POST")
}
