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

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.Comment
	req.Content = r.PostFormValue("content")
	req.PostId, _ = strconv.Atoi(r.PostFormValue("postId"))

	user := AuthenticationService.GetUserFromCookie(r)

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	req.AuthorId = user.Id
	req.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	comment, err := dbfunctions.InsertCommentData(req)
	if err != nil {
		http.Error(w, "Insert comment data failed: "+err.Error(), http.StatusInternalServerError)
		return
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

			if postImg, err := dbfunctions.InsertCommentImageConnection(models.CommentImage{CommentId: comment.Id, ImagePath: filePath}); err != nil {
				http.Error(w, "Error storing file: "+err.Error(), http.StatusInternalServerError)
				return
			} else {
				comment.ImageArray = append(comment.ImageArray, postImg)
			}
		}
	}

	res, _ := json.Marshal(comment)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
