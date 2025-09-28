package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"tchat.dev/commerce/models"
	"tchat.dev/commerce/repository"
	sharedModels "tchat.dev/shared/models"
)

type ReviewService interface {
	CreateReview(ctx context.Context, req models.CreateReviewRequest) (*models.Review, error)
	GetReview(ctx context.Context, id uuid.UUID) (*models.Review, error)
	UpdateReview(ctx context.Context, id uuid.UUID, req models.UpdateReviewRequest) (*models.Review, error)
	DeleteReview(ctx context.Context, id uuid.UUID) error
	ListReviews(ctx context.Context, filters ReviewFilters, pagination models.Pagination) (*models.ReviewResponse, error)
	GetReviewsByProduct(ctx context.Context, productID uuid.UUID, pagination models.Pagination) (*models.ReviewResponse, error)
	GetReviewsByBusiness(ctx context.Context, businessID uuid.UUID, pagination models.Pagination) (*models.ReviewResponse, error)
	MarkReviewHelpful(ctx context.Context, reviewID, userID uuid.UUID, helpful bool) error
	ReportReview(ctx context.Context, reviewID, userID uuid.UUID, reason, comment string) error
	ModerateReview(ctx context.Context, reviewID uuid.UUID, status models.ReviewStatus, moderatorID uuid.UUID, notes string) error
	GetAverageRating(ctx context.Context, targetType models.ReviewType, targetID uuid.UUID) (decimal.Decimal, int64, error)
}

type ReviewFilters struct {
	Type       *models.ReviewType   `json:"type,omitempty"`
	Status     *models.ReviewStatus `json:"status,omitempty"`
	ProductID  *uuid.UUID           `json:"productId,omitempty"`
	BusinessID *uuid.UUID           `json:"businessId,omitempty"`
	OrderID    *uuid.UUID           `json:"orderId,omitempty"`
	UserID     *uuid.UUID           `json:"userId,omitempty"`
	MinRating  *decimal.Decimal     `json:"minRating,omitempty"`
	MaxRating  *decimal.Decimal     `json:"maxRating,omitempty"`
	Search     *string              `json:"search,omitempty"`
}

type reviewService struct {
	reviewRepo repository.ReviewRepository
	db         *gorm.DB
}

func NewReviewService(reviewRepo repository.ReviewRepository, db *gorm.DB) ReviewService {
	return &reviewService{
		reviewRepo: reviewRepo,
		db:         db,
	}
}

func (s *reviewService) CreateReview(ctx context.Context, req models.CreateReviewRequest) (*models.Review, error) {
	// Validate that user hasn't already reviewed this target
	var existingCount int64
	query := s.db.Model(&models.Review{}).Where("user_id = ? AND type = ?", req.UserID, req.Type)

	if req.ProductID != nil {
		query = query.Where("product_id = ?", *req.ProductID)
	}
	if req.BusinessID != nil {
		query = query.Where("business_id = ?", *req.BusinessID)
	}
	if req.OrderID != nil {
		query = query.Where("order_id = ?", *req.OrderID)
	}

	if err := query.Count(&existingCount).Error; err != nil {
		return nil, fmt.Errorf("failed to check existing reviews: %w", err)
	}

	if existingCount > 0 {
		return nil, fmt.Errorf("user has already reviewed this %s", req.Type)
	}

	review := &models.Review{
		Type:       req.Type,
		ProductID:  req.ProductID,
		BusinessID: req.BusinessID,
		OrderID:    req.OrderID,
		UserID:     req.UserID,
		UserName:   req.UserName,
		UserEmail:  req.UserEmail,
		Rating:     req.Rating,
		Title:      req.Title,
		Content:    req.Content,
		Images:     req.Images,
		Status:     models.ReviewStatusPending,
		IsVerified: false,
	}

	if err := s.db.WithContext(ctx).Create(review).Error; err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	// Publish review created event
	eventData := map[string]interface{}{
		"review_id":   review.ID,
		"product_id":  review.ProductID,
		"business_id": review.BusinessID,
		"rating":      review.Rating,
		"type":        review.Type,
	}

	event := &sharedModels.Event{
		Type:          sharedModels.EventTypeReviewCreated,
		Category:      sharedModels.EventCategoryDomain,
		Severity:      sharedModels.SeverityInfo,
		Subject:       "Review Created",
		Description:   "A new review has been created",
		AggregateType: "review",
		AggregateID:   review.ID.String(),
	}

	if err := event.MarshalData(eventData); err != nil {
		fmt.Printf("Failed to marshal event data: %v\n", err)
	}

	if err := s.db.WithContext(ctx).Create(event).Error; err != nil {
		// Log error but don't fail the review creation
		fmt.Printf("Failed to create review event: %v\n", err)
	}

	return review, nil
}

func (s *reviewService) GetReview(ctx context.Context, id uuid.UUID) (*models.Review, error) {
	var review models.Review
	if err := s.db.WithContext(ctx).First(&review, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("review not found")
		}
		return nil, fmt.Errorf("failed to get review: %w", err)
	}

	return &review, nil
}

func (s *reviewService) UpdateReview(ctx context.Context, id uuid.UUID, req models.UpdateReviewRequest) (*models.Review, error) {
	review, err := s.GetReview(ctx, id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.Rating != nil {
		updates["rating"] = *req.Rating
	}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Images != nil {
		updates["images"] = req.Images
	}

	if len(updates) > 0 {
		if err := s.db.WithContext(ctx).Model(review).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("failed to update review: %w", err)
		}
	}

	return review, nil
}

func (s *reviewService) DeleteReview(ctx context.Context, id uuid.UUID) error {
	result := s.db.WithContext(ctx).Delete(&models.Review{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete review: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("review not found")
	}
	return nil
}

func (s *reviewService) ListReviews(ctx context.Context, filters ReviewFilters, pagination models.Pagination) (*models.ReviewResponse, error) {
	query := s.db.WithContext(ctx).Model(&models.Review{})

	// Apply filters
	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}
	if filters.ProductID != nil {
		query = query.Where("product_id = ?", *filters.ProductID)
	}
	if filters.BusinessID != nil {
		query = query.Where("business_id = ?", *filters.BusinessID)
	}
	if filters.OrderID != nil {
		query = query.Where("order_id = ?", *filters.OrderID)
	}
	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}
	if filters.MinRating != nil {
		query = query.Where("rating >= ?", *filters.MinRating)
	}
	if filters.MaxRating != nil {
		query = query.Where("rating <= ?", *filters.MaxRating)
	}
	if filters.Search != nil && *filters.Search != "" {
		searchTerm := "%" + *filters.Search + "%"
		query = query.Where("title ILIKE ? OR content ILIKE ?", searchTerm, searchTerm)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count reviews: %w", err)
	}

	// Apply pagination and get results
	var reviews []*models.Review
	offset := (pagination.Page - 1) * pagination.PageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pagination.PageSize).Find(&reviews).Error; err != nil {
		return nil, fmt.Errorf("failed to list reviews: %w", err)
	}

	totalPages := (total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize)

	return &models.ReviewResponse{
		Reviews:    reviews,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *reviewService) GetReviewsByProduct(ctx context.Context, productID uuid.UUID, pagination models.Pagination) (*models.ReviewResponse, error) {
	filters := ReviewFilters{
		Type:      &[]models.ReviewType{models.ReviewTypeProduct}[0],
		ProductID: &productID,
		Status:    &[]models.ReviewStatus{models.ReviewStatusApproved}[0],
	}
	return s.ListReviews(ctx, filters, pagination)
}

func (s *reviewService) GetReviewsByBusiness(ctx context.Context, businessID uuid.UUID, pagination models.Pagination) (*models.ReviewResponse, error) {
	filters := ReviewFilters{
		Type:       &[]models.ReviewType{models.ReviewTypeBusiness}[0],
		BusinessID: &businessID,
		Status:     &[]models.ReviewStatus{models.ReviewStatusApproved}[0],
	}
	return s.ListReviews(ctx, filters, pagination)
}

func (s *reviewService) MarkReviewHelpful(ctx context.Context, reviewID, userID uuid.UUID, helpful bool) error {
	// Check if user has already voted
	var existing models.ReviewHelpful
	err := s.db.WithContext(ctx).Where("review_id = ? AND user_id = ?", reviewID, userID).First(&existing).Error

	if err == nil {
		// Update existing vote
		if err := s.db.WithContext(ctx).Model(&existing).Update("is_helpful", helpful).Error; err != nil {
			return fmt.Errorf("failed to update helpful vote: %w", err)
		}
	} else if err == gorm.ErrRecordNotFound {
		// Create new vote
		vote := &models.ReviewHelpful{
			ReviewID:  reviewID,
			UserID:    userID,
			IsHelpful: helpful,
		}
		if err := s.db.WithContext(ctx).Create(vote).Error; err != nil {
			return fmt.Errorf("failed to create helpful vote: %w", err)
		}
	} else {
		return fmt.Errorf("failed to check existing vote: %w", err)
	}

	// Update review helpful/not helpful counts
	return s.updateReviewHelpfulCounts(ctx, reviewID)
}

func (s *reviewService) ReportReview(ctx context.Context, reviewID, userID uuid.UUID, reason, comment string) error {
	report := &models.ReviewReport{
		ReviewID: reviewID,
		UserID:   userID,
		Reason:   reason,
		Comment:  comment,
		Status:   "pending",
	}

	if err := s.db.WithContext(ctx).Create(report).Error; err != nil {
		return fmt.Errorf("failed to create review report: %w", err)
	}

	// Update review report count
	if err := s.db.WithContext(ctx).Model(&models.Review{}).Where("id = ?", reviewID).
		UpdateColumn("report_count", gorm.Expr("report_count + 1")).Error; err != nil {
		return fmt.Errorf("failed to update report count: %w", err)
	}

	return nil
}

func (s *reviewService) ModerateReview(ctx context.Context, reviewID uuid.UUID, status models.ReviewStatus, moderatorID uuid.UUID, notes string) error {
	updates := map[string]interface{}{
		"status":           status,
		"moderation_notes": notes,
		"moderated_by":     moderatorID,
		"moderated_at":     "NOW()",
	}

	if err := s.db.WithContext(ctx).Model(&models.Review{}).Where("id = ?", reviewID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to moderate review: %w", err)
	}

	return nil
}

func (s *reviewService) GetAverageRating(ctx context.Context, targetType models.ReviewType, targetID uuid.UUID) (decimal.Decimal, int64, error) {
	var result struct {
		AvgRating decimal.Decimal
		Count     int64
	}

	query := s.db.WithContext(ctx).Model(&models.Review{}).
		Select("COALESCE(AVG(rating), 0) as avg_rating, COUNT(*) as count").
		Where("type = ? AND status = ?", targetType, models.ReviewStatusApproved)

	switch targetType {
	case models.ReviewTypeProduct:
		query = query.Where("product_id = ?", targetID)
	case models.ReviewTypeBusiness:
		query = query.Where("business_id = ?", targetID)
	case models.ReviewTypeOrder:
		query = query.Where("order_id = ?", targetID)
	}

	if err := query.Scan(&result).Error; err != nil {
		return decimal.Zero, 0, fmt.Errorf("failed to get average rating: %w", err)
	}

	return result.AvgRating, result.Count, nil
}

func (s *reviewService) updateReviewHelpfulCounts(ctx context.Context, reviewID uuid.UUID) error {
	var helpfulCount, notHelpfulCount int64

	if err := s.db.WithContext(ctx).Model(&models.ReviewHelpful{}).
		Where("review_id = ? AND is_helpful = ?", reviewID, true).Count(&helpfulCount).Error; err != nil {
		return fmt.Errorf("failed to count helpful votes: %w", err)
	}

	if err := s.db.WithContext(ctx).Model(&models.ReviewHelpful{}).
		Where("review_id = ? AND is_helpful = ?", reviewID, false).Count(&notHelpfulCount).Error; err != nil {
		return fmt.Errorf("failed to count not helpful votes: %w", err)
	}

	updates := map[string]interface{}{
		"helpful_count":     helpfulCount,
		"not_helpful_count": notHelpfulCount,
	}

	if err := s.db.WithContext(ctx).Model(&models.Review{}).Where("id = ?", reviewID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update helpful counts: %w", err)
	}

	return nil
}