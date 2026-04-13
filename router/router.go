package router

import (
	"shortener/handlers"

	"github.com/gin-gonic/gin"
)

func Setup(h *handlers.ShortLinkHandler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/shortlinks", h.Create)
		api.GET("/shortlinks/:id", h.GetByID)
	}

	r.GET("/shortlinks/:id", h.Redirect)

	return r
}
