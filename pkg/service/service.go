package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/artemivashinasv/music-api/models"
	"github.com/artemivashinasv/music-api/pkg/repository"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllSongs(page, limit int) ([]models.Song, error) {
	offset := (page - 1) * limit
	return s.repo.GetAllSongs(offset, limit)
}

func (s *Service) GetSongByID(id uint) (*models.Song, error) {
	return s.repo.GetSongByID(id)
}

func (s *Service) SaveSong(song *models.Song) error {
	// Запрос к внешнему API
	apiURL := fmt.Sprintf("http://external-api/info?group=%s&song=%s", song.Group, song.Song)
	resp, err := http.Get(apiURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("external API returned an error")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var detail models.SongDetail
	if err := json.Unmarshal(body, &detail); err != nil {
		return err
	}

	// Обогащение данных
	song.ReleaseDate = detail.ReleaseDate
	song.Text = detail.Text
	song.Link = detail.Link

	return s.repo.SaveSong(song)
}

func (s *Service) UpdateSong(song *models.Song) error {
	return s.repo.UpdateSong(song)
}

func (s *Service) DeleteSong(id uint) error {
	return s.repo.DeleteSong(id)
}
