package handler

import (
	"github.com/artemivashinasv/music-api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Получение всех песен
func (h *Handler) GetAllSongs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	songs, err := h.services.GetAllSongs(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch songs"})
		return
	}

	c.JSON(http.StatusOK, songs)
}

// Получение песни по ID
func (h *Handler) GetSongByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	song, err := h.services.GetSongByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	c.JSON(http.StatusOK, song)
}

// Сохранение новой песни
func (h *Handler) SaveSong(c *gin.Context) {
	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.services.SaveSong(&song); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save song"})
		return
	}

	c.JSON(http.StatusCreated, song)
}

// Обновление песни
func (h *Handler) UpdateSong(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	song.ID = uint(id)
	if err := h.services.UpdateSong(&song); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update song"})
		return
	}

	c.JSON(http.StatusOK, song)
}

// Удаление песни
func (h *Handler) DeleteSong(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.services.DeleteSong(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Song deleted"})
}