package models

import (
	"github.com/google/uuid"
	sharedModels "tchat.dev/shared/models"
)

// Re-export shared types for convenience
type ProductType = sharedModels.ProductType
type ProductPricing = sharedModels.ProductPricing
type ProductImage = sharedModels.ProductImage
type ProductVariant = sharedModels.ProductVariant

const (
	ProductTypePhysical = sharedModels.ProductTypePhysical
	ProductTypeDigital  = sharedModels.ProductTypeDigital
	ProductTypeService  = sharedModels.ProductTypeService
	ProductTypeMedia    = sharedModels.ProductTypeDigital // Alias for backward compatibility

	ProductStatusDraft    = sharedModels.ProductStatusDraft
	ProductStatusActive   = sharedModels.ProductStatusActive
	ProductStatusInactive = sharedModels.ProductStatusInactive
)

// Product is an alias to shared Product model
type Product = sharedModels.Product

// Helper functions for Product
func ProductIsPhysical(p *Product) bool {
	return p.Type == ProductTypePhysical
}

func ProductIsMedia(p *Product) bool {
	return p.Type == ProductTypeDigital
}

// GetDefaultPrice returns the first price in the pricing array (for backward compatibility)
func GetDefaultPrice(p *Product) (float64, string) {
	if len(p.Pricing) > 0 {
		priceFloat, _ := p.Pricing[0].Price.Float64()
		return priceFloat, p.Pricing[0].Currency
	}
	return 0.0, "USD"
}

// GetMediaContentID returns a UUID for media content (for backward compatibility)
// Since shared Product doesn't have MediaContentID, we return the product ID itself
func GetMediaContentID(p *Product) *uuid.UUID {
	id := p.ID
	return &id
}