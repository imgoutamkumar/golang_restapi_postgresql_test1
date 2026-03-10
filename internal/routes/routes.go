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

		product.GET("/filters", handlers.GetFilters)
		productProtected := product.Group("/")
		productProtected.Use(middleware.AuthMiddleware())
		{
			product.GET("/allProducts", handlers.GetAllProductsForAdmin)
			productProtected.GET("/:id", handlers.GetProductById)
			productProtected.POST("/create", handlers.CreateNewProduct)
			productProtected.PUT("/:id", handlers.UpdateProduct)
			productProtected.DELETE("/:id", handlers.DeleteProduct)
			productProtected.PATCH("/reorder-images", handlers.ProductImagesReorder)
			productProtected.POST("/brand/create", handlers.CreateNewBrand)
			productProtected.GET("/brand/all", handlers.GetAllBrands)
			productProtected.POST("/category/create", handlers.CreateNewCategory)
			productProtected.GET("/category/all", handlers.GetAllCategories)
			productProtected.GET("/attribute/all", handlers.GetAllAttributes)
			productProtected.GET("/attribute-value/:attributeId", handlers.GetAttributeValuesByAttributeId)
			productProtected.POST("/attribute/create", handlers.CreateNewAttribute)
			productProtected.POST("/attribute/value/create/:attributeId", handlers.CreateNewAttributeValue)
		}
	}

	// cart routes

	cart := r.Group(("/cart"))
	cartUserProtected := cart.Group("/")
	cartUserProtected.Use(middleware.AuthMiddleware())
	{
		cartUserProtected.GET("/items", handlers.GetAllCartItems)
		cartUserProtected.POST("/item", handlers.AddOrUpdateCartItem)
		cartUserProtected.POST("/checkout", handlers.Checkout)
		cartUserProtected.POST("/verify-payment", handlers.VerifyPayment)
	}

	// Admin routes (admin authorized routes)
	cartAdminProtected := cart.Group("/")
	cartAdminProtected.Use(middleware.AuthMiddleware(), middleware.IsAuthorized("Admin"))
	{
		cartAdminProtected.GET("/:userId", handlers.GetCart)
		cartAdminProtected.DELETE("/:userId", handlers.DeleteCart)
	}
}
