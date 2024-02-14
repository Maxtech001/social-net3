package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
	"social-network/pkg/models"
)

func UnfollowUserHandler(w http.ResponseWriter, r *http.Request) {
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

	follower := AuthenticationService.GetUserFromCookie(r)

	if follower == nil || req.FollowerId != follower.Id {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	unFollowedUser, err := dbfunctions.GetAccountDataById(req.FollowedId)
	if err != nil || unFollowedUser == nil {
		http.Error(w, "Failed to get followed user data", http.StatusUnprocessableEntity)
		return
	}

	if req.FollowerId == req.FollowedId {
		http.Error(w, "Account id's cannot match", http.StatusUnprocessableEntity)
		return
	}

	unFollowedUser, err = dbfunctions.RemoveFollowConnection(req)
	if err != nil {
		http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		return
	}

	isExisting, err := dbfunctions.IsExistingFollowConnection(models.FollowConnection{FollowedId: unFollowedUser.Id, FollowerId: follower.Id})
	if err != nil {
		http.Error(w, "Failed to get existing follow connection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	unFollowedUser.IsFollowed = isExisting

	res, _ := json.Marshal(unFollowedUser)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
