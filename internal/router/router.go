package router

import (
	"ramah-disabilitas-be/internal/handler"
	"ramah-disabilitas-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// Serve static files from storage directory
	r.Static("/storage", "./storage")

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Halo! Aplikasi Go berhasil jalan di Koyeb.")
	})

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", handler.Register)
			auth.POST("/login", handler.Login)
		}

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/upload", handler.UploadFile)
			protected.GET("/auth/me", handler.GetMe)
			protected.POST("/auth/logout", handler.Logout)

			admin := protected.Group("/admin")
			admin.Use(middleware.AdminMiddleware())
			{
				admin.GET("/health", func(c *gin.Context) {
					c.JSON(200, gin.H{"status": "ok", "role": "admin"})
				})
			}

			protected.POST("/user/accessibility", handler.UpdateAccessibility)
			protected.POST("/courses/join", handler.JoinCourse)
			protected.GET("/courses/joined", handler.GetMyJoinedCourses)

			lecturer := protected.Group("/lecturer")
			lecturer.Use(middleware.LecturerMiddleware())
			{
				lecturer.POST("/courses", handler.CreateCourse)
				lecturer.GET("/courses", handler.GetMyCourses)
				lecturer.PUT("/courses/:id", handler.UpdateCourse)
				lecturer.DELETE("/courses/:id", handler.DeleteCourse)
			}
		}
	}

	return r
}
