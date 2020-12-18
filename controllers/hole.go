package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pingcap/log"
	"github.com/bobbaicloudwithpants/ios_backend/models"
	"github.com/bobbaicloudwithpants/ios_backend/service"
)

// 用户创建 HOLE
// CreateHole godoc
// @Summary CreateHole
// @Description	CreateHole
// @Tags Holes
// @Accept	mpfd
// @Produce	json
// @Param token header string true "将token放在请求头部的‘Authorization‘字段中，并以‘Bearer ‘开头""
// @Param title formData string true "Hole 的标题"
// @Param content formData string true "Hole 的内容"
// @Success 200 {object} responses.StatusOKResponse "创建 hole 成功"
// @Failure 400 {object} responses.StatusBadRequestResponse "树洞的标题或者内容不得为空"
// @Failure 500 {object} responses.StatusInternalServerError "插入用户创建的hole失败"
// @Router /forums/{forum_id}/holes [post]
func CreateHole(c *gin.Context) {
	log.Info("user create hole")
	var data interface{}
	forum_id, _ := strconv.Atoi(c.Param("forum_id"))
	title, content := c.PostForm("title"), c.PostForm("content")
	user_id := service.GetUserFromContext(c).UserId

	if title == "" || content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "树洞的标题或者内容不得为空", "data": data})
		return
	}

	_, err := models.CreateHole(models.Hole{ForumID: forum_id, UserID: user_id, Title: title, Content: content}) //HoleID与CreateAt由数据库生成
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "插入用户创建的hole失败", "data": data})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建 hole 成功", "data": data})
	return
}

// 获取某个 forum 下的所有hole
// GetAllHolesByForumID godoc
// @Summary GetAllHolesByForumID
// @Description GetAllHolesByForumID
// @Tags Holes
// @Accept json
// @Produce json
// @Param token header string true "将token放在请求头部的‘Authorization‘字段中，并以‘Bearer ‘开头""
// @Success 200 {object} responses.StatusOKResponse{data=[]models.Hole} "获取树洞帖子成功"
// @Failure 500 {object} responses.StatusInternalServerError "查询数据库出现异常"
// @Router /forums/{forum_id}/holes [get]
func GetAllHolesByForumID(c *gin.Context) {
	log.Info("get all holes by forum_id controller")
	forum_id, _ := strconv.Atoi(c.Param("forum_id"))

	data, err := service.GetAllHolesByForumID(forum_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询数据库出现异常" + err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  fmt.Sprintf("获取论坛 %d 下的全部树洞帖子成功", forum_id),
		"data": data,
	})
}

// 根据 id 获取某个hole的详情
// GetOneHoleDetailByHoleID godoc
// @Summary GetOneHoleDetailByHoleID
// @Description GetOneHoleDetailByHoleID
// @Tags Holes
// @Accept json
// @Produce json
// @Param token header string true "将token放在请求头部的‘Authorization‘字段中，并以‘Bearer ‘开头""
// @Success 200 {object} responses.StatusOKResponse{data=[]models.HoleDetail} "获取树洞帖子成功"
// @Failure 400 {object} responses.StatusInternalServerError "数据库查询异常，或者该hole不存在"
// @Router /forums/{forum_id}/holes/{hole_id} [get]
func GetOneHoleDetailByHoleID(c *gin.Context) {
	log.Info("get one hole detail by hole_id")
	hole_id, _ := strconv.Atoi(c.Param("hole_id"))

	data, err := service.GetOneHoleDetailByHoleID(hole_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 403, "msg": "数据库查询异常，或者该hole不存在：" + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  fmt.Sprintf("获取第 %d 号树洞帖子成功", hole_id),
		"data": data,
	})
}
