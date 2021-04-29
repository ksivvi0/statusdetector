package store

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx"
	"github.com/ksivvi0/statusdetector/internal/helper"
)

type Store struct {
	Conn             *pgx.Conn
	operationTimeout time.Duration
}

type StoreHistory struct {
	CheckTime   time.Time
	CheckStatus int
}

func NewStore(connString string) (*Store, error) {
	if !helper.CheckNullString(connString) {
		return nil, errors.New("пустая строка подключения")
	}
	connCfg, err := pgx.ParseConnectionString(connString)
	helper.IsError(err, true)

	conn, err := pgx.Connect(connCfg)
	helper.IsError(err, true)

	return &Store{Conn: conn, operationTimeout: time.Second * 30}, nil
}

func (s *Store) GetStatus(url string, _time int) (*StoreHistory, error) {
	if !helper.CheckNullString(url) {
		return nil, errors.New("некорректный адрес")
	}
	history := &StoreHistory{}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	history.CheckStatus = resp.StatusCode

	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()

	_, err = s.Conn.QueryEx(ctx, "SELECT add_status($1::text, $2::integer)", nil, url, history.CheckStatus)
	if err != nil {
		return history, err
	}

	return history, nil
}

func (s *Store) GetHistory(url string) ([]StoreHistory, error) {
	if !helper.CheckNullString(url) {
		return nil, errors.New("пустой url")
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.operationTimeout)
	defer cancel()
	rows, err := s.Conn.QueryEx(ctx, "SELECT get_history($1::text)", nil, url)
	if err != nil {
		return nil, err
	}
	history := make([]StoreHistory, 0)

	for rows.Next() {
		var tmp StoreHistory
		if err = rows.Scan(&tmp); err != nil {
			return nil, err
		}
		history = append(history, tmp)
	}
	return history, nil
}
