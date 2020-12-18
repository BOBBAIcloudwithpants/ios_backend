package middlewares

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pingcap/log"
	"github.com/bobbaicloudwithpants/ios_backend/service"
)

// 检查用户能否查看当前论坛
func CanUserWatchTheForum() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info("can user watch the forum middleware")
		user := service.GetUserFromContext(c)
		forum_id, _ := strconv.Atoi(c.Param("forum_id"))

		ok, err := service.IsUserInForum(user.UserId, forum_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "中间件查询错误 " + err.Error(), "data": nil})
			c.Abort()
			return
		}

		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "该用户没有查看当前论坛的权限", "data": nil})
			c.Abort()
			return

		}
		c.Next()
	}

}
