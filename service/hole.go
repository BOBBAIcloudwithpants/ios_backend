package service

import (
	"github.com/pkg/errors"
	"github.com/bobbaicloudwithpants/ios_backend/models"
)

func GetAllHolesByForumID(forum_id int) ([]models.Hole, error) {
	return models.GetAllHolesByForumID(forum_id)
}

// 根据 hole_id 获取一个hole的详情
func GetOneHoleDetailByHoleID(hole_id int) ([]models.HoleDetail, error) {
	var holeDetails []models.HoleDetail

	holes, err := models.GetOneHoleByHoleID(hole_id)
	if err != nil {
		return nil, err
	}
	if len(holes) == 0 {
		return nil, errors.New("该hole_id对应的hole不存在")
	}
	hole := holes[0]

	holeDetail := models.HoleDetail{Hole: hole}
	holeDetails = []models.HoleDetail{holeDetail}
	return holeDetails, nil
}
