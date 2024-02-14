package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
	"social-network/pkg/models"
)

func AcceptGroupInviteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.GroupInvite
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

	if err := dbfunctions.RemoveGroupInvite(req); err != nil {
		http.Error(w, "Deleting group invite failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = dbfunctions.InsertGroupMember(models.GroupMemberConnection{GroupId: req.GroupId, MemberId: req.InvitedId})
	if err != nil {
		http.Error(w, "Failed to insert user to group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
