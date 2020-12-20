package controllers

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"

	"github.com/pingcap/log"
	"github.com/bobbaicloudwithpants/ios_backend/service"

	"github.com/gin-gonic/gin"
	"github.com/bobbaicloudwithpants/ios_backend/models"
)

type ForumResponse struct {
	Forums     []models.Forum    `json:"forums"`
	UserDetail models.UserDetail `json:"user_detail"`
}

// GetAllPublicFroums godoc
// @Summary GetAllPublicFroums
// @Description GetAllPublicFroums
// @Tags Forums
// @Accept  json
// @Produce  json
// @Success 200 {object} responses.StatusOKResponse{data=ForumResponse} "获取全部公开论坛"
// @Failure 500 {object} responses.StatusInternalServerError "服务器错误"
// @Router /forums [get]
func GetAllPublicFroums(c *gin.Context) {
	log.Info("get all public forims controller")
	var data ForumResponse

	forums, err := models.GetAllPublicForums()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: " + err.Error(), "data": data})
		return
	}
	user := service.GetUserFromContext(c)
	userDetail, err := service.GetOneUserDetail(user.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询用户信息错误: " + err.Error(), "data": data})
		return
	}

	data.UserDetail = userDetail
	data.Forums = forums
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取全部公开论坛", "data": data})
}

// GetForumByID godoc
// @Summary GetForumByID
// @Description GetForumByID
// @Tags Forums
// @Accept  json
// @Produce  json
// @Success 200 {object} responses.StatusOKResponse{data=ForumResponse} "获取全部公开论坛"
// @Failure 500 {object} responses.StatusInternalServerError "服务器错误"
// @Router /forums/{forum_id} [get]
func GetForumByID(c *gin.Context) {
	log.Info("get forum by id controller")
	var data ForumResponse
	forum_id, _ := strconv.Atoi(c.Param("forum_id"))
	forums, err := models.GetForumByID(forum_id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: " + err.Error(), "data": data})
		return
	}
	user := service.GetUserFromContext(c)
	userDetail, err := service.GetOneUserDetail(user.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询用户信息错误: " + err.Error(), "data": data})
		return
	}

	data.UserDetail = userDetail
	data.Forums = forums
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取论坛成功", "data": data})
}

// 论坛参数
type ForumParam struct {
	ForumName   string `json:"forum_name"`
	IsPublic    bool   `json:"is_public"`
	Description string `json:"description"`
}

// CreateForum godoc
// @Summary CreateForum
// @Description CreateForum
// @Tags Forums
// @Accept  json
// @Produce  json
// @Param forum_name body string true "论坛名"
// @Param is_public body bool true "是否公开"
// @Param description body string true "论坛描述"
// @Success 200 {object} responses.StatusOKResponse "论坛创建成功"
// @Failure 400 {object} responses.StatusBadRequestResponse "参数不合法"
// @Failure 500 {object} responses.StatusInternalServerError "服务器错误"
// @Router /forums [post]
func CreateForum(c *gin.Context) {
	log.Info("create forum controller")
	var param ForumParam
	data := make(map[string]string)
	err := c.BindJSON(&param)
	fmt.Println(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数不合法: " + err.Error(), "data": data})
		return
	}
	// TODO:
	// 1.检查论坛同名问题？好像没必要
	forum_id, err := service.CreateForum(param.ForumName, param.Description, param.IsPublic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: " + err.Error()})
		return
	}
	// 2.写入forum_user的关系
	user := service.GetUserFromContext(c)
	err = models.AddRoleInForum(int(forum_id), user.UserId, "owner")
	if err != nil {
		// TODO:
		// 撤销论坛的建立？
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "论坛创建成功"})
}

// 上传论坛封面
// UploadCover godoc
// @Summary UploadCover
// @Description UploadCover
// @Tags Forums
// @Accept  json
// @Produce  json
// @Param cover formData file true "论坛封面"
// @Param token header string true "将token放在请求头部的‘Authorization‘字段中，并以‘Bearer ‘开头""
// @Success 200 {object} responses.StatusOKResponse "上传封面成功"
// @Failure 400 {object} responses.StatusBadRequestResponse "请求格式不正确"
// @Failure 403 {object} responses.StatusForbiddenResponse "禁止更改他人资源"
// @Failure 500 {object} responses.StatusInternalServerError "文件服务错误"
// @Router /forums/{forum_id}/cover [post]
func UploadCover(c *gin.Context) {
	log.Info("upload cover controller")
	data := make(map[string]string)
	// TODO:
	// 查表当前的用户是否是论坛的主人或者管理员
	file, header, err := c.Request.FormFile("cover")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求格式不正确: " + err.Error(), "data": data})
	} else {
		fmt.Println(c.Request.URL.String())
		fmt.Println("base", path.Base(c.Request.URL.Path))
		fmt.Println("not base", c.Request.URL.Path)
		_, err := service.FileUpload(file, header, path.Base(c.Request.URL.Path), c.Request.URL.Path, ".png")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "文件服务错误: " + err.Error(), "data": data})
		} else {
			forumID, err := strconv.Atoi(c.Param("forum_id"))
			if err != nil {
				// id错误?
			}
			// TODO:
			err = models.UpdateCover(c.Request.URL.Path+".png", forumID)
			if err != nil {
				// 更新数据库失败
			}
			c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "上传封面成功", "data": data})
		}
	}
}

// GetCover godoc
// @Summary GetCover
// @Description GetCover
// @Tags Forums
// @Accept  json
// @Produce  image/jpeg
// @Success 200 {object} responses.StatusOKResponse{data=[]byte} "读取封面成功，data为字节数组"
// @Failure 404 {object} responses.StatusForbiddenResponse "获取封面失败，下载时出错"
// @Failure 500 {object} responses.StatusInternalServerError "读取图片失败，处理时出错"
// @Header 200 {string} Content-Disposition "attachment; filename=hello.txt"
// @Header 200 {string} Content-Type "image/jpeg"
// @Header 200 {string} Accept-Length "image's length"
// @Router /forums/{froum_id}/cover [get]
func GetCover(c *gin.Context) {
	log.Info("get cover controller")
	var data interface{}
	forum_id, _ := strconv.Atoi(c.Param("forum_id"))
	forums, _ := models.GetForumByID(forum_id)
	fmt.Println("forum is", forums[0])
	filename := forums[0].Cover

	fmt.Println("filename", filename)
	rawImage, err := service.FileDownloadByName(filename, "cover")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "获取封面失败" + err.Error(), "data": data})
	} else {
		// 图片最多2个M
		image := make([]byte, 2000000)
		len, err := rawImage.Read(image)
		if err != nil {
			if err != io.EOF && err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "读取图片失败 " + err.Error(), "data": data})
			} else {
				c.Writer.WriteHeader(http.StatusOK)
				c.Header("Content-Disposition", "attachment; filename=hello.txt")
				c.Header("Content-Type", "image/jpeg")
				c.Header("Accept-Length", fmt.Sprintf("%d", len))
				c.Writer.Write(image)
			}
		}
	}
}

// 获取用户身份 注意不需要jwt
// GetRoleInForum godoc
// @Summary GetRoleInForum
// @Description GetRoleInForum
// @Tags Role
// @Accept  json
// @Produce  json
// @Success 200 {object} responses.StatusOKResponse "获取角色成功"
// @Failure 400 {object} responses.StatusBadRequestResponse "该用户不再此论坛下"
// @Failure 500 {object} responses.StatusInternalServerError "服务器错误"
// @Router /forums/{forum_id}/role/{user_id} [get]
func GetRoleInForum(c *gin.Context) {
	log.Info("get role in forum controller")
	forum_id, _ := strconv.Atoi(c.Param("forum_id"))
	user_id, _ := strconv.Atoi(c.Param("user_id"))
	role, err := models.FindRoleInForum(forum_id, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: " + err.Error(), "data": nil})
		return
	} else if role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "该用户不再此论坛下", "data": role})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取角色成功", "data": role})
}

type RoleParam struct {
	Role string `json:"role" example:"user/admin/owner/null"`
}

// 修改用户身份 需要jwt以及json带上role
// UpdateRoleInForum godoc
// @Summary UpdateRoleInForum
// @Description UpdateRoleInForum
// @Tags Role
// @Accept  json
// @Produce  json
// @Param token header string true "将token放在请求头部的‘Authorization‘字段中，并以‘Bearer ‘开头""
// @Param role body string true "目标用户的身份"
// @Success 200 {object} responses.StatusOKResponse "所有者转让成功"
// @Success 200 {object} responses.StatusOKResponse "授予管理员成功"
// @Failure 400 {object} responses.StatusBadRequestResponse "请求格式不正确"
// @Failure 403 {object} responses.StatusForbiddenResponse "操作者非本论坛成员"
// @Failure 403 {object} responses.StatusForbiddenResponse "操作者身份权限不足"
// @Router /forums/{forum_id}/role/{user_id} [patch]
func UpdateRoleInForum(c *gin.Context) {
	log.Info("update role in forum controller")

	var roleParam RoleParam
	err := c.BindJSON(&roleParam)
	if err != nil {
		// 参数错误
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求格式不正确: " + err.Error(), "data": nil})
		return
	}

	forum_id, _ := strconv.Atoi(c.Param("forum_id"))
	// 需要更新身份的用户id
	update_user_id, _ := strconv.Atoi(c.Param("user_id"))

	// 将要修改的角色
	update_role := roleParam.Role
	// 执行修改操作的用户id
	user_id := service.GetUserFromContext(c).UserId
	// 操作者的角色
	role, err := models.FindRoleInForum(forum_id, user_id)
	if err != nil {
		// 操作者身份不明确
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "操作者非本论坛成员", "data": nil})
		return
	}
	// fmt.Println("操作者为", role)
	if role == "user" {
		// 如果操作者是用户
		// 权限不够
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "操作者身份权限不足", "data": nil})
		return
	} else if role == "owner" {
		// 如果操作者是owner，可以指派管理员，也可以指派所有者，但这时是转移
		if update_role == "owner" {
			// 所有者转让
			err = models.UpdateRoleInForum(forum_id, user_id, "admin")
			if err != nil {

			}
			err := models.UpdateRoleInForum(forum_id, update_user_id, "owner")
			if err != nil {

			}
			c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "所有者转让成功", "data": nil})
			return
		} else if update_role == "admin" {
			err := models.UpdateRoleInForum(forum_id, update_user_id, "admin")
			if err != nil {

			}
			c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "授予管理员成功", "data": nil})
		}
	} else {
		// 如果操作者是管理员，可以踢人
		// TODO
		if update_role == "null" {
			//踢人
		}
	}
}

// 用户订阅forum
// SubscribeForum godoc
// @Summary SubscribeForum
// @Description SubscribeForum
// @Tags Role
// @Accept  json
// @Produce  json
// @Param token header string true "将token放在请求头部的‘Authorization‘字段中，并以‘Bearer ‘开头""
// @Success 200 {object} responses.StatusOKResponse "订阅成功"
// @Failure 403 {object} responses.StatusForbiddenResponse "不可重复订阅"
// @Failure 500 {object} responses.StatusInternalServerError "服务器错误"
// @Router /forums/{forum_id}/role [post]
func SubscribeForum(c *gin.Context) {
	log.Info("subscribe forum controller")
	forum_id, _ := strconv.Atoi(c.Param("forum_id"))
	// 用户id
	user_id := service.GetUserFromContext(c).UserId
	// 不能重复订阅
	role, _ := models.FindRoleInForum(forum_id, user_id)
	if role != "" {
		// 重复订阅失败
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "不可重复订阅", "data": nil})
		return
	}

	err := models.AddRoleInForum(forum_id, user_id, "user")
	if err != nil {
		// 服务器错误
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: " + err.Error(), "data": nil})
		return
	}
	// 订阅成功
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "订阅成功", "data": nil})
}

// UnSubscribeForum godoc
// @Summary UnSubscribeForum
// @Description UnSubscribeForum
// @Tags Role
// @Accept  json
// @Produce  json
// @Param token header string true "将token放在请求头部的‘Authorization‘字段中，并以‘Bearer ‘开头""
// @Success 200 {object} responses.StatusOKResponse "取消订阅成功"
// @Failure 500 {object} responses.StatusInternalServerError "服务器错误"
// @Router /forums/{forum_id}/role [delete]
func UnSubscribeForum(c *gin.Context) {
	log.Info("unsubscribe forum controller")
	forum_id, _ := strconv.Atoi(c.Param("forum_id"))
	// 用户id
	user_id := service.GetUserFromContext(c).UserId
	// TODO：记得删除对应的数据库记录
	err := models.DeleteRoleInForum(forum_id, user_id)
	if err != nil {
		// 服务器错误
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: " + err.Error(), "data": nil})
		return
	}
	// 订阅成功
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "取消订阅成功", "data": nil})
}

type UsersListParam struct {
	Users []int `json:"users"`
}

// 批量添加用户
// AddUsersToForum godoc
// @Summary AddUsersToForum
// @Description AddUsersToForum
// @Tags Role
// @Accept  json
// @Produce  json
// @Param token header string true "将token放在请求头部的‘Authorization‘字段中，并以‘Bearer ‘开头""
// @Param users body UsersListParam true "用户id列表"
// @Success 200 {object} responses.StatusOKResponse "订阅成功"
// @Failure 403 {object} responses.StatusForbiddenResponse "不可重复订阅"
// @Failure 500 {object} responses.StatusInternalServerError "服务器错误"
// @Router /forums/{forum_id}/role [put]
func AddUsersToForum(c *gin.Context) {
	log.Info("add users to forum controller")
	forum_id, _ := strconv.Atoi(c.Param("forum_id"))
	// 用户id
	// user_id := service.GetUserFromContext(c).UserId
	// 不能重复订阅
	// role, _ := models.FindRoleInForum(forum_id, user_id)
	// if role != "" {
	// 	// 重复订阅失败
	// 	c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "不可重复订阅", "data": nil})
	// 	return
	// }
	var idList UsersListParam
	err := c.BindJSON(&idList)
	if err != nil {
		// 参数错误
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求格式不正确: " + err.Error(), "data": nil})
		return
	}
	for _, r := range idList.Users {
		// fmt.Println("user", r)
		err := models.AddRoleInForum(forum_id, r, "user")
		if err != nil {
			// 服务器错误
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: " + err.Error(), "data": nil})
			return
		}
	}
	// 订阅成功
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "订阅成功", "data": nil})
}
