package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/handlers"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/middleware"
)

func SetRoutes(r *gin.Engine) {

	user := r.Group("/users")
	{
		user.POST("/register", handlers.Register)
		user.POST("/login", handlers.Login)

		user.GET("/", middleware.AuthMiddleware(), handlers.GetAllUsers)
	}

	product := r.Group("/products")
	{
		product.GET("/", handlers.GetAllProducts)

		productProtected := product.Group("/")
		productProtected.Use(middleware.AuthMiddleware())
		{
			productProtected.GET("/:id", handlers.GetProduct)
			productProtected.POST("/", handlers.CreateNewProduct)
			productProtected.PUT("/:id", handlers.UpdateProduct)
			productProtected.DELETE("/:id", handlers.DeleteProduct)
		}
	}
}
