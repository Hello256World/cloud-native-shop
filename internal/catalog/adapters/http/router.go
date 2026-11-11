package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup, ch *CategoryHandler, vh *VariantTypeHandler, ph *ProductHandler, pub *PublicProductHandler, ah *AttributeHandler, authMW, adminMW gin.HandlerFunc) {
	api.GET("/categories", ch.ListPublic)
	api.GET("/products", pub.List)
	api.GET("/products/:slug", pub.Get)

	admin := api.Group("/admin", authMW, adminMW)
	admin.GET("/categories", ch.ListAdmin)
	admin.POST("/categories", ch.Create)
	admin.PUT("/categories/:id", ch.Update)
	admin.DELETE("/categories/:id", ch.Delete)

	admin.GET("/variant-types", vh.List)
	admin.POST("/variant-types", vh.CreateType)
	admin.POST("/variant-types/:id/values", vh.CreateValue)

	admin.POST("/products", ph.Create)
	admin.GET("/products/:id", ph.Get)
	admin.PUT("/products/:id", ph.Update)
	admin.DELETE("/products/:id", ph.Delete)
	admin.POST("/products/:id/variants", ph.AddVariant)
	admin.PUT("/products/:id/variants/:variantId", ph.UpdateVariant)
	admin.POST("/products/:id/images", ph.AddImage)
	admin.DELETE("/products/:id/images/:imageId", ph.DeleteImage)
	admin.GET("/products", ph.List)
	admin.GET("/attributes", ah.List)
	admin.POST("/attributes", ah.Create)
}
