package middlewares

import (
	"github.com/gin-gonic/gin"
	"goApiByGin/helpers"
	"net/http"
)

func RoleValidateFactory(allowances []string) func(*gin.Context) {
	return func(c *gin.Context) {
		if roles, ok := c.Get("roles"); ok {
			for _, r := range roles.([]string) {
				if helpers.IsElementExist(allowances, r) {
					c.Next()
					return
				}
			}
		}

		c.JSON(http.StatusForbidden, helpers.Response{Code: http.StatusForbidden, Message: "User is not allow to access this api", Data: nil})
		c.Abort()
		return
	}

}
