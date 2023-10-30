package middlewares

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"goApiByGin/db"
	"goApiByGin/ent/role"
	"goApiByGin/ent/token"
	"goApiByGin/ent/user"
	"goApiByGin/helpers"
	"net/http"
)

func AuthValidateMiddleware(c *gin.Context) {
	authorization := c.Request.Header.Get("X-API-TOKEN")
	if len(authorization) == 0 {
		c.JSON(http.StatusForbidden, helpers.Response{Code: http.StatusForbidden, Message: "Unauthorized", Data: nil})
		c.Abort()
		return
	}

	authToken := authorization[7:]

	if tokenObj, err := db.Client().Token.Query().Where(token.Token(authToken)).Only(context.Background()); err == nil {
		token, parseTokenErr := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(tokenObj.Secret), nil
		})

		if parseTokenErr != nil {
			c.JSON(http.StatusForbidden, helpers.Response{Code: http.StatusForbidden, Message: "Unauthorized", Data: nil})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			c.JSON(http.StatusForbidden, helpers.Response{Code: http.StatusForbidden, Message: "Invalid token", Data: nil})
			c.Abort()
			return
		}

		username := claims["username"]

		fmt.Println(username)

		userData, userErr := db.Client().User.Query().Where(user.Username(username.(string))).Only(context.Background())
		if userErr != nil {
			c.JSON(http.StatusForbidden, helpers.Response{Code: http.StatusForbidden, Message: "User not found", Data: nil})
			c.Abort()
			return
		}

		roleData, roleErr := db.Client().Role.Query().Where(role.User(userData.ID)).All(context.Background())

		if roleErr != nil {
			c.JSON(http.StatusForbidden, helpers.Response{Code: http.StatusForbidden, Message: "Roles not found", Data: nil})
			c.Abort()
			return
		}

		roles := make([]string, len(roleData))
		for i, r := range roleData {
			roles[i] = r.Type.String()
		}

		c.Set("userId", userData.ID.String())
		c.Set("username", username)
		c.Set("roles", roles)
		c.Next()
		return
	}

	c.JSON(http.StatusForbidden, helpers.Response{Code: http.StatusForbidden, Message: "Unauthorized", Data: nil})
	c.Abort()
	return
}
