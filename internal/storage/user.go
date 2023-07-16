package storage

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/triabokon/gotagv/internal/model"
	"github.com/triabokon/gotagv/internal/postgresql"
)

const userTable = "users"

func (s *Storage) GetUser(ctx context.Context, id string) (string, error) {
	sql, params, err := postgresql.StatementBuilder.
		Select("id").
		From(userTable).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("failed to build query: %w", err)
	}

	row := s.client.DB.QueryRow(ctx, sql, params...)
	var userID string
	if rErr := row.Scan(&userID); rErr != nil {
		return "", fmt.Errorf("failed to scan user id: %w", rErr)
	}
	switch errors.Cause(err) {
	case nil:
		return userID, nil
	case pgx.ErrNoRows:
		return "", model.ErrNotFound
	default:
		return "", err
	}
}

func (s *Storage) InsertUser(ctx context.Context, id string) error {
	query, args, err := postgresql.StatementBuilder.
		Insert(userTable).
		SetMap(map[string]interface{}{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	if _, qErr := s.client.DB.Exec(ctx, query, args...); qErr != nil {
		pgErr, ok := qErr.(*pgconn.PgError)
		if ok && pgErr.Code == uniqueViolation {
			return model.ErrAlreadyExists
		}
		return fmt.Errorf("failed to insert: %w", qErr)
	}
	return nil
}
