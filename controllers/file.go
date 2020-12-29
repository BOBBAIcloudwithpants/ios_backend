package controllers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"path/filepath"
	"github.com/gin-gonic/gin"
	"github.com/pingcap/log"
	"github.com/bobbaicloudwithpants/ios_backend/service"
)

// 获取某个 POST 下的全部文件
// GetFilesByPostID godoc
// @Summary GetFilesByPostID
// @Description GetFilesByPostID
// @Tags Files
// @Accept json
// @Produce json
// @Param token header string true "将token放在请求头部的‘Authorization‘字段中，并以‘Bearer ‘开头""
// @Success 200 {object} responses.StatusOKResponse{data=[]models.ExtendedFile}
// @Failure 500 {object} responses.StatusInternalServerError "服务器错误"
// @Router /forums/{forum_id}/posts/{post_id}/files [get]
func GetFilesByPostID(c *gin.Context) {
	log.Info("get files by post_id")
	post_id, _ := strconv.Atoi(c.Param("post_id"))

	files, err := service.GetFilesByPostID(post_id)
	if err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "数据库查询错误 " + err.Error(), "data": nil})
	    return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg": fmt.Sprintf("获取第 %d 帖子下的所有文件成功", post_id),
		"data": files,
	})
}

// 获取某个文件的内容
// GetOneFile godoc
// @Summary GetOneFile
// @Description GetOneFile
// @Tags Files
// @Accept json
// @Produce  image/jpeg
// @Success 200 {object} responses.StatusOKResponse{data=[]byte} "读取文件成功"
// @Failure 404 {object} responses.StatusForbiddenResponse "获取文件失败"
// @Failure 404 {object} responses.StatusInternalServerError "参数不能为空"
// @Header 200 {string} Content-Disposition "attachment; filename=hello.txt"
// @Header 200 {string} Content-Type "image/jpeg"
// @Header 200 {string} Accept-Length "image's length"
// @Router /files/{filename} [get]
func GetOneFile(c *gin.Context){
	log.Info("get one file controller")
	filename := c.Param("filename")
	query := c.Query("help")

	if filename == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg": "参数不能为空",
			"data": nil,
		})
		return
	}
	bucket := "posts"
	if query != "" {
		bucket = "help"
	}


	rawFile, err := service.FileDownloadByName(filename, bucket)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg": "获取文件失败 " + err.Error(),
			"data": nil,
		})
	} else {
		image := make([]byte, 3000000)
		t := "image/jpeg"
		if filepath.Ext(filename) == ".mp3" {
			t = "audio/mpeg"
		}
		len, err := rawFile.Read(image)
		if err != nil {
			if err != io.EOF && err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "读取图片失败 " + err.Error(), "data": nil})
			} else {
				c.Writer.WriteHeader(http.StatusOK)

				c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
				c.Header("Content-Type", t)
				c.Header("Accept-Length", fmt.Sprintf("%d", len))
				c.Writer.Write(image)
			}
		}
	}
}

