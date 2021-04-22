package service

import (
	"log"
	"os"
	"os/signal"

	"github.com/rlaskowski/go-email/config"
	"github.com/rlaskowski/go-email/email"
	"github.com/rlaskowski/go-email/queue"
	"github.com/rlaskowski/go-email/registries"
	"github.com/rlaskowski/go-email/router"
)

type Service struct {
	http         *router.HttpServer
	grpc         *router.GrpcServer
	email        *email.Email
	queueFactory *queue.QueueFactory
	registries   *registries.Registries
}

func New() *Service {
	email := email.NewEmail()
	queueFactory := queue.NewFactory(email)

	registries := registries.NewRegistries(queueFactory, email)

	return &Service{
		http:         router.NewHttpServer(registries),
		grpc:         router.NewGrpcServer(registries),
		email:        email,
		queueFactory: queueFactory,
	}
}

func (s *Service) Start() error {
	log.Printf("Email http service working in directory %s with pid %d", config.GetExecutableDirectory(), os.Getpid())
	log.Printf("Path to store temporary file before send: %s", config.FileStorePath)

	if err := s.http.Start(); err != nil {
		log.Printf("Could not start HTTP server: %s", err)
	}

	if err := s.grpc.Start(); err != nil {
		log.Printf("Could not start GRPC server: %s", err)
	}

	if err := s.email.Start(); err != nil {
		log.Printf("Could not start Email: %s", err)
	}

	if err := s.queueFactory.Start(); err != nil {
		log.Printf("Could not start Queue factory: %s", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	<-sigChan

	return s.Stop()

}

func (s *Service) Stop() error {
	if err := s.http.Stop(); err != nil {
		log.Printf("Could not stop HTTP server: %s", err)
	}

	if err := s.grpc.Stop(); err != nil {
		log.Printf("Could not stop GRPC server: %s", err)
	}

	if err := s.email.Stop(); err != nil {
		log.Printf("Could not stop Email: %s", err)
	}

	if err := s.queueFactory.Stop(); err != nil {
		log.Printf("Could not stop Queue factory: %s", err)
	}

	return nil
}
