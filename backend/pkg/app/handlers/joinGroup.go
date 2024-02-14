package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
	"social-network/pkg/models"
	"time"
)

const GroupRequestNotificationContent = " has requested to join your group "

func JoinGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.Group
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

	group, err = dbfunctions.GetGroupDataById(req.Id)
	if err != nil {
		http.Error(w, "Failed to get group data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := dbfunctions.InsertGroupRequest(req.Id, user.Id); err != nil {
		http.Error(w, "Failed to request to joing group: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = dbfunctions.InsertGroupNotification(
		models.Notification{
			UserId:           group.AdminId,
			Content:          user.FirstName + " " + user.LastName + GroupRequestNotificationContent + group.GroupName,
			NotificationType: 2,
			IsRead:           false,
			CreatedAt:        time.Now().Format("2006-01-02 15:04:05"),
			GroupId:          group.Id,
		})
	if err != nil {
		http.Error(w, "Insert notification failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
