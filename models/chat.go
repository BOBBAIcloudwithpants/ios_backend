package models

import "strconv"

type Chat struct {
	ChatID int			`json:"chat_id"`
	User1 User			`json:"user1"`
	User2 User			`json:"user2"`
	LastModified string	`json:"last_modified"`
}

type Message struct {
	MessageID int	`json:"message_id"`
	ChatID	int		`json:"chat_id"`
	SenderID int	`json:"sender_id"`
	ReceiverID int	`json:"receiver_id"`
	Content string`json:"content"`
	CreateAt string	`json:"create_at"`
	IsRead	bool	`json:"is_read"`
}

type ChatDetail struct {
	Chat
	Messages []Message	`json:"messages"`
}

func convertMapToMessage(message map[string]string) Message {
	message_id, _ := strconv.Atoi(message["message_id"])
	chat_id, _ := strconv.Atoi(message["chat_id"])
	sender_id, _ := strconv.Atoi(message["sender_id"])
	receiver_id, _ := strconv.Atoi(message["receiver_id"])

	is_read_int, _ := strconv.Atoi(message["is_read"])
	is_read := false
	if is_read_int > 0 {
		is_read = true
	}
	return Message{
		MessageID: message_id,
		ChatID: chat_id,
		IsRead: is_read,
		SenderID: sender_id,
		ReceiverID: receiver_id,
		CreateAt: message["create_at"],
		Content: message["content"],
	}
}


func convertMapToChat(chat map[string]string) Chat {
	chat_id, _ := strconv.Atoi(chat["chat_id"])
	user1_id, _ := strconv.Atoi(chat["user1_id"])
	user2_id, _ := strconv.Atoi(chat["user2_id"])
	return Chat{
		ChatID: chat_id,
		User1: User{
			UserId: user1_id,
			Username: chat["user1_name"],
			Email: chat["user1_email"],
		},
		User2: User{
			UserId: user2_id,
			Username: chat["user2_name"],
			Email: chat["user2_email"],
		},
		LastModified: chat["last_modified"],
	}
}

func GetChatsByUserID(user_id int)([]Chat, error) {
	var ret []Chat
	sql :=
		`
		SELECT 
				chat.chat_id,
				user1.user_id AS user1_id,
				user1.username AS user1_name,
				user1.email	AS user1_email,
				user2.user_id AS user2_id,
				user2.username AS user2_name,
				user2.email AS user2_email,
				chat.last_modified
		FROM chat
				INNER JOIN user AS user1 ON chat.user1_id = user1.user_id
				INNER JOIN user AS user2 ON chat.user2_id = user2.user_id
		WHERE
				chat.user1_id = ? OR chat.user2_id = ?
		`
	res, err := QueryRows(sql, user_id, user_id)
	if err != nil {
		return nil, err
	}
	for _, val := range res {
		ret = append(ret, convertMapToChat(val))
	}
	return ret, nil
}



func GetMessageByChatID(chat_id int)([]Message, error) {
	var ret []Message
	sql :=
		`
		SELECT
				message_id,
				chat_id,
				sender_id,
				receiver_id,
				create_at,
				content,
				is_read,
		FROM message
		WHERE chat_id = ?
		`
	res, err := QueryRows(sql, chat_id)
	if err != nil {
		return ret, err
	}
	for _, val := range res {
		ret = append(ret, convertMapToMessage(val))
	}
	return ret, nil
}

