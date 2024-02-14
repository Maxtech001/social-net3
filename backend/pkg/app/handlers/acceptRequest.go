package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
	"social-network/pkg/models"
)

func AcceptRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.FollowConnection
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

	if err := dbfunctions.RemoveFollowRequest(req); err != nil {
		http.Error(w, "Deleting follow request failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	followedUser, err := dbfunctions.InsertFollowConnection(req)
	if err != nil {
		http.Error(w, "Failed to follow user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	isExisting, err := dbfunctions.IsExistingFollowConnection(models.FollowConnection{FollowedId: followedUser.Id, FollowerId: user.Id})
	if err != nil {
		http.Error(w, "Failed to get existing follow connection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	followedUser.IsFollowed = isExisting

	res, _ := json.Marshal(followedUser)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
