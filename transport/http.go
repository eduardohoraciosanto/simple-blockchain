package transport

import (
	"net/http"

	"github.com/eduardohoraciosanto/simple-blockchain/controller"
	"github.com/eduardohoraciosanto/simple-blockchain/pkg/blockchain"
	"github.com/eduardohoraciosanto/simple-blockchain/pkg/health"
	"github.com/gorilla/mux"
)

func NewHTTPRouter(svc blockchain.Service, hsvc health.Service) *mux.Router {
	hc := controller.HealthController{
		Service: hsvc,
	}
	bc := controller.BlockController{
		Service: svc,
	}

	r := mux.NewRouter()

	r.HandleFunc("/health", hc.Health).Methods(http.MethodGet)
	r.HandleFunc("/blockchain", bc.GetBlockchain).Methods(http.MethodGet)
	r.HandleFunc("/blockchain", bc.InsertDataBlock).Methods(http.MethodPost)

	r.PathPrefix("/swagger").Handler(http.StripPrefix("/swagger", http.FileServer(http.Dir("./swagger"))))
	return r
}
