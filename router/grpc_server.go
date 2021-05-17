package router

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/rlaskowski/go-email/config"
	"github.com/rlaskowski/go-email/grpc/protobuf/emailservice"
	"github.com/rlaskowski/go-email/registries"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	grpc       *grpc.Server
	context    context.Context
	cancel     context.CancelFunc
	registries *registries.Registries
}

func NewGrpcServer(registries *registries.Registries) *GrpcServer {
	ctx, cancel := context.WithCancel(context.Background())

	g := &GrpcServer{
		grpc:       grpc.NewServer(),
		context:    ctx,
		cancel:     cancel,
		registries: registries,
	}

	return g
}

func (g *GrpcServer) configureGrpc() {
	reflection.Register(g.grpc)

	emailservice.RegisterEmailServiceServer(g.grpc, g.registries.GrpcEmailQueue)
}

func (g *GrpcServer) Start() error {

	g.configureGrpc()

	go g.start()

	return nil
}

func (g *GrpcServer) Stop() error {
	g.cancel()

	log.Print("Stopping GRPC Server...")

	g.grpc.Stop()

	return nil
}

func (g *GrpcServer) start() error {
	log.Printf("Starting GRPC Server on %d port", config.GrpcListenPort)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GrpcListenPort))
	if err != nil {
		log.Fatalf("Caught error while starting listen on port %d, %s", config.GrpcListenPort, err.Error())
	}

	if err := g.grpc.Serve(listener); err != nil {
		log.Fatalf("Caught error while starting grpc server: %s", err.Error())
	}

	return nil
}
