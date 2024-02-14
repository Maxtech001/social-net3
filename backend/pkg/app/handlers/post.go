package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
	"social-network/pkg/models"
	"strconv"
	"strings"
	"time"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.Post
	req.Content = r.PostFormValue("content")
	req.GroupId, _ = strconv.Atoi(r.PostFormValue("groupId"))
	req.PostType, _ = strconv.Atoi(r.PostFormValue("privacyType"))
	if r.PostFormValue("specifiedUsers") != "" {
		err := json.Unmarshal([]byte(r.PostFormValue("specifiedUsers")), &req.SpecifiedUsers)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	user := AuthenticationService.GetUserFromCookie(r)

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	req.AuthorId = user.Id
	req.TotalComments = 0
	req.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	post, err := dbfunctions.InsertPostData(req)
	if err != nil {
		http.Error(w, "Insert post data failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, u := range req.SpecifiedUsers {
		if _, err := dbfunctions.InsertPostUserConnection(models.PostUser{PostId: post.Id, UserId: u.Id}); err != nil {
			http.Error(w, "Error inserting post user connection: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	for _, fs := range r.MultipartForm.File {
		for _, f := range fs {
			file, _ := f.Open()
			fileNameSplit := strings.Split(f.Filename, ".")
			filePath := fmt.Sprintf("/static/%d_%d.%s", user.Id, time.Now().UnixNano(), fileNameSplit[len(fileNameSplit)-1])
			dst, err := os.Create("./pkg" + filePath)
			if err != nil {
				fmt.Println("Failed to create file: " + err.Error())
				http.Error(w, "Error storing file: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer dst.Close()

			if _, err := io.Copy(dst, file); err != nil {
				fmt.Println("Failed to copy file: " + err.Error())
				http.Error(w, "Error storing file: "+err.Error(), http.StatusInternalServerError)
				return
			}
			file.Close()

			if postImg, err := dbfunctions.InsertPostImageConnection(models.PostImage{PostId: post.Id, ImagePath: filePath}); err != nil {
				http.Error(w, "Error inserting post image connection: "+err.Error(), http.StatusInternalServerError)
				return
			} else {
				post.ImageArray = append(post.ImageArray, postImg)
			}
		}
	}

	res, _ := json.Marshal(post)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
