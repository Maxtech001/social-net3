package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
	"social-network/pkg/models"
)

func GetUserFollowingHandler(w http.ResponseWriter, r *http.Request) {
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

	user := AuthenticationService.GetUserFromCookie(r)

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	following := []*models.User{}

	isAllowed, err := isAllowedToView(user.Id, req.Id)
	if err != nil {
		http.Error(w, "Failed to get existing follow connection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if isAllowed {
		following, err = dbfunctions.GetFollowedusersByFollowerId(req.Id)
		if err != nil {
			http.Error(w, "Querying followers failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	res, _ := json.Marshal(following)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
