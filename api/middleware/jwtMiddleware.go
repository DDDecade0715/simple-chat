package middleware

import (
	"gin-derived/pkg/app/response"
	jwtPkg "gin-derived/pkg/jwt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		//只允许 header Authorization 参数提交token
		auth := c.GetHeader("Authorization")
		if auth == "" {
			response.FailWithMessage("header缺少Authorization参数", c)
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(auth, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.FailWithMessage("请求头中auth格式有误", c)
			c.Abort()
			return
		}
		claims, err := jwtPkg.ParseToken(parts[1])
		if err != nil {
			var errMsg string
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				errMsg = "鉴权失败，Token 超时"
			default:
				errMsg = "鉴权失败，Token 错误"
			}
			response.FailWithMessage(errMsg, c)
			c.Abort()
			return
		}

		c.Set("UserName", claims.UserName)
		c.Set("UserID", claims.UserID)
		c.Next()
	}
}
