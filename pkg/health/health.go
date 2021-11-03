package health

import "github.com/sirupsen/logrus"

//Service is the interface for the health
type Service interface {
	HealthCheck() (service bool, err error)
}
type svc struct {
	log *logrus.Entry
}

//NewService gives a new Service
func NewService(log *logrus.Entry) Service {
	return &svc{
		log: log,
	}
}

//HealthCheck returns the status of the API and it's components
func (s *svc) HealthCheck() (service bool, err error) {
	s.log.Info("Performing Healthcheck")
	return true, nil
}
