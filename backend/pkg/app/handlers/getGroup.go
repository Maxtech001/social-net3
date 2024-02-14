package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
	"social-network/pkg/models"
)

func GetGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.Group
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

	group, err := dbfunctions.GetGroupDataById(req.Id)
	if err != nil {
		http.Error(w, "Getting group data failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	group.IsMember, err = dbfunctions.ExistsGroupMemberByPK(models.GroupMemberConnection{GroupId: req.Id, MemberId: user.Id})
	if err != nil {
		http.Error(w, "Getting group request failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	group.IsRequested, err = dbfunctions.ExistsGroupRequestByPK(models.GroupRequest{GroupId: req.Id, RequesterId: user.Id})
	if err != nil {
		http.Error(w, "Getting group request failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res, _ := json.Marshal(group)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
