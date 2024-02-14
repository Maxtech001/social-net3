package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
	"social-network/pkg/models"
	"time"
)

func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
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

	req.AdminId = user.Id
	req.MemberCount = 1
	req.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	group, err := dbfunctions.InsertGroupData(req)
	if err != nil {
		http.Error(w, "Inserting group failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := dbfunctions.InsertGroupMember(models.GroupMemberConnection{GroupId: group.Id, MemberId: group.AdminId}); err != nil {
		http.Error(w, "Inserting group member failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res, _ := json.Marshal(group)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
