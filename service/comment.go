package service

import (
	"github.com/pkg/errors"
	"github.com/bobbaicloudwithpants/ios_backend/models"
)

// 根据 post_id 获取全部comment的详情
func GetAllCommentsByPostID(post_id int) ([]models.Comment, error) {
	return models.GetAllCommentsByPostID(post_id)
}

// 根据 comment_id 获取一个comment的详情
func GetOneCommentDetailByCommentID(comment_id int) ([]models.CommentDetail, error) {
	var commentDetails []models.CommentDetail

	comments, err := models.GetOneCommentByCommentID(comment_id)
	if err != nil {
		return nil, err
	}
	if len(comments) == 0 {
		return nil, errors.New("该comment_id对应的comment不存在")
	}
	comment := comments[0]

	commentDetail := models.CommentDetail{Comment: comment}
	commentDetails = []models.CommentDetail{commentDetail}
	return commentDetails, nil
}
