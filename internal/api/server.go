package api

import (
	v1 "github.com/Alina9496/library/pkg/api/v1"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"net/http"

	"github.com/gin-contrib/cors"

	"github.com/Alina9496/tool/pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	service Service
	l       *logger.Logger
}

func NewServer(handler *gin.Engine, l *logger.Logger, t Service) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowCredentials = true
	handler.Use(cors.New(corsConfig))

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))
	s := &Server{t, l}

	h := handler.Group("/api/v1")
	{
		h.POST("/song", s.Create)
		h.PATCH("/song/:id", s.Update)
		h.DELETE("/song/:id", s.Delete)
		h.GET("/song", s.GetTextSong)
		h.GET("/songs", s.GetSongs)
	}
}

func (s *Server) Create(c *gin.Context) {
	var sg v1.Song
	err := c.BindJSON(&sg)
	if err != nil {
		s.errorResponse(c, errToHttpStatus(err), err)
		return
	}

	song, err := toDomainSong(sg)
	if err != nil {
		s.errorResponse(c, errToHttpStatus(err), err)
		return
	}

	id, err := s.service.Create(c.Request.Context(), song)
	if err != nil {
		s.errorResponse(c, errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, toRespID(id))
}

func (s *Server) Update(c *gin.Context) {
	var sg v1.Song
	err := c.BindJSON(&sg)
	if err != nil {
		s.errorResponse(c, errToHttpStatus(err), err)
		return
	}

	song, err := toDomainSong(sg)
	if err != nil {
		s.errorResponse(c, errToHttpStatus(err), err)
		return
	}

	song.ID, err = uuid.Parse(c.Param("id"))
	if err != nil {
		s.errorResponse(c, errToHttpStatus(ErrParsingID), ErrParsingID)
		return
	}

	err = s.service.Update(c.Request.Context(), song)
	if err != nil {
		s.errorResponse(c, errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, map[string]string{"response": "ok"})
}

func (s *Server) Delete(c *gin.Context) {
	uuid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		s.errorResponse(c, errToHttpStatus(ErrParsingID), ErrParsingID)
		return
	}

	err = s.service.Delete(c.Request.Context(), &uuid)
	if err != nil {
		s.errorResponse(c, errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, map[string]string{"response": "ok"})
}

func (s *Server) GetTextSong(c *gin.Context) {
	param, err := toGetTextSongRequest(c)
	if err != nil {
		s.errorResponse(c, errToHttpStatus(err), err)
		return
	}

	text, err := s.service.GetTextSong(c, param)
	if err != nil {
		s.errorResponse(c, errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, map[string]string{"response": text})
}

func (s *Server) GetSongs(c *gin.Context) {
	filter, err := toGetSongsRequest(c)
	if err != nil {
		s.errorResponse(c, errToHttpStatus(err), err)
		return
	}

	songs, err := s.service.GetSongs(c, filter)
	if err != nil {
		s.errorResponse(c, errToHttpStatus(err), err)
		return
	}

	c.JSON(http.StatusOK, map[string]any{"response": toGetSongsResponse(songs)})
}
