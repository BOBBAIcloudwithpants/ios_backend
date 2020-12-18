package service

import "github.com/bobbaicloudwithpants/ios_backend/models"

func GetFilesByPostID(post_id int) ([]models.ExtendedFile, error) {
	files, err := models.GetFilesByPostID(post_id)
	if err != nil {
		return nil, err
	}
	return files, nil
}
