package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup, h *Handler, authMW gin.HandlerFunc) {
	auth := api.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.Refresh)
	auth.POST("/logout", h.Logout)

	me := api.Group("/me", authMW)
	me.GET("", h.Me)
}
