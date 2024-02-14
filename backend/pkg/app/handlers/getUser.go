package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
	"social-network/pkg/models"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
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

	existing, err := dbfunctions.GetAccountDataById(req.Id)
	if err != nil {
		http.Error(w, "Failed to get user data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if existing == nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	user := AuthenticationService.GetUserFromCookie(r)

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	isExisting, err := dbfunctions.IsExistingFollowConnection(models.FollowConnection{FollowedId: existing.Id, FollowerId: user.Id})
	if err != nil {
		http.Error(w, "Failed to get existing follow connection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	existing.IsFollowed = isExisting

	isExisting, err = dbfunctions.IsExistingFollowRequest(models.FollowConnection{FollowedId: existing.Id, FollowerId: user.Id})
	if err != nil {
		http.Error(w, "Failed to get existing follow request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	existing.IsRequested = isExisting

	res, _ := json.Marshal(existing)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
