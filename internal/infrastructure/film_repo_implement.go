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

// CreateFilm creates a new given film in the database connection
func (pg *FilmRepo) CreateFilm(film *domain.Film) error {
	if err := pg.dbConnection.Create(film).Error; err != nil {
		return err
	}
	return nil
}

// GetFilmByTitle gets a film details (this is all the fields for the film) for the given title
func (pg *FilmRepo) GetFilmByTitle(title string) error {
	var film domain.Film
	if err := pg.dbConnection.First(&film, title).Error; err != nil {
		return err
	}
	return nil
}

// PatchFilm updates a film details for the given film by title. Only updates the fields defined in the map newFilmFields
// Uses the GORM `Where` and `Updates` method to update the film fields
func (pg *FilmRepo) PatchFilm(filmTitleToUpdate string, newFilmFields *map[string]interface{}) (*domain.Film, error) {
	var film domain.Film
	if err := pg.dbConnection.Model(&film).Where("title = ?", filmTitleToUpdate).Updates(newFilmFields).Error; err != nil {
		return nil, err
	}
	return &film, nil
}

// PutFilm puts and saves a film details. Updates all the row for the given film, if there are missing fields, then default values will be filled for those missing fields.
// Uses the GORM `Where` and `Save` method to put the film update
// If the value doesn't have primary key, will insert it.
func (pg *FilmRepo) PutFilm(film *domain.Film) (*domain.Film, error) {
	if err := pg.dbConnection.Where("title = ?", film.Title).Save(&film).Error; err != nil {
		return nil, err
	}
	return film, nil
}

// GetAllFilms retrieves a list of films from the database, optionally filtered by the provided filter.
// The function uses the GORM `Where` method to apply any filters provided in the `optionalFilter` parameter.
// If no filters are provided, it retrieves all films from the database.
func (pg *FilmRepo) GetAllFilms(optionalFilter *domain.FilmFilter) ([]domain.Film, error) {
	var films []domain.Film
	if err := pg.dbConnection.Where(optionalFilter).Find(&films).Error; err != nil {
		return nil, err
	}
	return films, nil
}

// DeleteFilm removes a film from the database for the given title.
func (pg *FilmRepo) DeleteFilm(title string) (*domain.Film, error) {
	var film domain.Film
	if err := pg.dbConnection.Where("title = ?", title).Delete(&film).Error; err != nil {
		return nil, err
	}
	return &film, nil
}

// GetCreatorIdByTitle retrieves the creator user ID associated with a specific film title.
// It selects the `creator_user_id` column from the `films` table where the `title` matches the provided title.
func (pg *FilmRepo) GetCreatorIdByTitle(title string) (int, error) {
	var creatorId int
	if err := pg.dbConnection.Table("films").Select("creator_user_id").Where("title = ?", title).Scan(&creatorId).Error; err != nil {
		return -1, err
	}
	return creatorId, nil
}
