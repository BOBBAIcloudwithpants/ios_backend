package models

import (
	"strconv"
)

type Comment struct {
	CommentID int    `json:"comment_id"`
	PostID    int    `json:"post_id"`
	UserID    int    `json:"user_id"`
	UserName  string `json:"username"`
	Content   string `json:"content"`
	CreateAt  string `json:"create_at"`
}

type CommentDetail struct {
	Comment
}

// 将数据库查询结果转换为 Comment
func convertMapToComment(comment map[string]string) Comment {
	comment_id, _ := strconv.Atoi(comment["comment_id"])
	post_id, _ := strconv.Atoi(comment["post_id"])
	user_id, _ := strconv.Atoi(comment["user_id"])
	return Comment{
		CommentID: comment_id,
		PostID:    post_id,
		UserID:    user_id,
		UserName:  comment["username"],
		Content:   comment["content"],
		CreateAt:  comment["create_at"],
	}
}

// 在给定的PostID中创建一个评论Comment
func CreateComment(comment Comment) (int64, error) {
	sentence := "INSERT INTO comment(post_id, user_id, content) VALUES (?, ?, ?)"
	return Execute(sentence, comment.PostID, comment.UserID, comment.Content)
}

// 获取给定的 PostID 下的全部 comments
func GetAllCommentsByPostID(post_id int) ([]Comment, error) {
	var ret []Comment
	res, err := QueryRows("SELECT comment_id, post_id, comment.user_id, content, comment.create_at, comment.like, username FROM comment LEFT JOIN user ON comment.user_id=user.user_id WHERE post_id=?", post_id)
	if err != nil {
		return ret, err
	}

	for _, p := range res {
		ret = append(ret, convertMapToComment(p))
	}
	return ret, nil
}

// 根据id获取某个 Comment
func GetOneCommentByCommentID(comment_id int) ([]Comment, error) {
	var ret []Comment
	res, err := QueryRows("SELECT comment_id, post_id, comment.user_id, content, comment.create_at, comment.like, username FROM comment LEFT JOIN user ON comment.user_id=user.user_id WHERE comment_id=?", comment_id)
	if err != nil {
		return ret, err
	}

	for _, p := range res {
		ret = append(ret, convertMapToComment(p))
	}
	return ret, nil
}
