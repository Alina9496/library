package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Alina9496/library/internal/domain"
	v1 "github.com/Alina9496/library/pkg/api/v1"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_toDomainSong(t *testing.T) {
	dateString := "2006-01-02"
	date, _ := time.Parse(time.DateOnly, dateString)

	tests := []struct {
		name    string
		song    v1.Song
		want    *domain.Song
		wantErr error
	}{
		{
			name: "name is empty",
			song: v1.Song{
				Group:       "group",
				Name:        "",
				Link:        "link",
				ReleaseDate: dateString,
				Text: []v1.SongItem{
					{
						Type: "verse",
						Text: "text",
					},
				},
			},
			want:    nil,
			wantErr: ErrNameIsEmpty,
		},
		{
			name: "group is empty",
			song: v1.Song{
				Group:       "",
				Name:        "name",
				Link:        "link",
				ReleaseDate: dateString,
				Text: []v1.SongItem{
					{
						Type: "verse",
						Text: "text",
					},
				},
			},
			want:    nil,
			wantErr: ErrGroupIsEmpty,
		},
		{
			name: "link is not correct",
			song: v1.Song{
				Group:       "group",
				Name:        "name",
				Link:        "link",
				ReleaseDate: dateString,
				Text: []v1.SongItem{
					{
						Type: "verse",
						Text: "text",
					},
				},
			},
			want:    nil,
			wantErr: ErrLinkNotCorrect,
		},
		{
			name: "error parsing date",
			song: v1.Song{
				Group:       "group",
				Name:        "name",
				Link:        "https://example.org/",
				ReleaseDate: "2006-0102",
				Text: []v1.SongItem{
					{
						Type: "verse",
						Text: "text",
					},
				},
			},
			want:    nil,
			wantErr: ErrParsingCreateDate,
		},
		{
			name: "text is empty",
			song: v1.Song{
				Group:       "group",
				Name:        "name",
				Link:        "https://example.org/",
				ReleaseDate: dateString,
				Text: []v1.SongItem{
					{
						Type: "",
						Text: "text",
					},
				},
			},
			want:    nil,
			wantErr: ErrTextIsEmpty,
		},
		{
			name: "text incorrect",
			song: v1.Song{
				Group:       "group",
				Name:        "name",
				Link:        "https://example.org/",
				ReleaseDate: dateString,
				Text: []v1.SongItem{
					{
						Type: "test",
						Text: "text",
					},
				},
			},
			want:    nil,
			wantErr: errInvalidText,
		},
		{
			name: "conversion from v1 in domain",
			song: v1.Song{
				Group:       "group",
				Name:        "name",
				Link:        "https://example.org/",
				ReleaseDate: dateString,
				Text: []v1.SongItem{
					{
						Type: "verse",
						Text: "text",
					},
				},
			},
			want: &domain.Song{
				Group:       "group",
				Name:        "name",
				Link:        "https://example.org/",
				ReleaseDate: date,
				Text: domain.SongText{
					{
						Type: "verse",
						Text: "text",
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toDomainSong(tt.song)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_toRespID(t *testing.T) {
	id := uuid.New()
	tests := []struct {
		name string
		id   *uuid.UUID
		want v1.RespID
	}{
		{
			name: "conversion from id in v1.RespID",
			id:   &id,
			want: v1.RespID{
				ID: id.String(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toRespID(tt.id)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_toGetTextSongRequest(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  *domain.SongRequest
		err   error
	}{
		{
			name:  "group is empty",
			query: "/test?name=name&offset=1",
			want:  nil,
			err:   errInvalidRequest,
		},
		{
			name:  "name is empty",
			query: "/test?group=group&offset=1",
			want:  nil,
			err:   errInvalidRequest,
		},
		{
			name:  "offset is empty",
			query: "/test?group=group&name=name",
			want:  nil,
			err:   errInvalidRequest,
		},
		{
			name:  "error parsing offset",
			query: "/test?group=group&name=name&offset=e",
			want:  nil,
			err:   ErrParsingNumber,
		},
		{
			name:  "conversion in domain.SongRequest",
			query: "/test?group=group&name=name&offset=1",
			want: &domain.SongRequest{
				Group:  "group",
				Name:   "name",
				Offset: 1,
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request, _ = http.NewRequest(http.MethodGet, tt.query, nil)
			got, err := toGetTextSongRequest(c)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.err, err)
		})
	}
}

func Test_toGetSongsRequest(t *testing.T) {
	date, _ := time.Parse(time.DateOnly, "2006-01-02")

	tests := []struct {
		name    string
		query   string
		want    *domain.SongRequest
		wantErr error
	}{
		{
			name:    "error parsing offset",
			query:   "/test?group=group&name=name&link=link&release_date=2006-01-02&offset=e&limit=1",
			want:    nil,
			wantErr: ErrParsingNumber,
		},
		{
			name:    "error parsing limit",
			query:   "/test?group=group&name=name&link=link&release_date=2006-01-02&offset=1&limit=e",
			want:    nil,
			wantErr: ErrParsingNumber,
		},
		{
			name:    "error parsing date",
			query:   "/test?group=group&name=name&link=link&release_date=2006-0102&offset=1&limit=1",
			want:    nil,
			wantErr: ErrParsingCreateDate,
		},
		{
			name:  "conversion in domain.SongRequest",
			query: "/test?group=group&name=name&link=link&release_date=2006-01-02&offset=1&limit=1",
			want: &domain.SongRequest{
				Group:       "group",
				Name:        "name",
				Link:        "link",
				ReleaseDate: date,
				Limit:       1,
				Offset:      1,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request, _ = http.NewRequest(http.MethodGet, tt.query, nil)
			got, err := toGetSongsRequest(c)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_toGetSongsResponse(t *testing.T) {
	dateString := "2006-01-02"
	date, _ := time.Parse(time.DateOnly, dateString)

	tests := []struct {
		name string
		s    []domain.Song
		want []v1.Song
	}{
		{
			name: "conversion from domain.Song in v1.Song",
			s: []domain.Song{
				{
					Group:       "group",
					Name:        "name",
					Link:        "link",
					ReleaseDate: date,
				},
			},
			want: []v1.Song{
				{
					Group:       "group",
					Name:        "name",
					Link:        "link",
					ReleaseDate: dateString,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toGetSongsResponse(tt.s)
			assert.Equal(t, tt.want, got)
		})
	}
}
