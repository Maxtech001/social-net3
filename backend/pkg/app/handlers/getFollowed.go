package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
)

func GetFollowedHandler(w http.ResponseWriter, r *http.Request) {
	user := AuthenticationService.GetUserFromCookie(r)

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	followers, err := dbfunctions.GetFollowedusersByFollowerId(user.Id)
	if err != nil {
		http.Error(w, "Querying followed users failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res, _ := json.Marshal(followers)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
