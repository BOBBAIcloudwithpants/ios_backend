package controllers

import (
	"net/http"
	"path"
	"strconv"

	"github.com/bobbaicloudwithpants/ios_backend/models"
	"github.com/bobbaicloudwithpants/ios_backend/service"
	"github.com/gin-gonic/gin"
	"github.com/pingcap/log"
)

// 用户创建 HELP
// CreateHelp godoc
// @Summary CreateHelp
// @Description	CreateHelp
// @Tags Helps
// @Accept	mpfd
// @Produce	json
// @Param token header string true "将token放在请求头部的‘Authorization‘字段中，并以‘Bearer ‘开头""
// @Param title formData string true "Help 的标题"
// @Param content formData string true "Help 的内容"
// @Param point formData string true "Help 的悬赏点数"
// @Param file formData file true "一个多媒体文件"
// @Success 200 {object} responses.StatusOKResponse "创建 Help 成功"
// @Failure 403 {object} responses.StatusBadRequestResponse "标题或者内容不得为空"
// @Failure 403 {object} responses.StatusBadRequestResponse "您所上传的文件无法打开"
// @Failure 500 {object} responses.StatusInternalServerError "服务器错误"
// @Router /forums/{forum_id}/helps [post]
func CreateHelp(c *gin.Context) {
	log.Info("create help controller")
	forum_id, _ := strconv.Atoi(c.Param("forum_id"))
	form, _ := c.MultipartForm()
	file := form.File["file"][0]	// 仅为 录音 或者 视频
	title, content := c.PostForm("title"), c.PostForm("content")
	bonus, _ := strconv.Atoi(c.PostForm("point"))
	user_id := service.GetUserFromContext(c).UserId

	if title == "" || content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "标题或者内容不得为空"})
		return
	}
	var fileToBeUpload service.File

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "您所上传的文件无法打开"})
		return
	}
	fileToBeUpload = service.File{F:f, H:file}
	bucketName := "help"

	ext := path.Ext(file.Filename)
	name, err := service.FileUpload(fileToBeUpload.F, fileToBeUpload.H, bucketName, c.Request.URL.Path, ext)
	if err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "文件系统出错："+err.Error()})
	    return
	}
	_, err = models.CreateHelp(forum_id, user_id, title, content, bonus, name)

	if err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error()})
	    return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建 help 成功"})
	return
}


// 获取全部的没有人应答的Help
// GetAllUnfinishedHelpByForumID godoc
// @Summary GetAllUnfinishedHelpByForumID
// @Description	GetAllUnfinishedHelpByForumID
// @Tags Helps
// @Accept	mpfd
// @Produce	json
// @Param token header string true "将token放在请求头部的‘Authorization‘字段中，并以‘Bearer ‘开头""
// @Success 200 {object} responses.StatusOKResponse{data=[]models.UnfinishedHelpDetail}
// @Failure 500 {object} responses.StatusInternalServerError "服务器错误"
// @Router /forums/{forum_id}/helps/unfinished [get]
func GetAllUnfinishedHelpByForumID(c *gin.Context) {
	log.Info("get all unfinished help by forum ID controller")
	var unfinishedHelp []models.UnfinishedHelpDetail
	forum_id, _ := strconv.Atoi(c.Param("forum_id"))
	unfinishedHelp, err := models.GetUnfinishedHelpsByForumID(forum_id)
	if err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error(), "data": unfinishedHelp})
	    return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取所有的没有应答的 help 成功", "data": unfinishedHelp})
	return
}

// 获取所有的已经应答但是还没有完成的 help
// GetAllPendingHelpByForumID godoc
// @Summary GetAllPendingHelpByForumID
// @Description	GetAllPendingHelpByForumID
// @Tags Helps
// @Accept	mpfd
// @Produce	json
// @Param token header string true "将token放在请求头部的‘Authorization‘字段中，并以‘Bearer ‘开头""
// @Success 200 {object} responses.StatusOKResponse{data=[]models.PendingOrFinishedHelpDetail}
// @Failure 500 {object} responses.StatusInternalServerError "服务器错误"
// @Router /forums/{forum_id}/helps/pending [get]
func GetAllPendingHelpByForumID(c *gin.Context) {
	log.Info("get all pending helps by forum_id controller")
	var pendingHelps []models.PendingOrFinishedHelpDetail
	forum_id, _ := strconv.Atoi(c.Param("forum_id"))
	pendingHelps, err := models.GetPendingHelpsByForumID(forum_id)
	if err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error(), "data": pendingHelps})
	    return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取所有的已经应答但是还没有完成的 help 成功", "data": pendingHelps})
	return
}

// 获取所有的已经应答但是还没有完成的 help
// GetAllFinishedHelpByForumID godoc
// @Summary GetAllFinishedHelpByForumID
// @Description	GetAllFinishedHelpByForumID
// @Tags Helps
// @Accept	mpfd
// @Produce	json
// @Param token header string true "将token放在请求头部的‘Authorization‘字段中，并以‘Bearer ‘开头""
// @Success 200 {object} responses.StatusOKResponse{data=[]models.PendingOrFinishedHelpDetail}
// @Failure 500 {object} responses.StatusInternalServerError "服务器错误"
// @Router /forums/{forum_id}/helps/finished [get]
func GetAllFinishedHelpByForumID(c *gin.Context) {
	log.Info("get all finished helps by forum_id controller")
	forum_id, _ := strconv.Atoi(c.Param("forum_id"))
	pendingHelps, err := models.GetFinishedHelpsByForumID(forum_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error(), "data": pendingHelps})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取所有的已经应答但是还没有完成的 help 成功", "data": pendingHelps})
	return
}



