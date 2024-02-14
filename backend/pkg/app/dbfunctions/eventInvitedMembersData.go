package dbfunctions

import (
	"database/sql"
	"fmt"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/models"
)

func InsertEventInvitedMember(invitedMember models.InvitedMember) error {
	_, err := sqlite.Db.Exec("INSERT INTO event_invited_members (EVENTID, USERID, ANSWERTYPE) VALUES (?, ?, ?)", invitedMember.EventId, invitedMember.MemberId, invitedMember.Answer)
	if err != nil {
		fmt.Println("Inserting event invited members failed: " + err.Error())
	}
	return err
}

func GetEventInvitedMembersByEventId(eventId int) ([]*models.InvitedMember, error) {
	members := []*models.InvitedMember{}
	rows, err := sqlite.Db.Query("SELECT EVENTID, USERID, ANSWERTYPE, users.ID, users.EMAIL, users.FIRSTNAME, users.LASTNAME, users.BIRTHDATE, users.AVATARPATH, users.NICKNAME, users.ABOUTME, users.FOLLOWERS, users.ISPUBLIC FROM event_invited_members INNER JOIN users ON users.ID = USERID WHERE EVENTID = ?", eventId)
	if err != nil {
		fmt.Println("Getting invited members failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var member models.InvitedMember
		if err := rows.Scan(&member.EventId, &member.MemberId, &member.Answer, &member.Member.Id, &member.Member.Email, &member.Member.FirstName, &member.Member.LastName, &member.Member.BirthDate, &member.Member.AvatarPath, &member.Member.Nickname, &member.Member.AboutMe, &member.Member.Followers, &member.Member.IsPublic); err != nil {
			fmt.Println("Scanning invited members failed: " + err.Error())
			return nil, err
		}
		members = append(members, &member)
	}
	return members, nil
}

func GetGoingMembersByEventId(eventId int) ([]*models.InvitedMember, error) {
	members := []*models.InvitedMember{}
	rows, err := sqlite.Db.Query("SELECT EVENTID, USERID, ANSWERTYPE, users.ID, users.EMAIL, users.FIRSTNAME, users.LASTNAME, users.BIRTHDATE, users.AVATARPATH, users.NICKNAME, users.ABOUTME, users.FOLLOWERS, users.ISPUBLIC FROM event_invited_members INNER JOIN users ON users.ID = USERID WHERE EVENTID = ? AND ANSWERTYPE = 2", eventId)
	if err != nil {
		fmt.Println("Getting going members failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var member models.InvitedMember
		if err := rows.Scan(&member.EventId, &member.MemberId, &member.Answer, &member.Member.Id, &member.Member.Email, &member.Member.FirstName, &member.Member.LastName, &member.Member.BirthDate, &member.Member.AvatarPath, &member.Member.Nickname, &member.Member.AboutMe, &member.Member.Followers, &member.Member.IsPublic); err != nil {
			fmt.Println("Scanning going members failed: " + err.Error())
			return nil, err
		}
		members = append(members, &member)
	}
	return members, nil
}

func GetNotGoingMembersByEventId(eventId int) ([]*models.InvitedMember, error) {
	members := []*models.InvitedMember{}
	rows, err := sqlite.Db.Query("SELECT EVENTID, USERID, ANSWERTYPE, users.ID, users.EMAIL, users.FIRSTNAME, users.LASTNAME, users.BIRTHDATE, users.AVATARPATH, users.NICKNAME, users.ABOUTME, users.FOLLOWERS, users.ISPUBLIC FROM event_invited_members INNER JOIN users ON users.ID = USERID WHERE EVENTID = ? AND ANSWERTYPE = 1", eventId)
	if err != nil {
		fmt.Println("Getting going members failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var member models.InvitedMember
		if err := rows.Scan(&member.EventId, &member.MemberId, &member.Answer, &member.Member.Id, &member.Member.Email, &member.Member.FirstName, &member.Member.LastName, &member.Member.BirthDate, &member.Member.AvatarPath, &member.Member.Nickname, &member.Member.AboutMe, &member.Member.Followers, &member.Member.IsPublic); err != nil {
			fmt.Println("Scanning going members failed: " + err.Error())
			return nil, err
		}
		members = append(members, &member)
	}
	return members, nil
}

func GetUserEventAnswer(invitedMember models.InvitedMember) (int, error) {
	var res int
	err := sqlite.Db.QueryRow("SELECT ANSWERTYPE FROM event_invited_members WHERE EVENTID = ? AND USERID = ?", invitedMember.EventId, invitedMember.MemberId).Scan(&res)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Querying is user going failed: " + err.Error())
			return 0, err
		}
		return 0, nil
	}
	return res, err
}

func SetUserEventAnswer(invitedMember models.InvitedMember) error {
	_, err := sqlite.Db.Exec("UPDATE event_invited_members SET ANSWERTYPE = ? WHERE EVENTID = ? AND USERID = ?", invitedMember.Answer, invitedMember.EventId, invitedMember.MemberId)
	if err != nil {
		fmt.Println("Updating user event answer failed: " + err.Error())
	}
	return err
}
