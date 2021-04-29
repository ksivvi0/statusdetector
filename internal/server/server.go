package server

import (
	"context"
	"github.com/golang/protobuf/ptypes/timestamp"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ksivvi0/statusdetector/config"
	"github.com/ksivvi0/statusdetector/internal/helper"
	"github.com/ksivvi0/statusdetector/internal/store"
)

type Server struct {
	Logger  *log.Logger
	Store   *store.Store
	TimeOut time.Duration
}

func NewServer(cfg *config.Config, store *store.Store) (*Server, error) {
	logFile, err := os.OpenFile(cfg.LogPath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		helper.IsError(err, true)
	}

	return &Server{
		Logger:  log.New(logFile, "DETECTOR:", log.Ldate|log.Ltime|log.Lshortfile),
		Store:   store,
		TimeOut: time.Second * 30,
	}, nil
}

func (d *Server) GetHistory(ctx context.Context, req *URLRequest) (*HistoryResponse, error) {
	d.Logger.Printf("Запрос GetHistory: %s %d\n", req.Url, req.Time)
	history, err := d.Store.GetHistory(req.Url)
	if err != nil {
		d.Logger.Println(err)
		return nil, err
	}
	historyResponse := &HistoryResponse{Url: req.Url}

	for _, v := range history {
		historyResponse.History = append(historyResponse.History, &History{
			CheckTime: &timestamp.Timestamp{Seconds: v.CheckTime.UnixNano()},
			Code:      int32(v.CheckStatus)})
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.TimeOut)
	defer cancel()
	rows, err := d.Store.Conn.QueryEx(ctx, "SELECT get_history($1::text)", nil, req.Url)
	if err != nil {
		return nil, err
	}
	result := &HistoryResponse{}
	for rows.Next() {
		var tmp History
		if err = rows.Scan(&tmp); err != nil {
			return nil, err
		}
		result.History = append(result.History, &tmp)
	}

	return result, nil
}

func (d *Server) GetStatus(ctx context.Context, req *URLRequest) (*URLResponse, error) {
	d.Logger.Printf("Запрос GetStatus: %s %d\n", req.Url, req.Time)

	resp, err := http.Get(req.Url)
	if err != nil {
		return nil, err
	}
	history := History{Code: int32(resp.StatusCode)}
	ctx, cancel := context.WithTimeout(context.Background(), d.TimeOut)
	defer cancel()

	_, err = d.Store.Conn.QueryEx(ctx, "SELECT add_status($1::text, $2::integer)", nil, req.Url, resp.StatusCode)
	if err != nil {
		return nil, err
	}

	result := &URLResponse{}
	result.History = append(result.History, &history)
	return result, err
}

func (d *Server) DropURL(ctx context.Context, req *URLRequest) (*VoidResponse, error) {
	return nil, nil
}
