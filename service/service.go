package service

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rlaskowski/go-email/config"
	"github.com/rlaskowski/go-email/grpc"
	"github.com/rlaskowski/go-email/queue"
	"github.com/rlaskowski/go-email/rest"
	"github.com/rlaskowski/go-email/router"
)

type Service struct {
	http             *router.HttpServer
	grpc             *router.GrpcServer
	queueBox         *queue.QueueBox
	emailService     *grpc.EmailService
	emailRestService *rest.EmailService
	serviceConfig    config.ServiceConfig
}

func NewService() *Service {
	c := config.DefaultServiceConfig
	return newService(c)
}

func ServiceWithConfig(config config.ServiceConfig) *Service {
	return newService(config)
}

func newService(serviceConfig config.ServiceConfig) *Service {
	s := &Service{
		queueBox:      queue.NewQueuBox(serviceConfig),
		serviceConfig: serviceConfig,
	}

	s.emailService = grpc.NewEmailService(s.queueBox)
	s.emailRestService = rest.NewEmailService(s.queueBox)
	s.http = router.NewHttpServer(s)
	s.grpc = router.NewGrpcServer(s)

	return s
}

func (s *Service) Start() error {
	log.Printf("Email http service working in directory %s with pid %d", s.serviceConfig.FileStorePath, os.Getpid())
	log.Printf("Path to store temporary file before send: %s", s.serviceConfig.FileStorePath)

	if err := s.queueBox.Start(); err != nil {
		log.Printf("Could not start QueueBox: %s", err)
	}

	if err := s.http.Start(); err != nil {
		log.Printf("Could not start HTTP server: %s", err)
	}

	if err := s.grpc.Start(); err != nil {
		log.Printf("Could not start GRPC server: %s", err)
	}

	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, syscall.SIGTERM)

	<-sigChan

	return s.Stop()

}

func (s *Service) Stop() error {
	if err := s.queueBox.Stop(); err != nil {
		log.Printf("Could not stop QueueBox: %s", err)
	}

	if err := s.http.Stop(); err != nil {
		log.Printf("Could not stop HTTP server: %s", err)
	}

	if err := s.grpc.Stop(); err != nil {
		log.Printf("Could not stop GRPC server: %s", err)
	}

	return nil
}

func (s *Service) EmailService() *grpc.EmailService {
	return s.emailService
}

func (s *Service) EmailRestService() *rest.EmailService {
	return s.emailRestService
}

func (s *Service) QueueBox() *queue.QueueBox {
	return s.queueBox
}

func (s *Service) ServiceConfig() config.ServiceConfig {
	return s.serviceConfig
}
