package service

import "github.com/bobbaicloudwithpants/ios_backend/models"

func GetPopularPosts() ([] models.PostDetail, error){
	var postDetails []models.PostDetail

	posts, err := models.GetPopularPosts()
	if err != nil {
		return postDetails, err
	}
	for _, post := range posts {
		files, err := models.GetFilesByPostID(post.PostID)
		if err != nil {
			return nil, err
		}
		postDetail := models.PostDetail{Files: files, Post: post}
		postDetails = append(postDetails, postDetail)
	}
	return postDetails, nil
}