package dbfunctions

import (
	"database/sql"
	"fmt"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/models"
)

func InsertCommentData(comment models.Comment) (*models.Comment, error) {
	var res models.Comment
	err := sqlite.Db.QueryRow("INSERT INTO comments (POSTID, USERID, CONTENT, CREATEDAT) VALUES (?, ?, ?, ?) RETURNING ID, POSTID, USERID, (SELECT FIRSTNAME || ' ' || LASTNAME FROM users WHERE users.ID = USERID), (SELECT AVATARPATH FROM users WHERE users.ID = USERID), CONTENT, CREATEDAT", comment.PostId, comment.AuthorId, comment.Content, comment.CreatedAt).Scan(&res.Id, &res.PostId, &res.AuthorId, &res.AuthorName, &res.AuthorAvatar, &res.Content, &res.CreatedAt)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Insert comment data failed: " + err.Error())
			return nil, err
		}
		return nil, nil
	}
	err = UpdatePostCommentCount(res.PostId)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func GetCommentsByPostId(postId int) ([]*models.Comment, error) {
	comments := []*models.Comment{}
	rows, err := sqlite.Db.Query("SELECT comments.ID, POSTID, USERID, users.FIRSTNAME || ' ' || users.LASTNAME, users.AVATARPATH, CONTENT, CREATEDAT FROM comments INNER JOIN users ON users.ID = comments.USERID WHERE POSTID = ? ORDER BY comments.ID DESC", postId)
	if err != nil {
		fmt.Println("Getting post comments failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.Id, &comment.PostId, &comment.AuthorId, &comment.AuthorName, &comment.AuthorAvatar, &comment.Content, &comment.CreatedAt); err != nil {
			fmt.Println("Scanning post comments failed: " + err.Error())
			return nil, err
		}
		if comment.ImageArray, err = GetCommentImageConnectionByCommentId(comment.Id); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}
