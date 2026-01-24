package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/repository"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/utils"
	"gorm.io/gorm"
)

type AddCartItemRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
}

func CreateCart(c *gin.Context) {
	val, ok := c.Get("userId")
	if !ok {
		utils.ResponseError(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	userId, err := uuid.Parse(val.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid userid": err})
		return
	}
	cart, err := repository.GetCartByUserId(userId)

	// cart already exists
	if err == nil && cart != nil {
		utils.ResponseSuccess(c, http.StatusOK, "Cart already exists", cart)
		return
	}

	// real DB error (not record not found)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to fetch cart", nil)
		return
	}

	// create cart
	err = repository.CreateCart(userId)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to create cart", nil)
		return
	}
	utils.ResponseSuccess(c, http.StatusCreated, "Cart created", nil)

}

func AddOrUpdateCartItem(c *gin.Context) {
	var req AddCartItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	val, ok := c.Get("userId")
	if !ok {
		utils.ResponseError(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	userId, err := uuid.Parse(val.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	cart, err := repository.GetCartByUserId(userId)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Cart Does not exist", nil)

		err = repository.CreateCart(userId)
		if err != nil {
			utils.ResponseError(c, http.StatusInternalServerError, "Failed to create cart", nil)
			return
		}
		// fetch newly created cart
		cart, err = repository.GetCartByUserId(userId)
		if err != nil {
			utils.ResponseError(c, http.StatusInternalServerError, "Failed to fetch cart", err)
			return
		}
	}
	product, err := repository.GetProductByUUID(req.ProductID)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Product does not exist", nil)
		return
	}
	if product.NumberOfStock < req.Quantity {
		utils.ResponseError(c, http.StatusBadRequest, "Product out of stock", nil)
		return
	}
	cartItem, err := repository.GetCartItem(cart.ID, req.ProductID)
	if err == nil {
		// Update quantity
		cartItem.Quantity += req.Quantity
		if err := repository.UpdateCartItem(cartItem); err != nil {
			utils.ResponseError(c, http.StatusInternalServerError, "Failed to update cart item", err)
			return
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Add new cart item
		err := repository.CreateCartItem(&models.CartItems{
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		})
		if err != nil {
			utils.ResponseError(c, http.StatusInternalServerError, "Failed to add cart item", err)
			return
		}
	} else {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to fetch cart item", err)
		return
	}

	// 1. Bind request body (product_id, quantity)
	// 2. Check if item exists in cart
	// 3. Update quantity OR insert new row
}
