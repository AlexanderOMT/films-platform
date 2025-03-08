package domain

type FilmRepository interface {
	CreateFilm(film *Film) error
	PatchFilm(filmTitleToUpdate string, newFilmFields *map[string]interface{}) (*Film, error)
	PutFilm(film *Film) (*Film, error)
	DeleteFilm(title string) (*Film, error)

	GetFilmByTitle(title string) error
	GetAllFilms(optionalFilter *FilmFilter) ([]Film, error)
	GetCreatorIdByTitle(title string) (int, error)
}
