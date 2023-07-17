package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/triabokon/gotagv/internal/model"
	"github.com/triabokon/gotagv/internal/postgresql"
)

const videoTable = "videos"

func (s *Storage) ListVideos(ctx context.Context) ([]*model.Video, error) {
	sql, params, err := postgresql.StatementBuilder.
		Select(videoColumns()...).
		From(videoTable).
		OrderBy("updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := s.client.DB.Query(ctx, sql, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to exec query: %w", err)
	}
	defer rows.Close()

	var result []*model.Video
	for rows.Next() {
		v, sErr := scanVideo(rows)
		if sErr != nil {
			return nil, fmt.Errorf("scan failed: %w", sErr)
		}
		result = append(result, v)
	}
	if rErr := rows.Err(); rErr != nil {
		return nil, rErr
	}
	return result, nil
}

func (s *Storage) GetVideo(ctx context.Context, id string) (*model.Video, error) {
	sql, params, err := postgresql.StatementBuilder.
		Select(videoColumns()...).
		From(videoTable).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	row := s.client.DB.QueryRow(ctx, sql, params...)
	v, sErr := scanVideo(row)
	if errors.Is(sErr, pgx.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if sErr != nil {
		return nil, fmt.Errorf("failed to get video: %w", sErr)
	}
	return v, nil
}

func (s *Storage) InsertVideo(ctx context.Context, video *model.Video) error {
	query, args, err := postgresql.StatementBuilder.
		Insert(videoTable).
		SetMap(map[string]interface{}{
			"id":         video.ID,
			"user_id":    video.UserID,
			"url":        video.URL,
			"duration":   video.Duration.Seconds(),
			"created_at": video.CreatedAt,
			"updated_at": video.UpdatedAt,
		}).ToSql()
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

func (s *Storage) DeleteVideo(ctx context.Context, id string) error {
	deleteBuilder := postgresql.StatementBuilder.Delete(videoTable).
		Where(squirrel.Eq{"id": id})

	sql, params, err := deleteBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	_, err = s.client.DB.Exec(ctx, sql, params...)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}
	return nil
}

func videoColumns() []string {
	columns := []string{
		"id", "user_id", "url", "duration", "created_at", "updated_at",
	}
	return columns
}

func scanVideo(row pgx.Row) (*model.Video, error) {
	var durationSeconds int
	var v model.Video
	if rErr := row.Scan(
		&v.ID, &v.UserID, &v.URL,
		&durationSeconds, &v.CreatedAt, &v.UpdatedAt,
	); rErr != nil {
		return nil, fmt.Errorf("failed to scan video: %w", rErr)
	}
	v.Duration = time.Duration(durationSeconds) * time.Second
	return &v, nil
}
