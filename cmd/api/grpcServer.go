package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"news_service.andreyklimov.net/internal/data"
)

type GRPCServer struct {
	addr   string
	model  data.Models
	server *grpc.Server
}

func NewGRPCServer(addr string, models data.Models) *GRPCServer {
	return &GRPCServer{addr: addr,
		model: models}
}

func (s *GRPCServer) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.server = grpc.NewServer()

	// register our grpc services
	newsService := s.model
	NewNewsService(s.server, newsService)

	log.Println("Starting gRPC server on", s.addr)

	return s.server.Serve(lis)
}

// waitForShutdown блокирует до получения системного сигнала и корректно останавливает gRPC-сервер.
func waitForShutdown(server *grpc.Server) {
	// Создаём канал, в который пойдут сигналы ОС
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Received shutdown signal, gracefully stopping gRPC server...")
	server.GracefulStop()
}
