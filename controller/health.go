package controller

import (
	"net/http"

	"github.com/eduardohoraciosanto/simple-blockchain/pkg/health"
	"github.com/eduardohoraciosanto/simple-blockchain/viewmodels"
)

type HealthController struct {
	Service health.Service
}

//Health is the handler for the health endpoint
func (c *HealthController) Health(w http.ResponseWriter, r *http.Request) {

	//using lower level pkg to do the logic
	service, err := c.Service.HealthCheck()
	if err != nil {
		viewmodels.RespondInternalServerError(w)
		return
	}
	hr := viewmodels.HealthResponse{
		Services: []viewmodels.Health{
			{
				Name:  "service",
				Alive: service,
			},
		},
	}
	viewmodels.RespondWithData(w, http.StatusOK, hr)
}
