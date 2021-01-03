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
	var ret []models.PendingOrFinishedHelpDetail
	pendingHelps, err := models.GetFinishedHelpsByForumID(forum_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error(), "data": ret})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取所有的已经应答但是还没有完成的 help 成功", "data": pendingHelps})
	return
}


// 获取所有的已经应答但是还没有完成的 help
// GetAllHelpedPeople godoc
// @Summary GetAllHelpedPeople
// @Description	GetAllHelpedPeople
// @Tags Users
// @Accept	json
// @Produce	json
// @Param token header string true "将token放在请求头部的‘Authorization‘字段中，并以‘Bearer ‘开头""
// @Success 200 {object} responses.StatusOKResponse{data=[]models.User}
// @Failure 500 {object} responses.StatusInternalServerError "服务器错误"
// @Router /users/{user_id}/helped [get]
func GetAllHelpedPeople(c *gin.Context) {
	var ret []models.User
	log.Info("get all helped people")
	user_id, _ := strconv.Atoi(c.Param("user_id"))
	res, err := models.GetAllHelpedPeopleByUserID(user_id)
	if err != nil {
	    c.JSON(500, gin.H{"code": 500, "msg": "查询所有帮助过的用户信息异常 "+ err.Error(), "data": ret})
	    return
	}
	c.JSON(200, gin.H{"code": 200, "msg": "查询所有帮助过的用户信息成功", "data": res})
	return
}

type ModifyParam struct {
	UserID int	`json:"user_id"`
	IsFinish bool `json:"is_finish"`
}

// ModifyStatusOfOneHelp godoc
// @Summary ModifyStatusOfOneHelp
// @Description ModifyStatusOfOneHelp
// @Tags Helps
// @Accept  json
// @Produce  json
// @Param is_finish body bool true "表示本次修改的类型：为true则为完成该 Help, 为false 则为响应该help（即'伸出援手'）"
// @Param user_id body int true "帮助者的id"
// @Success 200 {object} responses.StatusOKResponse "修改 help 状态成功"
// @Router /helps/{help_id} [patch]
func ModifyStatusOfOneHelp(c *gin.Context) {
	log.Info("modify status of one help")
	var param ModifyParam
	err := c.BindJSON(&param)
	help_id, _ := strconv.Atoi(c.Param("help_id"))

	if err != nil {
	    c.JSON(403, gin.H{"code": 403, "msg": "请求参数不合法", "data": nil})
	    return
	}

	// 如果该字段为真，则完成某个请求
	if param.IsFinish {
		err := models.FinishHelpByHelpID(help_id)
		if err != nil {
		    c.JSON(500, gin.H{"code": 500, "msg": "数据库修改错误 "+err.Error(), "data": nil})
		    return
		}
	} else {
		err := models.AnswerHelpByHelpIDAndUserID(help_id, param.UserID)
		if err != nil {
			c.JSON(500, gin.H{"code": 500, "msg": "数据库修改错误 "+err.Error(), "data": nil})
			return
		}
	}
	c.JSON(200, gin.H{"code": 200, "msg": "修改 help 状态成功", "data": nil})
	return
}






