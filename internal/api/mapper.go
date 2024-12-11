package api

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Alina9496/library/internal/domain"
	service "github.com/Alina9496/library/internal/service"
	v1 "github.com/Alina9496/library/pkg/api/v1"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func toDomainSong(song v1.Song) (*domain.Song, error) {
	if song.Name == "" {
		return nil, ErrNameIsEmpty
	}
	if song.Group == "" {
		return nil, ErrGroupIsEmpty
	}

	u, err := url.Parse(song.Link)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return nil, ErrLinkNotCorrect
	}

	parsedDate, err := time.Parse(time.DateOnly, song.ReleaseDate)
	if err != nil {
		return nil, ErrParsingCreateDate
	}

	songText := make(domain.SongText, 0, len(song.Text))
	for _, val := range song.Text {
		if val.Type == "" || val.Text == "" {
			return nil, ErrTextIsEmpty
		}
		songText = append(songText, domain.SongItem{
			Type: domain.TypeSongItem(val.Type),
			Text: val.Text,
		})
	}

	if !songText.IsValidType() {
		return nil, errInvalidText
	}

	return &domain.Song{
		Text:        songText,
		Name:        song.Name,
		Group:       song.Group,
		ReleaseDate: parsedDate,
		Link:        song.Link,
	}, nil
}

func toRespID(id *uuid.UUID) v1.RespID {
	return v1.RespID{
		ID: id.String(),
	}
}

func toGetTextSongRequest(c *gin.Context) (*domain.SongRequest, error) {
	if c.Query("group") == "" || c.Query("name") == "" || c.Query("offset") == "" {
		return nil, errInvalidRequest
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		return nil, ErrParsingNumber
	}

	return &domain.SongRequest{
		Group:  c.Query("group"),
		Name:   c.Query("name"),
		Offset: offset,
	}, nil
}

func toGetSongsRequest(c *gin.Context) (*domain.SongRequest, error) {
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		return nil, ErrParsingNumber
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		return nil, ErrParsingNumber
	}

	filter := domain.SongRequest{
		Group:  c.Query("group"),
		Name:   c.Query("name"),
		Link:   c.Query("link"),
		Offset: offset,
		Limit:  limit,
	}

	if c.Query("release_date") != "" {
		filter.ReleaseDate, err = time.Parse(time.DateOnly, c.Query("release_date"))
		if err != nil {
			return nil, ErrParsingCreateDate
		}
	}

	return &filter, nil
}

func toGetSongsResponse(s []domain.Song) []v1.Song {
	songs := make([]v1.Song, 0, len(s))
	for _, value := range s {
		songs = append(songs, v1.Song{
			Name:        value.Name,
			Group:       value.Group,
			ReleaseDate: value.ReleaseDate.Format(time.DateOnly),
			Link:        value.Link,
		})
	}
	return songs
}

func errToHttpStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch {
	case errors.Is(err, ErrParsingCreateDate),
		errors.Is(err, ErrParsingID),
		errors.Is(err, ErrNameIsEmpty),
		errors.Is(err, ErrGroupIsEmpty),
		errors.Is(err, ErrLinkNotCorrect),
		errors.Is(err, ErrTextIsEmpty):
		return http.StatusBadRequest
	case errors.Is(err, service.ErrSongNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
