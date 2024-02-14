package dbfunctions

import (
	"database/sql"
	"fmt"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/models"
)

func GetGroupMembersByGroupId(groupId int) ([]*models.User, error) {
	members := []*models.User{}
	rows, err := sqlite.Db.Query("SELECT users.ID, users.EMAIL, users.FIRSTNAME, users.LASTNAME, users.BIRTHDATE, users.AVATARPATH, users.NICKNAME, users.ABOUTME, users.FOLLOWERS, users.ISPUBLIC FROM group_members INNER JOIN users ON users.ID = group_members.MEMBERID WHERE GROUPID = ?", groupId)
	if err != nil {
		fmt.Println("Getting post comments failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var member models.User
		if err := rows.Scan(&member.Id, &member.Email, &member.FirstName, &member.LastName, &member.BirthDate, &member.AvatarPath, &member.Nickname, &member.AboutMe, &member.Followers, &member.IsPublic); err != nil {
			fmt.Println("Scanning post comments failed: " + err.Error())
			return nil, err
		}
		members = append(members, &member)
	}
	return members, nil
}

func InsertGroupMember(conn models.GroupMemberConnection) error {
	_, err := sqlite.Db.Exec("INSERT INTO group_members (GROUPID, MEMBERID) VALUES (?, ?)", conn.GroupId, conn.MemberId)
	if err != nil {
		fmt.Println("Inserting group member failed: " + err.Error())
	}
	if err = UpdateGroupMemberCount(conn.GroupId, conn.MemberId); err != nil {
		return err
	}
	return err
}

func GetGroupsByGroupMemberId(memberId int) ([]*models.Group, error) {
	groups := []*models.Group{}
	rows, err := sqlite.Db.Query("SELECT groups.ID, ADMINID, GROUPNAME, GROUPDESCRIPTION, MEMBERCOUNT, CREATEDAT FROM group_members INNER JOIN groups ON groups.ID = GROUPID WHERE MEMBERID = ?", memberId)
	if err != nil {
		fmt.Println("Getting user groups failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.Id, &group.AdminId, &group.GroupName, &group.Description, &group.MemberCount, &group.CreatedAt); err != nil {
			fmt.Println("Scanning user groups failed: " + err.Error())
			return nil, err
		}
		groups = append(groups, &group)
	}
	return groups, nil
}

func RemoveGroupMember(groupId, userId int) error {
	_, err := sqlite.Db.Exec("DELETE FROM group_members WHERE GROUPID = ? AND MEMBERID = ? ", groupId, userId)
	if err != nil {
		fmt.Println("Failed to remove group member: " + err.Error())
	}
	if err = UpdateGroupMemberCount(groupId, userId); err != nil {
		return err
	}
	return err
}

func ExistsGroupMemberByPK(groupMember models.GroupMemberConnection) (bool, error) {
	var res bool
	err := sqlite.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM group_members WHERE GROUPID = ? AND MEMBERID = ?)", groupMember.GroupId, groupMember.MemberId).Scan(&res)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Querying is exisiting group request failed: " + err.Error())
			return false, err
		}
		return false, nil
	}
	return res, err
}
