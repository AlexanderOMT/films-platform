package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"

	"golang-api-film-management/internal/domain"
	"golang-api-film-management/internal/usecase"
)

type FilmHandler struct {
	filmService usecase.FilmService
}

func NewFilmHandler(filmService usecase.FilmService) *FilmHandler {
	return &FilmHandler{filmService: filmService}
}

// These are protected routes. Should be assume that there is always a token in it header?

// #TODO: enhance: implements better logger

// #TODO: fix: http.Error should return a proper status depending on the situation. Check in all the API

// #TODO: enhance: implements input validations

func (f *FilmHandler) CreateFilm(w http.ResponseWriter, r *http.Request) {
	subjectID, ok := extractSubjectIdFromContext(r)
	if !ok || subjectID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var filmToCreate domain.Film
	if err := deserializeJSONFromRequest(r, &filmToCreate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filmToCreate.CreatorUserID = subjectID

	err := f.filmService.CreateFilm(&filmToCreate)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	writeJSONResponse(w, http.StatusOK, filmToCreate)
}

func (f *FilmHandler) GetAllFilms(w http.ResponseWriter, r *http.Request) {
	var customFilter domain.FilmFilter
	if err := schema.NewDecoder().Decode(&customFilter, r.URL.Query()); err != nil {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	films, err := f.filmService.GetAllFilms(&customFilter)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	writeJSONResponse(w, http.StatusOK, films)
}

func (f *FilmHandler) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	subjectID, ok := extractSubjectIdFromContext(r)
	if !ok || subjectID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	title := r.URL.Query().Get("title")

	film, err := f.filmService.DeleteFilm(title, subjectID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	writeJSONResponse(w, http.StatusOK, film)
}

func (f *FilmHandler) PatchFilm(w http.ResponseWriter, r *http.Request) {
	subjectID, ok := extractSubjectIdFromContext(r)
	if !ok || subjectID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	title := r.URL.Query().Get("title")

	var filmToUpdateFields map[string]interface{}
	if err := deserializeJSONFromRequest(r, &filmToUpdateFields); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Update is: %v", filmToUpdateFields)
	film, err := f.filmService.PatchFilm(title, &filmToUpdateFields, subjectID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	writeJSONResponse(w, http.StatusOK, film)
}

func (f *FilmHandler) PutFilm(w http.ResponseWriter, r *http.Request) {
	subjectID, ok := r.Context().Value("subjectId").(int)
	if !ok || subjectID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	title := r.URL.Query().Get("title")

	var filmToSave domain.Film
	if err := deserializeJSONFromRequest(r, &filmToSave); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var validate = validator.New()
	if err := validate.Struct(filmToSave); err != nil {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	film, err := f.filmService.PutFilm(title, &filmToSave, subjectID)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	writeJSONResponse(w, http.StatusOK, film)
}

// #TODO: style: These methods should be reachable from others handlers. Should reconsider the location and implementations of thsis
// Another idea is to encapsulate this logic into a web struct
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
	}
}

func deserializeJSONFromRequest(r *http.Request, data interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return fmt.Errorf("error deserializing the object from json: %v", err)
	}
	return nil
}

// TODO: style :refactor this. This seems not to be a not good practice as is not flexible
func extractSubjectIdFromContext(r *http.Request) (int, bool) {
	subjectID, ok := r.Context().Value("subjectId").(int)
	return subjectID, ok
}
