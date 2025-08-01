package main

import "net/http"

func (cfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	if cfg.env != "dev" {
		ResponseWithError(w, http.StatusForbidden, "It is not allowed to perform reset except for the local dev environment")
	}
	cfg.fileserverHits.Store(0)

	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		ResponseWithError(w, 500, "An error haapened when deleting users..")
	}
	RespondWithJson(w, http.StatusOK, regularReponse{
		Valid: true,
	})
}
