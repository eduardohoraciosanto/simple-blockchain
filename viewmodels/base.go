package viewmodels

import (
	"encoding/json"
	"net/http"

	"github.com/eduardohoraciosanto/simple-blockchain/config"
)

type Meta struct {
	Version string `json:"version"`
}

type BaseResponse struct {
	Meta  Meta        `json:"meta"`
	Data  interface{} `json:"data,omitempty"`
	Error interface{} `json:"error,omitempty"`
}

func newBaseResponseWithData(data interface{}) BaseResponse {
	return BaseResponse{
		Meta: Meta{
			Version: config.GetVersion(),
		},
		Data: data,
	}
}

func newBaseResponseWithError(err interface{}) BaseResponse {
	return BaseResponse{
		Meta: Meta{
			Version: config.GetVersion(),
		},
		Error: err,
	}
}

func RespondWithData(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(newBaseResponseWithData(data))
}

func RespondBadRequest(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	return json.NewEncoder(w).Encode(newBaseResponseWithError(StandardBadBodyRequest))
}

func RespondInternalServerError(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	return json.NewEncoder(w).Encode(newBaseResponseWithError(StandardInternalServerError))
}
