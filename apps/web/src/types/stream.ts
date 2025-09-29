// Stream Store TypeScript interfaces
// Generated for Stream Store Tabs feature implementation

export interface StreamCategory {
  id: string;
  name: string;
  displayOrder: number;
  iconName: string;
  isActive: boolean;
  featuredContentEnabled: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface StreamSubtab {
  id: string;
  categoryId: string;
  name: string;
  displayOrder: number;
  filterCriteria: Record<string, any>;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export type StreamContentType =
  | 'book'
  | 'podcast'
  | 'cartoon'
  | 'short_movie'
  | 'long_movie'
  | 'music'
  | 'art';

export type StreamAvailabilityStatus =
  | 'available'
  | 'coming_soon'
  | 'unavailable';

export interface StreamContentItem {
  id: string;
  categoryId: string;
  title: string;
  description: string;
  thumbnailUrl: string;
  contentType: StreamContentType;
  duration?: number; // in seconds, null for books
  price: number;
  currency: string;
  availabilityStatus: StreamAvailabilityStatus;
  isFeatured: boolean;
  featuredOrder?: number;
  metadata: Record<string, any>;
  createdAt: string;
  updatedAt: string;
}

export interface StreamProduct {
  id: string;
  name: string;
  description: string;
  price: number;
  currency: string;
  productType: 'physical' | 'media';
  mediaContentId?: string;
  mediaMetadata?: {
    contentType: StreamContentType;
    duration?: number;
    format?: string;
    license?: string;
  };
  category: string;
  isActive: boolean;
  stockQuantity?: number;
  createdAt: string;
  updatedAt: string;
}

export interface StreamCartItem {
  id: string;
  cartId: string;
  productId: string;
  mediaContentId?: string;
  quantity: number;
  unitPrice: number;
  totalPrice: number;
  mediaLicense?: 'personal' | 'family';
  downloadFormat?: 'PDF' | 'EPUB' | 'MP3' | 'MP4' | 'FLAC';
  createdAt: string;
  updatedAt: string;
}

export interface StreamOrder {
  id: string;
  userId: string;
  status: 'pending' | 'processing' | 'completed' | 'cancelled';
  totalPhysicalAmount: number;
  totalStreamAmount: number;
  totalAmount: number;
  currency: string;
  containsMediaItems: boolean;
  mediaDeliveryStatus: 'pending' | 'delivered' | 'failed';
  shippingAddress?: string;
  items: StreamOrderItem[];
  createdAt: string;
  updatedAt: string;
}

export interface StreamOrderItem {
  id: string;
  orderId: string;
  productId: string;
  mediaContentId?: string;
  quantity: number;
  unitPrice: number;
  totalPrice: number;
  mediaLicense?: 'personal' | 'family';
  downloadFormat?: 'PDF' | 'EPUB' | 'MP3' | 'MP4' | 'FLAC';
  deliveryStatus?: 'pending' | 'delivered' | 'failed';
  downloadUrl?: string;
  createdAt: string;
  updatedAt: string;
}

export interface ContentCollection {
  id: string;
  name: string;
  categoryId: string;
  collectionType: 'featured' | 'new_releases' | 'trending' | 'curated';
  displayOrder: number;
  isActive: boolean;
  itemIds: string[];
  maxItems: number;
  createdAt: string;
  updatedAt: string;
}

// Navigation state management
export interface TabNavigationState {
  userId: string;
  currentCategoryId: string;
  currentSubtabId?: string;
  lastVisitedAt: string;
  sessionId: string;
}

// API Response types
export interface StreamCategoriesResponse {
  categories: StreamCategory[];
  total: number;
}

export interface StreamContentResponse {
  items: StreamContentItem[];
  page: number;
  limit: number;
  total: number;
  hasMore: boolean;
}

export interface StreamFeaturedResponse {
  items: StreamContentItem[];
  total: number;
  hasMore: boolean;
}

export interface StreamSubtabsResponse {
  subtabs: StreamSubtab[];
  defaultSubtab: string;
}

export interface StreamSearchResponse {
  items: StreamContentItem[];
  query: string;
  total: number;
  page: number;
}

// Store integration types
export interface AddStreamToCartRequest {
  mediaContentId: string;
  quantity: number;
  mediaLicense: 'personal' | 'family';
  downloadFormat: 'PDF' | 'EPUB' | 'MP3' | 'MP4' | 'FLAC';
}

export interface AddStreamToCartResponse {
  cartId: string;
  itemsCount: number;
  totalAmount: number;
  currency: string;
  addedItem: StreamCartItem;
}

export interface UnifiedCartResponse {
  cartId: string;
  physicalItems: StreamCartItem[];
  streamItems: StreamCartItem[];
  totalPhysicalAmount: number;
  totalStreamAmount: number;
  totalAmount: number;
  currency: string;
  itemsCount: number;
}

export interface StreamCheckoutValidationRequest {
  cartId: string;
  streamItems: StreamCartItem[];
}

export interface StreamCheckoutValidationResponse {
  isValid: boolean;
  validItems: StreamCartItem[];
  invalidItems: StreamCartItem[];
  totalStreamAmount: number;
  estimatedDeliveryTime: string;
}

export interface StreamOrdersResponse {
  orders: StreamOrder[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    hasMore: boolean;
  };
}

// Error types
export interface StreamApiError {
  error: string;
  message: string;
  details?: Record<string, any>;
}

// Filter and search types
export interface StreamFilters {
  categoryId?: string;
  contentType?: StreamContentType;
  priceMin?: number;
  priceMax?: number;
  isFeatured?: boolean;
  availabilityStatus?: StreamAvailabilityStatus;
  durationMin?: number;
  durationMax?: number;
}

export interface StreamSortOptions {
  field: 'title' | 'price' | 'createdAt' | 'featuredOrder';
  order: 'asc' | 'desc';
}

// Redux state types
export interface StreamState {
  categories: StreamCategory[];
  currentCategoryId?: string;
  currentSubtabId?: string;
  content: Record<string, StreamContentItem[]>;
  featuredContent: Record<string, StreamContentItem[]>;
  loading: {
    categories: boolean;
    content: boolean;
    featured: boolean;
  };
  error: {
    categories?: string;
    content?: string;
    featured?: string;
  };
  cache: {
    lastUpdated: Record<string, string>;
    ttl: number;
  };
}