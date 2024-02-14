package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
)

func GetGroupRequestsHandler(w http.ResponseWriter, r *http.Request) {
	user := AuthenticationService.GetUserFromCookie(r)

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	groupRequests, err := dbfunctions.GetGroupRequestsByUserId(user.Id)
	if err != nil {
		http.Error(w, "Getting user groups requests failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res, _ := json.Marshal(groupRequests)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
