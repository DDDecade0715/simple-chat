package controller

import (
	"gin-derived/pkg/app/response"
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	response.OkWithMessage("hello gin-derived", c)
}
