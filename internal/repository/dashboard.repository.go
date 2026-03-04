package repository

import (
	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/config"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/dto"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/utils"
)

func GetDashboardCardData() (map[string]interface{}, error) {

	type DashboardStats struct {
		CurrentOrders      int64
		PreviousOrders     int64
		CurrentRevenue     float64
		PreviousRevenue    float64
		CancelledOrders    int64
		NewCustomers       int64
		LowStockProducts   int64
		TopSellingProducts int64
	}

	var stats DashboardStats

	err := config.DB.Raw(`
		SELECT
			COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '30 days') AS current_orders,
			COUNT(*) FILTER (
				WHERE created_at BETWEEN NOW() - INTERVAL '60 days' AND NOW() - INTERVAL '30 days'
			) AS previous_orders,

			COALESCE(SUM(total_amount) FILTER (
				WHERE status = 'completed' AND created_at >= NOW() - INTERVAL '30 days'
			), 0) AS current_revenue,

			COALESCE(SUM(total_amount) FILTER (
				WHERE status = 'completed' AND created_at BETWEEN 
					NOW() - INTERVAL '60 days' AND NOW() - INTERVAL '30 days'
			), 0) AS previous_revenue,

			COUNT(*) FILTER (
				WHERE status = 'cancelled' AND created_at >= NOW() - INTERVAL '30 days'
			) AS cancelled_orders,
	(
			SELECT COUNT(*) 
			FROM users 
			WHERE created_at >= NOW() - INTERVAL '30 days'
		) AS new_customers,

  (
    SELECT COUNT(*) 
    FROM product_variants 
    WHERE stock <= 10
  ) AS low_stock_products,

		FROM orders;
	`).Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	// Growth calculations
	orderGrowth := utils.CalculateGrowth(
		float64(stats.CurrentOrders),
		float64(stats.PreviousOrders),
	)

	revenueGrowth := utils.CalculateGrowth(
		stats.CurrentRevenue,
		stats.PreviousRevenue,
	)

	// Final response
	return map[string]interface{}{
		"orders": map[string]interface{}{
			"current":  stats.CurrentOrders,
			"previous": stats.PreviousOrders,
			"growth":   orderGrowth,
		},
		"revenue": map[string]interface{}{
			"current":  stats.CurrentRevenue,
			"previous": stats.PreviousRevenue,
			"growth":   revenueGrowth,
		},
		"customers": map[string]interface{}{
			"new": stats.NewCustomers,
		},
		"cancelled_orders":   stats.CancelledOrders,
		"low_stock_products": stats.LowStockProducts,
	}, nil
}

func TopSellingProducts() ([]dto.TopSellingProductResponse, error) {
	type Result struct {
		ProductID   uuid.UUID
		ProductName string
		TotalSold   int64
	}

	var results []Result
	err := config.DB.Raw(`
		SELECT pv.product_id, p.name as product_name, SUM(oi.quantity) as total_sold
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		JOIN product_variants pv ON pv.id = oi.variant_id
		JOIN products p ON p.id = pv.product_id
		WHERE o.created_at >= NOW() - INTERVAL '30 days'
		GROUP BY pv.product_id, p.name
		ORDER BY total_sold DESC
		LIMIT 10;
	`).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	var response []dto.TopSellingProductResponse

	for _, r := range results {
		response = append(response, dto.TopSellingProductResponse{
			ProductID:   r.ProductID,
			ProductName: r.ProductName,
			TotalSold:   r.TotalSold,
		})
	}

	return response, nil
}
