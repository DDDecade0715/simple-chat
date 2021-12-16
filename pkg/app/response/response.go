package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var (
	SUCCESS = 0
	ERROR   = 7
)

func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"msg":     msg,
		"data":    data,
		"elapsed": GetElapsed(c),
	})
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(message string, data interface{}, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

func GetElapsed(c *gin.Context) float64 {
	elapsed := 0.00
	if requestTime := Get(c, "beginTime"); requestTime != nil {
		elapsed = float64(time.Since(requestTime.(time.Time))) / 1e9
	}
	return elapsed
}

func Get(c *gin.Context, key string) interface{} {
	val, _ := c.Get(key)
	return val
}
