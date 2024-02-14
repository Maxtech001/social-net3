package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
	"social-network/pkg/models"
)

func GetUserPostsHandler(w http.ResponseWriter, r *http.Request) {
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

	if user.Id != existing.Id && !isExisting && !existing.IsPublic {
		res, _ := json.Marshal([]models.Post{})
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	}

	posts, err := dbfunctions.GetUserPosts(existing.Id, user.Id)
	if err != nil {
		http.Error(w, "Getting user posts failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	res, _ := json.Marshal(posts)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
