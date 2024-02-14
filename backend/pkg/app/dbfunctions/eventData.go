package dbfunctions

import (
	"database/sql"
	"fmt"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/models"
)

func InsertEventData(event models.Event) (*models.Event, error) {
	var res models.Event
	var author *models.User
	err := sqlite.Db.QueryRow("INSERT INTO events (GROUPID, AUTHORID, EVENTTITLE, EVENTDESCRIPTION, STARTDATE, ENDDATE, CREATEDAT) VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING events.ID, events.GROUPID, events.AUTHORID, events.EVENTTITLE, events.EVENTDESCRIPTION, events.STARTDATE, events.ENDDATE, events.CREATEDAT", event.GroupId, event.AuthorId, event.Title, event.Description, event.StartDate, event.EndDate, event.CreatedAt).Scan(&res.Id, &res.GroupId, &res.AuthorId, &res.Title, &res.Description, &res.StartDate, &res.EndDate, &res.CreatedAt)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Insert event data failed: " + err.Error())
			return nil, err
		}
		return nil, nil
	}
	if author, err = GetAccountDataById(res.AuthorId); err != nil {
		return nil, err
	}
	res.Author = *author
	return &res, nil
}

func GetEventById(eventId int) (*models.Event, error) {
	var event models.Event
	err := sqlite.Db.QueryRow("SELECT events.ID, events.GROUPID, events.AUTHORID, events.EVENTTITLE, events.EVENTDESCRIPTION, events.STARTDATE, events.ENDDATE, events.CREATEDAT, users.ID, users.EMAIL, users.FIRSTNAME, users.LASTNAME, users.BIRTHDATE, users.AVATARPATH, users.NICKNAME, users.ABOUTME, users.FOLLOWERS, users.ISPUBLIC, groups.ID, groups.ADMINID, groups.GROUPNAME, groups.GROUPDESCRIPTION, groups.MEMBERCOUNT, groups.CREATEDAT FROM events INNER JOIN users ON users.ID = events.AUTHORID INNER JOIN groups ON groups.ID = events.GROUPID WHERE events.ID = ?", eventId).Scan(&event.Id, &event.GroupId, &event.AuthorId, &event.Title, &event.Description, &event.StartDate, &event.EndDate, &event.CreatedAt, &event.Author.Id, &event.Author.Email, &event.Author.FirstName, &event.Author.LastName, &event.Author.BirthDate, &event.Author.AvatarPath, &event.Author.Nickname, &event.Author.AboutMe, &event.Author.Followers, &event.Author.IsPublic, &event.Group.Id, &event.Group.AdminId, &event.Group.GroupName, &event.Group.Description, &event.Group.MemberCount, &event.Group.CreatedAt)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Getting event data failed: " + err.Error())
			return nil, err
		}
		return nil, nil
	}
	return &event, nil
}

func GetUpcomingEventsByGroupId(groupId int) ([]*models.Event, error) {
	events := []*models.Event{}
	rows, err := sqlite.Db.Query("SELECT events.ID, events.GROUPID, events.AUTHORID, events.EVENTTITLE, events.EVENTDESCRIPTION, events.STARTDATE, events.ENDDATE, events.CREATEDAT, users.ID, users.EMAIL, users.FIRSTNAME, users.LASTNAME, users.BIRTHDATE, users.AVATARPATH, users.NICKNAME, users.ABOUTME, users.FOLLOWERS, users.ISPUBLIC FROM events INNER JOIN users ON users.ID = events.AUTHORID WHERE events.GROUPID = ? AND STARTDATE >= CURRENT_DATE ORDER BY STARTDATE DESC", groupId)
	if err != nil {
		fmt.Println("Getting upcoming events failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.Id, &event.GroupId, &event.AuthorId, &event.Title, &event.Description, &event.StartDate, &event.EndDate, &event.CreatedAt, &event.Author.Id, &event.Author.Email, &event.Author.FirstName, &event.Author.LastName, &event.Author.BirthDate, &event.Author.AvatarPath, &event.Author.Nickname, &event.Author.AboutMe, &event.Author.Followers, &event.Author.IsPublic); err != nil {
			fmt.Println("Scanning upcoming events failed: " + err.Error())
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}

func GetPastEventsByGroupId(groupId int) ([]*models.Event, error) {
	events := []*models.Event{}
	rows, err := sqlite.Db.Query("SELECT events.ID, events.GROUPID, events.AUTHORID, events.EVENTTITLE, events.EVENTDESCRIPTION, events.STARTDATE, events.ENDDATE, events.CREATEDAT, users.ID, users.EMAIL, users.FIRSTNAME, users.LASTNAME, users.BIRTHDATE, users.AVATARPATH, users.NICKNAME, users.ABOUTME, users.FOLLOWERS, users.ISPUBLIC FROM events INNER JOIN users ON users.ID = events.AUTHORID WHERE events.GROUPID = ? AND STARTDATE < CURRENT_DATE ORDER BY STARTDATE DESC", groupId)
	if err != nil {
		fmt.Println("Getting past events failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.Id, &event.GroupId, &event.AuthorId, &event.Title, &event.Description, &event.StartDate, &event.EndDate, &event.CreatedAt, &event.Author.Id, &event.Author.Email, &event.Author.FirstName, &event.Author.LastName, &event.Author.BirthDate, &event.Author.AvatarPath, &event.Author.Nickname, &event.Author.AboutMe, &event.Author.Followers, &event.Author.IsPublic); err != nil {
			fmt.Println("Scanning past events failed: " + err.Error())
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}
