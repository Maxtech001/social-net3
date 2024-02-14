package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
	"social-network/pkg/models"
)

func GetPastEventsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.Event
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

	events, err := dbfunctions.GetPastEventsByGroupId(req.GroupId)
	if err != nil {
		http.Error(w, "Querying past events failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res, _ := json.Marshal(events)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
