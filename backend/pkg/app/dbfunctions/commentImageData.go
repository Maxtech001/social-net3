package dbfunctions

import (
	"fmt"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/models"
)

func InsertCommentImageConnection(commentImg models.CommentImage) (*models.CommentImage, error) {
	var ret models.CommentImage
	err := sqlite.Db.QueryRow("INSERT INTO comment_images_connections (COMMENTID, IMAGEPATH) VALUES(?, ?) RETURNING COMMENTID, IMAGEPATH", commentImg.CommentId, commentImg.ImagePath).Scan(&ret.CommentId, &ret.ImagePath)
	if err != nil {
		fmt.Println("Inserting comment image connection failed: " + err.Error())
		return nil, err
	}
	return &ret, err
}

func GetCommentImageConnectionByCommentId(commentId int) ([]*models.CommentImage, error) {
	imgArray := []*models.CommentImage{}
	rows, err := sqlite.Db.Query("SELECT COMMENTID, IMAGEPATH FROM comment_images_connections WHERE COMMENTID = ?", commentId)
	if err != nil {
		fmt.Println("Getting comment images failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var postImg models.CommentImage
		if err := rows.Scan(&postImg.CommentId, &postImg.ImagePath); err != nil {
			fmt.Println("Scanning comment images failed: " + err.Error())
			return nil, err
		}
		imgArray = append(imgArray, &postImg)
	}
	return imgArray, nil
}
