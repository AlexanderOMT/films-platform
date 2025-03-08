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

// CreateFilm creates a new film in the system and make the relationship with their creator
// The relationship of the film with their creator is the subjecy user id extracted from the request context
// Its response is the film created or any error if encountred
func (f *FilmHandler) CreateFilm(w http.ResponseWriter, r *http.Request) {
	subjectID, ok := extractSubjectIdFromContext(r)
	if !ok || subjectID == 0 {
		log.Println("Error creating film: No subject ID found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var filmToCreate domain.Film
	if err := deserializeJSONFromRequest(r, &filmToCreate); err != nil {
		log.Printf("Error deserializing JSON for creating film: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filmToCreate.CreatorUserID = subjectID

	err := f.filmService.CreateFilm(&filmToCreate)
	if err != nil {
		log.Printf("Error creating film: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("Successfully created film | Film: %v", filmToCreate)
	writeJSONResponse(w, http.StatusOK, filmToCreate)
}

// GetAllFilms retrieves all films based on a custom filter from query parameters.
// It maps the query parameters with a `FilmFilter` struct
// Calls a service so delegates the responsability of applying the filter for retrieving the list of film
// Its response is the list of all the films with a opctional filter applied or any error if encountred
func (f *FilmHandler) GetAllFilms(w http.ResponseWriter, r *http.Request) {
	var customFilter domain.FilmFilter
	if err := schema.NewDecoder().Decode(&customFilter, r.URL.Query()); err != nil {
		log.Printf("Invalid query parameters for GetAllFilms: %v", err)
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	films, err := f.filmService.GetAllFilms(&customFilter)
	if err != nil {
		log.Printf("Error retrieving films list: %v", err)
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	log.Printf("Successfully retrieved all the films list | Films: %v", films)
	writeJSONResponse(w, http.StatusOK, films)
}

// DeleteFilm deletes a film based on the provided title.
// Calls a service so delegates the responsability if the subject user id is allowed (or not allowed) for deleting the given film
// Its response is the film removed or any error if encountred
func (f *FilmHandler) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	subjectID, ok := extractSubjectIdFromContext(r)
	if !ok || subjectID == 0 {
		log.Printf("Error deleting film: No subject ID found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	title := r.URL.Query().Get("title")

	film, err := f.filmService.DeleteFilm(title, subjectID)
	if err != nil {
		log.Printf("Error deleting film with title %v: %v", title, err)
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	log.Printf("Successfully deleted film | Film title: %v", title)
	writeJSONResponse(w, http.StatusOK, film)
}

// PatchFilm updates a film fields based on the provided JSON payload.
// It decoded the fields from the json and call a service to update the fields
// The service called has the responsability if the subject user id is allowed (or not allowed) for updating the given film
// Its response is the film updated or any error if encountred
func (f *FilmHandler) PatchFilm(w http.ResponseWriter, r *http.Request) {
	subjectID, ok := extractSubjectIdFromContext(r)
	if !ok || subjectID == 0 {
		log.Println("Error updating (PATCH) the film: No subject ID found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	title := r.URL.Query().Get("title")

	var filmToUpdateFields map[string]interface{}
	if err := deserializeJSONFromRequest(r, &filmToUpdateFields); err != nil {
		log.Printf("Error deserializing JSON for updating (PATCH) film: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	film, err := f.filmService.PatchFilm(title, &filmToUpdateFields, subjectID)
	if err != nil {
		log.Printf("Error updating (PATCH) film with title %s: %v", title, err)
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	log.Printf("Successfully updated (PATCH) film with title: %s . This are the fields provided: %v", title, filmToUpdateFields)
	writeJSONResponse(w, http.StatusOK, film)
}

// PutFilm puts and saves a film details.
// If there is not all the fields of the film given, then will return a error response
// The service called has the responsability if the subject user id is allowed (or not allowed) for updating the given film
// Its response is the film updated or any error if encountred
func (f *FilmHandler) PutFilm(w http.ResponseWriter, r *http.Request) {
	subjectID, ok := r.Context().Value("subjectId").(int)
	if !ok || subjectID == 0 {
		log.Println("Error updating (PUT) the film: No subject ID found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	title := r.URL.Query().Get("title")

	var filmToSave domain.Film
	if err := deserializeJSONFromRequest(r, &filmToSave); err != nil {
		log.Printf("Error deserializing JSON for saving (PUT) film: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var validate = validator.New()
	if err := validate.Struct(filmToSave); err != nil {
		log.Printf("Missing required fields for updating (PUT) film: %v", err)
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	film, err := f.filmService.PutFilm(title, &filmToSave, subjectID)
	if err != nil {
		log.Printf("Error saving (PUT) film with title %s: %v", title, err)
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	log.Printf("Successfully saved (PUT) film | Film title: %s", title)
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
