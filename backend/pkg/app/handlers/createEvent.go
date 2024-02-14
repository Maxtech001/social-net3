package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
	"social-network/pkg/models"
	"time"
)

const EventNotificationContent = "You've been invited to the event "

func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
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

	req.AuthorId = user.Id
	req.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	var group *models.Group
	event, err := dbfunctions.InsertEventData(req)
	if err != nil {
		http.Error(w, "Inserting event failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	group, err = dbfunctions.GetGroupDataById(req.GroupId)
	if err != nil {
		http.Error(w, "Failed to get group data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var groupMembers []*models.User

	if groupMembers, err = dbfunctions.GetGroupMembersByGroupId(req.GroupId); err != nil {
		fmt.Println("Getting group members failed: " + err.Error())
		return
	}

	for _, member := range groupMembers {
		if err = dbfunctions.InsertEventInvitedMember(models.InvitedMember{EventId: event.Id, MemberId: member.Id, Answer: 0}); err != nil {
			fmt.Println("Inserting event invited members failed: " + err.Error())
			return
		}

		if member.Id != user.Id {
			_, err = dbfunctions.InsertEventNotification(
				models.Notification{
					UserId:           member.Id,
					Content:          EventNotificationContent + event.Title + " in the group " + group.GroupName,
					NotificationType: 3,
					IsRead:           false,
					CreatedAt:        time.Now().Format("2006-01-02 15:04:05"),
					EventId:          event.Id,
				})
			if err != nil {
				http.Error(w, "Insert notification failed: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	res, _ := json.Marshal(event)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
