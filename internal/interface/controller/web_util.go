package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// #TODO: style: These methods should be reachable from others handlers. Should reconsider the location and implementations of thsis
// Another idea is to encapsulate this logic into a web struct

func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
	}
}

func DeserializeJSONFromRequest(r *http.Request, data interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return fmt.Errorf("error deserializing the object from json: %v", err)
	}
	return nil
}

// extractSubjectIdFromContext extracts the subject id from the request context
// It assumes that is always with the key `subjectId` as `int`
// TODO: style: refactor this: seems not to be a not good practice as is not flexible
func ExtractSubjectIdFromContext(r *http.Request) (int, bool) {
	subjectID, ok := r.Context().Value("subjectId").(int)
	if !ok {
		log.Println("error extracting subjectId from context: not found or not an int")
		return 0, false
	}
	return subjectID, ok
}
