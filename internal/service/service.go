package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alina9496/library/internal/domain"
	"github.com/Alina9496/library/internal/repo"
	"github.com/Alina9496/tool/pkg/logger"
	"github.com/google/uuid"
)

type Service struct {
	repo Repository
	log  *logger.Logger
}

func New(
	r Repository,
	log *logger.Logger,
) *Service {
	return &Service{
		repo: r,
		log:  log,
	}
}

func (s *Service) Create(ctx context.Context, song *domain.Song) (*uuid.UUID, error) {
	l := s.log.WithField("service_method", "Create")
	if song == nil {
		l.Debug(ErrSongIsNil.Error())
		return nil, ErrSongIsNil
	}

	id, err := s.repo.Create(ctx, song)
	if err != nil {
		l.WithError(err).Error("error when create")
		return nil, fmt.Errorf("error when create: %w", ErrCreateSong)
	}

	l.WithField("id", id).Info("create song was successfully")
	return id, nil
}

func (s *Service) Update(ctx context.Context, song *domain.Song) error {
	l := s.log.WithField("service_method", "Update")
	if song == nil {
		l.Debug(ErrSongIsNil.Error())
		return ErrSongIsNil
	}

	err := s.repo.Update(ctx, song)
	if err != nil {
		if errors.Is(err, repo.ErrSongNotFound) {
			return ErrSongNotFound
		}
		l.WithError(err).Error("error when update")
		return fmt.Errorf("error when update: %w", ErrUpdateSong)
	}

	l.Info("update song was successfully")
	return nil
}

func (s *Service) Delete(ctx context.Context, id *uuid.UUID) error {
	l := s.log.WithField("service_method", "Delete")
	if id == nil {
		l.Debug(ErrIDIsNil.Error())
		return ErrIDIsNil
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrSongNotFound) {
			return ErrSongNotFound
		}
		l.WithError(err).Error("error when delete")
		return fmt.Errorf("error when delete: %w", ErrDeleteSong)
	}

	l.WithField("id", id).Info("delete song was successfully")
	return nil
}

func (s *Service) GetTextSong(ctx context.Context, filter *domain.SongRequest) (string, error) {
	l := s.log.WithField("service_method", "GetTextSong")
	if filter == nil {
		l.Debug(ErrFilterIsNil.Error())
		return "", ErrFilterIsNil
	}

	songText, err := s.repo.GetTextSong(ctx, filter)
	if err != nil {
		if errors.Is(err, repo.ErrSongNotFound) {
			return "", ErrSongNotFound
		}
		l.WithError(err).Error("error when getTextSong")
		return "", fmt.Errorf("error when getTextSong: %w", ErrGetSong)
	}

	count := 0
	for _, val := range songText {
		if val.Type == domain.Verse {
			count++
			if count == filter.Offset {
				l.Info("the song text was found successfully")
				return val.Text, nil
			}
		}
	}

	return "", ErrSongNotFound
}

func (s *Service) GetSongs(ctx context.Context, filter *domain.SongRequest) ([]domain.Song, error) {
	l := s.log.WithField("service_method", "GetSongs")
	if filter == nil {
		l.Debug(ErrFilterIsNil.Error())
		return nil, ErrFilterIsNil
	}

	songs, err := s.repo.GetSongs(ctx, filter)
	if err != nil {
		if errors.Is(err, repo.ErrSongNotFound) {
			return nil, ErrSongNotFound
		}
		l.WithError(err).Error("error when getSongs")
		return nil, fmt.Errorf("error when getSongs: %w", ErrGetSong)
	}

	l.Info("the songs was found successfully")
	return songs, nil
}
