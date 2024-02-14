package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
)

func GetUserGroupsHandler(w http.ResponseWriter, r *http.Request) {
	user := AuthenticationService.GetUserFromCookie(r)

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	groups, err := dbfunctions.GetGroupsByGroupMemberId(user.Id)
	if err != nil {
		http.Error(w, "Getting user groups failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res, _ := json.Marshal(groups)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
