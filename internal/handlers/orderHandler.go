package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/repository"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/utils"
)

// use goroutines and channels to handle multiple tasks concurrently
// Scenario: A user places an Order. The Problem: You need to:
// Save Order to DB.
// Send an Email Confirmation.
// Send a WhatsApp Notification.
// Update Analytics Dashboard.

func Checkout(c *gin.Context) {
	val, ok := c.Get("userId")
	if !ok {
		utils.ResponseError(c, http.StatusBadRequest, "Unauthorized", nil)
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
		return
	}
	var total float64
	var finalOrderItems []models.OrderItem

	for _, item := range cart.CartItems {

		product, err := repository.GetProductByUUID(item.ProductID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Product does not exist": err})
			return
		}
		if product.Stock < item.Quantity {
			utils.ResponseError(c, http.StatusBadRequest, "Product is ut of stock", nil)
		}

		// Deduct Stock
		product.Stock = product.Stock - item.Quantity
		// repository.UpdateProduct(product)

		finalOrderItems = append(finalOrderItems, models.OrderItem{
			ProductID:    item.ProductID,
			Quantity:     item.Quantity,
			ProductPrice: product.BasePrice,
		})

		// calculate total
		total += product.BasePrice * float64(item.Quantity)

	}
	// 3. Create Order
	order := models.Order{
		UserID:      userId,
		OrderNumber: GenerateOrderNumber(),
		Subtotal:    total,
		Status:      models.OrderPending,
		OrderItems:  finalOrderItems,
		TotalAmount: total,
	}
	createdOrder, err := repository.CreateOrder(&order)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Order Failed", nil)
		return
	}
	utils.ResponseSuccess(c, http.StatusBadRequest, "Order Placed", createdOrder)
}

func CreateOrder(order *models.Order) {

}

func GenerateOrderNumber() string {
	return "5"
}
