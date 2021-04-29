package main

import (
	"fmt"
	"net"

	"github.com/ksivvi0/statusdetector/config"
	"github.com/ksivvi0/statusdetector/internal/helper"
	"github.com/ksivvi0/statusdetector/internal/server"
	"github.com/ksivvi0/statusdetector/internal/store"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.NewConfig("127.0.0.1", 8000, "./detector.log", "postgres://postgres:pwd@localhost:5432/detector")
	helper.IsError(err, true)

	dbStore, err := store.NewStore(cfg.ConnectionString)
	helper.IsError(err, true)

	srv, err := server.NewServer(cfg, dbStore)
	helper.IsError(err, true)
	grpcServer := grpc.NewServer()

	socket, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", cfg.ListenAddress, cfg.ListenPort))
	helper.IsError(err, true)

	server.RegisterDetectorServer(grpcServer, srv)

	// go listenSignals(grpcServer)

	err = grpcServer.Serve(socket)
	helper.IsError(err, true)
}

// func listenSignals(grpcServer *grpc.Server) {
// 	signalsChan := make(chan os.Signal, 1)
// 	signal.Notify(signalsChan, syscall.SIGINT|syscall.SIGTERM)
// 	<-signalsChan
// 	grpcServer.GracefulStop()
// }
