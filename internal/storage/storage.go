package storage

import (
	"context"

	"github.com/jackc/pgx/v4"

	"github.com/triabokon/gotagv/internal/postgresql"
)

const uniqueViolation = "23505"

type Storage struct {
	client *postgresql.Client
}

func New(client *postgresql.Client) *Storage {
	return &Storage{
		client: client,
	}
}

func (s *Storage) BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(tx pgx.Tx) error) error {
	return s.client.DB.BeginTxFunc(ctx, txOptions, f)
}
