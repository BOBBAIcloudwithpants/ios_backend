package controllers

import (
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pingcap/log"
	"github.com/bobbaicloudwithpants/ios_backend/models"
	"github.com/bobbaicloudwithpants/ios_backend/service"
)

// 用户创建 POST
// CreatePost godoc
// @Summary CreatePost
// @Description	CreatePost
// @Tags Posts
// @Accept	mpfd
// @Produce	json
// @Param token header string true "将token放在请求头部的‘Authorization‘字段中，并以‘Bearer ‘开头""
// @Param title formData string true "Post 的标题"
// @Param content formData string true "Post 的内容"
// @Param files[] formData file true "文件内容"
// @Success 200 {object} responses.StatusOKResponse "创建 Post 成功"
// @Failure 403 {object} responses.StatusBadRequestResponse "标题或者内容不得为空"
// @Failure 403 {object} responses.StatusBadRequestResponse "您所上传的文件无法打开"
// @Failure 500 {object} responses.StatusInternalServerError "服务器错误"
// @Router /forums/{forum_id}/posts [post]
func CreatePost(c *gin.Context) {
	log.Info("user create post")
	var data interface{}
	forum_id, _ := strconv.Atoi(c.Param("forum_id"))
	form, _ := c.MultipartForm()
	files := form.File["files[]"]
	title, content := c.PostForm("title"), c.PostForm("content")
	user_id := service.GetUserFromContext(c).UserId

	if title == "" || content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "标题或者内容不得为空", "data": data})
		return
	}

	var filesToBeUpload []service.File
	for _, fileHeader := range files {
		f, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "您所上传的文件无法打开", "data": data})
			return
		}
		filesToBeUpload = append(filesToBeUpload, service.File{F: f, H: fileHeader})
	}
	bucketName := path.Base(c.Request.URL.Path)
	names, err := service.MultipleFilesUpload(filesToBeUpload, "posts", c.Request.URL.Path, ".png")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "上传文件失败，服务器内部错误" + err.Error(), "data": data})
		log.Info("上传文件失败，服务器内部错误" + err.Error())
		return
	}
	// 首先插入 post, 获取post_id
	post_id, err := models.CreatePost(models.Post{ForumID: forum_id, UserID: user_id, Title: title, Content: content})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "插入用户创建的post失败", "data": data})
		log.Info("插入用户创建的post失败" + err.Error())
		return
	}

	for _, name := range names {
		_, err := models.CreateFile(models.ExtendedFile{PostID: int(post_id), Bucket: bucketName, FileName: name})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "数据库插入异常 " + err.Error(), "data": data})
			log.Info("数据库插入异常 " + err.Error())
			return
		}
		fmt.Println(name)
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建 post 成功", "data": data})
	return
}


type PostsAndUserDetail struct {
	PostDetails []models.PostDetail	`json:"posts"`
	UserDetail models.UserDetail	`json:"user"`
}
// 获取某个 forum 下的所有post
// GetAllPostsByForumID godoc
// @Summary GetAllPostsByForumID
// @Description GetAllPostsByForumID
// @Tags Posts
// @Accept json
// @Produce json
// @Param token header string true "将token放在请求头部的‘Authorization‘字段中，并以‘Bearer ‘开头""
// @Success 200 {object} responses.StatusOKResponse{data=PostsAndUserDetail}
// @Failure 500 {object} responses.StatusInternalServerError "服务器错误"
// @Router /forums/{forum_id}/posts [get]
func GetAllPostsByForumID(c *gin.Context) {
	log.Info("get all posts by forum_id controller")
	forum_id, _ := strconv.Atoi(c.Param("forum_id"))
	user := service.GetUserFromContext(c)
	user_id := user.UserId
	userDetail, err := service.GetOneUserDetail(user_id)
	if err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取用户详情出错", "data": nil})
	    return
	}

	data, err := service.GetAllPostDetailsByForumID(forum_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询数据库出现异常" + err.Error(), "data": nil})
		return
	}
	for _, post := range data {
		post.IsLikeByCurrentUser = false
		if userDetail.LikeList != nil {
			for _, id := range userDetail.LikeList {
				if id == post.PostID {
					post.IsLikeByCurrentUser = true
				}
			}
		}
	}
	ret := PostsAndUserDetail{UserDetail: userDetail, PostDetails: data}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  fmt.Sprintf("获取论坛 %d 下的全部帖子成功", forum_id),
		"data": ret,
	})
}

// 根据 id 获取某个post的详情
// GetOnePostDetailByPostID godoc
// @Summary GetOnePostDetailByPostID
// @Description GetOnePostDetailByPostID
// @Tags Posts
// @Accept json
// @Produce json
// @Param token header string true "将token放在请求头部的‘Authorization‘字段中，并以‘Bearer ‘开头""
// @Success 200 {object} responses.StatusOKResponse{data=[]models.PostDetail}
// @Failure 400 {object} responses.StatusInternalServerError "数据库查询异常，或者该post不存在"
// @Router /forums/{forum_id}/posts/{post_id} [get]
func GetOnePostDetailByPostID(c *gin.Context) {
	log.Info("get one post detail by post_id")
	post_id, _ := strconv.Atoi(c.Param("post_id"))

	data, err := service.GetOnePostDetailByPostID(post_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "数据库查询异常，或者该post不存在：" + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  fmt.Sprintf("获取第 %d 号帖子成功", post_id),
		"data": data,
	})
}


func LikeOnePostByPostID(c *gin.Context) {
	log.Info("like one post by post id contaroller")
	post_id, _ := strconv.Atoi(c.Param("post_id"))
	user := service.GetUserFromContext(c)
	user_id := user.UserId

	err := models.LikeOnePostByUserIDAndPostID(user_id, post_id)
	if err != nil {
	    c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "点赞失败，您已经点过赞啦 " + err.Error(), "data": nil})
	    return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "点赞成功", "data": nil})
	return
}
