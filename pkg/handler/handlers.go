package handler

import (
	_ "github.com/artemivashinasv/music-api/docs" // Подключение документации
	"github.com/artemivashinasv/music-api/pkg/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

// @Summary Get all songs
// @Tags Songs
// @Description Retrieve all songs with pagination
// @Accept  json
// @Produce  json
// @Param   page   query     int     false  "Page number"
// @Param   limit  query     int     false  "Items per page"
// @Success 200 {array} models.Song
// @Failure 500 {object} gin.H "Internal Server Error"
// @Router /songs [get]

// Инициализация маршрутов
func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.New()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/songs", h.GetAllSongs)
	r.GET("/songs/:id", h.GetSongByID)
	r.POST("/songs", h.SaveSong)
	r.PUT("/songs/:id", h.UpdateSong)
	r.DELETE("/songs/:id", h.DeleteSong)

	return r
}
