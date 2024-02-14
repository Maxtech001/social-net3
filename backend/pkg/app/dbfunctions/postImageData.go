package dbfunctions

import (
	"fmt"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/models"
)

func InsertPostImageConnection(postImg models.PostImage) (*models.PostImage, error) {
	var ret models.PostImage
	err := sqlite.Db.QueryRow("INSERT INTO post_images_connections (POSTID, IMAGEPATH) VALUES(?, ?) RETURNING POSTID, IMAGEPATH", postImg.PostId, postImg.ImagePath).Scan(&ret.PostId, &ret.ImagePath)
	if err != nil {
		fmt.Println("Inserting post image connection failed: " + err.Error())
		return nil, err
	}
	return &ret, err
}

func GetPostImageConnectionByPostId(postId int) ([]*models.PostImage, error) {
	imgArray := []*models.PostImage{}
	rows, err := sqlite.Db.Query("SELECT POSTID, IMAGEPATH FROM post_images_connections WHERE POSTID = ?", postId)
	if err != nil {
		fmt.Println("Getting post images failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var postImg models.PostImage
		if err := rows.Scan(&postImg.PostId, &postImg.ImagePath); err != nil {
			fmt.Println("Scanning post images failed: " + err.Error())
			return nil, err
		}
		imgArray = append(imgArray, &postImg)
	}
	return imgArray, nil
}
