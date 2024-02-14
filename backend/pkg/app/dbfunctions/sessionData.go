package dbfunctions

import (
	"database/sql"
	"fmt"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/models"
)

func RemoveCookieData(userId int) error {
	_, err := sqlite.Db.Exec("DELETE FROM sessions WHERE USERID = ?", userId)
	if err != nil {
		fmt.Println("failed to remove cookie data: " + err.Error())
	}
	return err
}

func GetSessionData(sessionId string) (*models.Session, error) {
	var session models.Session
	err := sqlite.Db.QueryRow("SELECT ID, USERID, EXPIRY FROM sessions WHERE ID = ?", sessionId).Scan(&session.Id, &session.AccountId, &session.Expiry)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Get session data failed: " + err.Error())
			return nil, err
		}
		return nil, nil
	}
	return &session, nil
}

func GetUserBySessionId(sessionId string) (*models.User, error) {
	var user models.User
	err := sqlite.Db.QueryRow("SELECT users.ID, users.EMAIL, users.UPASSWORD, users.FIRSTNAME, users.LASTNAME, users.BIRTHDATE, users.AVATARPATH, users.NICKNAME, users.ABOUTME, users.FOLLOWERS, users.ISPUBLIC FROM sessions INNER JOIN users ON users.ID = sessions.USERID WHERE sessions.ID = ?", sessionId).Scan(&user.Id, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.BirthDate, &user.AvatarPath, &user.Nickname, &user.AboutMe, &user.Followers, &user.IsPublic)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Get all post data failed: " + err.Error())
			return nil, err
		}
		return nil, nil
	}
	return &user, nil
}

func InsertCookieData(session models.Session) error {
	_, err := sqlite.Db.Exec("INSERT INTO sessions (ID, USERID, EXPIRY) VALUES (?, ?, ?)", session.Id, session.AccountId, session.Expiry)
	if err != nil {
		fmt.Println("Inserting cookie failed: ", err.Error())
	}
	return err
}
