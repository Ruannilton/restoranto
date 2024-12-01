package httpapi

import (
	httputils "costumers-api/presentation/http_utils"
	"net/http"

	"github.com/gorilla/mux"
)

type HealthController struct{}

type IndexResponse struct {
	Message string
}

func NewHealthController() HealthController {
	return HealthController{}
}

func (controller *HealthController) Inject(mux *mux.Router) {
	mux.HandleFunc("/", controller.getIndex).Methods("GET")
	mux.HandleFunc("/health", controller.getHealth).Methods("GET")
}

func (controller *HealthController) getIndex(w http.ResponseWriter, r *http.Request) {
	response := IndexResponse{Message: "Running"}
	httputils.JsonResponse(w, response, http.StatusOK)
}

func (controller *HealthController) getHealth(w http.ResponseWriter, r *http.Request) {
	response := IndexResponse{Message: "Running"}
	httputils.JsonResponse(w, response, http.StatusOK)
}
