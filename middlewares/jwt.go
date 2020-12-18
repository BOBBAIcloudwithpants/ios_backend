package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pingcap/log"
	"github.com/bobbaicloudwithpants/ios_backend/service"
)

// 提取token
func ExtractToken(r *http.Request) string {
	// 'Authorization' : 'Bearer token'
	strArr := strings.Split(r.Header.Get("Authorization"), " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func VerifyJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info("verify jwt middleware")
		isError := false
		var data interface{}
		token := ExtractToken(c.Request)
		if token == "" {
			isError = true
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "缺少token",
				"data": data,
			})
		} else {
			claims, err := service.ParseToken(token)
			if err != nil {
				isError = true
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": 500,
					"msg":  "token校验发生错误",
					"data": data,
				})
			} else if time.Now().Unix() > claims.ExpiresAt {
				isError = true
				c.JSON(http.StatusForbidden, gin.H{
					"code": 403,
					"msg":  "token已过期",
					"data": data,
				})
			} else {
				// TODO:
				// 如何优雅防止拿着自己的token，却修改别人的资料？
				// 临时解决方案:通过context传递
				c.Set("Claims", claims)
				// fmt.Println("middleware.userid", claims.UserId)
			}

		}

		if isError {
			c.Abort()
			return
		}
		c.Next()
	}
}
