package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/justinyeh1995/go-web-server-template/internal/database"
)

func (cfg *apiConfig) handlerCreateChirpy(w http.ResponseWriter, r *http.Request) {
	type param struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}

	type Chirpy struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
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
		Body: sql.NullString{
			String: params.Body,
			Valid:  true,
		},
		UserID: sql.NullString{
			String: params.UserID,
			Valid:  true,
		},
	})
	if err != nil {
		// log.Fatalf("Error creating chirpy, with err msg %v", err)
		ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating chirpy, with err msg %v", err))
		return
	}

	chirpyID, err := uuid.Parse(chirpy.ID)
	if err != nil {
		// log.Fatalf("Error parsing chirpy id, with err msg %v", err)
		ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error parsing chirpy id, with err msg %v", err))
		return
	}
	userID, err := uuid.Parse(chirpy.UserID.String)
	if err != nil {
		// log.Fatalf("Error parsing chirpy's user id, with err msg %v", err)
		ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error parsing chirpy's user id, with err msg %v", err))
		return
	}
	// 3) Save it to the database
	RespondWithJson(w, http.StatusCreated, Chirpy{
		ID:        chirpyID,
		CreatedAt: chirpy.CreatedAt.Time,
		UpdatedAt: chirpy.UpdatedAt.Time,
		Body:      chirpy.Body.String,
		UserID:    userID,
	})
}
