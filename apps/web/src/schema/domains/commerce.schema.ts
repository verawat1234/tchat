/**
 * E-Commerce & Shopping Domain Schema
 *
 * Handles products, shops, orders, carts, payments, and marketplace features
 * Optimized for Southeast Asian markets with multi-currency and local payment methods
 */

import { UUID, Timestamp, Currency, CountryCode, Locale } from '../schema';

// =============================================================================
// PRODUCT MANAGEMENT
// =============================================================================

export interface Product {
  id: UUID;
  shopId: UUID;
  title: string;
  description: string;
  shortDescription?: string;
  images: ProductImage[];
  videos?: ProductVideo[];
  price: number;
  compareAtPrice?: number;
  currency: Currency;
  cost?: number;
  sku?: string;
  barcode?: string;
  inventory: ProductInventory;
  variants: ProductVariant[];
  category: ProductCategory;
  tags: string[];
  attributes: ProductAttribute[];
  seo: ProductSEO;
  status: ProductStatus;
  isDigital: boolean;
  weight?: number;
  dimensions?: ProductDimensions;
  shipping: ProductShipping;
  ratings: ProductRating;
  reviews: ProductReview[];
  localization: ProductLocalization;
  compliance: ProductCompliance;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export type ProductStatus = 'draft' | 'active' | 'out_of_stock' | 'discontinued' | 'archived';

export interface ProductImage {
  id: UUID;
  url: string;
  thumbnailUrl?: string;
  alt?: string;
  order: number;
  isPrimary: boolean;
  variants?: UUID[]; // Which variants this image applies to
}

export interface ProductVideo {
  id: UUID;
  url: string;
  thumbnailUrl?: string;
  duration: number;
  title?: string;
  order: number;
  type: 'demo' | 'review' | 'unboxing' | 'tutorial';
}

export interface ProductInventory {
  trackQuantity: boolean;
  quantity?: number;
  lowStockThreshold?: number;
  isInStock: boolean;
  allowBackorders: boolean;
  location?: string;
  reservedQuantity?: number;
  damagedQuantity?: number;
}

export interface ProductVariant {
  id: UUID;
  title: string;
  sku?: string;
  price?: number;
  compareAtPrice?: number;
  cost?: number;
  inventory: ProductInventory;
  options: VariantOption[];
  image?: string;
  weight?: number;
  dimensions?: ProductDimensions;
  isDefault: boolean;
}

export interface VariantOption {
  name: string; // e.g., "Color", "Size", "Material"
  value: string; // e.g., "Red", "Large", "Cotton"
  displayName?: string; // Localized display name
  colorCode?: string; // For color variants
  imageUrl?: string; // For image-based variants
}

export interface ProductDimensions {
  length: number;
  width: number;
  height: number;
  unit: 'cm' | 'in' | 'mm';
}

export interface ProductShipping {
  isShippingRequired: boolean;
  weight?: number;
  dimensions?: ProductDimensions;
  shippingClass?: string;
  processingTime?: number; // days
  shippingMethods: ShippingMethod[];
  restrictions?: ShippingRestriction[];
}

export interface ShippingMethod {
  id: UUID;
  name: string;
  description?: string;
  price: number;
  currency: Currency;
  estimatedDays: { min: number; max: number };
  regions: CountryCode[];
  carrier?: string;
  trackingEnabled: boolean;
}

export interface ShippingRestriction {
  countries: CountryCode[];
  reason: string;
  type: 'prohibited' | 'restricted' | 'requires_permit';
}

export interface ProductCategory {
  id: UUID;
  name: string;
  slug: string;
  parentId?: UUID;
  description?: string;
  image?: string;
  icon?: string;
  isActive: boolean;
  sortOrder: number;
  seoTitle?: string;
  seoDescription?: string;
  localization: CategoryLocalization;
}

export interface CategoryLocalization {
  [locale: string]: {
    name: string;
    description?: string;
    seoTitle?: string;
    seoDescription?: string;
  };
}

export interface ProductAttribute {
  name: string;
  value: string;
  displayName?: string;
  isVariant: boolean;
  isFilterable: boolean;
  isRequired: boolean;
  type: 'text' | 'number' | 'boolean' | 'date' | 'color' | 'image';
  unit?: string;
}

export interface ProductSEO {
  title?: string;
  description?: string;
  keywords: string[];
  slug: string;
  metaTags?: Record<string, string>;
  canonicalUrl?: string;
  openGraphImage?: string;
}

export interface ProductRating {
  averageRating: number;
  totalReviews: number;
  distribution: { [key: number]: number }; // rating -> count
  verifiedPurchaseRating?: number;
  recentRating?: number; // Last 30 days
}

export interface ProductReview {
  id: UUID;
  productId: UUID;
  userId: UUID;
  variantId?: UUID;
  orderId?: UUID;
  rating: number;
  title?: string;
  comment: string;
  pros?: string[];
  cons?: string[];
  images?: ReviewImage[];
  videos?: ReviewVideo[];
  isVerifiedPurchase: boolean;
  helpfulCount: number;
  reportCount: number;
  moderationStatus: 'pending' | 'approved' | 'rejected' | 'hidden';
  language: Locale;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export interface ReviewImage {
  id: UUID;
  url: string;
  thumbnailUrl?: string;
  caption?: string;
  order: number;
}

export interface ReviewVideo {
  id: UUID;
  url: string;
  thumbnailUrl?: string;
  duration: number;
  caption?: string;
  order: number;
}

export interface ProductLocalization {
  [locale: string]: {
    title: string;
    description: string;
    shortDescription?: string;
    tags: string[];
    attributes: { [key: string]: string };
  };
}

export interface ProductCompliance {
  certifications: ProductCertification[];
  warnings: string[];
  ageRestriction?: number;
  countryRestrictions: CountryCode[];
  requires_id_verification?: boolean;
}

export interface ProductCertification {
  type: string; // e.g., "CE", "FDA", "HALAL", "Organic"
  number?: string;
  issuer: string;
  issuedAt: Timestamp;
  expiresAt?: Timestamp;
  documentUrl?: string;
}

// =============================================================================
// SHOP MANAGEMENT
// =============================================================================

export interface Shop {
  id: UUID;
  ownerId: UUID;
  name: string;
  description?: string;
  avatar?: string;
  coverImage?: string;
  isVerified: boolean;
  verificationLevel: 'none' | 'basic' | 'premium' | 'enterprise';
  status: ShopStatus;
  settings: ShopSettings;
  contact: ShopContact;
  location?: ShopLocation;
  stats: ShopStats;
  policies: ShopPolicies;
  categories: string[];
  tags: string[];
  subscription: ShopSubscription;
  compliance: ShopCompliance;
  localization: ShopLocalization;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export type ShopStatus = 'active' | 'suspended' | 'under_review' | 'closed' | 'maintenance';

export interface ShopSettings {
  isPublic: boolean;
  allowReviews: boolean;
  autoApproveOrders: boolean;
  currency: Currency;
  timezone: string;
  businessHours: ShopBusinessHours[];
  minimumOrder?: number;
  freeShippingThreshold?: number;
  taxSettings: TaxSettings;
  returnSettings: ReturnSettings;
}

export interface ShopBusinessHours {
  dayOfWeek: number; // 0-6, Sunday = 0
  openTime: string; // HH:mm
  closeTime: string; // HH:mm
  isClosed: boolean;
  isDeliveryAvailable: boolean;
  specialNotes?: string;
}

export interface TaxSettings {
  includeTax: boolean;
  taxRate: number;
  taxLabel: string;
  exemptProducts: UUID[];
  countrySpecific: { [country: string]: number };
}

export interface ReturnSettings {
  allowReturns: boolean;
  returnPeriod: number; // days
  returnShipping: 'customer_pays' | 'shop_pays' | 'free';
  conditions: string[];
  restockingFee?: number;
}

export interface ShopContact {
  email?: string;
  phone?: string;
  whatsapp?: string;
  website?: string;
  socialMedia: { [platform: string]: string };
  supportHours?: ShopBusinessHours[];
  responseTime?: number; // hours
}

export interface ShopLocation {
  address: string;
  city: string;
  state: string;
  country: CountryCode;
  postalCode: string;
  coordinates?: { latitude: number; longitude: number };
  isPhysicalStore: boolean;
  storeHours?: ShopBusinessHours[];
}

export interface ShopStats {
  totalProducts: number;
  activeProducts: number;
  totalOrders: number;
  totalRevenue: number;
  averageOrderValue: number;
  averageRating: number;
  totalReviews: number;
  responseTime: number; // hours
  responseRate: number; // percentage
  returnRate: number; // percentage
  customerRetentionRate: number; // percentage
  monthlyStats: MonthlyShopStats[];
}

export interface MonthlyShopStats {
  month: string; // YYYY-MM
  orders: number;
  revenue: number;
  newCustomers: number;
  returningCustomers: number;
  averageOrderValue: number;
  topProducts: { productId: UUID; sales: number }[];
}

export interface ShopPolicies {
  returnPolicy?: string;
  shippingPolicy?: string;
  privacyPolicy?: string;
  termsOfService?: string;
  refundPolicy?: string;
  warrantyPolicy?: string;
  lastUpdated: Timestamp;
}

export interface ShopSubscription {
  plan: 'free' | 'basic' | 'professional' | 'enterprise';
  status: 'active' | 'cancelled' | 'expired' | 'suspended';
  billingCycle: 'monthly' | 'yearly';
  price: number;
  currency: Currency;
  features: ShopFeature[];
  limits: ShopLimits;
  startDate: Timestamp;
  endDate?: Timestamp;
  autoRenew: boolean;
}

export interface ShopFeature {
  id: string;
  name: string;
  isEnabled: boolean;
  limit?: number;
  usage?: number;
}

export interface ShopLimits {
  maxProducts: number;
  maxCategories: number;
  maxImages: number;
  maxBandwidth: number; // GB
  maxOrders: number; // per month
  customDomain: boolean;
  advancedAnalytics: boolean;
  prioritySupport: boolean;
}

export interface ShopCompliance {
  businessLicense?: BusinessLicense;
  taxRegistration?: TaxRegistration;
  certifications: ShopCertification[];
  insurancePolicies: InsurancePolicy[];
  complianceChecks: ComplianceCheck[];
}

export interface BusinessLicense {
  number: string;
  issuer: string;
  issuedAt: Timestamp;
  expiresAt: Timestamp;
  documentUrl?: string;
  status: 'valid' | 'expired' | 'suspended';
}

export interface TaxRegistration {
  number: string;
  country: CountryCode;
  registeredAt: Timestamp;
  isActive: boolean;
}

export interface ShopCertification {
  type: string;
  issuer: string;
  validFrom: Timestamp;
  validTo: Timestamp;
  documentUrl?: string;
}

export interface InsurancePolicy {
  type: 'liability' | 'product' | 'cyber' | 'general';
  provider: string;
  policyNumber: string;
  coverage: number;
  currency: Currency;
  validFrom: Timestamp;
  validTo: Timestamp;
}

export interface ComplianceCheck {
  id: UUID;
  type: string;
  status: 'pending' | 'passed' | 'failed' | 'review_required';
  checkedAt: Timestamp;
  nextCheckAt?: Timestamp;
  notes?: string;
  documentUrls: string[];
}

export interface ShopLocalization {
  [locale: string]: {
    name: string;
    description?: string;
    policies?: Partial<ShopPolicies>;
    categories: string[];
  };
}

// =============================================================================
// CART & ORDERING
// =============================================================================

export interface Cart {
  id: UUID;
  userId: UUID;
  sessionId?: string; // For guest users
  items: CartItem[];
  subtotal: number;
  discount: number;
  tax: number;
  shipping: number;
  total: number;
  currency: Currency;
  couponCode?: string;
  shippingAddress?: Address;
  billingAddress?: Address;
  paymentMethod?: PaymentMethodInfo;
  notes?: string;
  metadata?: CartMetadata;
  expiresAt?: Timestamp;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export interface CartItem {
  id: UUID;
  cartId: UUID;
  productId: UUID;
  variantId?: UUID;
  quantity: number;
  price: number;
  compareAtPrice?: number;
  title: string;
  image?: string;
  variant?: string;
  shopId: UUID;
  shopName: string;
  isAvailable: boolean;
  isDigital: boolean;
  shippingRequired: boolean;
  weight?: number;
  dimensions?: ProductDimensions;
  customizations?: ProductCustomization[];
  addedAt: Timestamp;
  updatedAt: Timestamp;
}

export interface ProductCustomization {
  type: 'text' | 'image' | 'color' | 'option';
  label: string;
  value: string;
  price?: number; // Additional cost
  displayOrder: number;
}

export interface CartMetadata {
  promoCode?: string;
  referralCode?: string;
  affiliateId?: UUID;
  utmSource?: string;
  utmMedium?: string;
  utmCampaign?: string;
  deviceInfo?: string;
  ipAddress?: string;
  geoLocation?: string;
}

export interface Order {
  id: UUID;
  orderNumber: string;
  userId: UUID;
  customerInfo: CustomerInfo;
  status: OrderStatus;
  items: OrderItem[];
  subtotal: number;
  discount: number;
  tax: number;
  shipping: number;
  total: number;
  currency: Currency;
  couponCode?: string;
  shippingAddress: Address;
  billingAddress?: Address;
  payment: OrderPayment;
  fulfillment: OrderFulfillment;
  communications: OrderCommunication[];
  timeline: OrderTimeline[];
  refunds: OrderRefund[];
  notes?: string;
  metadata?: OrderMetadata;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export type OrderStatus =
  | 'pending'
  | 'confirmed'
  | 'processing'
  | 'shipped'
  | 'delivered'
  | 'cancelled'
  | 'refunded'
  | 'disputed';

export interface CustomerInfo {
  name: string;
  email?: string;
  phone?: string;
  preferredLanguage: Locale;
  isGuest: boolean;
  customerNotes?: string;
}

export interface OrderItem {
  id: UUID;
  orderId: UUID;
  productId: UUID;
  variantId?: UUID;
  quantity: number;
  price: number;
  total: number;
  title: string;
  image?: string;
  variant?: string;
  shopId: UUID;
  shopName: string;
  sku?: string;
  weight?: number;
  customizations?: ProductCustomization[];
  fulfillmentStatus: 'pending' | 'processing' | 'shipped' | 'delivered' | 'cancelled' | 'returned';
  trackingNumber?: string;
  returnStatus?: 'none' | 'requested' | 'approved' | 'returned' | 'refunded';
}

export interface OrderPayment {
  method: PaymentMethodInfo;
  status: PaymentStatus;
  transactionId?: string;
  amount: number;
  currency: Currency;
  fees: PaymentFee[];
  processedAt?: Timestamp;
  failureReason?: string;
  authorization?: PaymentAuthorization;
  refunds: PaymentRefund[];
}

export type PaymentStatus = 'pending' | 'authorized' | 'captured' | 'failed' | 'refunded' | 'disputed';

export interface PaymentMethodInfo {
  type: 'wallet' | 'bank_transfer' | 'credit_card' | 'qr_payment' | 'cod' | 'installment';
  details: PaymentMethodDetails;
  isStored: boolean;
  metadata?: Record<string, any>;
}

export interface PaymentMethodDetails {
  // Credit Card
  cardLast4?: string;
  cardBrand?: string;
  cardExpiry?: string;

  // Bank Transfer
  bankName?: string;
  accountLast4?: string;

  // QR Payment
  qrProvider?: string;
  qrReference?: string;

  // COD
  codFee?: number;

  // Installment
  installmentPlan?: InstallmentPlan;
}

export interface InstallmentPlan {
  provider: string;
  months: number;
  monthlyAmount: number;
  totalAmount: number;
  interestRate: number;
  firstPayment: Timestamp;
}

export interface PaymentFee {
  type: 'processing' | 'gateway' | 'currency_conversion' | 'installment';
  amount: number;
  currency: Currency;
  description: string;
}

export interface PaymentAuthorization {
  authorizationId: string;
  amount: number;
  currency: Currency;
  expiresAt: Timestamp;
  capturedAmount?: number;
  capturedAt?: Timestamp;
}

export interface PaymentRefund {
  id: UUID;
  amount: number;
  currency: Currency;
  reason: string;
  status: 'pending' | 'completed' | 'failed';
  refundId?: string;
  processedAt?: Timestamp;
  failureReason?: string;
}

export interface OrderFulfillment {
  status: 'pending' | 'processing' | 'partially_shipped' | 'shipped' | 'delivered' | 'cancelled';
  trackingNumbers: TrackingInfo[];
  estimatedDelivery?: Timestamp;
  actualDelivery?: Timestamp;
  shippingCarrier?: string;
  shippingMethod?: string;
  shippingCost: number;
  packaging?: PackagingInfo;
  deliveryInstructions?: string;
  deliveryAttempts: DeliveryAttempt[];
}

export interface TrackingInfo {
  trackingNumber: string;
  carrier: string;
  url?: string;
  status: 'created' | 'picked_up' | 'in_transit' | 'out_for_delivery' | 'delivered' | 'failed';
  lastUpdate: Timestamp;
  estimatedDelivery?: Timestamp;
  items: UUID[]; // OrderItem IDs
}

export interface PackagingInfo {
  type: 'envelope' | 'box' | 'tube' | 'custom';
  dimensions?: ProductDimensions;
  weight: number;
  materials: string[];
  isEcoFriendly: boolean;
}

export interface DeliveryAttempt {
  attemptNumber: number;
  attemptedAt: Timestamp;
  status: 'failed' | 'delivered' | 'rescheduled';
  reason?: string;
  nextAttempt?: Timestamp;
  signature?: string;
  photo?: string;
}

export interface OrderCommunication {
  id: UUID;
  orderId: UUID;
  type: 'email' | 'sms' | 'push' | 'in_app';
  direction: 'outbound' | 'inbound';
  subject?: string;
  content: string;
  sentAt: Timestamp;
  deliveredAt?: Timestamp;
  readAt?: Timestamp;
  responseRequired: boolean;
  templateId?: string;
}

export interface OrderTimeline {
  id: UUID;
  orderId: UUID;
  status: string;
  description: string;
  timestamp: Timestamp;
  userId?: UUID;
  metadata?: Record<string, any>;
  isPublic: boolean;
  notification?: {
    sent: boolean;
    channels: string[];
    sentAt?: Timestamp;
  };
}

export interface OrderRefund {
  id: UUID;
  orderId: UUID;
  amount: number;
  currency: Currency;
  reason: RefundReason;
  status: 'requested' | 'approved' | 'processing' | 'completed' | 'rejected';
  items: RefundItem[];
  refundMethod: 'original_payment' | 'store_credit' | 'bank_transfer';
  processedAt?: Timestamp;
  notes?: string;
  attachments: string[];
}

export type RefundReason =
  | 'defective_product'
  | 'wrong_item'
  | 'not_as_described'
  | 'arrived_late'
  | 'customer_changed_mind'
  | 'duplicate_order'
  | 'fraud'
  | 'other';

export interface RefundItem {
  orderItemId: UUID;
  quantity: number;
  amount: number;
  reason?: string;
  condition?: 'new' | 'opened' | 'used' | 'damaged';
  photos?: string[];
}

export interface OrderMetadata {
  source: 'web' | 'mobile' | 'api' | 'admin';
  referrer?: string;
  utmParams?: Record<string, string>;
  deviceInfo?: string;
  ipAddress?: string;
  fraudScore?: number;
  riskLevel?: 'low' | 'medium' | 'high';
  affiliateId?: UUID;
  promotionIds: UUID[];
}

export interface Address {
  id?: UUID;
  firstName: string;
  lastName: string;
  company?: string;
  address1: string;
  address2?: string;
  city: string;
  province: string;
  country: CountryCode;
  postalCode: string;
  phone?: string;
  email?: string;
  isDefault: boolean;
  type: 'shipping' | 'billing' | 'both';
  coordinates?: { latitude: number; longitude: number };
  deliveryInstructions?: string;
  accessCodes?: string;
  validatedAt?: Timestamp;
  validationStatus?: 'valid' | 'invalid' | 'unverified';
}

// =============================================================================
// BUSINESS LOGIC CONSTANTS
// =============================================================================

/**
 * Product status transitions
 */
export const PRODUCT_STATUS_TRANSITIONS: Record<ProductStatus, ProductStatus[]> = {
  draft: ['active', 'archived'],
  active: ['out_of_stock', 'discontinued', 'archived'],
  out_of_stock: ['active', 'discontinued', 'archived'],
  discontinued: ['active', 'archived'],
  archived: ['draft']
};

/**
 * Order status transitions
 */
export const ORDER_STATUS_TRANSITIONS: Record<OrderStatus, OrderStatus[]> = {
  pending: ['confirmed', 'cancelled'],
  confirmed: ['processing', 'cancelled'],
  processing: ['shipped', 'cancelled'],
  shipped: ['delivered', 'cancelled'],
  delivered: ['refunded', 'disputed'],
  cancelled: [],
  refunded: ['disputed'],
  disputed: []
};

/**
 * Currency configurations for SEA markets
 */
export const CURRENCY_CONFIG: Record<Currency, {
  symbol: string;
  decimals: number;
  placement: 'before' | 'after';
  countries: CountryCode[];
}> = {
  THB: { symbol: '฿', decimals: 2, placement: 'before', countries: ['TH'] },
  IDR: { symbol: 'Rp', decimals: 0, placement: 'before', countries: ['ID'] },
  MYR: { symbol: 'RM', decimals: 2, placement: 'before', countries: ['MY'] },
  SGD: { symbol: 'S$', decimals: 2, placement: 'before', countries: ['SG'] },
  PHP: { symbol: '₱', decimals: 2, placement: 'before', countries: ['PH'] },
  VND: { symbol: '₫', decimals: 0, placement: 'after', countries: ['VN'] },
  USD: { symbol: '$', decimals: 2, placement: 'before', countries: [] }
};

/**
 * Default shipping methods by country
 */
export const DEFAULT_SHIPPING_METHODS: Record<CountryCode, string[]> = {
  TH: ['Thailand Post', 'Kerry Express', 'J&T Express', 'Flash Express'],
  ID: ['JNE', 'TIKI', 'Pos Indonesia', 'J&T Express'],
  MY: ['Pos Malaysia', 'City-Link', 'Ninja Van', 'J&T Express'],
  SG: ['SingPost', 'Ninja Van', 'Qxpress', 'J&T Express'],
  PH: ['LBC', 'J&T Express', '2GO Express', 'Ninja Van'],
  VN: ['Vietnam Post', 'Giao Hang Nhanh', 'Viettel Post', 'J&T Express']
};