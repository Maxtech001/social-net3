package dbfunctions

import (
	"database/sql"
	"fmt"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/models"
)

func InsertFollowNotification(notification models.Notification) (*models.Notification, error) {
	var res models.Notification
	err := sqlite.Db.QueryRow("INSERT INTO notifications (USERID, CONTENT, NOTIFICATIONTYPE, ISREAD, CREATEDAT, FOLLOWERID) VALUES (?, ?, ?, ?, ?, ?) RETURNING ID, USERID, CONTENT, NOTIFICATIONTYPE, ISREAD, CREATEDAT, FOLLOWERID", notification.UserId, notification.Content, notification.NotificationType, notification.IsRead, notification.CreatedAt, notification.FollowerId).Scan(&res.Id, &res.UserId, &res.Content, &res.NotificationType, &res.IsRead, &res.CreatedAt, &res.FollowerId)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Insert notification failed: " + err.Error())
			return nil, err
		}
		return nil, nil
	}
	var user *models.User
	if user, err = GetAccountDataById(res.UserId); err != nil {
		return nil, err
	}
	res.User = *user
	if user, err = GetAccountDataById(res.FollowerId); err != nil {
		return nil, err
	}
	res.Follower = user
	return &res, nil
}

func InsertEventNotification(notification models.Notification) (*models.Notification, error) {
	var res models.Notification
	err := sqlite.Db.QueryRow("INSERT INTO notifications (USERID, CONTENT, NOTIFICATIONTYPE, ISREAD, CREATEDAT, EVENTID) VALUES (?, ?, ?, ?, ?, ?) RETURNING ID, USERID, CONTENT, NOTIFICATIONTYPE, ISREAD, CREATEDAT, EVENTID", notification.UserId, notification.Content, notification.NotificationType, notification.IsRead, notification.CreatedAt, notification.EventId).Scan(&res.Id, &res.UserId, &res.Content, &res.NotificationType, &res.IsRead, &res.CreatedAt, &res.EventId)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Insert notification failed: " + err.Error())
			return nil, err
		}
		return nil, nil
	}
	var user *models.User
	if user, err = GetAccountDataById(res.UserId); err != nil {
		return nil, err
	}
	res.User = *user
	var event *models.Event
	if event, err = GetEventById(res.EventId); err != nil {
		return nil, err
	}
	res.Event = event
	return &res, nil
}

func InsertGroupNotification(notification models.Notification) (*models.Notification, error) {
	var res models.Notification
	err := sqlite.Db.QueryRow("INSERT INTO notifications (USERID, CONTENT, NOTIFICATIONTYPE, ISREAD, CREATEDAT, GROUPID) VALUES (?, ?, ?, ?, ?, ?) RETURNING ID, USERID, CONTENT, NOTIFICATIONTYPE, ISREAD, CREATEDAT, GROUPID", notification.UserId, notification.Content, notification.NotificationType, notification.IsRead, notification.CreatedAt, notification.GroupId).Scan(&res.Id, &res.UserId, &res.Content, &res.NotificationType, &res.IsRead, &res.CreatedAt, &res.GroupId)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Insert notification failed: " + err.Error())
			return nil, err
		}
		return nil, nil
	}
	var user *models.User
	if user, err = GetAccountDataById(res.UserId); err != nil {
		return nil, err
	}
	res.User = *user
	var group *models.Group
	if group, err = GetGroupDataById(res.GroupId); err != nil {
		return nil, err
	}
	res.Group = group
	return &res, nil
}

func GetAllNotificationsByUserId(userId int) ([]*models.Notification, error) {
	notifications := []*models.Notification{}
	rows, err := sqlite.Db.Query("SELECT n.ID, n.USERID, n.CONTENT, n.NOTIFICATIONTYPE, n.ISREAD, n.CREATEDAT, n.FOLLOWERID, n.GROUPID, n.EVENTID, u.ID, u.EMAIL, u.FIRSTNAME, u.LASTNAME, u.BIRTHDATE, u.AVATARPATH, u.NICKNAME, u.ABOUTME, u.FOLLOWERS, u.ISPUBLIC FROM notifications as n INNER JOIN users as u on u.ID = n.USERID WHERE n.USERID = ?", userId)
	if err != nil {
		fmt.Println("Getting notifications failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var notification models.Notification
		if err := rows.Scan(&notification.Id, &notification.UserId, &notification.Content, &notification.NotificationType, &notification.IsRead, &notification.CreatedAt, &notification.NullFollowerId, &notification.NullGroupId, &notification.NullEventId, &notification.User.Id, &notification.User.Email, &notification.User.FirstName, &notification.User.LastName, &notification.User.BirthDate, &notification.User.AvatarPath, &notification.User.Nickname, &notification.User.AboutMe, &notification.User.Followers, &notification.User.IsPublic); err != nil {
			fmt.Println("Scanning notifications failed: " + err.Error())
			return nil, err
		}
		var group *models.Group
		notification.GroupId = int(notification.NullGroupId.Int32)
		if group, err = GetGroupDataById(notification.GroupId); err != nil {
			return nil, err
		}
		notification.Group = group
		var event *models.Event
		notification.EventId = int(notification.NullEventId.Int32)
		if event, err = GetEventById(notification.EventId); err != nil {
			return nil, err
		}
		notification.Event = event
		var follower *models.User
		notification.FollowerId = int(notification.NullFollowerId.Int32)
		if follower, err = GetAccountDataById(notification.FollowerId); err != nil {
			return nil, err
		}
		notification.Follower = follower
		notifications = append(notifications, &notification)
	}
	return notifications, nil
}

func SetNotificationIsRead(notificationId int, val bool) error {
	_, err := sqlite.Db.Exec("UPDATE notifications SET ISREAD = ? WHERE ID = ?", val, notificationId)
	if err != nil {
		fmt.Println("Failed to set isRead status: " + err.Error())
	}
	return err
}
