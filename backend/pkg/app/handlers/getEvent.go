package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
	"social-network/pkg/models"
)

func GetEventHandler(w http.ResponseWriter, r *http.Request) {
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

	event, err := dbfunctions.GetEventById(req.Id)
	if err != nil {
		http.Error(w, "Querying event data failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if event != nil {
		if event.InvitedMembers, err = dbfunctions.GetEventInvitedMembersByEventId(event.Id); err != nil {
			http.Error(w, "Querying invited members failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if event.GoingMembers, err = dbfunctions.GetGoingMembersByEventId(event.Id); err != nil {
			http.Error(w, "Querying going members failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if event.NotGoingMembers, err = dbfunctions.GetNotGoingMembersByEventId(event.Id); err != nil {
			http.Error(w, "Querying not going members failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if event.UserAnswer, err = dbfunctions.GetUserEventAnswer(models.InvitedMember{EventId: event.Id, MemberId: user.Id}); err != nil {
			http.Error(w, "Querying is user going failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	res, _ := json.Marshal(event)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
