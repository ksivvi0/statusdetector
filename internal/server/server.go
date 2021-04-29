package server

import (
	"context"
	"log"
	"os"

	"github.com/ksivvi0/statusdetector/config"
	"github.com/ksivvi0/statusdetector/internal/helper"
	"github.com/ksivvi0/statusdetector/internal/store"
)

type Server struct {
	Logger *log.Logger
	Store  *store.Store
}

func NewServer(cfg *config.Config, store *store.Store) (*Server, error) {

	logFile, err := os.OpenFile(cfg.LogPath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		helper.IsError(err, true)
	}

	return &Server{
		Logger: log.New(logFile, "DETECTOR:", log.Ldate|log.Ltime|log.Lshortfile),
		Store:  store,
	}, nil
}

func (d *Server) GetHistory(ctx context.Context, req *URLRequest) (*HistoryResponse, error) {
	return nil, nil
}

func (d *Server) GetStatus(ctx context.Context, req *URLRequest) (*URLResponse, error) {
	d.Logger.Printf("Запрос: %s %d\n", req.Url, req.Time)
	code, err := d.Store.GetStatus(req.Url, int(req.Time))
	if err != nil {
		d.Logger.Println(err)
		return nil, err
	}
	return &URLResponse{Url: req.Url, Code: int32(code)}, nil
}

func (d *Server) DropURL(ctx context.Context, req *URLRequest) (*VoidResponse, error) {
	return nil, nil
}
