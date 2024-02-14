package dbfunctions

import (
	"database/sql"
	"fmt"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/models"
)

func InsertGroupRequest(groupId, userId int) error {
	_, err := sqlite.Db.Exec("INSERT INTO group_requests (GROUPID, REQUESTERID) VALUES (?, ?)", groupId, userId)
	if err != nil {
		fmt.Println("Inserting group request failed: " + err.Error())
	}
	return err
}

func RemoveGroupRequest(groupRequest models.GroupRequest) error {
	_, err := sqlite.Db.Exec("DELETE FROM group_requests WHERE GROUPID = ? AND REQUESTERID = ?; DELETE FROM group_invites WHERE GROUPID = ? AND INVITEDID = ?;", groupRequest.GroupId, groupRequest.RequesterId, groupRequest.GroupId, groupRequest.RequesterId)
	if err != nil {
		fmt.Println("Failed to remove group member: " + err.Error())
	}
	return err
}

func ExistsGroupRequestByPK(groupRequest models.GroupRequest) (bool, error) {
	var res bool
	err := sqlite.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM group_requests WHERE GROUPID = ? AND REQUESTERID = ?)", groupRequest.GroupId, groupRequest.RequesterId).Scan(&res)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Querying is exisiting group request failed: " + err.Error())
			return false, err
		}
		return false, nil
	}
	return res, err
}

func GetGroupRequestsByUserId(userId int) ([]*models.GroupRequest, error) {
	groupRequests := []*models.GroupRequest{}
	rows, err := sqlite.Db.Query("SELECT g.ID, g.ADMINID, g.GROUPNAME, g.GROUPDESCRIPTION, g.MEMBERCOUNT, g.CREATEDAT, u.ID, u.EMAIL, u.UPASSWORD, u.FIRSTNAME, u.LASTNAME, u.BIRTHDATE, u.AVATARPATH, u.NICKNAME, u.ABOUTME, u.FOLLOWERS, u.ISPUBLIC FROM groups as g INNER JOIN group_requests as gp ON g.ID = gp.GROUPID INNER JOIN users as u ON u.ID = gp.REQUESTERID WHERE g.ADMINID = ? ", userId)
	if err != nil {
		fmt.Println("Getting user group requests failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var groupRequest models.GroupRequest
		if err := rows.Scan(&groupRequest.Group.Id, &groupRequest.Group.AdminId, &groupRequest.Group.GroupName, &groupRequest.Group.Description, &groupRequest.Group.MemberCount, &groupRequest.Group.CreatedAt, &groupRequest.User.Id, &groupRequest.User.Email, &groupRequest.User.Password, &groupRequest.User.FirstName, &groupRequest.User.LastName, &groupRequest.User.BirthDate, &groupRequest.User.AvatarPath, &groupRequest.User.Nickname, &groupRequest.User.AboutMe, &groupRequest.User.Followers, &groupRequest.User.IsPublic); err != nil {
			fmt.Println("Scanning user group requests failed: " + err.Error())
			return nil, err
		}
		groupRequests = append(groupRequests, &groupRequest)
	}
	return groupRequests, nil
}
