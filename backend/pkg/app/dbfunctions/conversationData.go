package dbfunctions

import (
	"database/sql"
	"fmt"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/models"
)

func InsertConversation(conversationIN models.Conversation) (*models.Conversation, error) {
	var ret models.Conversation
	err := sqlite.Db.QueryRow("INSERT INTO conversations (USER1ID, USER2ID, GROUPID, CREATEDAT) VALUES(?, ?, ?, ?) RETURNING ID, USER1ID, USER2ID, GROUPID",
		conversationIN.User1ID, conversationIN.User2ID, conversationIN.GroupId, conversationIN.CreatedAt).Scan(&ret.Id, &ret.User1ID, &ret.User2ID, &ret.GroupId)
	if err != nil {
		fmt.Println("Inserting conversation failed: " + err.Error())
		return nil, err
	}
	return &ret, nil
}

func GetConversationBySenderReceiverIDs(User1ID int, User2ID int) (*models.Conversation, error) {
	var ret models.Conversation
	err := sqlite.Db.QueryRow("SELECT ID, USER1ID, USER2ID, GROUPID FROM conversations WHERE USER1ID in (?,?) AND USER2ID in (?,?)",
		User1ID, User2ID, User1ID, User2ID).Scan(&ret.Id, &ret.User1ID, &ret.User2ID, &ret.GroupId)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("getting conversation failed: " + err.Error())
			return nil, err
		}
		return nil, nil
	}
	return &ret, nil
}

func GetConversationByGroupId(groupId int) (*models.Conversation, error) {
	var ret models.Conversation
	err := sqlite.Db.QueryRow("SELECT ID, GROUPID FROM conversations WHERE GROUPID = ?",
		groupId).Scan(&ret.Id, &ret.GroupId)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("getting conversation failed: " + err.Error())
			return nil, err
		}
		return nil, nil
	}
	return &ret, nil
}

func GetConversationsByUserID(UserID int) ([]*models.Conversation, error) {
	rows, err := sqlite.Db.Query(`
	SELECT
		c.ID,
		c.USER1ID,
		c.USER2ID,
		c.GROUPID,
		m.MESSAGE,
		m.SENDERID,
		m.ISREAD,
		max(m.CREATEDAT)
	FROM messages AS m
	JOIN conversations as c on c.ID = m.CONVERSATIONID
	WHERE c.USER1ID = ? OR c.USER2ID = ? OR c.GROUPID IN (SELECT group_members.GROUPID FROM group_members WHERE MEMBERID = ?)
	GROUP BY m.CONVERSATIONID
	ORDER BY m.CREATEDAT DESC`, UserID, UserID, UserID)
	if err != nil {
		fmt.Println("Getting conversation by accountId failed: " + err.Error())
		return nil, err
	}
	var retAll []*models.Conversation
	defer rows.Close()
	for rows.Next() {
		var conversation models.Conversation
		if err := rows.Scan(&conversation.Id, &conversation.User1ID, &conversation.User2ID, &conversation.GroupId, &conversation.LastMessage.Message, &conversation.LastMessage.SenderId, &conversation.LastMessage.IsRead, &conversation.LastMessage.CreatedAt); err != nil {
			fmt.Println(err)
			continue
		}
		if conversation.GroupId != 0 {
			conversation.GroupObj, _ = GetGroupDataById(conversation.GroupId)
		} else {
			if conversation.User1ID == UserID {
				conversation.UserObj, _ = GetAccountDataById(conversation.User2ID)
			} else {
				conversation.UserObj, _ = GetAccountDataById(conversation.User1ID)
			}
		}

		conversation.LastMessage.IsRead = conversation.LastMessage.IsRead || conversation.LastMessage.SenderId == UserID

		retAll = append(retAll, &conversation)
	}
	return retAll, nil
}
