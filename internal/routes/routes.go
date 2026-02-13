package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/handlers"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/middleware"
)

func SetRoutes(r *gin.Engine) {

	//user routes

	user := r.Group("/users")
	{
		user.POST("/register", handlers.Register)
		user.POST("/login", handlers.Login)
		user.GET("/all", middleware.AuthMiddleware(), handlers.GetAllUsers)
		// user.GET("/:id", middleware.AuthMiddleware(), handlers.GetUser)
		user.GET("/user", middleware.AuthMiddleware(), handlers.GetUserByEmail)
		user.POST("/password/forgot", handlers.SendOtpRequest)
		user.POST("/password/verify", handlers.VerifyOtpRequest)
		user.POST("/password/reset", handlers.PasswordReset)
	}
	productUser := user.Group("/")
	productUser.Use(middleware.AuthMiddleware())
	{
		productUser.GET("/profile", handlers.GetUser)
	}

	// product routes

	product := r.Group("/products")
	{
		product.GET("/all", handlers.GetAllProducts)

		productProtected := product.Group("/")
		productProtected.Use(middleware.AuthMiddleware())
		{
			productProtected.GET("/:id", handlers.GetProductById)
			productProtected.POST("/create", handlers.CreateNewProduct)
			productProtected.PUT("/:id", handlers.UpdateProduct)
			productProtected.DELETE("/:id", handlers.DeleteProduct)
			productProtected.PATCH("/reorder-images", handlers.ProductImagesReorder)
		}
	}

	// cart routes

	cart := r.Group(("/cart"))
	cartUserProtected := cart.Group("/")
	cartUserProtected.Use(middleware.AuthMiddleware())
	{
		cartUserProtected.GET("/items", handlers.GetAllCartItems)
		cartUserProtected.POST("/item", handlers.AddOrUpdateCartItem)
	}

	// Admin routes (admin authorized routes)
	cartAdminProtected := cart.Group("/")
	cartAdminProtected.Use(middleware.AuthMiddleware(), middleware.IsAuthorized("Admin"))
	{
		cartAdminProtected.GET("/:userId", handlers.GetCart)
		cartAdminProtected.DELETE("/:userId", handlers.DeleteCart)
	}
}
