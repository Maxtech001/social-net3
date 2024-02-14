package dbfunctions

import (
	"fmt"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/models"
)

func InsertGroupInvite(groupInvite models.GroupInvite) error {
	_, err := sqlite.Db.Exec("INSERT INTO group_invites (GROUPID, INVITEDID, INVITERID) VALUES (?, ?, ?)", groupInvite.GroupId, groupInvite.InvitedId, groupInvite.InviterId)
	if err != nil {
		fmt.Println("Inserting group invite failed: " + err.Error())
	}
	return err
}

func RemoveGroupInvite(groupInvite models.GroupInvite) error {
	_, err := sqlite.Db.Exec("DELETE FROM group_invites WHERE GROUPID = ? AND INVITEDID = ? AND INVITERID = ?; DELETE FROM group_requests WHERE GROUPID = ? AND REQUESTERID = ?", groupInvite.GroupId, groupInvite.InvitedId, groupInvite.InviterId, groupInvite.GroupId, groupInvite.InvitedId)
	if err != nil {
		fmt.Println("Failed to remove group invite: " + err.Error())
	}
	return err
}

func GetGroupInvitesByUserId(userId int) ([]*models.GroupInvite, error) {
	groupInvites := []*models.GroupInvite{}
	rows, err := sqlite.Db.Query("SELECT g.ID, g.ADMINID, g.GROUPNAME, g.GROUPDESCRIPTION, g.MEMBERCOUNT, g.CREATEDAT, u.ID, u.EMAIL, u.UPASSWORD, u.FIRSTNAME, u.LASTNAME, u.BIRTHDATE, u.AVATARPATH, u.NICKNAME, u.ABOUTME, u.FOLLOWERS, u.ISPUBLIC, ui.ID, ui.EMAIL, ui.UPASSWORD, ui.FIRSTNAME, ui.LASTNAME, ui.BIRTHDATE, ui.AVATARPATH, ui.NICKNAME, ui.ABOUTME, ui.FOLLOWERS, ui.ISPUBLIC FROM groups as g INNER JOIN group_invites as gi ON g.ID = gi.GROUPID INNER JOIN users as u ON u.ID = gi.INVITEDID INNER JOIN users as ui ON ui.ID = gi.INVITERID WHERE gi.INVITEDID = ? ", userId)
	if err != nil {
		fmt.Println("Getting user group invites failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var groupInvite models.GroupInvite
		if err := rows.Scan(&groupInvite.Group.Id, &groupInvite.Group.AdminId, &groupInvite.Group.GroupName, &groupInvite.Group.Description, &groupInvite.Group.MemberCount, &groupInvite.Group.CreatedAt, &groupInvite.InvitedUser.Id, &groupInvite.InvitedUser.Email, &groupInvite.InvitedUser.Password, &groupInvite.InvitedUser.FirstName, &groupInvite.InvitedUser.LastName, &groupInvite.InvitedUser.BirthDate, &groupInvite.InvitedUser.AvatarPath, &groupInvite.InvitedUser.Nickname, &groupInvite.InvitedUser.AboutMe, &groupInvite.InvitedUser.Followers, &groupInvite.InvitedUser.IsPublic, &groupInvite.InviterUser.Id, &groupInvite.InviterUser.Email, &groupInvite.InviterUser.Password, &groupInvite.InviterUser.FirstName, &groupInvite.InviterUser.LastName, &groupInvite.InviterUser.BirthDate, &groupInvite.InviterUser.AvatarPath, &groupInvite.InviterUser.Nickname, &groupInvite.InviterUser.AboutMe, &groupInvite.InviterUser.Followers, &groupInvite.InviterUser.IsPublic); err != nil {
			fmt.Println("Scanning user group invites failed: " + err.Error())
			return nil, err
		}
		groupInvites = append(groupInvites, &groupInvite)
	}
	return groupInvites, nil
}
