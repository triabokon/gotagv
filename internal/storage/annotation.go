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

const annotationTable = "annotations"

func (s *Storage) GetAnnotationWithDuration(ctx context.Context, id string) (*model.Annotation, error) {
	sql, params, err := postgresql.StatementBuilder.
		Select(append(annotationColumns(), "videos.duration")...).
		From(annotationTable).
		Where(squirrel.Eq{"annotations.id": id}).
		Join(fmt.Sprintf("%s ON %s.video_id = %s.id", videoTable, annotationTable, videoTable)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	row := s.client.DB.QueryRow(ctx, sql, params...)

	a, sErr := scanAnnotation(row, true)
	if errors.Is(sErr, pgx.ErrNoRows) {
		return nil, model.ErrNotFound
	}
	if sErr != nil {
		return nil, fmt.Errorf("failed to get annotation: %w", sErr)
	}
	return a, nil
}

func (s *Storage) ListAnnotations(ctx context.Context, videoID string) ([]*model.Annotation, error) {
	sql, params, err := postgresql.StatementBuilder.
		Select(annotationColumns()...).
		Where(squirrel.Eq{"video_id": videoID}).
		From(annotationTable).
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

	var result []*model.Annotation
	for rows.Next() {
		a, sErr := scanAnnotation(rows, false)
		if sErr != nil {
			return nil, fmt.Errorf("scan failed: %w", sErr)
		}
		result = append(result, a)
	}
	if rErr := rows.Err(); rErr != nil {
		return nil, rErr
	}
	return result, nil
}

func (s *Storage) InsertAnnotation(ctx context.Context, a *model.Annotation) error {
	query, args, err := postgresql.StatementBuilder.
		Insert(annotationTable).
		SetMap(map[string]interface{}{
			"id":         a.ID,
			"video_id":   a.VideoID,
			"user_id":    a.UserID,
			"start_time": a.StartTime.Seconds(),
			"end_time":   a.EndTime.Seconds(),
			"type":       a.Type,
			"message":    a.Message,
			"url":        a.URL,
			"title":      a.Title,
			"created_at": a.CreatedAt,
			"updated_at": a.UpdatedAt,
		}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, qErr := s.client.DB.Exec(ctx, query, args...)
	if qErr != nil {
		pgErr, ok := qErr.(*pgconn.PgError)
		if ok && pgErr.Code == uniqueViolation {
			return model.ErrAlreadyExists
		}
		return fmt.Errorf("failed to insert: %w", qErr)
	}
	return nil
}

func (s *Storage) UpdateAnnotation(ctx context.Context, id string, p *model.UpdateAnnotationParams) error {
	if p.NoUpdates() {
		return fmt.Errorf("no updates")
	}
	builder := postgresql.StatementBuilder.
		Update(annotationTable).
		Where(squirrel.Eq{"id": id}).
		Set("updated_at", time.Now())

	if p.StartTime != nil {
		builder = builder.Set("start_time", p.StartTime.Seconds())
	}
	if p.EndTime != nil {
		builder = builder.Set("end_time", p.EndTime.Seconds())
	}
	if p.Type != nil {
		builder = builder.Set("type", *p.Type)
	}
	if p.Message != nil {
		builder = builder.Set("message", *p.Message)
	}
	if p.URL != nil {
		builder = builder.Set("url", *p.URL)
	}
	if p.Title != nil {
		builder = builder.Set("title", *p.Title)
	}
	sql, params, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	ct, err := s.client.DB.Exec(ctx, sql, params...)
	if err != nil {
		return fmt.Errorf("failed to execute: %w", err)
	}
	switch ra := ct.RowsAffected(); ra {
	case 0:
		return model.ErrNotFound
	case 1:
		return nil
	default:
		return fmt.Errorf("update more than one (%d) annotation with id %q", ra, id)
	}
}

func (s *Storage) DeleteAnnotation(ctx context.Context, id string) error {
	deleteBuilder := postgresql.StatementBuilder.Delete(annotationTable).
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

func annotationColumns() []string {
	columns := []string{
		"annotations.id", "annotations.video_id", "annotations.user_id", "annotations.start_time",
		"annotations.end_time", "annotations.type", "annotations.message", "annotations.url",
		"annotations.title", "annotations.created_at", "annotations.updated_at",
	}
	return columns
}

func scanAnnotation(row pgx.Row, withDuration bool) (*model.Annotation, error) {
	var a model.Annotation
	var startTime, endTime, vidDuration int
	var rErr error
	if withDuration {
		rErr = row.Scan(
			&a.ID, &a.VideoID, &a.UserID, &startTime,
			&endTime, &a.Type, &a.Message, &a.URL, &a.Title,
			&a.CreatedAt, &a.UpdatedAt, &vidDuration,
		)
		a.VideoDuration = time.Duration(vidDuration) * time.Second
	} else {
		rErr = row.Scan(
			&a.ID, &a.VideoID, &a.UserID, &startTime,
			&endTime, &a.Type, &a.Message, &a.URL, &a.Title,
			&a.CreatedAt, &a.UpdatedAt,
		)
	}
	if rErr != nil {
		return nil, fmt.Errorf("failed to scan annotation: %w", rErr)
	}
	a.StartTime = time.Duration(startTime) * time.Second
	a.EndTime = time.Duration(endTime) * time.Second
	return &a, nil
}
