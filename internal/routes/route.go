package routes

import (
	"github.com/EputraP/kfc_be/internal/handler"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth *handler.AuthHandler
}

type Middlewares struct {
	Auth gin.HandlerFunc
}

func Build(srv *gin.Engine, h *Handlers, middlewares *Middlewares) {

	auth := srv.Group("/auth")
	auth.POST("/register", h.Auth.CreateUser)
	auth.POST("/login", h.Auth.Login)
	auth.GET("/refresh", middlewares.Auth, h.Auth.Refresh)
	auth.GET("/logout", h.Auth.Logout)

}
