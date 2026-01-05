package router

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Clash BE API",
		})
	})

	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		})
	}

	return r
}
