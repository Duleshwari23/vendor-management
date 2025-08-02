package main

import (
	"log"
	"time"
	"vendor-management/handlers"
	"vendor-management/middleware"
	"vendor-management/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize cron job
	c := cron.New()
	// Run at 12 PM every day
	_, err := c.AddFunc("0 12 * * *", utils.UpdateAttendance)
	if err != nil {
		log.Fatal("Error setting up cron job:", err)
	}
	c.Start()

	// Auth routes
	r.POST("/api/auth/login", handlers.HandleLogin)
	r.POST("/api/auth/signup", handlers.HandleSignup)

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// Vendor routes
		api.GET("/profile", handlers.GetProfile)
		api.GET("/my-attendance", handlers.GetMyAttendance)
		api.GET("/my-assets", handlers.GetMyAssets)
		api.GET("/my-documents", handlers.GetMyDocuments)

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			// Vendor management
			admin.POST("/vendors", handlers.CreateVendor)
			admin.GET("/vendors", handlers.ListVendors)
			admin.GET("/vendors/:id", handlers.GetVendor)
			admin.PUT("/vendors/:id", handlers.UpdateVendor)

			// Asset management
			admin.POST("/assets", handlers.CreateAsset)
			admin.GET("/assets", handlers.ListAssets)
			admin.PUT("/assets/:id", handlers.UpdateAsset)
			admin.POST("/assets/:id/assign", handlers.AssignAsset)
			admin.POST("/assets/:id/return", handlers.ReturnAsset)

			// Document management
			admin.POST("/documents", handlers.UploadDocument)
			admin.GET("/documents", handlers.ListDocuments)
			admin.GET("/documents/:id", handlers.GetDocument)
			admin.DELETE("/documents/:id", handlers.DeleteDocument)

			// Attendance management
			admin.GET("/attendance", handlers.ListAttendance)
			admin.GET("/attendance/:vendorId", handlers.GetVendorAttendance)
		}
	}

	log.Println("Server starting on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
