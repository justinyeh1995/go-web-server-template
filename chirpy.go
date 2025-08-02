package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/justinyeh1995/go-web-server-template/internal/database"
)

type Chirpy struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirpy(w http.ResponseWriter, r *http.Request) {
	type param struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := param{}
	err := decoder.Decode(&params)

	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}

	if len(params.Body) > 140 {
		ResponseWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	chirpy, err := cfg.db.CreateChirpy(r.Context(), database.CreateChirpyParams{
		Body:   params.Body,
		UserID: params.UserID,
	})
	if err != nil {
		// log.Fatalf("Error creating chirpy, with err msg %v", err)
		ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating chirpy, with err msg %v", err))
		return
	}

	// 3) Save it to the database
	RespondWithJson(w, http.StatusCreated, Chirpy{
		ID:        chirpy.ID,
		CreatedAt: chirpy.CreatedAt.Time,
		UpdatedAt: chirpy.UpdatedAt.Time,
		Body:      chirpy.Body,
		UserID:    chirpy.UserID,
	})
}

func (cfg *apiConfig) handlerListChirpies(w http.ResponseWriter, r *http.Request) {
	chirpies, err := cfg.db.ListChirpies(r.Context())
	if err != nil {
		ResponseWithError(w, http.StatusInternalServerError, "Error listing all chirpies!")
		return
	}
	res := []Chirpy{}
	for _, chirpy := range chirpies {
		res = append(res, Chirpy{
			ID:        chirpy.ID,
			CreatedAt: chirpy.CreatedAt.Time,
			UpdatedAt: chirpy.UpdatedAt.Time,
			Body:      chirpy.Body,
			UserID:    chirpy.UserID,
		})
	}

	RespondWithJson(w, http.StatusOK, res)
}

func (cfg *apiConfig) handlerGetChirpyByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	if id == "" {
		ResponseWithError(w, http.StatusInternalServerError, "Path value is empty!")
		return
	}
	chirpyID, err := uuid.Parse(id)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirpy, err := cfg.db.GetChirpyByID(r.Context(), chirpyID)
	if err != nil {
		ResponseWithError(w, http.StatusNotFound, fmt.Sprintf("Error getting ChirpyID %s, it does not exist in the databse!", chirpyID))
		return
	}

	RespondWithJson(w, http.StatusOK, Chirpy{
		ID:        chirpy.ID,
		CreatedAt: chirpy.CreatedAt.Time,
		UpdatedAt: chirpy.UpdatedAt.Time,
		Body:      chirpy.Body,
		UserID:    chirpy.UserID,
	})
}
