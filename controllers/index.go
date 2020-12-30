package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pingcap/log"
)

type IndexImage struct {
	ImageName string
	Description string
}
func GetIndexImages(c *gin.Context) {
	log.Info("get index images")

	var ret []IndexImage

	ret = append(ret, IndexImage{
		ImageName: "programming.jpg", Description: "程序猿饲养地",
	})
	ret = append(ret, IndexImage{
		ImageName: "movie.jpg", Description: "电影天堂",
	})
	ret = append(ret, IndexImage{
		ImageName: "rock.jpg", Description: "摇滚圣经",
	})

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取头图成功", "data": ret})
	return

}
