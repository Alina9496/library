package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Alina9496/library/internal/domain"
	"github.com/Alina9496/library/internal/repo"
	"github.com/Alina9496/tool/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
	repo    *MockRepository
	service Service
}

func (s *ServiceSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.repo = NewMockRepository(ctrl)
	s.service = *New(s.repo, logger.New(""))
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func (s *ServiceSuite) Test_Create() {
	ctx := context.Background()
	id := uuid.New()
	song := &domain.Song{
		Group:       "group",
		Name:        "name",
		Link:        "link",
		ReleaseDate: time.Now(),
		Text: domain.SongText{
			{
				Type: "verse",
				Text: "text",
			},
		},
	}
	tests := []struct {
		name  string
		ctx   context.Context
		song  *domain.Song
		want  *uuid.UUID
		err   error
		calls func()
	}{
		{
			name:  "song equal nil",
			ctx:   ctx,
			song:  nil,
			want:  nil,
			err:   ErrSongIsNil,
			calls: func() {},
		},
		{
			name: "error create song",
			ctx:  ctx,
			song: song,
			want: nil,
			err:  fmt.Errorf("error when create: %w", ErrCreateSong),
			calls: func() {
				s.repo.EXPECT().Create(ctx, song).Return(nil, errors.ErrUnsupported)
			},
		},
		{
			name: "song was successfully created",
			ctx:  ctx,
			song: song,
			want: &id,
			err:  nil,
			calls: func() {
				s.repo.EXPECT().Create(ctx, song).Return(&id, nil)
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.calls()
			got, err := s.service.Create(tt.ctx, tt.song)
			s.Equal(tt.want, got)
			s.Equal(tt.err, err)
		})
	}
}

func (s *ServiceSuite) Test_Update() {
	ctx := context.Background()
	song := &domain.Song{
		Group:       "group",
		Name:        "name",
		Link:        "link",
		ReleaseDate: time.Now(),
		Text: domain.SongText{
			{
				Type: "verse",
				Text: "text",
			},
		},
	}
	tests := []struct {
		name  string
		ctx   context.Context
		song  *domain.Song
		err   error
		calls func()
	}{
		{
			name:  "song equal nil",
			ctx:   ctx,
			song:  nil,
			err:   ErrSongIsNil,
			calls: func() {},
		},
		{
			name: "song not found",
			ctx:  ctx,
			song: song,
			err:  fmt.Errorf("error when update: %w", ErrSongNotFound),
			calls: func() {
				s.repo.EXPECT().Update(ctx, song).Return(repo.ErrSongNotFound)
			},
		},
		{
			name: "error update song",
			ctx:  ctx,
			song: song,
			err:  fmt.Errorf("error when update: %w", ErrUpdateSong),
			calls: func() {
				s.repo.EXPECT().Update(ctx, song).Return(errors.ErrUnsupported)
			},
		},
		{
			name: "song was successfully update",
			ctx:  ctx,
			song: song,
			err:  nil,
			calls: func() {
				s.repo.EXPECT().Update(ctx, song).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.calls()
			err := s.service.Update(tt.ctx, tt.song)
			s.Equal(tt.err, err)
		})
	}
}

func (s *ServiceSuite) Test_Delete() {
	ctx := context.Background()
	id := uuid.New()
	tests := []struct {
		name  string
		ctx   context.Context
		id    *uuid.UUID
		err   error
		calls func()
	}{
		{
			name:  "id equal nil",
			ctx:   ctx,
			id:    nil,
			err:   ErrSongIsNil,
			calls: func() {},
		},
		{
			name: "song not found",
			ctx:  ctx,
			id:   &id,
			err:  ErrSongNotFound,
			calls: func() {
				s.repo.EXPECT().Delete(ctx, &id).Return(repo.ErrSongNotFound)
			},
		},
		{
			name: "error delete song",
			ctx:  ctx,
			id:   &id,
			err:  fmt.Errorf("error when delete: %w", ErrDeleteSong),
			calls: func() {
				s.repo.EXPECT().Delete(ctx, id).Return(errors.ErrUnsupported)
			},
		},
		{
			name: "song was successfully delete",
			ctx:  ctx,
			id:   &id,
			err:  nil,
			calls: func() {
				s.repo.EXPECT().Delete(ctx, id).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.calls()
			err := s.service.Delete(tt.ctx, &id)
			s.Equal(tt.err, err)
		})
	}
}

func (s *ServiceSuite) Test_GetTextSong() {
	ctx := context.Background()
	filter := &domain.SongRequest{
		Group:  "group",
		Name:   "name",
		Offset: 2,
	}
	filterNotFound := &domain.SongRequest{
		Group:  "group",
		Name:   "name",
		Offset: 5,
	}
	songText := domain.SongText{
		{Type: "verse", Text: "text1"},
		{Type: "chorus", Text: "text2"},
		{Type: "verse", Text: "text3"},
	}

	tests := []struct {
		name   string
		ctx    context.Context
		filter *domain.SongRequest
		wait   string
		err    error
		calls  func()
	}{
		{
			name:   "filter equal nil",
			ctx:    ctx,
			filter: nil,
			wait:   "",
			err:    ErrFilterIsNil,
			calls:  func() {},
		},
		{
			name:   "song not found",
			ctx:    ctx,
			filter: filter,
			wait:   "",
			err:    ErrSongNotFound,
			calls: func() {
				s.repo.EXPECT().GetTextSong(ctx, filter).Return(nil, repo.ErrSongNotFound)
			},
		},
		{
			name:   "error get song",
			ctx:    ctx,
			filter: filter,
			wait:   "",
			err:    fmt.Errorf("error when getTextSong: %w", ErrGetSong),
			calls: func() {
				s.repo.EXPECT().GetTextSong(ctx, filter).Return(nil, errors.ErrUnsupported)
			},
		},
		{
			name:   "get song",
			ctx:    ctx,
			filter: filter,
			wait:   "text3",
			err:    nil,
			calls: func() {
				s.repo.EXPECT().GetTextSong(ctx, filter).Return(songText, nil)
			},
		},
		{
			name:   "verse not found",
			ctx:    ctx,
			filter: filterNotFound,
			wait:   "",
			err:    ErrSongNotFound,
			calls: func() {
				s.repo.EXPECT().GetTextSong(ctx, filterNotFound).Return(songText, nil)
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.calls()
			text, err := s.service.GetTextSong(tt.ctx, tt.filter)
			s.Equal(tt.wait, text)
			s.Equal(tt.err, err)
		})
	}
}

func (s *ServiceSuite) Test_GetSongs() {
	ctx := context.Background()
	date, _ := time.Parse(time.DateOnly, "2006-01-02")
	filter := &domain.SongRequest{
		Group:       "group",
		Name:        "name",
		Link:        "https://lyrsense.com",
		ReleaseDate: date,
		Limit:       2,
		Offset:      2,
	}
	res := []domain.Song{
		{
			Name:        "name",
			Group:       "group",
			Link:        "https://lyrsense.com",
			ReleaseDate: date,
		},
	}

	tests := []struct {
		name   string
		ctx    context.Context
		filter *domain.SongRequest
		wait   []domain.Song
		err    error
		calls  func()
	}{
		{
			name:   "filter equal nil",
			ctx:    ctx,
			filter: nil,
			wait:   nil,
			err:    ErrFilterIsNil,
			calls:  func() {},
		},
		{
			name:   "song not found",
			ctx:    ctx,
			filter: filter,
			wait:   nil,
			err:    ErrSongNotFound,
			calls: func() {
				s.repo.EXPECT().GetSongs(ctx, filter).Return(nil, repo.ErrSongNotFound)
			},
		},
		{
			name:   "error get song",
			ctx:    ctx,
			filter: filter,
			wait:   nil,
			err:    fmt.Errorf("error when getSongs: %w", ErrGetSong),
			calls: func() {
				s.repo.EXPECT().GetSongs(ctx, filter).Return(nil, errors.ErrUnsupported)
			},
		},
		{
			name:   "get song",
			ctx:    ctx,
			filter: filter,
			wait:   res,
			err:    nil,
			calls: func() {
				s.repo.EXPECT().GetSongs(ctx, filter).Return(res, nil)
			},
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.calls()
			text, err := s.service.GetSongs(tt.ctx, tt.filter)
			s.Equal(tt.wait, text)
			s.Equal(tt.err, err)
		})
	}
}
