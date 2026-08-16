package router

import (
	controller "core/api/controller"

	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	// _ "core/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Config(db *gorm.DB) http.Handler {
	router := gin.Default()
	controllers := controller.Create(db)
	router.Use(cors.New(
		cors.Config{
			AllowCredentials: true,
			AllowBrowserExtensions: true,
			
			AllowHeaders: []string{"*"},
			AllowOrigins: []string{
				"http://localhost:3000",
			},
			AllowMethods: []string{
				"GET",
				"POST",
			},
		},
	))

	router.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(
			swaggerFiles.Handler,
		),
	)

	v1 := router.Group("/api/v1")

	// Health
	v1.GET(
		"/health",
		func(c *gin.Context) {
			c.JSON(
				http.StatusOK,
				gin.H{
					"status": "healthy",
				},
			)
		},
	)

	v1.GET("/chat", controllers.GetMessages)
	v1.POST("/widget", controllers.CreateWidget)
	v1.POST("/ask", controllers.Ask)
	return router
}
