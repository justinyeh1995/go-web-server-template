package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerCheckChirpy(w http.ResponseWriter, r *http.Request) {
	type param struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := param{}
	err := decoder.Decode(&params)

	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		ResponseWithError(w, 500, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}

	if len(params.Body) > 140 {
		ResponseWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	RespondWithJson(w, 200, regularReponse{
		Valid: true,
	})
}
