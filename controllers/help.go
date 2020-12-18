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

func CreateHelp(c *gin.Context) {
	log.Info("create help controller")
	forum_id, _ := strconv.Atoi(c.Param("forum_id"))
	form, _ := c.MultipartForm()
	file := form.File["file"][0]	// 仅为 录音 或者 视频
	title, content := c.PostForm("title"), c.PostForm("content")
	bonus, _ := strconv.Atoi(c.PostForm("bonus"))
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
	name, err := service.FileUpload(fileToBeUpload, file, bucketName, c.Request.URL.Path, ext)
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