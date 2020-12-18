package models

import (
	"fmt"
	"strconv"
)

type User struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
	Avatar   string `json:"avatar"`
	CreateAt string `json:"create_at"`
}

type SubscribeList struct {
	ParticipateList []int `json:"participate_list"`
	FocusList       []int `json:"star_list"`
}

type UserDetail struct {
	User
	SubscribeList
	LikeList []int `json:"like_list"`
}

// 将数据库查询的结果转换为 User
func convertMapToUser(user map[string]string) User {
	user_id, _ := strconv.Atoi(user["user_id"])
	is_admin := false
	if user["is_admin"] == "1" {
		is_admin = true
	}

	return User{UserId: user_id, Username: user["username"], Email: user["email"], Password: user["password"], IsAdmin: is_admin, Avatar: user["avatar"], CreateAt: user["create_at"]}
}

// 创建用户
func CreateUser(user User) error {
	sentence := "INSERT INTO user(username, password, email ,is_admin, avatar) VALUES(?, ?, ?, ?, ?)"
	_, err := Execute(sentence, user.Username, user.Password, user.Email, user.IsAdmin, user.Avatar)
	return err
}

// 根据用户id获取用户
func GetUserById(user_id int) ([]User, error) {
	var ret []User

	res, err := QueryRows("SELECT user_id, username, password, email ,is_admin, create_at, avatar FROM user WHERE user_id = ?", user_id)

	if err != nil {
		return nil, err
	}

	for _, r := range res {
		ret = append(ret, convertMapToUser(r))
	}

	return ret, nil
}

// 根据用户名获取用户
func GetUserByUsername(username string) ([]User, error) {
	var ret []User

	res, err := QueryRows("SELECT user_id, username, password, email, is_admin, create_at, avatar FROM user WHERE username = ?", username)

	if err != nil {
		return nil, err
	}

	for _, r := range res {
		ret = append(ret, convertMapToUser(r))
	}

	return ret, err
}

func GetUserByEmail(email string) ([]User, error) {
	var ret []User

	res, err := QueryRows("SELECT user_id, username, password, email,is_admin, create_at, avatar FROM user WHERE email = ?", email)

	if err != nil {
		return nil, err
	}

	for _, r := range res {
		ret = append(ret, convertMapToUser(r))
	}

	return ret, err
}

func GetAllUsers() ([]User, error) {
	var ret []User

	res, err := QueryRows("SELECT user_id, username, password, email,is_admin, create_at, avatar FROM user")

	if err != nil {
		return nil, err
	}

	for _, r := range res {
		ret = append(ret, convertMapToUser(r))
	}

	return ret, err
}

func GetAllUsersContains(str string) ([]User, error) {
	var ret []User
	query := fmt.Sprintf("SELECT user_id, username, password, email,is_admin, create_at, avatar FROM user WHERE username LIKE %s", strconv.Quote("%"+str+"%"))
	res, err := QueryRows(query)

	if err != nil {
		return nil, err
	}

	for _, r := range res {
		ret = append(ret, convertMapToUser(r))
	}

	return ret, err
}

// 根据用户id获取某个用户信息以及所参与的/关注的列表
func GetOneUserSubscribe(userID int) (SubscribeList, error) {
	var ret SubscribeList
	sql :=
		`
		SELECT forum.is_public, forum.forum_id FROM forum
			INNER JOIN forum_user ON forum.forum_id = forum_user.forum_id
			WHERE forum_user.user_id = ?
		`
	res, err := QueryRows(sql, userID)
	if err != nil {
		return ret, err
	}

	for _, val := range res {
		is_public, _ := strconv.Atoi(val["is_public"])
		forum_id, _ := strconv.Atoi(val["forum_id"])
		if is_public == 1 {
			ret.FocusList = append(ret.FocusList, forum_id)
		} else {
			ret.ParticipateList = append(ret.ParticipateList, forum_id)
		}

	}
	return ret, nil
}

func UpdateUserAvatarByUserId(userID int, avatar_path string) error {
	sql :=
		`
		UPDATE user SET avatar=? WHERE user_id=?
		`
	_, err := Execute(sql, avatar_path, userID)
	return err
}

func GetOneUserLikeListByUserID(userID int) ([]int, error) {
	sql :=
		`
		SELECT post_id FROM post_like WHERE user_id = ?;
		`
	data, err := QueryRows(sql, userID)
	if err != nil {
		return nil, err
	}
	var ret []int
	for _, val := range data {
		id, err := strconv.Atoi(val["post_id"])
		if err != nil {
			return nil, err
		}
		ret = append(ret, id)
	}
	return ret, nil
}

