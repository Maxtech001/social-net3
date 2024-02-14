package handlers

import (
	"encoding/json"
	"net/http"
	"social-network/pkg/app/AuthenticationService"
	"social-network/pkg/app/dbfunctions"
	"social-network/pkg/models"
	"sort"
)

func GetContactsHandler(w http.ResponseWriter, r *http.Request) {
	user := AuthenticationService.GetUserFromCookie(r)

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	allContacts := []models.Contact{}

	groups, err := dbfunctions.GetGroupsByGroupMemberId(user.Id)
	if err != nil {
		http.Error(w, "Getting groups failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	for _, group := range groups {
		allContacts = append(allContacts, models.Contact{Group: group})
	}

	userContacts, err := dbfunctions.GetFollowedAndFollowerUsersByUserId(user.Id)
	for _, user := range userContacts {
		allContacts = append(allContacts, models.Contact{User: user})
	}

	sort.Slice(allContacts, func(i, j int) bool {
		firstName := ""
		secondName := ""

		if allContacts[i].User != nil {
			firstName = allContacts[i].User.FirstName + " " + allContacts[i].User.LastName
		} else {
			firstName = allContacts[i].Group.GroupName
		}

		if allContacts[j].User != nil {
			secondName = allContacts[j].User.FirstName + " " + allContacts[j].User.LastName
		} else {
			secondName = allContacts[j].Group.GroupName
		}
		return firstName < secondName
	})

	res, _ := json.Marshal(allContacts)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
