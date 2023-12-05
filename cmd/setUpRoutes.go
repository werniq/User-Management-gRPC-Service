package main

import (
	"github.com/Chained/auth-service/cmd/middleware"
	"github.com/gin-gonic/gin"
)

func (app *Application) SetupRoutes(r *gin.Engine) {
	r.Use(middleware.CorsMiddleware())

	r.POST("/is-authenticated", app.CheckAuthentication)
	r.POST("/authenticate", app.CreateAuthenticationToken)
	// forgot password, reset password

	users := r.Group("/all-users")

	users.POST("/create", app.CreateUser)
	users.GET("/{id}", app.GetUser)
	users.PUT("/edit/{id}", app.EditUser)
	users.DELETE("/delete/{id}", app.DeleteUser)
}
