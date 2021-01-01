package models

import (
	"errors"
	"fmt"
	"strconv"
)

type Forum struct {
	ForumID      int    `json:"forum_id"`
	ForumName    string `json:"forum_name"`
	IsPublic     bool   `json:"is_public"`
	Description  string `json:"description"`
	CreateAt     string `json:"create_at"`
	Cover        string `json:"cover"`
	PostNum      int    `json:"post_num"`
	SubscribeNum int    `json:"subscribe_num"`
	AdminList    []int  `json:"admin_list"`
}

// 包含 forum 信息的用户
type ForumUser struct {
	User
	Role    string `json:"role"`
	ForumID int    `json:"forum_id"`
}

// 将数据库查询结果转换为Forum
func convertMapToForum(forum map[string]string) Forum {
	forum_id, _ := strconv.Atoi(forum["forum_id"])
	is_public := true
	if forum["is_public"] == "0" {
		is_public = false
	}
	return Forum{
		ForumID:     forum_id,
		ForumName:   forum["forum_name"],
		IsPublic:    is_public,
		Description: forum["description"],
		CreateAt:    forum["create_at"],
		Cover:       forum["cover"],
	}
}

func convertMapToForumUser(forumUser map[string]string) ForumUser {
	forum_id, _ := strconv.Atoi(forumUser["forum_id"])
	return ForumUser{
		ForumID: forum_id,
		Role:    forumUser["role"],
		User:    convertMapToUser(forumUser),
	}
}

// 创建论坛
func CreateForum(forum Forum) (int64, error) {
	sentence := "INSERT INTO forum(forum_name, is_public, description, cover) VALUES(?, ?, ?, ?)"
	id, err := Execute(sentence, forum.ForumName, forum.IsPublic, forum.Description, forum.Cover)
	if err != nil {
		return -1, err
	}
	return id, nil
}

// 往论坛写入封面的url
func UpdateCover(path string, forumID int) error {
	sentence := "UPDATE forum SET cover=? WHERE forum_id=?"
	_, err := Execute(sentence, path, forumID)
	return err
}

// 获取论坛封面的url
func GetCoverURL(forumID int) error {
	sentence := "SELECT cover FROM forum WHERE user_id=?"
	_, err := Execute(sentence, forumID)
	return err
}

// 获取所有公开论坛，用于未登陆时获取公开论坛
func GetAllPublicForums() ([]Forum, error) {
	var ret []Forum

	res, err := QueryRows("SELECT forum_id, forum_name, is_public, description, create_at, cover FROM forum WHERE is_public = 1 ORDER BY create_at DESC")

	if err != nil {
		return nil, err
	}

	for _, r := range res {
		// fmt.Println("forum_id is ", r["forum_id"])
		// 拿到post数量
		forum_id, _ := strconv.Atoi(r["forum_id"])
		postNum, err := GetPostNumInForum(forum_id)
		if err != nil {
			fmt.Println(err)
		}
		subscribeNUm, err := GetSubscribeNumInForum(forum_id)
		if err != nil {
			fmt.Println(err)
		}
		adminList, err := GetAllAdminsInForum(forum_id)
		if err != nil {
			fmt.Println(err)
		}
		forum := convertMapToForum(r)
		forum.PostNum = postNum
		forum.SubscribeNum = subscribeNUm
		forum.AdminList = adminList
		ret = append(ret, forum)
	}

	return ret, nil
}

// 将某用户"添加为"某论坛的"user/owner/admin"
func AddRoleInForum(forum_id, user_id int, role string) error {
	sentence := "INSERT INTO forum_user(forum_id, user_id, role) VALUES(?, ?, ?)"
	_, err := Execute(sentence, forum_id, user_id, role)
	return err
}

// 查找某用户在某论坛中的角色
func FindRoleInForum(forum_id, user_id int) (string, error) {
	res, err := QueryRows("SELECT role FROM forum_user WHERE forum_id=? AND user_id=?", forum_id, user_id)
	if err != nil {
		return "", err
	}
	if len(res) == 0 {
		return "", errors.New("此用户不在该论坛下")
	}
	role := res[0]["role"]
	// fmt.Println("role is", role)
	return role, nil
}

// 将某用户在某论坛的角色修改为role
func UpdateRoleInForum(forum_id, user_id int, role string) error {
	sentence := "UPDATE forum_user SET role=? WHERE forum_id=? AND user_id=?"
	_, err := Execute(sentence, role, forum_id, user_id)
	return err
}

// 将某用户在论坛中删除
func DeleteRoleInForum(forum_id, user_id int) error {
	sentence := "DELETE FROM forum_user WHERE forum_id=? AND user_id=?"
	_, err := Execute(sentence, forum_id, user_id)
	return err
}

// 根据论坛的 ID 获取论坛信息
func GetForumByID(forum_id int) ([]Forum, error) {
	var ret []Forum

	res, err := QueryRows("SELECT * FROM forum WHERE forum_id=?", forum_id)
	if err != nil {
		return nil, err
	}

	for _, r := range res {
		postNum, err := GetPostNumInForum(forum_id)
		if err != nil {
			fmt.Println(err)
		}
		subscribeNUm, err := GetSubscribeNumInForum(forum_id)
		if err != nil {
			fmt.Println(err)
		}
		adminList, err := GetAllAdminsInForum(forum_id)
		if err != nil {
			fmt.Println(err)
		}
		forum := convertMapToForum(r)
		forum.PostNum = postNum
		forum.SubscribeNum = subscribeNUm
		forum.AdminList = adminList
		ret = append(ret, forum)
	}

	return ret, nil
}

// 根据论坛id获取其成员信息/关注者信息
func GetForumUserByForumID(forum_id int) ([]ForumUser, error) {
	var users []ForumUser

	sentence :=
		`
SELECT user.user_id, user.username, user.email, user.is_admin, user.create_at, user.avatar, forum_user.role
			FROM user INNER JOIN forum_user
            ON user.user_id = forum_user.user_id
			WHERE forum_user.forum_id = ?;
		`
	res, err := QueryRows(sentence, forum_id)
	if err != nil {
		return nil, err
	}
	for _, u := range res {
		users = append(users, convertMapToForumUser(u))
	}

	return users, nil
}

func GetPostNumInForum(forum_id int) (int, error) {
	sentence := "SELECT COUNT(*) AS num FROM post WHERE forum_id=?"
	res, err := QueryRows(sentence, forum_id)
	if err != nil {
		return -1, err
	}
	num, err := strconv.Atoi(res[0]["num"])
	if err != nil {
		return -1, err
	}
	// fmt.Println("postnum ", num)
	return num, nil
}

func GetSubscribeNumInForum(forum_id int) (int, error) {
	sentence := "SELECT COUNT(*) AS num FROM forum_user WHERE forum_id=?"
	res, err := QueryRows(sentence, forum_id)
	if err != nil {
		return -1, err
	}
	num, err := strconv.Atoi(res[0]["num"])
	if err != nil {
		return -1, err
	}
	// fmt.Println("subscribe num ", num)
	return num, nil
}

func GetAllAdminsInForum(forum_id int) ([]int, error) {
	sentence := "SELECT user_id FROM forum_user WHERE forum_id=? AND role<>'user' ORDER BY role DESC"
	res, err := QueryRows(sentence, forum_id)
	if err != nil {
		return nil, err
	}
	var adminList []int
	for _, r := range res {
		user_id, _ := strconv.Atoi(r["user_id"])
		adminList = append(adminList, user_id)
	}
	return adminList, nil
}

// 根据用户id获取公开及其相关论坛