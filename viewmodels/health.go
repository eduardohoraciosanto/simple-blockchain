package viewmodels

type Health struct {
	Name  string `json:"name"`
	Alive bool   `json:"alive"`
}
type HealthResponse struct {
	Services []Health `json:"services"`
}
