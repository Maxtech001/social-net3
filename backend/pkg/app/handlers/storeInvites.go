package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
	"social-network/pkg/models"
	"time"
)

const GroupInviteNotificationContent = "You've been invited to join the group "

func StoreInvitesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	type Request struct {
		Users   []models.User `json:"users"`
		GroupId int           `json:"groupId"`
	}

	var req Request
	var group *models.Group
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

	group, err = dbfunctions.GetGroupDataById(req.GroupId)
	if err != nil {
		http.Error(w, "Failed to get group data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, invitedUser := range req.Users {
		if err := dbfunctions.InsertGroupInvite(models.GroupInvite{GroupId: req.GroupId, InvitedId: invitedUser.Id, InviterId: user.Id}); err != nil {
			http.Error(w, "Failed to insert group invite: "+err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = dbfunctions.InsertGroupNotification(
			models.Notification{
				UserId:           invitedUser.Id,
				Content:          GroupInviteNotificationContent + group.GroupName,
				NotificationType: 1,
				IsRead:           false,
				CreatedAt:        time.Now().Format("2006-01-02 15:04:05"),
				GroupId:          group.Id,
			})
		if err != nil {
			http.Error(w, "Insert notification failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
