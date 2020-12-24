package controllers

import (
	"net/http"

	"github.com/bobbaicloudwithpants/ios_backend/models"
	"github.com/bobbaicloudwithpants/ios_backend/service"
	"github.com/gin-gonic/gin"
	"github.com/pingcap/log"
)


// GetAllUnreadNotification godoc
// @Summary GetAllUnreadNotification
// @Description GetAllUnreadNotification
// @Tags Forums
// @Accept  json
// @Produce  json
// @Success 200 {object} responses.StatusOKResponse{data=[]models.NotificationDetail}
// @Failure 500 {object} responses.StatusInternalServerError "服务器错误"
// @Router /notifications [get]
func GetAllUnreadNotification(c *gin.Context) {
	log.Info("get all unread notification")
	user_id := service.GetUserFromContext(c).UserId

	ret, err := models.GetUnreadNotificationByUserID(user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询数据库 "+err.Error(), "data": ret})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "查询所有未读通知成功", "data": ret})
}