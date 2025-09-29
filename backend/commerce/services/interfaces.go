package services

import (
	"context"

	"github.com/google/uuid"
	"tchat.dev/commerce/models"
	sharedModels "tchat.dev/shared/models"
)

// BusinessService defines the interface for business operations
type BusinessService interface {
	GetBusinesses(ctx context.Context, filters models.BusinessFilters, pagination models.Pagination, sort models.SortOptions) (*models.BusinessResponse, error)
	GetBusiness(ctx context.Context, id uuid.UUID) (*sharedModels.Business, error)
	CreateBusiness(ctx context.Context, ownerID uuid.UUID, req *models.CreateBusinessRequest) (*sharedModels.Business, error)
	UpdateBusiness(ctx context.Context, id uuid.UUID, req *models.UpdateBusinessRequest) (*sharedModels.Business, error)
	DeleteBusiness(ctx context.Context, id uuid.UUID) error
	GetBusinessProducts(ctx context.Context, businessID uuid.UUID, pagination models.Pagination, sort models.SortOptions) (*models.ProductResponse, error)
}

// ProductService defines the interface for product operations
type ProductService interface {
	GetProducts(ctx context.Context, filters models.ProductFilters, pagination models.Pagination, sort models.SortOptions) (*models.ProductResponse, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error)
	CreateProduct(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req *models.UpdateProductRequest) (*models.Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

// BusinessRepository defines the interface for business data access
type BusinessRepository interface {
	FindBusinesses(ctx context.Context, filters models.BusinessFilters, pagination models.Pagination, sort models.SortOptions) ([]*sharedModels.Business, int64, error)
	FindBusinessByID(ctx context.Context, id uuid.UUID) (*sharedModels.Business, error)
	CreateBusiness(ctx context.Context, business *sharedModels.Business) error
	UpdateBusiness(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteBusiness(ctx context.Context, id uuid.UUID) error
}

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	FindProducts(ctx context.Context, filters models.ProductFilters, pagination models.Pagination, sort models.SortOptions) ([]*models.Product, int64, error)
	FindProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	FindProductsByBusinessID(ctx context.Context, businessID uuid.UUID, pagination models.Pagination, sort models.SortOptions) ([]*models.Product, int64, error)
	CreateProduct(ctx context.Context, product *models.Product) error
	UpdateProduct(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}