package handlers

import "github.com/gin-gonic/gin"

func GetDashboardCardData(c *gin.Context) {
	// Simulate fetching data for dashboard cards
	cardData := map[string]interface{}{
		"totalRevenue":    1500,
		"totalOrders":    25000,
		"newOrders": 120,
		"pendingOrders": 30,
		"completedOrders": 200,
		"cancelledOrders": 20,
		"newCustomers": 50,
		"lowStockProducts": 10,
	}
	c.JSON(200, cardData)
}
