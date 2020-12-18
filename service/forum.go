package service

import (
	"github.com/bobbaicloudwithpants/ios_backend/models"
)

// 创建论坛
func CreateForum(forumName, description string, isPublic bool) (int64, error) {
	forum := models.Forum{
		ForumName:   forumName,
		IsPublic:    isPublic,
		Description: description,
		Cover:       "cover.png",
	}
	return models.CreateForum(forum)
}

// 查找当前用户是否是论坛成员
// 如果该论坛是公开的，则直接返回true, 否则查看论坛的成员列表中是否有该用户
func IsUserInForum(user_id int, forum_id int) (bool, error) {
	forums, err := models.GetForumByID(forum_id)

	if err != nil {
		return false, err
	}

	// 证明论坛不存在
	if len(forums) == 0 {
		return false, nil
	}

	forum := forums[0]

	// 论坛公开，可以访问
	if forum.IsPublic {
		return true, nil
	}

	forumUsers, err := models.GetForumUserByForumID(forum_id)
	if err != nil {
		return false, err
	}

	for _, fu := range forumUsers {
		if user_id == fu.UserId {
			return true, nil
		}
	}
	return false, nil
}
