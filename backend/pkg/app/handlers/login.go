package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
	"social-network/pkg/models"
)

type Error struct {
	errorMessage string
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existing, err := dbfunctions.GetAccountDataByEmail(req.Email)

	if err != nil {
		http.Error(w, "Failed to login: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if existing == nil || req.Password != existing.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	cookie, err := AuthenticationService.GenerateCookies(*existing)
	if err != nil {
		http.Error(w, "Failed to create session: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}
