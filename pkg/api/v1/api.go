package v1

type SongItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Song struct {
	Text        []SongItem `json:"text,omitempty"`
	Name        string     `json:"name"`
	Group       string     `json:"group"`
	ReleaseDate string     `json:"release_date"`
	Link        string     `json:"link"`
}

type RespID struct {
	ID string `json:"id"`
}
