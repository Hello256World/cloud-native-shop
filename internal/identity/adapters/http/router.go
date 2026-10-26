package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup, authH *AuthHandler, addrH *AddressHandler, authMW gin.HandlerFunc) {
	auth := api.Group("/auth")
	auth.POST("/register", authH.Register)
	auth.POST("/login", authH.Login)
	auth.POST("/refresh", authH.Refresh)
	auth.POST("/logout", authH.Logout)

	me := api.Group("/me", authMW)
	me.GET("", authH.Me)

	addrs := me.Group("/addresses")
	addrs.POST("", addrH.Create)
	addrs.GET("", addrH.List)
	addrs.PUT("/:id", addrH.Update)
	addrs.DELETE("/:id", addrH.Delete)
}
