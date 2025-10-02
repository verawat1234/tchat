package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"tchat.dev/shared/logger"
)

// VideoCommerceHandler handles commerce operations for video content
type VideoCommerceHandler struct {
	logger *logger.TchatLogger
}

// NewVideoCommerceHandler creates a new video commerce handler instance
func NewVideoCommerceHandler(logger *logger.TchatLogger) *VideoCommerceHandler {
	return &VideoCommerceHandler{
		logger: logger,
	}
}

// VideoProduct represents a product associated with a video
type VideoProduct struct {
	ID          string  `json:"id"`
	VideoID     string  `json:"video_id"`
	ProductID   string  `json:"product_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	ImageURL    string  `json:"image_url"`
	StockStatus string  `json:"stock_status"` // in_stock, out_of_stock, pre_order
	Timestamp   int     `json:"timestamp"`    // Video timestamp where product appears
	CreatedAt   string  `json:"created_at"`
}

// VideoPurchase represents a video purchase transaction
type VideoPurchase struct {
	ID            string  `json:"id"`
	VideoID       string  `json:"video_id"`
	UserID        string  `json:"user_id"`
	Price         float64 `json:"price"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"` // pending, completed, failed
	PaymentMethod string  `json:"payment_method"`
	CreatedAt     string  `json:"created_at"`
	ExpiresAt     string  `json:"expires_at,omitempty"` // For rental videos
}

// VideoMonetization represents monetization settings for a video
type VideoMonetization struct {
	VideoID          string  `json:"video_id"`
	IsMonetized      bool    `json:"is_monetized"`
	PricingType      string  `json:"pricing_type"` // free, purchase, rental, subscription
	Price            float64 `json:"price"`
	Currency         string  `json:"currency"`
	RentalDuration   int     `json:"rental_duration,omitempty"` // Hours
	RevenueShare     float64 `json:"revenue_share"`             // Creator revenue share percentage
	AdEnabled        bool    `json:"ad_enabled"`
	SponsorshipDeals int     `json:"sponsorship_deals"`
}

// GetVideoPurchaseOptions retrieves purchase options for a video
// GET /api/v1/videos/:id/purchase
func (h *VideoCommerceHandler) GetVideoPurchaseOptions(c *gin.Context) {
	videoID := c.Param("id")
	if _, err := uuid.Parse(videoID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_video_id",
			"message": "Invalid video ID format",
		})
		return
	}

	// Mock data - would fetch from database
	monetization := VideoMonetization{
		VideoID:      videoID,
		IsMonetized:  true,
		PricingType:  "purchase",
		Price:        9.99,
		Currency:     "USD",
		RevenueShare: 70.0,
		AdEnabled:    true,
	}

	h.logger.WithFields(map[string]interface{}{
		"video_id":    videoID,
		"pricing_type": monetization.PricingType,
		"price":       monetization.Price,
	}).Info("Video purchase options retrieved")

	c.JSON(http.StatusOK, gin.H{
		"monetization": monetization,
		"purchase_url": "https://tchat.dev/purchase/video/" + videoID,
	})
}

// PurchaseVideo handles video purchase transaction
// POST /api/v1/videos/:id/purchase
func (h *VideoCommerceHandler) PurchaseVideo(c *gin.Context) {
	videoID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	var req struct {
		PaymentMethod string `json:"payment_method" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Payment method is required",
		})
		return
	}

	purchase := VideoPurchase{
		ID:            uuid.New().String(),
		VideoID:       videoID,
		UserID:        userID.(string),
		Price:         9.99,
		Currency:      "USD",
		Status:        "completed",
		PaymentMethod: req.PaymentMethod,
		CreatedAt:     "2024-01-01T00:00:00Z",
	}

	h.logger.WithFields(map[string]interface{}{
		"purchase_id":    purchase.ID,
		"video_id":       videoID,
		"user_id":        userID,
		"payment_method": req.PaymentMethod,
		"amount":         purchase.Price,
	}).Info("Video purchased")

	c.JSON(http.StatusCreated, gin.H{
		"purchase": purchase,
		"message":  "Video purchased successfully",
		"access":   "unlimited",
	})
}

// RentVideo handles video rental transaction
// POST /api/v1/videos/:id/rent
func (h *VideoCommerceHandler) RentVideo(c *gin.Context) {
	videoID := c.Param("id")
	userID, _ := c.Get("user_id")

	var req struct {
		PaymentMethod  string `json:"payment_method" binding:"required"`
		RentalDuration int    `json:"rental_duration"` // Hours (default 48)
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Payment method is required",
		})
		return
	}

	if req.RentalDuration == 0 {
		req.RentalDuration = 48 // Default 48 hours
	}

	purchase := VideoPurchase{
		ID:            uuid.New().String(),
		VideoID:       videoID,
		UserID:        userID.(string),
		Price:         3.99,
		Currency:      "USD",
		Status:        "completed",
		PaymentMethod: req.PaymentMethod,
		CreatedAt:     "2024-01-01T00:00:00Z",
		ExpiresAt:     "2024-01-03T00:00:00Z", // Would calculate based on rental duration
	}

	h.logger.WithFields(map[string]interface{}{
		"purchase_id":      purchase.ID,
		"video_id":         videoID,
		"user_id":          userID,
		"rental_duration":  req.RentalDuration,
		"payment_method":   req.PaymentMethod,
	}).Info("Video rented")

	c.JSON(http.StatusCreated, gin.H{
		"purchase":        purchase,
		"message":         "Video rented successfully",
		"rental_duration": req.RentalDuration,
		"expires_at":      purchase.ExpiresAt,
	})
}

// GetVideoProducts retrieves products featured in a video
// GET /api/v1/videos/:id/products
func (h *VideoCommerceHandler) GetVideoProducts(c *gin.Context) {
	videoID := c.Param("id")

	// Mock data - would fetch from database
	products := []VideoProduct{
		{
			ID:          uuid.New().String(),
			VideoID:     videoID,
			ProductID:   uuid.New().String(),
			Name:        "Product shown in video",
			Description: "High-quality product featured",
			Price:       29.99,
			Currency:    "USD",
			ImageURL:    "https://cdn.tchat.dev/products/sample.jpg",
			StockStatus: "in_stock",
			Timestamp:   120, // Appears at 2:00
			CreatedAt:   "2024-01-01T00:00:00Z",
		},
	}

	h.logger.WithFields(map[string]interface{}{
		"video_id":      videoID,
		"product_count": len(products),
	}).Info("Video products retrieved")

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    len(products),
	})
}

// AddVideoProduct associates a product with a video
// POST /api/v1/videos/:id/products
func (h *VideoCommerceHandler) AddVideoProduct(c *gin.Context) {
	videoID := c.Param("id")
	userID, _ := c.Get("user_id")

	var req struct {
		ProductID   string  `json:"product_id" binding:"required"`
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Price       float64 `json:"price" binding:"required"`
		Currency    string  `json:"currency" binding:"required"`
		ImageURL    string  `json:"image_url"`
		Timestamp   int     `json:"timestamp"` // Video timestamp
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Required fields missing",
		})
		return
	}

	product := VideoProduct{
		ID:          uuid.New().String(),
		VideoID:     videoID,
		ProductID:   req.ProductID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Currency:    req.Currency,
		ImageURL:    req.ImageURL,
		StockStatus: "in_stock",
		Timestamp:   req.Timestamp,
		CreatedAt:   "2024-01-01T00:00:00Z",
	}

	h.logger.WithFields(map[string]interface{}{
		"video_id":   videoID,
		"product_id": req.ProductID,
		"user_id":    userID,
		"timestamp":  req.Timestamp,
	}).Info("Product added to video")

	c.JSON(http.StatusCreated, gin.H{
		"product": product,
		"message": "Product added to video successfully",
	})
}

// GetUserPurchases retrieves a user's video purchases
// GET /api/v1/users/purchases
func (h *VideoCommerceHandler) GetUserPurchases(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not authenticated",
		})
		return
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Mock data - would fetch from database
	purchases := []VideoPurchase{
		{
			ID:            uuid.New().String(),
			VideoID:       uuid.New().String(),
			UserID:        userID.(string),
			Price:         9.99,
			Currency:      "USD",
			Status:        "completed",
			PaymentMethod: "credit_card",
			CreatedAt:     "2024-01-01T00:00:00Z",
		},
	}

	h.logger.WithFields(map[string]interface{}{
		"user_id": userID,
		"page":    page,
		"limit":   limit,
	}).Info("User video purchases retrieved")

	c.JSON(http.StatusOK, gin.H{
		"purchases": purchases,
		"page":      page,
		"limit":     limit,
		"total":     len(purchases),
	})
}

// UpdateVideoMonetization updates monetization settings for a video
// PUT /api/v1/videos/:id/monetization
func (h *VideoCommerceHandler) UpdateVideoMonetization(c *gin.Context) {
	videoID := c.Param("id")
	userID, _ := c.Get("user_id")

	var req VideoMonetization
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Invalid monetization settings",
		})
		return
	}

	req.VideoID = videoID

	h.logger.WithFields(map[string]interface{}{
		"video_id":     videoID,
		"user_id":      userID,
		"is_monetized": req.IsMonetized,
		"pricing_type": req.PricingType,
		"price":        req.Price,
	}).Info("Video monetization updated")

	c.JSON(http.StatusOK, gin.H{
		"monetization": req,
		"message":      "Monetization settings updated successfully",
	})
}

// GetVideoRevenue retrieves revenue statistics for a video
// GET /api/v1/videos/:id/revenue
func (h *VideoCommerceHandler) GetVideoRevenue(c *gin.Context) {
	videoID := c.Param("id")
	userID, _ := c.Get("user_id")

	// Mock data - would calculate from database
	revenue := gin.H{
		"video_id":            videoID,
		"total_revenue":       1458.75,
		"creator_revenue":     1021.13, // 70% revenue share
		"platform_revenue":    437.62,
		"currency":            "USD",
		"purchases":           147,
		"rentals":             89,
		"ad_revenue":          234.50,
		"sponsorship_revenue": 500.00,
		"period":              "all_time",
	}

	h.logger.WithFields(map[string]interface{}{
		"video_id":       videoID,
		"user_id":        userID,
		"total_revenue":  revenue["total_revenue"],
	}).Info("Video revenue retrieved")

	c.JSON(http.StatusOK, revenue)
}

// GetMarketplaceVideos retrieves videos available for purchase
// GET /api/v1/marketplace/videos
func (h *VideoCommerceHandler) GetMarketplaceVideos(c *gin.Context) {
	// Filters
	category := c.Query("category")
	priceMin, _ := strconv.ParseFloat(c.DefaultQuery("price_min", "0"), 64)
	priceMax, _ := strconv.ParseFloat(c.DefaultQuery("price_max", "100"), 64)
	pricingType := c.DefaultQuery("pricing_type", "all")

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Mock data - would fetch from database
	videos := []gin.H{
		{
			"video_id":     uuid.New().String(),
			"title":        "Premium Video Content",
			"price":        9.99,
			"currency":     "USD",
			"pricing_type": "purchase",
			"creator":      "Creator Name",
			"views":        15678,
			"rating":       4.7,
		},
	}

	h.logger.WithFields(map[string]interface{}{
		"category":     category,
		"price_range":  []float64{priceMin, priceMax},
		"pricing_type": pricingType,
		"page":         page,
		"limit":        limit,
	}).Info("Marketplace videos retrieved")

	c.JSON(http.StatusOK, gin.H{
		"videos": videos,
		"page":   page,
		"limit":  limit,
		"total":  len(videos),
		"filters": gin.H{
			"category":     category,
			"price_min":    priceMin,
			"price_max":    priceMax,
			"pricing_type": pricingType,
		},
	})
}