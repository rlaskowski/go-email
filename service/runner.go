package service

import (
	"log"
	"os"
	"os/signal"

	"github.com/rlaskowski/go-email/config"
	"github.com/rlaskowski/go-email/registries"
	"github.com/rlaskowski/go-email/router"
)

type Service struct {
	http       *router.HttpServer
	grpc       *router.GrpcServer
	registries *registries.Registries
}

func New() *Service {
	registries := registries.NewRegistries()

	return &Service{
		http: router.NewHttpServer(registries),
		grpc: router.NewGrpcServer(registries),
	}
}

func (s *Service) Start() error {
	log.Printf("Email http service working in directory %s with pid %d", config.GetWorkingDirectory(), os.Getpid())
	log.Printf("Path to store temporary file before send: %s", config.FileStorePath)

	if err := s.registries.Start(); err != nil {
		log.Printf("Could not start registries: %s", err)
	}

	if err := s.http.Start(); err != nil {
		log.Printf("Could not start HTTP server: %s", err)
	}

	if err := s.grpc.Start(); err != nil {
		log.Printf("Could not start GRPC server: %s", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	<-sigChan

	return s.Stop()

}

func (s *Service) Stop() error {
	if err := s.registries.Stop(); err != nil {
		log.Printf("Could not stop Registries: %s", err)
	}

	if err := s.http.Stop(); err != nil {
		log.Printf("Could not stop HTTP server: %s", err)
	}

	if err := s.grpc.Stop(); err != nil {
		log.Printf("Could not stop GRPC server: %s", err)
	}

	return nil
}
