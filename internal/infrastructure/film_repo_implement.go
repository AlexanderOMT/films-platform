package infrastructure

// TODO change name of database package

import (
	_ "github.com/lib/pq" // SQL driver
	"gorm.io/gorm"

	"golang-api-film-management/internal/domain"
)

type FilmRepo struct {
	dbConnection *gorm.DB
}

func NewFilmRepo(db *gorm.DB) *FilmRepo {
	return &FilmRepo{dbConnection: db}
}

// Create film
func (pg *FilmRepo) CreateFilm(film domain.Film) error {
	return nil
}

// Get a film details
func (pg *FilmRepo) GetFilm(title string) error {
	return nil
}

// Update a film details
func (pg *FilmRepo) UpdateFilm(title string) error {
	return nil
}

// Get films list
func (pg *FilmRepo) GetFilms(title string) error {
	return nil
}

// Delete film
func (pg *FilmRepo) DeleteFilm(title string) error {
	return nil
}
