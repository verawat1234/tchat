/**
 * Commerce API Types
 *
 * Type definitions for the commerce microservice API endpoints.
 * These types correspond to the Go backend models and API responses.
 */

// ===== Core Types =====

export interface UUID {
  value: string;
}

export interface Decimal {
  value: string;
}

export interface Timestamp {
  value: string;
}

// ===== Business Types =====

export interface BusinessAddress {
  street: string;
  city: string;
  state: string;
  postalCode: string;
  country: string;
  coordinates?: {
    latitude: number;
    longitude: number;
  };
}

export interface BusinessContactInfo {
  email: string;
  phone: string;
  website?: string;
  social?: {
    facebook?: string;
    twitter?: string;
    instagram?: string;
  };
}

export interface BusinessSettings {
  currency: string;
  timezone: string;
  language: string;
  supportedCurrencies: string[];
  supportedLanguages: string[];
  features: {
    allowReviews: boolean;
    allowWishlists: boolean;
    requireApproval: boolean;
  };
}

export enum BusinessVerificationStatus {
  PENDING = 'pending',
  VERIFIED = 'verified',
  REJECTED = 'rejected',
  SUSPENDED = 'suspended'
}

export interface Business {
  id: string;
  name: string;
  description: string;
  category: string;
  address: BusinessAddress;
  contact: BusinessContactInfo;
  settings: BusinessSettings;
  verificationStatus: BusinessVerificationStatus;
  rating: string; // Decimal as string
  reviewCount: number;
  isActive: boolean;
  isFeatured: boolean;
  createdAt: string;
  updatedAt: string;
}

// ===== Product Types =====

export enum ProductStatus {
  DRAFT = 'draft',
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  ARCHIVED = 'archived'
}

export enum ProductType {
  PHYSICAL = 'physical',
  DIGITAL = 'digital',
  SERVICE = 'service',
  SUBSCRIPTION = 'subscription'
}

export interface ProductInventory {
  trackQuantity: boolean;
  quantity: number;
  lowStockThreshold: number;
  allowBackorder: boolean;
  maxPerOrder: number;
}

export interface ProductShipping {
  weight: string; // Decimal as string
  dimensions: {
    length: string; // Decimal as string
    width: string;  // Decimal as string
    height: string; // Decimal as string
  };
  shippingClass: string;
  handlingTime: string;
}

export interface ProductSEO {
  metaTitle: string;
  metaDescription: string;
  slug: string;
  keywords: string[];
}

export interface Product {
  id: string;
  businessId: string;
  name: string;
  description: string;
  category: string;
  type: ProductType;
  status: ProductStatus;
  price: string; // Decimal as string
  currency: string;
  compareAtPrice?: string; // Decimal as string
  costPrice?: string; // Decimal as string
  sku: string;
  barcode?: string;
  images: string[];
  tags: string[];
  variants: Record<string, any>;
  inventory: ProductInventory;
  shipping: ProductShipping;
  seo: ProductSEO;
  rating: string; // Decimal as string
  reviewCount: number;
  salesCount: number;
  viewCount: number;
  isActive: boolean;
  isFeatured: boolean;
  createdAt: string;
  updatedAt: string;
}

// ===== Cart Types =====

export enum CartStatus {
  ACTIVE = 'active',
  ABANDONED = 'abandoned',
  CONVERTED = 'converted',
  EXPIRED = 'expired'
}

export interface CartItem {
  id: string;
  cartId: string;
  productId: string;
  variantId?: string;
  quantity: number;
  unitPrice: string; // Decimal as string
  totalPrice: string; // Decimal as string
  productName: string;
  productImage?: string;
  productSku: string;
  variantName?: string;
  isGift: boolean;
  giftMessage?: string;
  isAvailable: boolean;
  stockQuantity: number;
  maxQuantity: number;
  addedAt: string;
}

export interface Cart {
  id: string;
  userId?: string;
  sessionId?: string;
  businessId: string;
  status: CartStatus;
  items: CartItem[];
  itemCount: number;
  subtotalAmount: string; // Decimal as string
  taxAmount: string; // Decimal as string
  shippingAmount: string; // Decimal as string
  discountAmount: string; // Decimal as string
  totalAmount: string; // Decimal as string
  currency: string;
  couponCode?: string;
  couponDiscount?: string; // Decimal as string
  notes?: string;
  expiresAt?: string;
  lastActivity: string;
  createdAt: string;
  updatedAt: string;
}

export interface CartValidationIssue {
  type: string;
  message: string;
  productId: string;
  severity: 'error' | 'warning' | 'info';
}

export interface CartValidation {
  isValid: boolean;
  issues: CartValidationIssue[];
  totalItems: number;
  totalValue: string; // Decimal as string
  currency: string;
  estimatedShipping: string; // Decimal as string
  estimatedTax: string; // Decimal as string
  estimatedTotal: string; // Decimal as string
  validatedAt: string;
}

export interface CartAbandonmentTracking {
  id: string;
  cartId: string;
  abandonmentStage: string;
  lastPageVisited: string;
  emailSent: boolean;
  smsent: boolean;
  isRecovered: boolean;
  recoveredAt?: string;
  createdAt: string;
  updatedAt: string;
}

// ===== Category Types =====

export enum CategoryStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  ARCHIVED = 'archived'
}

export enum CategoryType {
  GLOBAL = 'global',
  BUSINESS = 'business'
}

export interface CategoryImage {
  url: string;
  alt: string;
  width: number;
  height: number;
}

export interface CategorySEO {
  metaTitle: string;
  metaDescription: string;
  slug: string;
  keywords: string[];
}

export interface CategoryAttribute {
  name: string;
  type: string;
  required: boolean;
  options?: string[];
}

export interface Category {
  id: string;
  businessId?: string;
  parentId?: string;
  name: string;
  description: string;
  shortDescription: string;
  status: CategoryStatus;
  type: CategoryType;
  icon: string;
  image?: CategoryImage;
  color: string;
  level: number;
  sortOrder: number;
  productCount: number;
  activeProductCount: number;
  childrenCount: number;
  isVisible: boolean;
  isFeatured: boolean;
  allowProducts: boolean;
  seo: CategorySEO;
  attributes: CategoryAttribute[];
  createdAt: string;
  updatedAt: string;
}

export interface CategoryView {
  id: string;
  categoryId: string;
  userId?: string;
  sessionId: string;
  ipAddress?: string;
  userAgent?: string;
  referrer?: string;
  viewedAt: string;
}

export interface CategoryAnalytics {
  totalViews: number;
  uniqueVisitors: number;
  topReferrers: Array<{
    referrer: string;
    count: number;
  }>;
  viewsByDay: Array<{
    date: string;
    views: number;
  }>;
  productClicks: number;
  conversionRate: string; // Decimal as string
}

// ===== Review Types =====

export enum ReviewType {
  PRODUCT = 'product',
  BUSINESS = 'business',
  ORDER = 'order'
}

export enum ReviewStatus {
  PENDING = 'pending',
  APPROVED = 'approved',
  REJECTED = 'rejected',
  FLAGGED = 'flagged'
}

export interface ReviewImage {
  url: string;
  alt: string;
  width: number;
  height: number;
}

export interface Review {
  id: string;
  type: ReviewType;
  productId?: string;
  businessId?: string;
  orderId?: string;
  userId: string;
  userName: string;
  userEmail: string;
  userAvatar?: string;
  rating: string; // Decimal as string
  title: string;
  content: string;
  images: ReviewImage[];
  status: ReviewStatus;
  helpfulCount: number;
  reportCount: number;
  verifiedPurchase: boolean;
  response?: string;
  respondedAt?: string;
  createdAt: string;
  updatedAt: string;
}

// ===== Wishlist Types =====

export enum WishlistType {
  PERSONAL = 'personal',
  SHARED = 'shared',
  PUBLIC = 'public'
}

export enum WishlistPrivacy {
  PRIVATE = 'private',
  SHARED = 'shared',
  PUBLIC = 'public'
}

export interface WishlistItem {
  id: string;
  wishlistId: string;
  productId: string;
  variantId?: string;
  quantity: number;
  note: string;
  priority: number;
  addedAt: string;
}

export interface Wishlist {
  id: string;
  userId: string;
  name: string;
  description: string;
  type: WishlistType;
  privacy: WishlistPrivacy;
  items: WishlistItem[];
  itemCount: number;
  shareToken?: string;
  isDefault: boolean;
  createdAt: string;
  updatedAt: string;
}

// ===== Request/Response Types =====

export interface Pagination {
  page: number;
  pageSize: number;
}

export interface SortOptions {
  field: string;
  order: 'asc' | 'desc';
}

// Business Requests
export interface CreateBusinessRequest {
  name: string;
  description: string;
  category: string;
  address: BusinessAddress;
  contact: BusinessContactInfo;
  settings: BusinessSettings;
}

export interface UpdateBusinessRequest {
  name?: string;
  description?: string;
  category?: string;
  address?: BusinessAddress;
  contact?: BusinessContactInfo;
  settings?: BusinessSettings;
  isActive?: boolean;
}

export interface BusinessFilters {
  country?: string;
  category?: string;
  status?: BusinessVerificationStatus;
  search?: string;
}

export interface BusinessResponse {
  businesses: Business[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

// Product Requests
export interface CreateProductRequest {
  businessId: string;
  name: string;
  description: string;
  category: string;
  type: ProductType;
  status: ProductStatus;
  price: string; // Decimal as string
  currency: string;
  sku: string;
  images: string[];
  tags: string[];
  variants: Record<string, any>;
  inventory: ProductInventory;
  shipping: ProductShipping;
  seo: ProductSEO;
}

export interface UpdateProductRequest {
  name?: string;
  description?: string;
  category?: string;
  type?: ProductType;
  status?: ProductStatus;
  price?: string; // Decimal as string
  currency?: string;
  sku?: string;
  images?: string[];
  tags?: string[];
  variants?: Record<string, any>;
  inventory?: ProductInventory;
  shipping?: ProductShipping;
  seo?: ProductSEO;
  isActive?: boolean;
}

export interface ProductFilters {
  businessId?: string;
  category?: string;
  type?: ProductType;
  status?: ProductStatus;
  search?: string;
}

export interface ProductResponse {
  products: Product[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

// Cart Requests
export interface AddToCartRequest {
  productId: string;
  variantId?: string;
  quantity: number;
  isGift?: boolean;
  giftMessage?: string;
}

export interface UpdateCartItemRequest {
  quantity?: number;
  isGift?: boolean;
  giftMessage?: string;
}

export interface ApplyCouponRequest {
  couponCode: string;
}

export interface MergeCartRequest {
  userId: string;
  sessionId: string;
}

export interface CreateAbandonmentTrackingRequest {
  cartId: string;
  abandonmentStage: string;
  lastPageVisited: string;
}

export interface CartResponse {
  carts: Cart[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

export interface AbandonmentTrackingResponse {
  tracking: CartAbandonmentTracking[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

// Category Requests
export interface CreateCategoryRequest {
  businessId?: string;
  parentId?: string;
  name: string;
  description: string;
  shortDescription: string;
  type: CategoryType;
  icon: string;
  image?: CategoryImage;
  color: string;
  sortOrder: number;
  isVisible: boolean;
  isFeatured: boolean;
  allowProducts: boolean;
  seo: CategorySEO;
  attributes: CategoryAttribute[];
}

export interface UpdateCategoryRequest {
  name?: string;
  description?: string;
  shortDescription?: string;
  status?: CategoryStatus;
  type?: CategoryType;
  icon?: string;
  image?: CategoryImage;
  color?: string;
  sortOrder?: number;
  isVisible?: boolean;
  isFeatured?: boolean;
  allowProducts?: boolean;
  seo?: CategorySEO;
  attributes?: CategoryAttribute[];
}

export interface CategoryResponse {
  categories: Category[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

export interface AddProductToCategoryRequest {
  productId: string;
  categoryId: string;
  isPrimary: boolean;
  sortOrder: number;
}

export interface TrackCategoryViewRequest {
  categoryId: string;
  userId?: string;
  sessionId: string;
  ipAddress: string;
  userAgent: string;
  referrer: string;
}

// Review Requests
export interface CreateReviewRequest {
  type: ReviewType;
  productId?: string;
  businessId?: string;
  orderId?: string;
  userId: string;
  userName: string;
  userEmail: string;
  rating: string; // Decimal as string
  title: string;
  content: string;
  images?: ReviewImage[];
}

export interface UpdateReviewRequest {
  rating?: string; // Decimal as string
  title?: string;
  content?: string;
  images?: ReviewImage[];
}

export interface ReviewResponse {
  reviews: Review[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

// Wishlist Requests
export interface CreateWishlistRequest {
  name: string;
  description: string;
  type: WishlistType;
  privacy: WishlistPrivacy;
}

export interface UpdateWishlistRequest {
  name?: string;
  description?: string;
  privacy?: WishlistPrivacy;
}

export interface AddToWishlistRequest {
  productId: string;
  variantId?: string;
  quantity: number;
  note: string;
  priority: number;
}

export interface WishlistResponse {
  wishlists: Wishlist[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

// ===== API Hooks Types =====

export interface CommerceApiState {
  businesses: Business[];
  products: Product[];
  carts: Cart[];
  categories: Category[];
  reviews: Review[];
  wishlists: Wishlist[];
  isLoading: boolean;
  error: string | null;
}