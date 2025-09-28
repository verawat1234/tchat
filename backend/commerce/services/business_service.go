package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"tchat.dev/commerce/models"
	"tchat.dev/commerce/repository"
	sharedModels "tchat.dev/shared/models"
)

// businessService implements the BusinessService interface
type businessService struct {
	businessRepo repository.BusinessRepository
	productRepo  repository.ProductRepository
}

// NewBusinessService creates a new business service
func NewBusinessService(businessRepo repository.BusinessRepository, productRepo repository.ProductRepository) BusinessService {
	return &businessService{
		businessRepo: businessRepo,
		productRepo:  productRepo,
	}
}

// GetBusinesses retrieves businesses with filters and pagination
func (s *businessService) GetBusinesses(ctx context.Context, filters models.BusinessFilters, pagination models.Pagination, sort models.SortOptions) (*models.BusinessResponse, error) {
	businesses, total, err := s.businessRepo.FindBusinesses(ctx, filters, pagination, sort)
	if err != nil {
		return nil, fmt.Errorf("failed to find businesses: %w", err)
	}

	// Convert to shared models
	sharedBusinesses := make([]*sharedModels.Business, len(businesses))
	for i, business := range businesses {
		sharedBusinesses[i] = business
	}

	return &models.BusinessResponse{
		Businesses: sharedBusinesses,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: (total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize),
	}, nil
}

// GetBusiness retrieves a single business by ID
func (s *businessService) GetBusiness(ctx context.Context, id uuid.UUID) (*models.Business, error) {
	business, err := s.businessRepo.FindBusinessByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("business not found")
	}
	return business, nil
}

// CreateBusiness creates a new business
func (s *businessService) CreateBusiness(ctx context.Context, ownerID uuid.UUID, req *models.CreateBusinessRequest) (*models.Business, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, fmt.Errorf("business name is required")
	}

	if req.Address.Country == "" {
		return nil, fmt.Errorf("country is required")
	}

	// Create business entity
	business := &models.Business{
		ID:          uuid.New(),
		OwnerID:     ownerID,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Address: sharedModels.BusinessAddress{
			Street:     req.Address.Street,
			City:       req.Address.City,
			State:      req.Address.State,
			PostalCode: req.Address.PostalCode,
			Country:    req.Address.Country,
		},
		ContactInfo: sharedModels.BusinessContactInfo{
			Phone:   req.Contact.Phone,
			Email:   req.Contact.Email,
			Website: req.Contact.Website,
		},
		BusinessSettings: sharedModels.BusinessSettings{
			SupportedCurrencies: []string{"USD", "THB", "SGD"},
			SupportedLanguages:  []string{"en", "th", "zh"},
			ShippingCountries:   []string{req.Address.Country},
		},
		VerificationStatus: sharedModels.BusinessVerificationPending,
		IsActive:          true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Save to database
	err := s.businessRepo.CreateBusiness(ctx, business)
	if err != nil {
		return nil, fmt.Errorf("failed to create business: %w", err)
	}

	return business, nil
}

// UpdateBusiness updates an existing business
func (s *businessService) UpdateBusiness(ctx context.Context, id uuid.UUID, req *models.UpdateBusinessRequest) (*models.Business, error) {
	// Check if business exists
	_, err := s.businessRepo.FindBusinessByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("business not found")
	}

	// Prepare updates
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.Address != nil {
		updates["address_street"] = req.Address.Street
		updates["address_city"] = req.Address.City
		updates["address_state"] = req.Address.State
		updates["address_postal_code"] = req.Address.PostalCode
		updates["address_country"] = req.Address.Country
	}
	if req.Contact != nil {
		updates["contact_phone"] = req.Contact.Phone
		updates["contact_email"] = req.Contact.Email
		updates["contact_website"] = req.Contact.Website
	}
	if req.Settings != nil {
		if len(req.Settings.SupportedCurrencies) > 0 {
			updates["settings_currencies"] = req.Settings.SupportedCurrencies
		}
		if len(req.Settings.SupportedLanguages) > 0 {
			updates["settings_languages"] = req.Settings.SupportedLanguages
		}
		if len(req.Settings.ShippingCountries) > 0 {
			updates["settings_shipping"] = req.Settings.ShippingCountries
		}
		if req.Settings.TaxSettings != nil {
			updates["settings_tax"] = req.Settings.TaxSettings
		}
		if req.Settings.BusinessHours != nil {
			updates["settings_hours"] = req.Settings.BusinessHours
		}
		if len(req.Settings.PaymentMethods) > 0 {
			updates["settings_payments"] = req.Settings.PaymentMethods
		}
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	updates["updated_at"] = time.Now()

	// Update business
	err = s.businessRepo.UpdateBusiness(ctx, id, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update business: %w", err)
	}

	// Return updated business
	return s.businessRepo.FindBusinessByID(ctx, id)
}

// DeleteBusiness deletes a business (soft delete)
func (s *businessService) DeleteBusiness(ctx context.Context, id uuid.UUID) error {
	// Check if business exists
	_, err := s.businessRepo.FindBusinessByID(ctx, id)
	if err != nil {
		return fmt.Errorf("business not found")
	}

	// Soft delete
	err = s.businessRepo.DeleteBusiness(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete business: %w", err)
	}

	return nil
}

// GetBusinessProducts retrieves products for a specific business
func (s *businessService) GetBusinessProducts(ctx context.Context, businessID uuid.UUID, pagination models.Pagination, sort models.SortOptions) (*models.ProductResponse, error) {
	// Check if business exists
	_, err := s.businessRepo.FindBusinessByID(ctx, businessID)
	if err != nil {
		return nil, fmt.Errorf("business not found")
	}

	products, total, err := s.productRepo.FindProductsByBusinessID(ctx, businessID, pagination, sort)
	if err != nil {
		return nil, fmt.Errorf("failed to find products: %w", err)
	}

	// Convert to shared models
	sharedProducts := make([]*sharedModels.Product, len(products))
	for i, product := range products {
		sharedProducts[i] = product
	}

	return &models.ProductResponse{
		Products:   sharedProducts,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: (total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize),
	}, nil
}