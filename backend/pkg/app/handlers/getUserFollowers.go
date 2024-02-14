package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
	"social-network/pkg/models"
)

func GetUserFollowersHandler(w http.ResponseWriter, r *http.Request) {
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

	followers := []*models.User{}

	isAllowed, err := isAllowedToView(user.Id, req.Id)
	if err != nil {
		http.Error(w, "Failed to get existing follow connection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if isAllowed {
		followers, err = dbfunctions.GetFollowersByFollowedId(req.Id)
		if err != nil {
			http.Error(w, "Querying followers failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	res, _ := json.Marshal(followers)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func isAllowedToView(userId, reqId int) (bool, error) {
	if userId == reqId {
		return true, nil
	}

	reqUser, err := dbfunctions.GetAccountDataById(reqId)
	if err != nil {
		return false, err
	}

	if reqUser.IsPublic {
		return true, nil
	}

	isExisting, err := dbfunctions.IsExistingFollowConnection(models.FollowConnection{FollowedId: reqId, FollowerId: userId})
	if err != nil {
		return false, err
	}

	return isExisting, nil
}
