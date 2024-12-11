package api

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func (s *Server) errorResponse(c *gin.Context, code int, err error) {
	c.JSON(code, map[string]string{"error": err.Error()})
}

var (
	ErrParsingCreateDate = errors.New("Error parsing date")
	ErrParsingNumber     = errors.New("Error parsing number")
	ErrParsingID         = errors.New("Error parsing id")
	ErrNameIsEmpty       = errors.New("Name is empty")
	ErrGroupIsEmpty      = errors.New("Group is empty")
	ErrLinkNotCorrect    = errors.New("Link is not correct")
	ErrTextIsEmpty       = errors.New("Text is empty")
	errInvalidRequest    = errors.New("Incorrect parameters")
	errInvalidText       = errors.New("Incorrect text")
)
