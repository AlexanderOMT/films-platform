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

func NewFilmService(filmRepo domain.FilmRepository) FilmService {
	return &FilmServiceImpl{filmRepo: filmRepo}
}

func (f *FilmServiceImpl) CreateFilm(film *domain.Film) error {
	err := f.filmRepo.CreateFilm(film)
	if err != nil {
		return fmt.Errorf("error creating new film '%s': %w", film.Title, err)
	}
	return nil
}
func (f *FilmServiceImpl) GetAllFilms(optionalFilter *domain.FilmFilter) ([]domain.Film, error) {
	films, err := f.filmRepo.GetAllFilms(optionalFilter)
	if err != nil {
		return nil, fmt.Errorf("error retrieving films: %w", err)
	}
	return films, nil
}

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
