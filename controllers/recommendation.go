package controllers

import (
	"net/http"

	"github.com/bobbaicloudwithpants/ios_backend/models"
	"github.com/bobbaicloudwithpants/ios_backend/service"
	"github.com/gin-gonic/gin"
)

// 获取给用户的推荐
// GetRecommendations godoc
// @Summary GetRecommendations
// @Description GetRecommendations
// @Tags Posts
// @Accept json
// @Produce json
// @Param token header string true "将token放在请求头部的‘Authorization‘字段中，并以‘Bearer ‘开头""
// @Success 200 {object} responses.StatusOKResponse{data=[]models.PostDetail}
// @Failure 400 {object} responses.StatusInternalServerError "数据库查询异常"
// @Router /recommendation [get]
func GetRecommendations(c *gin.Context) {
	var ret []models.PostDetail
	postDetails, err := service.GetPopularPosts()
	if err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "数据库查询出错 "+err.Error(), "data": ret})
	    return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取推荐成功", "data": postDetails})
	return
}


func GetRecomCover(c *gin.Context) {
	ret := [3]string{"movie.jpg", "programming.jpg", "rock.jpg"}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取封面照片名成功", "data": ret})
	return
}
