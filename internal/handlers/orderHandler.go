package handlers

import (
	"net/http"
	// razorpay "github.com/razorpay/razorpay-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/config"
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

		productVariant, err := repository.GetVariantByVariantId(item.VariantID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Product does not exist": err})
			return
		}
		if productVariant.Stock < item.Quantity {
			utils.ResponseError(c, http.StatusBadRequest, "Product is out of stock", nil)
			return
		}

		// Deduct Stock
		productVariant.Stock = productVariant.Stock - item.Quantity
		// repository.UpdateProduct(productVariant)

		finalOrderItems = append(finalOrderItems, models.OrderItem{
			VariantID:    item.VariantID,
			Quantity:     item.Quantity,
			ProductPrice: productVariant.Price,
		})

		// calculate total
		total += productVariant.Price * float64(item.Quantity)

	}
	// Razorpay amount must be in paise
	amount := int(total * 100)

	data := map[string]interface{}{
		"amount":   amount,
		"currency": "INR",
		"receipt":  GenerateOrderNumber(),
	}

	body, err := config.RazorpayClient.Order.Create(data, nil)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, "Failed to create Razorpay order", nil)
		return
	}

	utils.ResponseSuccess(c, http.StatusOK, "Razorpay order created", body)
}

func VerifyPayment(c *gin.Context) {

	var req struct {
		RazorpayOrderID   string `json:"razorpay_order_id"`
		RazorpayPaymentID string `json:"razorpay_payment_id"`
		RazorpaySignature string `json:"razorpay_signature"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	body := req.RazorpayOrderID + "|" + req.RazorpayPaymentID

	if !utils.VerifySignature(body, req.RazorpaySignature) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
		return
	}

	val, _ := c.Get("userId")
	userId, _ := uuid.Parse(val.(string))

	cart, err := repository.GetCartByUserId(userId)
	if err != nil {
		utils.ResponseError(c, http.StatusBadRequest, "Cart not found", nil)
		return
	}

	var total float64
	var orderItems []models.OrderItem

	for _, item := range cart.CartItems {

		variant, _ := repository.GetVariantByVariantId(item.VariantID)

		orderItems = append(orderItems, models.OrderItem{
			VariantID:    item.VariantID,
			Quantity:     item.Quantity,
			ProductPrice: variant.Price,
		})

		total += variant.Price * float64(item.Quantity)

		variant.Stock -= item.Quantity
		repository.UpdateVariant(variant)
	}

	order := models.Order{
		UserID:      userId,
		OrderNumber: GenerateOrderNumber(),
		Status:      models.OrderPaid,
		Subtotal:    total,
		TotalAmount: total,
		OrderItems:  orderItems,
	}

	repository.CreateOrder(&order)

	repository.ClearCart(userId)

	payment := models.Payment{
		UserID:            userId,
		RazorpayOrderID:   req.RazorpayOrderID,
		RazorpayPaymentID: req.RazorpayPaymentID,
		Amount:            total,
		Status:            "success",
	}

	config.DB.Create(&payment)

	utils.ResponseSuccess(c, http.StatusOK, "Payment successful", order)
}

func CreateOrder(order *models.Order) {

}

func GenerateOrderNumber() string {
	return "5"
}
