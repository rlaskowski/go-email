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
	registries *registries.Registries
}

func New() *Service {
	registries := registries.NewRegistries()

	return &Service{
		http: router.NewHttpServer(registries),
	}
}

func (s *Service) Start() error {
	log.Printf("Email http service working in directory %s with pid %d", config.GetExecutableDirectory(), os.Getpid())

	if err := s.http.Start(); err != nil {
		log.Printf("Could not start HTTP server: %s", err)
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

	return nil
}
