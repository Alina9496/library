package api

import (
	"context"

	"github.com/Alina9496/library/internal/domain"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, song *domain.Song) (*uuid.UUID, error)
	Update(ctx context.Context, song *domain.Song) error
	Delete(ctx context.Context, id *uuid.UUID) error
	GetTextSong(ctx context.Context, filter *domain.SongRequest) (string, error)
	GetSongs(ctx context.Context, filter *domain.SongRequest) ([]domain.Song, error)
}
