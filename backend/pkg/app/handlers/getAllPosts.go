package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
)

func GetAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	user := AuthenticationService.GetUserFromCookie(r)

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	posts, err := dbfunctions.GetVisiblePosts(*user)
	if err != nil {
		http.Error(w, "Getting all posts failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	res, _ := json.Marshal(posts)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
