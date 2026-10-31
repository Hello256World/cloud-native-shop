package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup, ch *CategoryHandler, authMW, adminMW gin.HandlerFunc) {
	api.GET("/categories", ch.ListPublic)

	admin := api.Group("/admin", authMW, adminMW)
	admin.GET("/categories", ch.ListAdmin)
	admin.POST("/categories", ch.Create)
	admin.PUT("/categories/:id", ch.Update)
	admin.DELETE("/categories/:id", ch.Delete)
}
