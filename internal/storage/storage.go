package storage

import (
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
