package dbfunctions

import (
	"database/sql"
	"fmt"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/models"
)

func InsertFollowConnection(followConn models.FollowConnection) (*models.User, error) {
	_, err := sqlite.Db.Exec("INSERT INTO follower_connections (FOLLOWEDID, FOLLOWERID) VALUES (?, ?)", followConn.FollowedId, followConn.FollowerId)
	if err != nil {
		fmt.Println("Inserting follow connection failed: " + err.Error())
		return nil, err
	}
	err = UpdateFollowCount(followConn.FollowedId)
	if err != nil {
		return nil, err
	}
	return GetAccountDataById(followConn.FollowedId)
}
func GetFollowersByFollowedId(followedId int) ([]*models.User, error) {
	users := []*models.User{}
	rows, err := sqlite.Db.Query("SELECT users.ID, users.EMAIL, users.UPASSWORD, users.FIRSTNAME, users.LASTNAME, users.BIRTHDATE, users.AVATARPATH, users.NICKNAME, users.ABOUTME, users.FOLLOWERS, users.ISPUBLIC FROM follower_connections INNER JOIN users ON users.ID = FOLLOWERID WHERE FOLLOWEDID = ?", followedId)
	if err != nil {
		fmt.Println("Getting followers failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.BirthDate, &user.AvatarPath, &user.Nickname, &user.AboutMe, &user.Followers, &user.IsPublic); err != nil {
			fmt.Println("Scanning followers failed: " + err.Error())
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func GetFollowedusersByFollowerId(followerId int) ([]*models.User, error) {
	users := []*models.User{}
	rows, err := sqlite.Db.Query("SELECT users.ID, users.EMAIL, users.UPASSWORD, users.FIRSTNAME, users.LASTNAME, users.BIRTHDATE, users.AVATARPATH, users.NICKNAME, users.ABOUTME, users.FOLLOWERS, users.ISPUBLIC FROM follower_connections INNER JOIN users ON users.ID = FOLLOWEDID WHERE FOLLOWERID = ?", followerId)
	if err != nil {
		fmt.Println("Getting followed users failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.BirthDate, &user.AvatarPath, &user.Nickname, &user.AboutMe, &user.Followers, &user.IsPublic); err != nil {
			fmt.Println("Scanning followed users failed: " + err.Error())
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func GetFollowedAndFollowerUsersByUserId(userId int) ([]*models.User, error) {
	users := []*models.User{}
	rows, err := sqlite.Db.Query("SELECT users.ID, users.EMAIL, users.UPASSWORD, users.FIRSTNAME, users.LASTNAME, users.BIRTHDATE, users.AVATARPATH, users.NICKNAME, users.ABOUTME, users.FOLLOWERS, users.ISPUBLIC FROM follower_connections INNER JOIN users ON CASE WHEN FOLLOWEDID = ? THEN users.ID = FOLLOWERID ELSE users.ID = FOLLOWEDID END WHERE FOLLOWERID = ? OR FOLLOWEDID = ?", userId, userId, userId)
	if err != nil {
		fmt.Println("Getting followed users failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.BirthDate, &user.AvatarPath, &user.Nickname, &user.AboutMe, &user.Followers, &user.IsPublic); err != nil {
			fmt.Println("Scanning followed users failed: " + err.Error())
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func IsExistingFollowConnection(followConn models.FollowConnection) (bool, error) {
	var res bool
	err := sqlite.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM follower_connections WHERE FOLLOWEDID = ? AND FOLLOWERID = ?)", followConn.FollowedId, followConn.FollowerId).Scan(&res)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Querying is exisiting follow connection failed: " + err.Error())
			return false, err
		}
		return false, nil
	}
	return res, err
}

func RemoveFollowConnection(followConn models.FollowConnection) (*models.User, error) {
	_, err := sqlite.Db.Exec("DELETE FROM follower_connections WHERE FOLLOWEDID = ? AND FOLLOWERID = ?", followConn.FollowedId, followConn.FollowerId)
	if err != nil {
		fmt.Println("Deleting follow connection failed: " + err.Error())
	}
	err = UpdateFollowCount(followConn.FollowedId)
	if err != nil {
		return nil, err
	}
	return GetAccountDataById(followConn.FollowedId)
}

func GetMutualFollowers(userId int) ([]*models.User, error) {
	users := []*models.User{}
	rows, err := sqlite.Db.Query("SELECT users.ID, users.EMAIL, users.UPASSWORD, users.FIRSTNAME, users.LASTNAME, users.BIRTHDATE, users.AVATARPATH, users.NICKNAME, users.ABOUTME, users.FOLLOWERS, users.ISPUBLIC FROM follower_connections INNER JOIN users ON users.ID = FOLLOWERID WHERE FOLLOWEDID = ? AND EXISTS(SELECT 1 FROM follower_connections AS fw2 WHERE fw2.FOLLOWEDID = follower_connections.FOLLOWERID AND fw2.FOLLOWERID = follower_connections.FOLLOWEDID)", userId)
	if err != nil {
		fmt.Println("Getting all mutual followers failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.BirthDate, &user.AvatarPath, &user.Nickname, &user.AboutMe, &user.Followers, &user.IsPublic); err != nil {
			fmt.Println("Scanning all mutual followers failed: " + err.Error())
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func GetMutualFollowersNotInGroup(userId, groupId int) ([]*models.User, error) {
	users := []*models.User{}
	rows, err := sqlite.Db.Query(`
	SELECT u.id,
       u.email,
       u.upassword,
       u.firstname,
       u.lastname,
       u.birthdate,
       u.avatarpath,
       u.nickname,
       u.aboutme,
       u.followers,
       u.ispublic
	FROM   follower_connections AS fc1
		INNER JOIN follower_connections AS fc2
				ON fc1.followedid = fc2.followerid
					AND fc2.followedid = fc1.followerid
		INNER JOIN users AS u
				ON fc1.followerid = u.id
	WHERE  fc1.followedid = ?
		AND fc1.followerid NOT IN (
			SELECT MEMBERID 
			FROM group_members as gp
			WHERE gp.groupid = ?
		)
		AND fc1.followerid NOT IN (
			SELECT INVITEDID 
			FROM group_invites as gi 
			WHERE gi.groupid = ?
		)
			`, userId, groupId, groupId)
	if err != nil {
		fmt.Println("Getting all mutual followers failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.BirthDate, &user.AvatarPath, &user.Nickname, &user.AboutMe, &user.Followers, &user.IsPublic); err != nil {
			fmt.Println("Scanning all mutual followers failed: " + err.Error())
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func GetFollowersNotInGroup(userId, groupId int) ([]*models.User, error) {
	users := []*models.User{}
	rows, err := sqlite.Db.Query(`
	SELECT u.id,
       u.email,
       u.upassword,
       u.firstname,
       u.lastname,
       u.birthdate,
       u.avatarpath,
       u.nickname,
       u.aboutme,
       u.followers,
       u.ispublic
	FROM   follower_connections AS fc1
		INNER JOIN users AS u
				ON fc1.followerid = u.id
	WHERE  fc1.followedid = ?
		AND fc1.followerid NOT IN (
			SELECT MEMBERID 
			FROM group_members as gp
			WHERE gp.groupid = ?
		)
		AND fc1.followerid NOT IN (
			SELECT INVITEDID 
			FROM group_invites as gi 
			WHERE gi.groupid = ?
		)
			`, userId, groupId, groupId)
	if err != nil {
		fmt.Println("Getting all mutual followers failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.BirthDate, &user.AvatarPath, &user.Nickname, &user.AboutMe, &user.Followers, &user.IsPublic); err != nil {
			fmt.Println("Scanning all followers not in group failed: " + err.Error())
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
