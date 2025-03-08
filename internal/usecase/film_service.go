package usecase

import (
	"fmt"

	"golang-api-film-management/internal/domain"
)

type FilmService interface {
	CreateFilm(film *domain.Film) error
	DeleteFilm(title string, creatorUserId int) (*domain.Film, error)
	PatchFilm(filmTitle string, filmField *map[string]interface{}, userdId int) (*domain.Film, error)
	PutFilm(filmTitle string, film *domain.Film, userdId int) (*domain.Film, error)

	GetAllFilms(optionalFilter *domain.FilmFilter) ([]domain.Film, error)
}

type FilmServiceImpl struct {
	filmRepo domain.FilmRepository
}

// NewFilmService creates a new instance of FilmService with the provided FilmRepository.
func NewFilmService(filmRepo domain.FilmRepository) FilmService {
	return &FilmServiceImpl{filmRepo: filmRepo}
}

// CreateFilm creates a new film record in the repository.
// It takes a pointer to a domain.Film object as input and returns an error if the creation fails.
// If the film is successfully created, it returns nil.
func (f *FilmServiceImpl) CreateFilm(film *domain.Film) error {
	err := f.filmRepo.CreateFilm(film)
	if err != nil {
		return fmt.Errorf("error creating new film '%s': %w", film.Title, err)
	}
	return nil
}

// GetAllFilms retrieves all films from the repository, optionally filtered by the provided FilmFilter.
// If an error occurs during retrieval, it returns an error.
func (f *FilmServiceImpl) GetAllFilms(optionalFilter *domain.FilmFilter) ([]domain.Film, error) {
	films, err := f.filmRepo.GetAllFilms(optionalFilter)
	if err != nil {
		return nil, fmt.Errorf("error retrieving films: %w", err)
	}
	return films, nil
}

// DeleteFilm deletes a film with the given title if the user with the specified userId is the owner of the film.
// It returns the deleted film and an error if any occurred during the process.
func (f *FilmServiceImpl) DeleteFilm(title string, userId int) (*domain.Film, error) {
	if err := f.validateFilmOwnershipForUserId(title, userId); err != nil {
		return nil, fmt.Errorf("user %d is not owner for the film %v", userId, title)
	}

	film, err := f.filmRepo.DeleteFilm(title)
	if err != nil {
		return nil, fmt.Errorf("error deleting film '%s': %w", title, err)
	}
	return film, nil
}

// PatchFilm updates the fields of a film identified by its title. The fields provided by the map can contains all or some of the film fields
// It first validates if the user with the given userId is the owner of the film.
// If the user is not the owner, it returns an error.
// If the user is the owner, it proceeds to update the film with the provided fields.
func (f *FilmServiceImpl) PatchFilm(filmTitle string, filmFields *map[string]interface{}, userId int) (*domain.Film, error) {
	if err := f.validateFilmOwnershipForUserId(filmTitle, userId); err != nil {
		return nil, fmt.Errorf("user %d is not owner for the film %v", userId, filmTitle)
	}

	film, err := f.filmRepo.PatchFilm(filmTitle, filmFields)
	if err != nil {
		return nil, fmt.Errorf("error updating film '%s': %w", filmTitle, err)
	}

	return film, nil
}

// PutFilm updates the film row if the user is the owner.
// It validates the film ownership and updates the film in the repository. For those missing fields, default values will be filled
// Returns the updated film or an error if the update fails.
func (f *FilmServiceImpl) PutFilm(filmTitle string, film *domain.Film, userId int) (*domain.Film, error) {
	if err := f.validateFilmOwnershipForUserId(filmTitle, userId); err != nil {
		return nil, fmt.Errorf("user %d is not owner for the film %v", userId, filmTitle)
	}

	film.CreatorUserID = userId
	film.Title = filmTitle

	updatedFilm, err := f.filmRepo.PutFilm(film)
	if err != nil {
		return nil, fmt.Errorf("error updating film '%s': %w", filmTitle, err)
	}

	return updatedFilm, nil
}

// validateFilmOwnershipForUserId checks if the given userId is the owner of the film with the specified filmTitle.
// It retrieves the creator user ID from the film repository and compares it with the provided userId.
// If the userId does not match the creator user ID, an error is returned indicating that the user is not the owner of the film.
func (f *FilmServiceImpl) validateFilmOwnershipForUserId(filmTitle string, userId int) error {
	creatorUserId, err := f.filmRepo.GetCreatorIdByTitle(filmTitle)
	if err != nil {
		return fmt.Errorf("error fetching creator ID for film '%s': %w", filmTitle, err)
	}

	if userId != creatorUserId {
		return fmt.Errorf("user %d is not owner for the film %v", userId, filmTitle)
	}

	return nil
}
