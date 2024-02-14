package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
)

func GetgroupInvitesHandler(w http.ResponseWriter, r *http.Request) {
	user := AuthenticationService.GetUserFromCookie(r)

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	groupInvites, err := dbfunctions.GetGroupInvitesByUserId(user.Id)
	if err != nil {
		http.Error(w, "Getting user group invites failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res, _ := json.Marshal(groupInvites)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
