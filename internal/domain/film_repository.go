package domain

type filmRepository interface {
	CreateFilm(film Film) error
	UpdateFilm(title string) error
	DeleteFilm(title string) error
	GetFilm(title string) error
	GetFilms(title string) error
}
