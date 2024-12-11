package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type TypeSongItem string

const Verse TypeSongItem = "verse"
const Chorus TypeSongItem = "chorus"

type SongItem struct {
	Type TypeSongItem `json:"type,omitempty"`
	Text string       `json:"text,omitempty"`
}

type SongText []SongItem

type Song struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ReleaseDate time.Time
	ID          uuid.UUID
	Text        SongText
	Name        string
	Group       string
	Link        string
}

func (s SongText) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *SongText) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &s)
}

func (s SongText) IsValidType() bool {
	for _, v := range s {
		if v.Type != Verse && v.Type != Chorus {
			return false
		}
	}
	return true
}

type SongRequest struct {
	ReleaseDate time.Time
	Limit       int
	Offset      int
	Name        string
	Group       string
	Link        string
}
