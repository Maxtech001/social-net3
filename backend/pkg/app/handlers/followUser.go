package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
	"social-network/pkg/models"
	"time"
)

const (
	FollowRequestNotificationContent = "You have a new follow request from "
)

func FollowUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.FollowConnection
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	follower := AuthenticationService.GetUserFromCookie(r)

	if follower == nil || req.FollowerId != follower.Id {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	followedUser, err := dbfunctions.GetAccountDataById(req.FollowedId)
	if err != nil || followedUser == nil {
		http.Error(w, "Failed to get followed user data: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if req.FollowerId == req.FollowedId {
		http.Error(w, "Account id's cannot match", http.StatusUnprocessableEntity)
		return
	}

	if followedUser.IsPublic {
		followedUser, err = dbfunctions.InsertFollowConnection(req)
		if err != nil {
			http.Error(w, "Failed to follow user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		isExisting, err := dbfunctions.IsExistingFollowConnection(models.FollowConnection{FollowedId: followedUser.Id, FollowerId: follower.Id})
		if err != nil {
			http.Error(w, "Failed to get existing follow connection: "+err.Error(), http.StatusInternalServerError)
			return
		}

		followedUser.IsFollowed = isExisting
	} else {
		followedUser, err = dbfunctions.InsertFollowRequest(req)
		if err != nil {
			http.Error(w, "Failed to request to follow user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = dbfunctions.InsertFollowNotification(
			models.Notification{
				UserId:           followedUser.Id,
				Content:          FollowRequestNotificationContent + follower.FirstName + " " + follower.LastName,
				NotificationType: 0,
				IsRead:           false,
				CreatedAt:        time.Now().Format("2006-01-02 15:04:05"),
				FollowerId:       follower.Id,
			})
		if err != nil {
			http.Error(w, "Insert notification failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		isExisting, err := dbfunctions.IsExistingFollowRequest(models.FollowConnection{FollowedId: followedUser.Id, FollowerId: follower.Id})
		if err != nil {
			http.Error(w, "Failed to get existing follow request: "+err.Error(), http.StatusInternalServerError)
			return
		}

		followedUser.IsRequested = isExisting
	}

	res, _ := json.Marshal(followedUser)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
