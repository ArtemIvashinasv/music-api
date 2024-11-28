package repository

import (
	"database/sql"

	"github.com/artemivashinasv/music-api/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveSong(song *models.Song) error {
	query := `
INSERT INTO songs (group_name, song_name, release_date, text, link)
VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRow(query, song.Group, song.Song, song.ReleaseDate, song.Text, song.Link).Scan(&song.ID)
}

func (r *Repository) GetAllSongs(offset, limit int) ([]models.Song, error) {
	query := `
SELECT id, group_name, song_name, release_date, text, link, created_at, updated_at 
FROM songs ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var song models.Song
		err := rows.Scan(
			&song.ID, &song.Group, &song.Song, &song.ReleaseDate,
			&song.Text, &song.Link, &song.CreatedAt, &song.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	return songs, nil
}

func (r *Repository) GetSongByID(id uint) (*models.Song, error) {
	query := `
SELECT id, group_name, song_name, release_date, text, link, created_at, updated_at 
FROM songs WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var song models.Song
	err := row.Scan(
		&song.ID, &song.Group, &song.Song, &song.ReleaseDate,
		&song.Text, &song.Link, &song.CreatedAt, &song.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &song, nil
}

func (r *Repository) UpdateSong(song *models.Song) error {
	query := `
UPDATE songs SET group_name = $1, song_name = $2, release_date = $3, text = $4, link = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $6`
	_, err := r.db.Exec(query, song.Group, song.Song, song.ReleaseDate, song.Text, song.Link, song.ID)
	return err
}

func (r *Repository) DeleteSong(id uint) error {
	query := `DELETE FROM songs WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
