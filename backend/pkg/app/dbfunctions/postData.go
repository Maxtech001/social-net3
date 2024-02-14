package dbfunctions

import (
	"database/sql"
	"fmt"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/models"
)

func InsertPostData(post models.Post) (*models.Post, error) {
	var res models.Post
	err := sqlite.Db.QueryRow("INSERT INTO posts (USERID, GROUPID ,CONTENT, TOTALCOMMENTS, CREATEDAT, POSTTYPE) VALUES (?, ?, ?, ?, ?, ?) RETURNING ID, USERID, GROUPID, (SELECT FIRSTNAME || ' ' || LASTNAME FROM users WHERE users.ID = USERID), (SELECT AVATARPATH FROM users WHERE users.ID = USERID), CONTENT, TOTALCOMMENTS, CREATEDAT, POSTTYPE", post.AuthorId, post.GroupId, post.Content, post.TotalComments, post.CreatedAt, post.PostType).Scan(&res.Id, &res.AuthorId, &res.GroupId, &res.AuthorName, &res.AuthorAvatar, &res.Content, &res.TotalComments, &res.CreatedAt, &res.PostType)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Insert post data failed: " + err.Error())
			return nil, err
		}
		return nil, nil
	}
	return &res, nil
}

func GetVisiblePosts(user models.User) ([]*models.Post, error) {
	posts := []*models.Post{}
	rows, err := sqlite.Db.Query(`
	SELECT posts.ID,
		USERID,
		users.FIRSTNAME || ' ' || users.LASTNAME,
		users.AVATARPATH,
		CONTENT,
		TOTALCOMMENTS,
		posts.CREATEDAT,
		GROUPID
	FROM posts
		INNER JOIN users ON users.ID = posts.USERID
	WHERE USERID = ? OR
		(USERID IN (SELECT FOLLOWEDID FROM follower_connections WHERE FOLLOWERID = ?) AND (POSTTYPE != 2 OR ? IN (SELECT USERID from post_user_connections AS puc WHERE puc.POSTID = posts.ID))) OR
		posts.GROUPID IN (SELECT gpm.GROUPID FROM group_members as gpm WHERE gpm.MEMBERID = ?)
	ORDER BY posts.ID DESC`, user.Id, user.Id, user.Id, user.Id)
	if err != nil {
		fmt.Println("Getting all posts failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.Id, &post.AuthorId, &post.AuthorName, &post.AuthorAvatar, &post.Content, &post.TotalComments, &post.CreatedAt, &post.GroupId); err != nil {
			fmt.Println("Scanning all posts failed: " + err.Error())
			return nil, err
		}
		if post.ImageArray, err = GetPostImageConnectionByPostId(post.Id); err != nil {
			return nil, err
		}
		if post.GroupId != 0 {
			if group, err := GetGroupDataById(post.GroupId); err != nil {
				return nil, err
			} else if group != nil {
				post.Group = *group
			}
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

func GetAllGroupPosts(user models.User) ([]*models.Post, error) {
	posts := []*models.Post{}
	rows, err := sqlite.Db.Query(`
		SELECT posts.id,
			userid,
			users.firstname
			|| ' '
			|| users.lastname,
			users.avatarpath,
			posts.content,
			posts.totalcomments,
			posts.createdat,
			g.id,
			g.adminid,
			g.groupname,
			g.groupdescription,
			g.membercount,
			g.createdat
		FROM group_members  as gpm
			INNER JOIN posts
				ON posts.groupID = gpm.groupID
			INNER JOIN groups AS g
				ON gpm.GROUPID = g.ID
			INNER JOIN users
				ON posts.USERID = users.Id
		WHERE  gpm.memberid = ?
		ORDER  BY posts.id DESC
	`, user.Id)
	if err != nil {
		fmt.Println("Getting all posts failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.Id, &post.AuthorId, &post.AuthorName, &post.AuthorAvatar, &post.Content, &post.TotalComments, &post.CreatedAt, &post.Group.Id, &post.Group.AdminId, &post.Group.GroupName, &post.Group.Description, &post.Group.MemberCount, &post.Group.CreatedAt); err != nil {
			fmt.Println("Scanning all posts failed: " + err.Error())
			return nil, err
		}
		if post.ImageArray, err = GetPostImageConnectionByPostId(post.Id); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

func GetPostsByGroupId(groupId int) ([]*models.Post, error) {
	posts := []*models.Post{}
	rows, err := sqlite.Db.Query("SELECT posts.ID, USERID, users.FIRSTNAME || ' ' || users.LASTNAME, users.AVATARPATH, CONTENT, TOTALCOMMENTS, CREATEDAT FROM posts INNER JOIN users ON users.ID = posts.USERID WHERE GROUPID = ? ORDER BY posts.ID DESC", groupId)
	if err != nil {
		fmt.Println("Getting user posts failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.Id, &post.AuthorId, &post.AuthorName, &post.AuthorAvatar, &post.Content, &post.TotalComments, &post.CreatedAt); err != nil {
			fmt.Println("Scanning user posts failed: " + err.Error())
			return nil, err
		}
		if post.ImageArray, err = GetPostImageConnectionByPostId(post.Id); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

func GetUserPosts(reqUserId, curUserId int) ([]*models.Post, error) {
	posts := []*models.Post{}
	rows, err := sqlite.Db.Query(`
	SELECT posts.ID, USERID, users.FIRSTNAME || ' ' || users.LASTNAME, users.AVATARPATH, CONTENT, TOTALCOMMENTS, CREATEDAT
	FROM posts
		INNER JOIN users ON users.ID = posts.USERID
	WHERE
		USERID = ? AND
		(USERID = ? OR (POSTTYPE = 0 AND EXISTS(SELECT FOLLOWEDID FROM follower_connections WHERE FOLLOWERID = ? AND FOLLOWEDID = ?)) OR (users.ISPUBLIC AND POSTTYPE = 1) OR (POSTTYPE = 2 AND EXISTS(SELECT USERID from post_user_connections AS puc WHERE puc.POSTID = posts.ID AND USERID = ?)))
	ORDER BY posts.ID DESC`, reqUserId, curUserId, curUserId, reqUserId, curUserId)
	if err != nil {
		fmt.Println("Getting user posts failed: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.Id, &post.AuthorId, &post.AuthorName, &post.AuthorAvatar, &post.Content, &post.TotalComments, &post.CreatedAt); err != nil {
			fmt.Println("Scanning user posts failed: " + err.Error())
			return nil, err
		}
		if post.ImageArray, err = GetPostImageConnectionByPostId(post.Id); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

func UpdatePostCommentCount(postId int) error {
	_, err := sqlite.Db.Exec("UPDATE posts SET TOTALCOMMENTS = (SELECT COUNT(*) FROM comments WHERE POSTID = ?) WHERE ID = ?", postId, postId)
	if err != nil {
		fmt.Println("Updating post comment count failed: " + err.Error())
	}
	return err
}
