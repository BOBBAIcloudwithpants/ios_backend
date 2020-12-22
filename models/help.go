package models

import (
	"strconv"

	"github.com/pkg/errors"
)

type Help struct {
	HelpID int			`json:"help_id"`
	ForumID int			`json:"forum_id"`
	UserID int			`json:"user_id"`
	Title string		`json:"title"`
	Content string		`json:"content"`
	CreateAt string		`json:"create_at"`
	Point int			`json:"point"`
	IsFinished bool		`json:"is_finished"`
	HelperID int		`json:"helper_id"`
	Filename string		`json:"filename"`
}

type UnfinishedHelpDetail struct {
	Help
	Creator User `json:"creator"`
}

type PendingOrFinishedHelpDetail struct {
	Help
	Creator User	`json:"creator"`
	Helper  User	`json:"helper"`
}

func convertMapToHelp(helpMap map[string]string) Help {
	help_id, _ := strconv.Atoi(helpMap["help_id"])
	forum_id, _ := strconv.Atoi(helpMap["forum_id"])
	is_finished, _ := strconv.Atoi(helpMap["is_finished"])
	help_point, _ := strconv.Atoi(helpMap["point"])
	is_finish := false
	user_id, _ := strconv.Atoi(helpMap["user_id"])
	helper_id, _ := strconv.Atoi(helpMap["helper_id"])
	if is_finished == 1 {
		is_finish = true
	}
	return Help{
		HelpID: help_id,
		ForumID: forum_id,
		UserID: user_id,
		Title: helpMap["title"],
		Content: helpMap["content"],
		CreateAt: helpMap["create_at"],
		Point: help_point,
		HelperID: helper_id,
		Filename: helpMap["filename"],
		IsFinished: is_finish,
	}
}

func convertMapToPendingOrFinishedHelpDetail (helpMap map[string]string) PendingOrFinishedHelpDetail {
	help_id, _ := strconv.Atoi(helpMap["help_id"])
	forum_id, _ := strconv.Atoi(helpMap["forum_id"])
	creator_id, _ := strconv.Atoi(helpMap["creator_id"])
	helper_id, _ := strconv.Atoi(helpMap["helper_id"])
	is_finished, _ := strconv.Atoi(helpMap["is_finished"])
	help_point, _ := strconv.Atoi(helpMap["help_point"])
	creator_point, _ := strconv.Atoi(helpMap["creator_point"])
	helper_point, _ := strconv.Atoi(helpMap["helper_point"])

	is_finish := false
	if is_finished == 1 {
		is_finish = true
	}
	creator := User{
		UserId: creator_id,
		Username: helpMap["creator_username"],
		Email: helpMap["creator_email"],
		Point: creator_point,
	}

	helper := User{
		UserId: helper_id,
		Username: helpMap["helper_username"],
		Email: helpMap["helper_email"],
		Point: helper_point,
	}
	help := Help{
		ForumID: forum_id,
		HelperID: help_id,
		UserID: creator_id,
		Title: helpMap["title"],
		Content: helpMap["content"],
		CreateAt: helpMap["create_at"],
		Point: help_point,
		IsFinished: is_finish,
		HelpID: helper_id,
		Filename: helpMap["filename"],
	}

	return PendingOrFinishedHelpDetail{
		Creator: creator,
		Helper: helper,
		Help: help,
	}

}

func convertMapToUnfinishedHelpDetail (helpMap map[string]string) UnfinishedHelpDetail {
	help_id, _ := strconv.Atoi(helpMap["help_id"])
	forum_id, _ := strconv.Atoi(helpMap["forum_id"])
	user_id, _ := strconv.Atoi(helpMap["user_id"])
	help_point, _ := strconv.Atoi(helpMap["help_point"])
	creator_point, _ := strconv.Atoi(helpMap["creator_point"])
	is_finished, _ := strconv.Atoi(helpMap["is_finished"])
	is_finish := false
	if is_finished == 1 {
		is_finish = true
	}
	helper_id, _ := strconv.Atoi(helpMap["helper_id"])

	user := User{
		UserId: user_id,
		Username: helpMap["username"],
		Email: helpMap["email"],
		Point: creator_point,
	}
	help := Help{
		ForumID: forum_id,
		HelperID: help_id,
		UserID: user_id,
		Title: helpMap["title"],
		Content: helpMap["content"],
		CreateAt: helpMap["create_at"],
		Point: help_point,
		IsFinished: is_finish,
		HelpID: helper_id,
		Filename: helpMap["filename"],
	}
	return UnfinishedHelpDetail{
		Help:    help,
		Creator: user,
	}
}

// 创建 help
func CreateHelp(forum_id int, user_id int, title string, content string, bonus int, filename string) (int64, error) {
	sql :=
		`
			INSERT INTO help(forum_id, user_id, title, content, point, filename) VALUES(?,?,?,?,?,?);
		`
	return Execute(sql, forum_id, user_id, title, content, bonus, filename)
}

// 根据论坛ID查询等待应答的 Help
func GetUnfinishedHelpsByForumID(forum_id int) ([]UnfinishedHelpDetail, error){
	var ret []UnfinishedHelpDetail
	sql :=
		`
			SELECT 
				help_id,username,email, forum_id, user.user_id, title, content, help.create_at, help.point AS help_point,user.point AS creator_point, is_finished, helper_id, filename
			From
				help INNER JOIN user ON help.user_id = user.user_id
			WHERE
				help.forum_id = ?;
		`
	res, err := QueryRows(sql, forum_id)
	if err != nil {
		return ret, err
	}

	for _, val := range res {
		ret = append(ret, convertMapToUnfinishedHelpDetail(val))
	}

	return ret, nil
}


// 根据论坛 ID 查询正在pending的Help(pending的意思是已经有别人应答，但是没有得到确认)
func GetPendingHelpsByForumID(forum_id int)([]PendingOrFinishedHelpDetail, error) {
	var ret []PendingOrFinishedHelpDetail
	sql :=
		`
			SELECT 
				help_id, helper.user_id AS helper_id,
                         helper.username AS helper_username, 
						helper.point AS helper_point,
						creator.user_id AS creator_id, 
						creator.username AS creator_username,
						creator.point AS creator_point,
						help_id, forum_id, title, content, help.create_at, 
						help.point AS help_point,
						is_finished,
						filename
			FROM
				help 
                INNER JOIN user AS helper ON help.helper_id = helper.user_id
				INNER JOIN user AS creator ON help.user_id = creator.user_id
			WHERE forum_id =? AND is_finished = 0;
		`
	res, err := QueryRows(sql, forum_id)
	if err != nil {
		return ret, err
	}

	for _, val := range res {
		ret = append(ret, convertMapToPendingOrFinishedHelpDetail(val))
	}

	return ret, nil
}

// 根据论坛 ID 查询已经完成的helps
func GetFinishedHelpsByForumID(forum_id int)([]PendingOrFinishedHelpDetail, error) {
	var ret []PendingOrFinishedHelpDetail
	sql :=
		`
			SELECT 
				help_id, helper.user_id AS helper_id,
                         helper.username AS helper_username, 
						helper.point AS helper_point,
						creator.user_id AS creator_id, 
						creator.username AS creator_username,
						creator.point AS creator_point,
						help_id, forum_id, title, content, help.create_at, 
						help.point AS help_point,
						is_finished,
						filename
			FROM
				help 
                INNER JOIN user AS helper ON help.helper_id = helper.user_id
				INNER JOIN user AS creator ON help.user_id = creator.user_id
			WHERE forum_id =? AND is_finished = 1;
		`
	res, err := QueryRows(sql, forum_id)
	if err != nil {
		return ret, err
	}

	for _, val := range res {
		ret = append(ret, convertMapToPendingOrFinishedHelpDetail(val))
	}

	return ret, nil
}

// 根据iD应答某一个help
func AnswerHelpByHelpIDAndUserID(helpID int, userID int) error {
	sql :=
		`
			UPDATE help SET help.helper_id=? WHERE help.help_id;
		`
	_, err := Execute(sql, userID, helpID)
	return err
}

func GetHelpByHelpID(helpID int)(Help, error) {
	var ret Help
	sql :=
		`
			SELECT * from help WHERE help_id = ?;
		`
	res, err := QueryRows(sql, helpID)
	if err != nil {
		return ret, err
	}
	if len(res) == 0 {
		return ret, errors.New("该 help 不存在")
	}
	ret = convertMapToHelp(res[0])
	return ret, nil

}


func FinishHelpByHelpID(helpID int) error {
	sql :=
		`
			UPDATE help SET help.is_finished = 1 WHERE help.help_id = ?;
		`
	_, err := Execute(sql, helpID)
	if err != nil {
		return err
	}

	help, err := GetHelpByHelpID(helpID)
	if err != nil {
		return err
	}

	updateCreatorPointSql :=
		`
			UPDATE user SET user.point = user.point - ? WHERE user_id = ?
		`
	_, err = Execute(updateCreatorPointSql, help.Point, help.UserID)
	if err != nil {
		return err
	}

	updateHelperPointSql :=
		`
			UPDATE user SET user.point = user.point + ? WHERE user_id = ?
		`
	_, err = Execute(updateHelperPointSql, help.Point, help.ForumID)
	if err != nil {
		return err
	}

	return nil
}




