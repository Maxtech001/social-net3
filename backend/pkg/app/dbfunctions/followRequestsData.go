package dbfunctions

import (
	"database/sql"
	"fmt"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/models"
)

func InsertFollowRequest(followConn models.FollowConnection) (*models.User, error) {
	_, err := sqlite.Db.Exec("INSERT INTO follow_requests (FOLLOWEDID, FOLLOWERID) VALUES (?, ?)", followConn.FollowedId, followConn.FollowerId)
	if err != nil {
		fmt.Println("Inserting follow request failed: " + err.Error())
		return nil, err
	}
	return GetAccountDataById(followConn.FollowedId)
}

func RemoveFollowRequest(followConn models.FollowConnection) error {
	_, err := sqlite.Db.Exec("DELETE FROM follow_requests WHERE FOLLOWEDID = ? AND FOLLOWERID = ?", followConn.FollowedId, followConn.FollowerId)
	if err != nil {
		fmt.Println("Deleting follow request failed: " + err.Error())
	}
	return err
}

func GetFollowRequestByFollowedId(followedId int) ([]*models.User, error) {
	users := []*models.User{}
	rows, err := sqlite.Db.Query("SELECT users.ID, users.EMAIL, users.UPASSWORD, users.FIRSTNAME, users.LASTNAME, users.BIRTHDATE, users.AVATARPATH, users.NICKNAME, users.ABOUTME, users.FOLLOWERS, users.ISPUBLIC FROM follow_requests INNER JOIN users ON users.ID = FOLLOWERID WHERE FOLLOWEDID = ?", followedId)
	if err != nil {
		fmt.Println("Getting follow requests failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.BirthDate, &user.AvatarPath, &user.Nickname, &user.AboutMe, &user.Followers, &user.IsPublic); err != nil {
			fmt.Println("Scanning follow requests failed: " + err.Error())
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func IsExistingFollowRequest(followConn models.FollowConnection) (bool, error) {
	var res bool
	err := sqlite.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM follow_requests WHERE FOLLOWEDID = ? AND FOLLOWERID = ?)", followConn.FollowedId, followConn.FollowerId).Scan(&res)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Querying is exisiting follow request failed: " + err.Error())
			return false, err
		}
		return false, nil
	}
	return res, err
}
