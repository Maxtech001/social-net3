package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
)

func GetFollowersHandler(w http.ResponseWriter, r *http.Request) {
	user := AuthenticationService.GetUserFromCookie(r)

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	followers, err := dbfunctions.GetFollowersByFollowedId(user.Id)
	if err != nil {
		http.Error(w, "Querying followers failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res, _ := json.Marshal(followers)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
