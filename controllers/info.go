package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestParam struct {
	Msg string  `json:"msg"`
}
func Test(c *gin.Context) {
	var param RequestParam
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数不合法: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "成功", "data": 2})
	return
}