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
		user.GET("/all", middleware.AuthMiddleware(), handlers.GetAllUsers)
		user.GET("/user/:id", middleware.AuthMiddleware(), handlers.GetUser)
		user.GET("/user", middleware.AuthMiddleware(), handlers.GetUserByEmail)
	}

	product := r.Group("/products")
	{
		product.GET("/all", handlers.GetAllProducts)

		productProtected := product.Group("/")
		productProtected.Use(middleware.AuthMiddleware())
		{
			productProtected.GET("/:id", handlers.GetProductById)
			productProtected.POST("/", handlers.CreateNewProduct)
			productProtected.PUT("/:id", handlers.UpdateProduct)
			productProtected.DELETE("/:id", handlers.DeleteProduct)
		}
	}
}
