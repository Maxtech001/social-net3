package dbfunctions

import (
	"database/sql"
	"fmt"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/models"
)

func InsertGroupData(group models.Group) (*models.Group, error) {
	var res models.Group
	err := sqlite.Db.QueryRow("INSERT INTO groups (ADMINID, GROUPNAME, GROUPDESCRIPTION, MEMBERCOUNT, CREATEDAT) VALUES (?, ?, ?, ?, ?) RETURNING ID, ADMINID, GROUPNAME, GROUPDESCRIPTION, MEMBERCOUNT, CREATEDAT", group.AdminId, group.GroupName, group.Description, group.MemberCount, group.CreatedAt).Scan(&res.Id, &res.AdminId, &res.GroupName, &res.Description, &res.MemberCount, &res.CreatedAt)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Insert group data failed: " + err.Error())
			return nil, err
		}
		return nil, nil
	}
	return &res, nil
}

func GetGroupDataById(groupId int) (*models.Group, error) {
	var group models.Group
	err := sqlite.Db.QueryRow("SELECT groups.ID, ADMINID, GROUPNAME, GROUPDESCRIPTION, MEMBERCOUNT, CREATEDAT, users.ID, users.EMAIL, users.FIRSTNAME, users.LASTNAME, users.BIRTHDATE, users.AVATARPATH, users.NICKNAME, users.ABOUTME, users.FOLLOWERS, users.ISPUBLIC FROM groups INNER JOIN users ON users.ID = groups.ADMINID WHERE groups.ID = ?", groupId).Scan(&group.Id, &group.AdminId, &group.GroupName, &group.Description, &group.MemberCount, &group.CreatedAt, &group.Admin.Id, &group.Admin.Email, &group.Admin.FirstName, &group.Admin.LastName, &group.Admin.BirthDate, &group.Admin.AvatarPath, &group.Admin.Nickname, &group.Admin.AboutMe, &group.Admin.Followers, &group.Admin.IsPublic)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Getting group data failed: " + err.Error())
			return nil, err
		}
		return nil, nil
	}
	if group.Members, err = GetGroupMembersByGroupId(groupId); err != nil {
		fmt.Println("Getting group members failed: " + err.Error())
		return nil, err
	}
	return &group, nil
}

func GetGroupsByGroupName(groupName string) ([]*models.Group, error) {
	groups := []*models.Group{}
	rows, err := sqlite.Db.Query("SELECT ID, ADMINID, GROUPNAME, GROUPDESCRIPTION, MEMBERCOUNT, CREATEDAT FROM groups WHERE GROUPNAME LIKE '%" + groupName + "%'")
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

func UpdateGroupMemberCount(groupId, userId int) error {
	_, err := sqlite.Db.Exec("UPDATE groups SET MEMBERCOUNT = (SELECT COUNT(*) FROM group_members WHERE GROUPID = ?) WHERE ID = ?", groupId, groupId)
	if err != nil {
		fmt.Println("Updating group member count failed: " + err.Error())
	}
	return err
}
