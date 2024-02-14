package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
	"social-network/pkg/models"
	"strings"
	"time"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 5*(1<<20))

	var req models.User
	req.Email = r.PostFormValue("email")
	req.Password = r.PostFormValue("password")
	req.FirstName = r.PostFormValue("firstName")
	req.LastName = r.PostFormValue("lastName")
	req.BirthDate = r.PostFormValue("birthDate")
	req.Nickname = r.PostFormValue("nickname")
	req.AboutMe = r.PostFormValue("aboutMe")
	req.Followers = 0
	req.IsPublic = false

	file, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		if err != http.ErrMissingFile {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		file = nil
	}

	if err := validateRegisterRequest(req); err != nil {
		http.Error(w, err.errorMessage, http.StatusUnprocessableEntity)
		return
	}

	if username, err := dbfunctions.GetAccountDataByEmail(req.Email); err != nil {
		http.Error(w, "Failed to register: "+err.Error(), http.StatusInternalServerError)
		return
	} else if username != nil {
		http.Error(w, "Email is already in use", http.StatusUnprocessableEntity)
		return
	}

	user, err := dbfunctions.InsertUserData(req)
	if err != nil {
		http.Error(w, "Failed to register: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if file != nil {
		fileNameSplit := strings.Split(fileHeader.Filename, ".")
		user.AvatarPath = fmt.Sprintf("/static/%d_%d.%s", user.Id, time.Now().UnixNano(), fileNameSplit[len(fileNameSplit)-1])
		dst, err := os.Create("./pkg" + user.AvatarPath)
		if err != nil {
			fmt.Println("Failed to create file: " + err.Error())
			dbfunctions.RemoveUserData(*user)
			http.Error(w, "Error storing file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			fmt.Println("Failed to copy file: " + err.Error())
			dbfunctions.RemoveUserData(*user)
			http.Error(w, "Error storing file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		file.Close()

		if err := dbfunctions.InsertAvatarPath(*user); err != nil {
			dbfunctions.RemoveUserData(*user)
			http.Error(w, "Error storing file: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	cookie, err := AuthenticationService.GenerateCookies(*user)
	if err != nil {
		http.Error(w, "Failed to create session: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}

func validateRegisterRequest(req models.User) *Error {
	if req.Email == "" {
		return &Error{"'email' missing from JSON input"}
	}
	if req.Password == "" {
		return &Error{"'password' missing from JSON input"}
	}
	if req.FirstName == "" {
		return &Error{"'firstName' missing from JSON input"}
	}
	if req.LastName == "" {
		return &Error{"'lastName' missing from JSON input"}
	}
	if req.BirthDate == "" {
		return &Error{"'birthDate' missing from JSON input"}
	}

	return nil
}
