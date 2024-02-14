package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
	"social-network/pkg/models"
)

func GetConversationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	type Request struct {
		UserId  int `json:"userId"`
		GroupId int `json:"groupId"`
	}

	var req Request
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

	var conversation *models.Conversation
	messages := []*models.Message{}

	if req.GroupId != 0 {
		conversation, err = dbfunctions.GetConversationByGroupId(req.GroupId)
		if err != nil {
			http.Error(w, "Getting conversation failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		conversation, err = dbfunctions.GetConversationBySenderReceiverIDs(user.Id, req.UserId)
		if err != nil {
			http.Error(w, "Getting conversation failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if conversation != nil {
		messages, err = dbfunctions.GetMessageByConversation(conversation.Id, user.Id)
		if err != nil {
			http.Error(w, "Getting conversation history failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	res, _ := json.Marshal(messages)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
