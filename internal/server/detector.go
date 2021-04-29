package server

//
//import (
//	"context"
//	"fmt"
//	"log"
//	"os"
//
//	"github.com/ksivvi0/statusdetector/config"
//	"github.com/ksivvi0/statusdetector/internal/helper"
//	"github.com/ksivvi0/statusdetector/internal/store"
//)
//
//type Detector struct {
//	Config *config.Config
//	Logger *log.Logger
//	Store  store.URLStore
//}
//
//func NewDetector(cfg *config.Config, store *store.Store) (*Detector, error) {
//
//	logFile, err := os.OpenFile(cfg.LogPath, os.O_WRONLY|os.O_CREATE, 0600)
//	if err != nil {
//		helper.IsError(err, true)
//	}
//
//	return &Detector{
//		Config: cfg,
//		Logger: log.New(logFile, "DETECTOR:", log.Ldate|log.Ltime|log.Lshortfile),
//	}, nil
//}
//
//// func (d *Detector) GetHistory(ctx context.Context, req *URLRequest) (*HistoryResponse, error) {
//// }
//
//func (d *Detector) GetStatus(ctx context.Context, req *URLRequest) (*URLResponse, error) {
//	code, err := d.Store.GetStatus(req.Url, int(req.Time))
//	if err != nil {
//		d.Logger.Println(err)
//		return nil, err
//	}
//	fmt.Println(code)
//	return &URLResponse{Url: req.Url, Code: int32(code)}, nil
//}
//
//// func (d *Detector) DropURL(ctx context.Context, req *URLRequest) (*VoidResponse, error) {}
