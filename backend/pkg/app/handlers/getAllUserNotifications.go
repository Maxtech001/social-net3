package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
)

func GetAllUserNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	user := AuthenticationService.GetUserFromCookie(r)

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	notifications, err := dbfunctions.GetAllNotificationsByUserId(user.Id)
	if err != nil {
		http.Error(w, "Getting all posts failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res, _ := json.Marshal(notifications)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
