package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup, authH *AuthHandler, addH *AddressHandler, authMW gin.HandlerFunc) {
	auth := api.Group("/auth")
	auth.POST("/register", authH.Register)
	auth.POST("/login", authH.Login)
	auth.POST("/refresh", authH.Refresh)
	auth.POST("/logout", authH.Logout)

	me := api.Group("/me", authMW)
	me.GET("", authH.Me)

	addrs := me.Group("/addresses")
	addrs.POST("", addH.Create)
	addrs.GET("", addH.List)
	addrs.PUT("/:id", addH.Update)
	addrs.DELETE("/:id", addH.Delete)
}
