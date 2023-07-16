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
		Where(squirrel.Eq{"id": id}).
		Join(fmt.Sprintf("%s ON %s.video_id = %s.id;", videoTable, annotationTable, videoTable)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	row := s.client.DB.QueryRow(ctx, sql, params...)

	a, sErr := scanAnnotation(row, true)
	if sErr != nil {
		return nil, fmt.Errorf("scan failed: %w", sErr)
	}
	switch errors.Cause(err) {
	case nil:
		return a, nil
	case pgx.ErrNoRows:
		return nil, model.ErrNotFound
	default:
		return nil, err
	}
}

func (s *Storage) ListAnnotations(ctx context.Context) ([]*model.Annotation, error) {
	sql, params, err := postgresql.StatementBuilder.
		Select(annotationColumns()...).
		From(annotationTable).
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
			"type":       a.VideoID,
			"notes":      a.VideoID,
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
		builder = builder.Set("start_time", *p.StartTime)
	}
	if p.EndTime != nil {
		builder = builder.Set("end_time", *p.EndTime)
	}
	if p.Type != nil {
		builder = builder.Set("type", *p.Type)
	}
	if p.Notes != nil {
		builder = builder.Set("notes", *p.Notes)
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
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func annotationColumns() []string {
	columns := []string{
		"id", "video_id", "user_id", "start_time",
		"end_time", "type", "notes", "created_at", "updated_at",
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
			&endTime, &a.Type, &a.Notes, &a.CreatedAt, &a.UpdatedAt, &vidDuration,
		)
		a.VideoDuration = time.Duration(vidDuration) * time.Second
	} else {
		rErr = row.Scan(
			&a.ID, &a.VideoID, &a.UserID, &startTime,
			&endTime, &a.Type, &a.Notes, &a.CreatedAt, &a.UpdatedAt,
		)
	}
	if rErr != nil {
		return nil, fmt.Errorf("failed to scan annotation: %w", rErr)
	}
	a.StartTime = time.Duration(startTime) * time.Second
	a.EndTime = time.Duration(endTime) * time.Second
	return &a, nil
}
