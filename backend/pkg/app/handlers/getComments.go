package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
	"social-network/pkg/models"
)

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.Comment
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

	comments, err := dbfunctions.GetCommentsByPostId(req.PostId)
	if err != nil {
		http.Error(w, "Getting post comments failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res, _ := json.Marshal(comments)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
