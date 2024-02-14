package dbfunctions

import (
	"fmt"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/models"
)

func InsertPostUserConnection(postUser models.PostUser) (*models.PostUser, error) {
	var ret models.PostUser
	err := sqlite.Db.QueryRow("INSERT INTO post_user_connections (POSTID, USERID) VALUES(?, ?) RETURNING POSTID, USERID", postUser.PostId, postUser.UserId).Scan(&ret.PostId, &ret.UserId)
	if err != nil {
		fmt.Println("Inserting post user connection failed: " + err.Error())
		return nil, err
	}
	return &ret, err
}
