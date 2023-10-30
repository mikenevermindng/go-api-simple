package router

import (
	"github.com/gin-gonic/gin"
	"goApiByGin/controllers"
	"goApiByGin/middlewares"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		user := v1.Group("/user")
		{
			user.POST("/signup", controllers.SignUp)
			user.POST("/signin", controllers.SignIn)
			user.POST("/refresh", controllers.Refresh)
		}
		todo := v1.Group("/todo")
		todo.Use(middlewares.AuthValidateMiddleware, middlewares.RoleValidateFactory([]string{"USER"}))
		{
			todo.POST("/", controllers.AddNote)
			todo.GET("/", controllers.GetNotes)
		}
	}
	return r
}
