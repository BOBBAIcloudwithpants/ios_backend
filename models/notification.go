package models

import "strconv"

type Notification struct {
	NotifID int			`json:"notif_id"`
	SenderID int		`json:"sender_id"`
	ReceiverID int		`json:"receiver_id"`
	Content string		`json:"content"`
	IsRead bool			`json:"is_read"`
}


type NotificationDetail struct {
	Notification
	Sender User			`json:"sender"`
	Receiver User		`json:"receiver"`
}

func convertMapToNotificationDetail(notification map[string]string) NotificationDetail {
	notif_id, _ := strconv.Atoi(notification["notif_id"])
	sender_id, _ := strconv.Atoi(notification["sender_id"])
	receiver_id, _ := strconv.Atoi(notification["receiver_id"])
	is_r, _ := strconv.Atoi(notification["is_read"])
	is_read := false
	if is_r == 1 {
		is_read = true
	}

	sender := User{
		UserId: sender_id,
		Username: notification["sender_username"],
		Email: notification["sender_email"],
	}
	receiver:= User{
		UserId: receiver_id,
		Username: notification["receiver_username"],
		Email: notification["receiver_email"],
	}

	notif := Notification{
		NotifID: notif_id,
		SenderID: sender_id,
		ReceiverID: receiver_id,
		Content: notification["content"],
		IsRead: is_read,
	}

	return NotificationDetail{
		Notification: notif,
		Sender: sender,
		Receiver: receiver,
	}
}

func GetUnreadNotificationByUserID(userID int) ([]NotificationDetail, error) {
	var ret []NotificationDetail
	sql :=
		`
			SELECT 
				sender.username AS sender_username,
				sender.email AS sender_email,
				sender.user_id AS sender_id,
				receiver.username AS receiver_username,
				receiver.email AS receiver_email,
				receiver.user_id AS receiver_id,
				notif_id, content, is_read
			FROM
				notification INNER JOIN user AS sender ON notification.sender_id = sender.user_id
  							 INNER JOIN user AS receiver ON notification.receiver_id = receiver.user_id
				WHERE receiver.user_id = ? AND is_read = 0;
		`
	res, err := QueryRows(sql, userID)
	if err != nil {
		return ret, err
	}
	for _, val := range res {
		ret = append(ret, convertMapToNotificationDetail(val))
	}
	return ret, nil
}


func CreateNotification(senderID int, receiverID int, content string) error {
	sql :=
		`
			INSERT INTO notification (sender_id, receiver_id, content) VALUES(?, ?, ?)
		`
	_, err := Execute(sql, senderID, receiverID, content)
	return err
}


func ReadNotificationByNotifID(notifID int) error {
	sql :=
		`
			UPDATE notification SET is_read = 1 WHERE notif_id = ?
		`

	_, err := Execute(sql, notifID)
	return err
}






