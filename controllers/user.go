package controllers

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"

	"github.com/bobbaicloudwithpants/ios_backend/models"
	"github.com/bobbaicloudwithpants/ios_backend/service"
	"github.com/gin-gonic/gin"
	"github.com/pingcap/log"
)

// 用户注册需要提供的字段
type RegisterParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// 用户注册控制器
// UserRegister godoc
// @Summary UserRegister
// @Description UserRegister
// @Tags Users
// @Accept  json
// @Produce  json
// @Param username body string true "用户名"
// @Param password body string true "密码"
// @Param email body string true "邮箱"
// @Success 200 {object} responses.StatusOKResponse "注册成功"
// @Failure 400 {object} responses.StatusBadRequestResponse "参数不合法"
// @Failure 403 {object} responses.StatusForbiddenResponse "该用户名已经被使用"
// @Failure 403 {object} responses.StatusForbiddenResponse "该邮箱已经被使用"
// @Failure 500 {object} responses.StatusInternalServerError "服务器错误"
// @Router /users [post]
func UserRegister(c *gin.Context) {
	log.Info("user register controller")
	var param RegisterParam
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数不合法: " + err.Error()})
		return
	}

	if ok, err := service.IsUsernameExist(param.Username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: " + err.Error()})
		return
	} else if ok {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "该用户名已经被使用"})
		return
	}

	if ok, err := service.IsEmailExist(param.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: " + err.Error()})
		return
	} else if ok {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "该邮箱已经被使用"})
		return
	}

	err = service.CreateUser(param.Username, param.Password, param.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "注册成功"})
}

// 登陆参数
type LoginParam struct {
	Input    string `json:"input"`
	Password string `json:"password"`
}

// 登录成功返回的结果
type LoginResponse struct {
	Token  string `json:"token"`
	UserId string `json:"user_id"`
}

// 用户登录控制器
// UserLogin godoc
// @Summary UserLogin
// @Description UserLogin
// @Tags Users
// @Accept  json
// @Produce  json
// @Param input body string true "用户名或者邮箱"
// @Param password body string true "密码"
// @Success 200 {object} responses.StatusOKResponse{data=LoginResponse} "正确登陆"
// @Failure 400 {object} responses.StatusBadRequestResponse "参数不合法"
// @Failure 403 {object} responses.StatusForbiddenResponse "密码错误"
// @Failure 403 {object} responses.StatusForbiddenResponse "该用户名或邮箱不存在"
// @Failure 500 {object} responses.StatusInternalServerError "服务器错误"
// @Router /users [put]
func UserLogin(c *gin.Context) {
	log.Info("user login controller")
	var param LoginParam
	data := make(map[string]string)
	//buf := make([]byte, 1024)
	//n, _ := c.Request.Body.Read(buf)
	//fmt.Println(string(buf[0:n]))
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数不合法: " + err.Error(), "data": data})
		log.Info("参数不合法 "+err.Error())
		return
	}
	if ok, err := service.IsUsernameExist(param.Input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: " + err.Error(), "data": data})
		return
	} else if ok {
		pass, err := service.VerifyByUsernameAndPassword(param.Input, param.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: " + err.Error(), "data": data})
			return
		}
		if pass {
			data["token"], err = service.ProduceTokenByUsernameAndPasword(param.Input, param.Password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: " + err.Error(), "data": data})
				return
			}
			userinfo, _ := service.ParseToken(data["token"])

			data["user_id"] = strconv.Itoa(userinfo.UserId)
			c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "登录成功", "data": data})
			return
		} else {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "密码错误", "data": data})
			return
		}
	}

	if ok, err := service.IsEmailExist(param.Input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: " + err.Error(), "data": data})
		return
	} else if ok {
		pass, err := service.VerifyByEmailAndPassword(param.Input, param.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: " + err.Error(), "data": data})
			return
		}
		if pass {
			data["token"], err = service.ProduceTokenByEmailAndPassword(param.Input, param.Password)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: " + err.Error(), "data": data})
				return
			}
			userinfo, _ := service.ParseToken(data["token"])
			data["user_id"] = strconv.Itoa(userinfo.UserId)
			c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "登录成功", "data": data})
			return
		} else {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "密码错误", "data": data})
			return
		}
	}

	c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "该用户名或邮箱不存在", "data": data})
}

// GetAllUsers godoc
// @Summary GetAllUsers
// @Description GetAllUsers
// @Tags Users
// @Accept  json
// @Produce  json
// @Param token header string true "将token放在请求头部的‘Authorization‘字段中，并以‘Bearer ‘开头""
// @Param username header string false "用户名的子串"
// @Success 200 {object} responses.StatusOKResponse{data=[]models.User} "获取全部用户"
// @Failure 500 {object} responses.StatusInternalServerError "服务器错误"
// @Router /users [get]
func GetAllUsers(c *gin.Context) {
	log.Info("get all users controller")

	var data []models.User
	var err error
	query := c.Query("username")
	// fmt.Println("query", query)
	if query != "" {
		data, err = models.GetAllUsersContains(query)
	} else {
		data, err = models.GetAllUsers()
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: " + err.Error(), "data": data})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取全部用户", "data": data})
}

// 上传用户头像图像
// UploadAvatar godoc
// @Summary UploadAvatar
// @Description UploadAvatar
// @Tags Users
// @Accept  json
// @Produce  json
// @Param avatar formData file true "用户头像"
// @Param token header string true "将token放在请求头部的‘Authorization‘字段中，并以‘Bearer ‘开头""
// @Success 200 {object} responses.StatusOKResponse "上传头像成功"
// @Failure 400 {object} responses.StatusBadRequestResponse "请求格式不正确"
// @Failure 403 {object} responses.StatusForbiddenResponse "禁止更改他人资源"
// @Failure 500 {object} responses.StatusInternalServerError "文件服务错误"
// @Router /users/{user_id}/avatar [post]
func UploadAvatar(c *gin.Context) {
	log.Info("upload user avatar controller")
	data := make(map[string]string)
	// 获取token的claim
	claims, _ := c.MustGet("Claims").(*service.Claims)
	user_id, _ := strconv.Atoi(c.Param("user_id"))
	if claims.UserId != user_id {
		// 不允许使用自己的token修改他人的资源
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "禁止更改他人资源", "data": data})
		return
	}
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求格式不正确: " + err.Error(), "data": data})
	} else {
		fmt.Println(c.Request.URL.String())
		// 图片统一改成png上传
		filename, err := service.FileUpload(file, header, path.Base(c.Request.URL.Path), c.Request.URL.Path, ".png")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "文件服务错误: " + err.Error(), "data": data})
		} else {
			// 写入数据库
			avatar := filename
			err := models.UpdateUserAvatarByUserId(user_id, avatar)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "数据库写入错误: " + err.Error(), "data": nil})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "上传头像成功", "data": data})
		}
	}
}

// 获取用户图片
// GetAvatar godoc
// @Summary GetAvatar
// @Description GetAvatar
// @Tags Users
// @Accept  json
// @Produce  image/jpeg
// @Success 200 {object} responses.StatusOKResponse{data=[]byte} "读取头像成功，data为字节数足"
// @Failure 404 {object} responses.StatusForbiddenResponse "获取头像失败，下载时出错"
// @Failure 500 {object} responses.StatusInternalServerError "读取图片失败，处理时出错"
// @Header 200 {string} Content-Disposition "attachment; filename=hello.txt"
// @Header 200 {string} Content-Type "image/jpeg"
// @Header 200 {string} Accept-Length "image's length"
// @Router /users/{user_id}/avatar [get]
func GetAvatar(c *gin.Context) {
	log.Info("get user's avatar controller")
	var data interface{}
	user_id, _ := strconv.Atoi(c.Param("user_id"))

	users, err := models.GetUserById(user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "数据库查询出错: " + err.Error(), "data": nil})
		return
	}
	if len(users) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "该用户不存在", "data": nil})
		return
	}
	user := users[0]
	if user.Avatar == "0.jpg" {
		// 下载默认头像
		rawImage, err := service.FileDownloadByName("0.jpg", "avatar")
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "获取默认头像失败" + err.Error(), "data": data})
			return
		}
		// 图片最多3个M
		image := make([]byte, 3000000)
		len, err := rawImage.Read(image)
		if err != nil {
			if err != io.EOF && err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "读取图片失败 " + err.Error(), "data": data})
			} else {
				c.Writer.WriteHeader(http.StatusOK)
				c.Header("Content-Disposition", "attachment; filename=0.png")
				c.Header("Content-Type", "image/jpeg")
				c.Header("Accept-Length", fmt.Sprintf("%d", len))
				c.Writer.Write(image)
			}
		}
	} else {
		rawImage, err := service.FileDownloadByName(user.Avatar, "avatar")
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "获取头像失败" + err.Error(), "data": data})
			return
		}
		// 图片最多3个M
		image := make([]byte, 1000000)
		len, err := rawImage.Read(image)
		if err != nil {
			if err != io.EOF && err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "读取图片失败 " + err.Error(), "data": data})
			} else {
				c.Writer.WriteHeader(http.StatusOK)
				c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", user.Avatar))
				c.Header("Content-Type", "image/jpeg")
				c.Header("Accept-Length", fmt.Sprintf("%d", len))
				c.Writer.Write(image)
			}
		}
	}

}

// 获取用户的关注订阅列表
// GetOneUserSubscribe godoc
// @Summary GetOneUserSubscribe
// @Description	GetOneUserSubscribe
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} responses.StatusOKResponse{data=models.SubscribeList} "获取第{user_id}号用户的关注订阅列表成功"
// @Failure 500 {object} responses.StatusInternalServerError "数据库查询出错"
// @Router /users/{user_id}/subscribe [get]
func GetOneUserSubscribe(c *gin.Context) {
	log.Info("get one user's subscribe controller")
	user_id, _ := strconv.Atoi(c.Param("user_id"))

	subscribe, err := service.GetOneUserSubscribe(user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "数据库查询出错", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": fmt.Sprintf("获取第 %d 号用户的关注订阅列表成功", user_id), "data": subscribe})
	return
}

// 获取某个用户的详情
// GetOneUserDetailByUserID
// @Summary GetOneUserDetailByUserID
// @Description GetOneUserDetailByUserID
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} responses.StatusOKResponse
//{data=models.UserDetail}
// @Failure 500 {object} responses.StatusInternalServerError "数据库查询出错"
// @Router /users/{user_id}/info [get]
func GetOneUserDetailByUserID(c *gin.Context) {
	log.Info("get one user's detail by user_id")
	var ret models.UserDetail
	user_id, _ := strconv.Atoi(c.Param("user_id"))

	userDetail, err := service.GetOneUserDetail(user_id)
	if err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error(), "data": ret})
	    return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg":fmt.Sprintf("获取第 %d 号用户的详情成功", user_id), "data": userDetail })

}

// 获取某个用户所发布的帖子
// GetOneUserPostsByUserID
// @Summary GetOneUserPostsByUserID
// @Description GetOneUserPostsByUserID
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} responses.StatusOKResponse{data=[]models.PostDetail}
// @Failure 500 {object} responses.StatusInternalServerError "数据库查询出错"
// @Router /users/{user_id}/posts [get]
func GetOneUserPostsByUserID(c *gin.Context) {
	log.Info("get one user's posts by user_id")
	var ret []models.PostDetail
	user_id, _ := strconv.Atoi(c.Param("user_id"))
	posts, err := service.GetOneUserPostsByUserID(user_id)
	if err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": err.Error(), "data": ret})
	    return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg":fmt.Sprintf("获取第 %d 号用户的所有帖子", user_id), "data": posts })
}