package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
	"social-network/pkg/models"
)

func GetUsersByNameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.NameSearch == "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := AuthenticationService.GetUserFromCookie(r)

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	users, err := dbfunctions.GetUsersByName(req.NameSearch)
	if err != nil {
		http.Error(w, "Getting user groups failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res, _ := json.Marshal(users)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
