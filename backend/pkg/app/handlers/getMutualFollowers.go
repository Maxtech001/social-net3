package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
)

func GetMutualFollowersHandler(w http.ResponseWriter, r *http.Request) {
	user := AuthenticationService.GetUserFromCookie(r)

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	users, err := dbfunctions.GetMutualFollowers(user.Id)
	if err != nil {
		http.Error(w, "Getting all mutual followers failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res, _ := json.Marshal(users)
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}
