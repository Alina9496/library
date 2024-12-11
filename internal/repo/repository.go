package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Alina9496/library/internal/domain"
	"github.com/Alina9496/tool/pkg/logger"
	"github.com/Alina9496/tool/pkg/postgres"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type Repository struct {
	pg *postgres.Postgres
	l  *logger.Logger
}

func New(pg *postgres.Postgres, l *logger.Logger) *Repository {
	return &Repository{
		pg: pg,
		l:  l,
	}
}

type db interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func (r *Repository) conn(ctx context.Context) db {
	if tx, ok := ctx.Value(tansactionKey).(pgx.Tx); ok {
		return tx
	}
	return r.pg.Pool
}

func (r *Repository) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := ctx.Value(tansactionKey).(pgx.Tx); ok {
		return fn(ctx)
	}

	tx, err := r.pg.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, tansactionKey, tx)

	defer func() {
		if p := recover(); p != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				r.l.Error("rollback err %s", errRollback)
			}
			err = fmt.Errorf("panic :%s", p)
			return
		}
		if errCommit := tx.Commit(ctx); errCommit != nil {
			r.l.Error("commit err %s", errCommit)
		}
	}()
	return fn(ctx)
}

func (r *Repository) Create(ctx context.Context, song *domain.Song) (*uuid.UUID, error) {
	now := time.Now()
	query, args, err := r.pg.Builder.
		Insert(tableSong).
		Columns(
			"name",
			"executor",
			"text",
			"link",
			"release_date",
			"created_at",
			"updated_at",
		).
		Values(
			song.Name,
			song.Group,
			song.Text,
			song.Link,
			song.ReleaseDate,
			now,
			now,
		).
		Suffix(suffixReturningID).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error build query: %w", err)
	}

	var id uuid.UUID
	err = r.conn(ctx).QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error create: %w", err)
	}

	return &id, nil
}

func (r *Repository) Update(ctx context.Context, song *domain.Song) error {
	valuesMap := map[string]any{
		"name":         song.Name,
		"executor":     song.Group,
		"text":         song.Text,
		"link":         song.Link,
		"release_date": song.ReleaseDate,
		"updated_at":   time.Now(),
	}

	query, args, err := r.pg.Builder.
		Update(tableSong).
		SetMap(valuesMap).
		Where(squirrel.Eq{"id": song.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("error build query: %w", err)
	}

	commandTag, err := r.conn(ctx).Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error update: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return ErrSongNotFound
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id *uuid.UUID) error {
	query, args, err := r.pg.Builder.
		Delete(tableSong).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("error build query: %w", err)
	}

	commandTag, err := r.conn(ctx).Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error delete: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return ErrSongNotFound
	}
	return nil
}

func (r *Repository) GetTextSong(ctx context.Context, filter *domain.SongRequest) (domain.SongText, error) {
	query, args, err := r.pg.Builder.
		Select("text").
		From(tableSong).
		Where(squirrel.Eq{"executor": filter.Group}).
		Where(squirrel.Eq{"name": filter.Name}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error build query: %w", err)
	}

	var text domain.SongText
	err = r.conn(ctx).QueryRow(ctx, query, args...).Scan(&text)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSongNotFound
		}
		return nil, fmt.Errorf("error get text song: %w", err)
	}

	return text, nil
}

func (r *Repository) GetSongs(ctx context.Context, filter *domain.SongRequest) ([]domain.Song, error) {
	where := squirrel.And{}
	if filter.Group != "" {
		where = append(where, squirrel.Like{"executor": "%" + filter.Group + "%"})
	}
	if filter.Name != "" {
		where = append(where, squirrel.Like{"name": "%" + filter.Name + "%"})
	}
	if filter.Link != "" {
		where = append(where, squirrel.Like{"link": "%" + filter.Link + "%"})
	}
	if !filter.ReleaseDate.IsZero() {
		where = append(where, squirrel.Eq{"release_date": filter.ReleaseDate})
	}

	query, args, err := r.pg.Builder.Select(
		"name",
		"executor",
		"link",
		"release_date",
	).From(tableSong).
		Where(where).
		Limit(uint64(filter.Limit)).
		Offset(uint64(filter.Offset)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error build query: %w", err)
	}

	songs := make([]domain.Song, 0, filter.Limit)
	rows, err := r.conn(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var s domain.Song
		err := rows.Scan(&s.Name, &s.Group, &s.Link, &s.ReleaseDate)
		if err != nil {
			return nil, err
		}
		songs = append(songs, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return songs, nil
}
