package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type param struct {
		Email string `json:"email"`
	}
	type User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := param{}
	err := decoder.Decode(&params)
	if err != nil {
		ResponseWithError(w, 500, "Error Decoding Request Body.")
		return
	}

	email := params.Email
	user, err := cfg.db.CreateUser(r.Context(), sql.NullString{
		String: params.Email,
		Valid:  true,
	})
	if err != nil {
		ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error Creating the User with email: %s. Error: %v", email, err))
		return
	}
	id, err := uuid.Parse(user.ID.String)
	if err != nil {
		ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error Parsing Users uuid: %v.", err))
		return
	}
	RespondWithJson(w, http.StatusCreated, User{
		ID:        id,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Email:     user.Email.String,
	})
}
