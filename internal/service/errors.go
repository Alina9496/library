package service

import "errors"

var (
	ErrSongNotFound = errors.New("song not found")
	ErrSongIsNil    = errors.New("song is nil")
	ErrIDIsNil      = errors.New("ID is nil")
	ErrFilterIsNil  = errors.New("filter is nil")
	ErrCreateSong   = errors.New("song not create")
	ErrUpdateSong   = errors.New("song not update")
	ErrDeleteSong   = errors.New("song not delete")
	ErrGetSong      = errors.New("error get song")
)
