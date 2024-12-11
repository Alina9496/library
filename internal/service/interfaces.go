//go:generate mockgen -source=interfaces.go -destination=./mock_interfaces.go -package=service
package service

import (
	"context"

	"github.com/Alina9496/library/internal/domain"
	"github.com/google/uuid"
)

type Repository interface {
	ExecTx(ctx context.Context, fn func(ctx context.Context) error) error
	Create(ctx context.Context, song *domain.Song) (*uuid.UUID, error)
	Update(ctx context.Context, song *domain.Song) error
	Delete(ctx context.Context, id *uuid.UUID) error
	GetTextSong(ctx context.Context, filter *domain.SongRequest) (domain.SongText, error)
	GetSongs(ctx context.Context, filter *domain.SongRequest) ([]domain.Song, error)
}
