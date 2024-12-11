package repo

import "errors"

type tansaction string

const (
	tableSong                    = "songs"
	suffixReturningID            = "RETURNING id"
	tansactionKey     tansaction = "tansactionSQL"
)

var (
	ErrSongNotFound = errors.New("song not found")
	ErrParserJsonb  = errors.New("error text parser jsonb")
)
