package dbfunctions

import (
	"database/sql"
	"fmt"
	"social-network/pkg/db/sqlite"
	"social-network/pkg/models"
)

func InsertMessage(messageIN models.Message) (*models.Message, error) {
	var ret models.Message
	err := sqlite.Db.QueryRow("INSERT INTO messages (conversationID, senderID, message, createdAt, isRead) VALUES(?, ?, ?, ?, ?) RETURNING ID, conversationID, senderID, message, CreatedAt, isRead",
		messageIN.ConversationId, messageIN.SenderId, messageIN.Message, messageIN.CreatedAt, messageIN.IsRead).Scan(&ret.Id, &ret.ConversationId, &ret.SenderId, &ret.Message, &ret.CreatedAt, &ret.IsRead)
	if err != nil {
		fmt.Println("Inserting message failed: " + err.Error())
		return nil, err
	}
	return &ret, err
}

func GetMessageByConversation(conversationID, userId int) ([]*models.Message, error) {
	rows, err := sqlite.Db.Query("SELECT m.ID, m.conversationID, m.senderID, m.message, m.createdAt, (a.FIRSTNAME || ' ' || a.LASTNAME), m.isread, a.AVATARPATH FROM messages AS m JOIN users AS a ON m.senderID = a.ID WHERE m.conversationID = ? ORDER BY m.ID DESC", conversationID)
	if err != nil {
		fmt.Println("Getting messages By conversationID failed: " + err.Error())
		return nil, err
	}
	retAll := []*models.Message{}
	defer rows.Close()
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.Id, &message.ConversationId, &message.SenderId, &message.Message, &message.CreatedAt, &message.SenderName, &message.IsRead, &message.SenderObj.AvatarPath); err != nil {
			fmt.Println(err)
			continue
		}

		message.IsRead = message.IsRead || message.SenderId == userId

		retAll = append(retAll, &message)
	}
	return retAll, nil
}

func GetLastMessageByConversation(conversationID int) (*models.Message, error) {
	var ret models.Message
	err := sqlite.Db.QueryRow("SELECT m.ID, m.conversationID, m.senderID, (a.FIRSTNAME || ' ' || a.LASTNAME), m.message, m.createdAt, m.isread FROM messages AS m JOIN users AS a ON m.senderID = a.ID WHERE m.conversationID = ?  ORDER BY m.ID DESC", conversationID).Scan(&ret.Id, &ret.ConversationId, &ret.SenderId, &ret.SenderName, &ret.Message, &ret.CreatedAt, &ret.IsRead)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Get last message data failed: " + err.Error())
			return nil, err
		}
		return nil, nil
	}
	return &ret, nil
}

func SetMessageIsRead(messageId int, val bool) error {
	_, err := sqlite.Db.Exec("UPDATE messages SET ISREAD = ? WHERE ID = ?", val, messageId)
	if err != nil {
		fmt.Println("Failed to set isRead status: " + err.Error())
	}
	return err
}
