package infrastructure

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
func (pg *FilmRepo) CreateFilm(film *domain.Film) error {
	if err := pg.dbConnection.Create(film).Error; err != nil {
		return err
	}
	return nil
}

// Get a film details by title
func (pg *FilmRepo) GetFilmByTitle(title string) error {
	var film domain.Film
	if err := pg.dbConnection.First(&film, title).Error; err != nil {
		return err
	}
	return nil
}

// Update a film details
func (pg *FilmRepo) PatchFilm(filmTitleToUpdate string, newFilmFields *map[string]interface{}) (*domain.Film, error) {
	var film domain.Film
	if err := pg.dbConnection.Model(&film).Where("title = ?", filmTitleToUpdate).Updates(newFilmFields).Error; err != nil {
		return nil, err
	}
	return &film, nil
}

// Put/save a film details
func (pg *FilmRepo) PutFilm(film *domain.Film) (*domain.Film, error) {
	if err := pg.dbConnection.Where("title = ?", film.Title).Save(&film).Error; err != nil {
		return nil, err
	}
	return film, nil
}

// Get films list
func (pg *FilmRepo) GetAllFilms(optionalFilter *domain.FilmFilter) ([]domain.Film, error) {
	var films []domain.Film
	if err := pg.dbConnection.Where(optionalFilter).Find(&films).Error; err != nil {
		return nil, err
	}
	return films, nil
}

// Delete film
func (pg *FilmRepo) DeleteFilm(title string) (*domain.Film, error) {
	var film domain.Film
	if err := pg.dbConnection.Where("title = ?", title).Delete(&film).Error; err != nil {
		return nil, err
	}
	return &film, nil
}

func (pg *FilmRepo) GetCreatorIdByTitle(title string) (int, error) {
	var creatorId int
	if err := pg.dbConnection.Table("films").Select("creator_user_id").Where("title = ?", title).Scan(&creatorId).Error; err != nil {
		return -1, err
	}
	return creatorId, nil
}
