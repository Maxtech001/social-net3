package dbfunctions

import (
	"database/sql"
	"fmt"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/models"
)

func InsertUserData(user models.User) (*models.User, error) {
	var res models.User
	err := sqlite.Db.QueryRow("INSERT INTO users (EMAIL, UPASSWORD, FIRSTNAME, LASTNAME, BIRTHDATE, AVATARPATH, NICKNAME, ABOUTME, FOLLOWERS, ISPUBLIC) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING ID, EMAIL, UPASSWORD, FIRSTNAME, LASTNAME, BIRTHDATE, AVATARPATH, NICKNAME, ABOUTME, FOLLOWERS, ISPUBLIC",
		user.Email, user.Password, user.FirstName, user.LastName, user.BirthDate, user.AvatarPath, user.Nickname, user.AboutMe, user.Followers, user.IsPublic).Scan(&res.Id, &res.Email, &res.Password, &res.FirstName, &res.LastName, &res.BirthDate, &res.AvatarPath, &res.Nickname, &res.AboutMe, &res.Followers, &res.IsPublic)
	if err != nil {
		fmt.Println("Failed to store user data: " + err.Error())
		return nil, err
	}
	return &res, nil
}

func RemoveUserData(user models.User) error {
	_, err := sqlite.Db.Exec("DELETE FROM users WHERE ID = ?", user.Id)
	if err != nil {
		fmt.Println("Failed to remove user data: " + err.Error())
	}
	return err
}

func GetAccountDataByEmail(email string) (*models.User, error) {
	var user models.User
	err := sqlite.Db.QueryRow("SELECT ID, EMAIL, UPASSWORD, FIRSTNAME, LASTNAME, BIRTHDATE, AVATARPATH, NICKNAME, ABOUTME, FOLLOWERS, ISPUBLIC FROM users WHERE EMAIL = ?", email).Scan(&user.Id, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.BirthDate, &user.AvatarPath, &user.Nickname, &user.AboutMe, &user.Followers, &user.IsPublic)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Failed to get user data by email: " + err.Error())
			return nil, err
		}
		return nil, nil
	}
	return &user, nil
}

func GetAccountDataById(userId int) (*models.User, error) {
	var user models.User
	err := sqlite.Db.QueryRow("SELECT ID, EMAIL, UPASSWORD, FIRSTNAME, LASTNAME, BIRTHDATE, AVATARPATH, NICKNAME, ABOUTME, FOLLOWERS, ISPUBLIC FROM users WHERE ID = ?", userId).Scan(&user.Id, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.BirthDate, &user.AvatarPath, &user.Nickname, &user.AboutMe, &user.Followers, &user.IsPublic)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Failed to get user data by id: " + err.Error())
			return nil, err
		}
		return nil, nil
	}
	return &user, nil
}

func GetUsersByName(name string) ([]*models.User, error) {
	users := []*models.User{}
	rows, err := sqlite.Db.Query("SELECT ID, EMAIL, UPASSWORD, FIRSTNAME, LASTNAME, BIRTHDATE, AVATARPATH, NICKNAME, ABOUTME, FOLLOWERS, ISPUBLIC FROM users WHERE FIRSTNAME || ' ' || LASTNAME LIKE '%" + name + "%'")
	if err != nil {
		fmt.Println("Getting user groups failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.BirthDate, &user.AvatarPath, &user.Nickname, &user.AboutMe, &user.Followers, &user.IsPublic); err != nil {
			fmt.Println("Scanning user groups failed: " + err.Error())
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func InsertAvatarPath(user models.User) error {
	_, err := sqlite.Db.Exec("UPDATE users SET AVATARPATH = ? WHERE ID = ?", user.AvatarPath, user.Id)
	if err != nil {
		fmt.Println("Failed to store avatar path: " + err.Error())
	}
	return err
}

func UpdateFollowCount(userId int) error {
	_, err := sqlite.Db.Exec("UPDATE users SET FOLLOWERS = (SELECT COUNT(*) FROM follower_connections WHERE FOLLOWEDID = ?) WHERE ID = ?", userId, userId)
	if err != nil {
		fmt.Println("Updating follower count failed: " + err.Error())
	}
	return err
}

func SetUserPublic(userId int, val bool) error {
	_, err := sqlite.Db.Exec("UPDATE users SET ISPUBLIC = ? WHERE ID = ?", val, userId)
	if err != nil {
		fmt.Println("Failed to set public status: " + err.Error())
	}
	return err
}
