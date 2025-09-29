// Media Store TypeScript interfaces
// Generated for Media Store Tabs feature implementation

export interface MediaCategory {
  id: string;
  name: string;
  displayOrder: number;
  iconName: string;
  isActive: boolean;
  featuredContentEnabled: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface MediaSubtab {
  id: string;
  categoryId: string;
  name: string;
  displayOrder: number;
  filterCriteria: Record<string, any>;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface MediaContentItem {
  id: string;
  categoryId: string;
  title: string;
  description: string;
  thumbnailUrl: string;
  contentUrl?: string;
  contentType: 'book' | 'podcast' | 'video' | 'cartoon';
  duration?: number;
  price: number;
  currency: string;
  availabilityStatus: 'available' | 'coming_soon' | 'unavailable';
  isFeatured: boolean;
  featuredOrder?: number;
  metadata: Record<string, any>;
  createdAt: string;
  updatedAt: string;
}

export interface MediaProduct {
  id: string;
  name: string;
  description: string;
  price: number;
  currency: string;
  productType: 'physical' | 'media';
  mediaContentId?: string;
  mediaMetadata?: {
    contentType: 'book' | 'podcast' | 'video' | 'cartoon';
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

export interface MediaCartItem {
  id: string;
  cartId: string;
  productId: string;
  mediaContentId?: string;
  quantity: number;
  unitPrice: number;
  totalPrice: number;
  mediaLicense?: 'personal' | 'commercial' | 'educational';
  downloadFormat?: 'PDF' | 'EPUB' | 'MP3' | 'MP4' | 'FLAC';
  createdAt: string;
  updatedAt: string;
}

export interface MediaOrder {
  id: string;
  userId: string;
  status: 'pending' | 'processing' | 'completed' | 'cancelled';
  totalPhysicalAmount: number;
  totalMediaAmount: number;
  totalAmount: number;
  currency: string;
  mediaDeliveryStatus: 'pending' | 'delivered' | 'failed';
  shippingAddress?: string;
  items: MediaOrderItem[];
  createdAt: string;
  updatedAt: string;
}

export interface MediaOrderItem {
  id: string;
  orderId: string;
  productId: string;
  mediaContentId?: string;
  quantity: number;
  unitPrice: number;
  totalPrice: number;
  mediaLicense?: 'personal' | 'commercial' | 'educational';
  downloadFormat?: 'PDF' | 'EPUB' | 'MP3' | 'MP4' | 'FLAC';
  deliveryStatus?: 'pending' | 'delivered' | 'failed';
  downloadUrl?: string;
  createdAt: string;
  updatedAt: string;
}

// API Response types
export interface MediaCategoriesResponse {
  categories: MediaCategory[];
  total: number;
}

export interface MediaContentResponse {
  items: MediaContentItem[];
  page: number;
  limit: number;
  total: number;
  hasMore: boolean;
}

export interface MediaFeaturedResponse {
  items: MediaContentItem[];
  total: number;
  hasMore: boolean;
}

export interface MediaSubtabsResponse {
  subtabs: MediaSubtab[];
  defaultSubtab: string;
}

export interface MediaSearchResponse {
  items: MediaContentItem[];
  query: string;
  total: number;
  page: number;
}

// Store integration types
export interface AddMediaToCartRequest {
  mediaContentId: string;
  quantity: number;
  mediaLicense: 'personal' | 'commercial' | 'educational';
  downloadFormat: 'PDF' | 'EPUB' | 'MP3' | 'MP4' | 'FLAC';
}

export interface AddMediaToCartResponse {
  cartId: string;
  itemsCount: number;
  totalAmount: number;
  currency: string;
  addedItem: MediaCartItem;
}

export interface UnifiedCartResponse {
  cartId: string;
  physicalItems: MediaCartItem[];
  mediaItems: MediaCartItem[];
  totalPhysicalAmount: number;
  totalMediaAmount: number;
  totalAmount: number;
  currency: string;
  itemsCount: number;
}

export interface MediaCheckoutValidationRequest {
  cartId: string;
  mediaItems: MediaCartItem[];
}

export interface MediaCheckoutValidationResponse {
  isValid: boolean;
  validItems: MediaCartItem[];
  invalidItems: MediaCartItem[];
  totalMediaAmount: number;
  estimatedDeliveryTime: string;
}

export interface MediaOrdersResponse {
  orders: MediaOrder[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    hasMore: boolean;
  };
}

// Error types
export interface MediaApiError {
  error: string;
  message: string;
  details?: Record<string, any>;
}