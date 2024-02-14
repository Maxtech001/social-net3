package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
)

func GetUserConversationsHandler(w http.ResponseWriter, r *http.Request) {
	user := AuthenticationService.GetUserFromCookie(r)

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	conversations, err := dbfunctions.GetConversationsByUserID(user.Id)
	if err != nil {
		http.Error(w, "Getting user conversations failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res, _ := json.Marshal(conversations)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
