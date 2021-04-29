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
	Conn *pgx.Conn
}

func NewStore(connString string) (*Store, error) {
	if !helper.CheckNullString(connString) {
		return nil, errors.New("пустая строка подключения")
	}
	connCfg, err := pgx.ParseConnectionString(connString)
	helper.IsError(err, true)

	conn, err := pgx.Connect(connCfg)

	helper.IsError(err, true)

	return &Store{Conn: conn}, nil
}

func (s *Store) GetStatus(url string, _time int) (int, error) {
	if len(url) < 0 {
		return -1, errors.New("некорректный адрес")
	}
	code := -1
	resp, err := http.Get(url)
	if err != nil {
		return code, err
	}
	code = resp.StatusCode

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, err = s.Conn.QueryEx(ctx, "SELECT add_status($1::text, $2::integer)", nil, url, code)
	if err != nil {
		return -1, err
	}

	return code, nil
}
