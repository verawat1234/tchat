/**
 * Microservices API Integration - Unified endpoints for all backend services
 *
 * This file provides RTK Query endpoints that integrate with all 6 microservices
 * running on ports 8080-8086. It handles service discovery, load balancing,
 * and provides a unified API interface for the React application.
 *
 * Services:
 * - Gateway (8080) - API Gateway and routing
 * - Auth (8081) - Authentication and user management
 * - Content (8082) - Content management system
 * - Messaging (8083) - Real-time messaging and chat
 * - Commerce (8084) - E-commerce and store functionality
 * - Payment (8085) - Multi-currency payments for SEA
 * - Notification (8086) - Push notifications and alerts
 */

import { api } from './api';
import { getServiceConfig } from './serviceConfig';
import type {
  User,
  Message,
  Chat,
  ApiResponse,
  PaginatedResponse
} from '../types/api';

// =============================================================================
// Service Configuration
// =============================================================================

// Check if we should use direct services or gateway
const serviceConfig = getServiceConfig();
const useDirectServices = serviceConfig.useDirect;

const SERVICES = {
  GATEWAY: import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1',
  AUTH: useDirectServices ? (import.meta.env.VITE_AUTH_URL || 'http://localhost:8081/api/v1') : (import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'),
  CONTENT: useDirectServices ? (import.meta.env.VITE_CONTENT_URL || 'http://localhost:8082/api/v1') : (import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'),
  MESSAGING: useDirectServices ? (import.meta.env.VITE_MESSAGING_URL || 'http://localhost:8083/api/v1') : (import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'),
  COMMERCE: useDirectServices ? (import.meta.env.VITE_COMMERCE_URL || 'http://localhost:8084/api/v1') : (import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'),
  PAYMENT: useDirectServices ? (import.meta.env.VITE_PAYMENT_URL || 'http://localhost:8085/api/v1') : (import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'),
  NOTIFICATION: useDirectServices ? (import.meta.env.VITE_NOTIFICATION_URL || 'http://localhost:8086/api/v1') : (import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1')
};

// =============================================================================
// Type Definitions for Southeast Asian Commerce
// =============================================================================

interface Product {
  id: string;
  name: string;
  description: string;
  price: number;
  currency: 'THB' | 'IDR' | 'SGD' | 'PHP' | 'MYR' | 'VND';
  images: string[];
  category: string;
  merchant: {
    id: string;
    name: string;
    country: 'TH' | 'ID' | 'SG' | 'PH' | 'MY' | 'VN';
  };
  rating: number;
  reviewCount: number;
  shipping: {
    available: boolean;
    cost: number;
    estimatedDays: number;
  };
  createdAt: string;
  updatedAt: string;
}

interface Order {
  id: string;
  userId: string;
  products: {
    productId: string;
    quantity: number;
    price: number;
    currency: string;
  }[];
  totalAmount: number;
  currency: string;
  status: 'pending' | 'confirmed' | 'shipped' | 'delivered' | 'cancelled';
  payment: {
    method: 'PromptPay' | 'GrabPay' | 'GoPay' | 'ShopeePay' | 'Cash' | 'BankTransfer';
    status: 'pending' | 'paid' | 'failed' | 'refunded';
  };
  shipping: {
    address: string;
    city: string;
    country: string;
    postalCode: string;
    trackingNumber?: string;
  };
  createdAt: string;
  updatedAt: string;
}

interface Wallet {
  id: string;
  userId: string;
  balances: {
    currency: 'THB' | 'IDR' | 'SGD' | 'PHP' | 'MYR' | 'VND';
    amount: number;
  }[];
  paymentMethods: {
    id: string;
    type: 'PromptPay' | 'GrabPay' | 'GoPay' | 'ShopeePay' | 'BankTransfer';
    details: Record<string, any>;
    isDefault: boolean;
  }[];
  transactions: {
    id: string;
    type: 'credit' | 'debit';
    amount: number;
    currency: string;
    description: string;
    timestamp: string;
  }[];
}

interface Notification {
  id: string;
  userId: string;
  type: 'message' | 'order' | 'payment' | 'system' | 'promotion';
  title: string;
  content: string;
  data?: Record<string, any>;
  read: boolean;
  createdAt: string;
}

// =============================================================================
// Messaging Service Endpoints
// =============================================================================

export const messagingApi = api.injectEndpoints({
  endpoints: (builder) => ({
    /**
     * Get all chats for current user
     */
    getUserChats: builder.query<Chat[], void>({
      query: () => ({
        url: `/chats`,
        baseUrl: SERVICES.MESSAGING
      }),
      providesTags: ['Chat'],
      transformResponse: (response: ApiResponse<Chat[]>) => response.data,
    }),

    /**
     * Get chat messages with pagination
     */
    getChatMessages: builder.query<PaginatedResponse<Message>, { chatId: string; page?: number; limit?: number }>({
      query: ({ chatId, page = 1, limit = 20 }) => ({
        url: `/chats/${chatId}/messages?page=${page}&limit=${limit}`,
        baseUrl: SERVICES.MESSAGING
      }),
      providesTags: (result, error, { chatId }) => [
        { type: 'Message', id: `CHAT-${chatId}` }
      ],
    }),

    /**
     * Send new message
     */
    sendChatMessage: builder.mutation<Message, { chatId: string; content: string; type?: string }>({
      query: ({ chatId, content, type = 'text' }) => ({
        url: `/chats/${chatId}/messages`,
        method: 'POST',
        body: { content, type },
        baseUrl: SERVICES.MESSAGING
      }),
      invalidatesTags: (result, error, { chatId }) => [
        { type: 'Message', id: `CHAT-${chatId}` },
        { type: 'Chat', id: chatId }
      ],
      transformResponse: (response: ApiResponse<Message>) => response.data,
    }),

    /**
     * Create new chat
     */
    createChat: builder.mutation<Chat, { participants: string[]; type?: 'direct' | 'group'; name?: string }>({
      query: (data) => ({
        url: '/chats',
        method: 'POST',
        body: data,
        baseUrl: SERVICES.MESSAGING
      }),
      invalidatesTags: ['Chat'],
      transformResponse: (response: ApiResponse<Chat>) => response.data,
    })
  }),
});

// =============================================================================
// Commerce Service Endpoints
// =============================================================================

export const commerceApi = api.injectEndpoints({
  endpoints: (builder) => ({
    /**
     * Get featured shops/merchants
     */
    getFeaturedShops: builder.query<any[], { country?: string; category?: string; limit?: number }>({
      query: ({ country, category, limit = 20 }) => {
        const params = new URLSearchParams();
        if (country) params.append('country', country);
        if (category) params.append('category', category);
        params.append('limit', limit.toString());
        return {
          url: `/shops/featured?${params}`,
          baseUrl: SERVICES.COMMERCE
        };
      },
      providesTags: ['Shop'],
      transformResponse: (response: ApiResponse<any[]>) => response.data,
    }),

    /**
     * Search shops with filters
     */
    searchShops: builder.query<any[], {
      query?: string;
      category?: string;
      country?: string;
      isVerified?: boolean;
      page?: number;
      limit?: number;
    }>({
      query: (filters) => {
        const params = new URLSearchParams();
        Object.entries(filters).forEach(([key, value]) => {
          if (value !== undefined) params.append(key, value.toString());
        });
        return {
          url: `/shops/search?${params}`,
          baseUrl: SERVICES.COMMERCE
        };
      },
      providesTags: ['Shop'],
      transformResponse: (response: ApiResponse<any[]>) => response.data,
    }),

    /**
     * Get featured products for Southeast Asian markets
     */
    getFeaturedProducts: builder.query<Product[], { country?: string; limit?: number }>({
      query: ({ country, limit = 20 }) => {
        const params = new URLSearchParams();
        if (country) params.append('country', country);
        params.append('limit', limit.toString());
        return {
          url: `/products/featured?${params}`,
          baseUrl: SERVICES.COMMERCE
        };
      },
      providesTags: ['Product'],
      transformResponse: (response: ApiResponse<Product[]>) => response.data,
    }),

    /**
     * Search products with filters
     */
    searchProducts: builder.query<PaginatedResponse<Product>, {
      query?: string;
      category?: string;
      country?: string;
      currency?: string;
      minPrice?: number;
      maxPrice?: number;
      page?: number;
      limit?: number;
    }>({
      query: (filters) => {
        const params = new URLSearchParams();
        Object.entries(filters).forEach(([key, value]) => {
          if (value !== undefined) params.append(key, value.toString());
        });
        return {
          url: `/products/search?${params}`,
          baseUrl: SERVICES.COMMERCE
        };
      },
      providesTags: ['Product'],
    }),

    /**
     * Get product details
     */
    getProduct: builder.query<Product, string>({
      query: (productId) => ({
        url: `/products/${productId}`,
        baseUrl: SERVICES.COMMERCE
      }),
      providesTags: (result, error, id) => [{ type: 'Product', id }],
      transformResponse: (response: ApiResponse<Product>) => response.data,
    }),

    /**
     * Create new order
     */
    createOrder: builder.mutation<Order, Partial<Order>>({
      query: (orderData) => ({
        url: '/orders',
        method: 'POST',
        body: orderData,
        baseUrl: SERVICES.COMMERCE
      }),
      invalidatesTags: ['Order'],
      transformResponse: (response: ApiResponse<Order>) => response.data,
    }),

    /**
     * Get user orders
     */
    getUserOrders: builder.query<Order[], { status?: string; page?: number; limit?: number }>({
      query: ({ status, page = 1, limit = 20 }) => {
        const params = new URLSearchParams();
        if (status) params.append('status', status);
        params.append('page', page.toString());
        params.append('limit', limit.toString());
        return {
          url: `/orders?${params}`,
          baseUrl: SERVICES.COMMERCE
        };
      },
      providesTags: ['Order'],
      transformResponse: (response: ApiResponse<Order[]>) => response.data,
    }),

    /**
     * Update order status
     */
    updateOrderStatus: builder.mutation<Order, { orderId: string; status: string }>({
      query: ({ orderId, status }) => ({
        url: `/orders/${orderId}/status`,
        method: 'PATCH',
        body: { status },
        baseUrl: SERVICES.COMMERCE
      }),
      invalidatesTags: (result, error, { orderId }) => [{ type: 'Order', id: orderId }],
      transformResponse: (response: ApiResponse<Order>) => response.data,
    }),

    /**
     * Get user cart items
     */
    getCartItems: builder.query<any[], void>({
      query: () => ({
        url: '/cart',
        baseUrl: SERVICES.COMMERCE
      }),
      providesTags: ['Cart'],
      transformResponse: (response: ApiResponse<any[]>) => response.data,
    }),

    /**
     * Add item to cart
     */
    addToCart: builder.mutation<any, { productId: string; quantity?: number }>({
      query: ({ productId, quantity = 1 }) => ({
        url: '/cart/items',
        method: 'POST',
        body: { productId, quantity },
        baseUrl: SERVICES.COMMERCE
      }),
      invalidatesTags: ['Cart'],
      transformResponse: (response: ApiResponse<any>) => response.data,
    }),

    /**
     * Update cart item quantity
     */
    updateCartItem: builder.mutation<any, { itemId: string; quantity: number }>({
      query: ({ itemId, quantity }) => ({
        url: `/cart/items/${itemId}`,
        method: 'PATCH',
        body: { quantity },
        baseUrl: SERVICES.COMMERCE
      }),
      invalidatesTags: ['Cart'],
      transformResponse: (response: ApiResponse<any>) => response.data,
    }),

    /**
     * Remove item from cart
     */
    removeFromCart: builder.mutation<void, string>({
      query: (itemId) => ({
        url: `/cart/items/${itemId}`,
        method: 'DELETE',
        baseUrl: SERVICES.COMMERCE
      }),
      invalidatesTags: ['Cart'],
    }),

    /**
     * Clear cart
     */
    clearCart: builder.mutation<void, void>({
      query: () => ({
        url: '/cart',
        method: 'DELETE',
        baseUrl: SERVICES.COMMERCE
      }),
      invalidatesTags: ['Cart'],
    }),

    /**
     * Get products by shop ID
     */
    getShopProducts: builder.query<Product[], {
      shopId: string;
      category?: string;
      page?: number;
      limit?: number;
    }>({
      query: ({ shopId, category, page = 1, limit = 20 }) => {
        const params = new URLSearchParams();
        if (category) params.append('category', category);
        params.append('page', page.toString());
        params.append('limit', limit.toString());
        return {
          url: `/shops/${shopId}/products?${params}`,
          baseUrl: SERVICES.COMMERCE
        };
      },
      providesTags: (result, error, { shopId }) => [
        { type: 'Product', id: `shop-${shopId}` },
        { type: 'Product', id: 'LIST' }
      ],
      transformResponse: (response: ApiResponse<Product[]>) => response.data,
    }),

    /**
     * Get categories
     */
    getCategories: builder.query<any[], void>({
      query: () => ({
        url: '/categories',
        baseUrl: SERVICES.COMMERCE
      }),
      providesTags: ['Category'],
      transformResponse: (response: ApiResponse<any[]>) => response.data,
    })
  }),
});

// Export the generated hooks
export const {
  useGetFeaturedShopsQuery,
  useGetFeaturedProductsQuery,
  useSearchProductsQuery,
  useAddToCartMutation,
  useGetShopProductsQuery,
  useGetCategoriesQuery
} = commerceApi;

// =============================================================================
// Payment Service Endpoints (Southeast Asian Payment Methods)
// =============================================================================

export const paymentApi = api.injectEndpoints({
  endpoints: (builder) => ({
    /**
     * Get user wallet with multi-currency support
     */
    getUserWallet: builder.query<Wallet, void>({
      query: () => ({
        url: '/wallet',
        baseUrl: SERVICES.PAYMENT
      }),
      providesTags: ['Wallet'],
      transformResponse: (response: ApiResponse<Wallet>) => response.data,
    }),

    /**
     * Add money to wallet
     */
    addMoneyToWallet: builder.mutation<Wallet, {
      amount: number;
      currency: string;
      paymentMethod: string;
    }>({
      query: (data) => ({
        url: '/wallet/topup',
        method: 'POST',
        body: data,
        baseUrl: SERVICES.PAYMENT
      }),
      invalidatesTags: ['Wallet'],
      transformResponse: (response: ApiResponse<Wallet>) => response.data,
    }),

    /**
     * Send money to another user (P2P payments)
     */
    sendMoney: builder.mutation<any, {
      recipientId: string;
      amount: number;
      currency: string;
      message?: string;
    }>({
      query: (data) => ({
        url: '/transfers',
        method: 'POST',
        body: data,
        baseUrl: SERVICES.PAYMENT
      }),
      invalidatesTags: ['Wallet'],
    }),

    /**
     * Process payment for order
     */
    processPayment: builder.mutation<any, {
      orderId: string;
      amount: number;
      currency: string;
      paymentMethod: string;
    }>({
      query: (data) => ({
        url: '/payments/process',
        method: 'POST',
        body: data,
        baseUrl: SERVICES.PAYMENT
      }),
      invalidatesTags: ['Wallet', 'Order'],
    }),

    /**
     * Add payment method
     */
    addPaymentMethod: builder.mutation<any, {
      type: string;
      details: Record<string, any>;
    }>({
      query: (data) => ({
        url: '/payment-methods',
        method: 'POST',
        body: data,
        baseUrl: SERVICES.PAYMENT
      }),
      invalidatesTags: ['Wallet'],
    }),

    /**
     * Get transaction history
     */
    getTransactionHistory: builder.query<any[], { page?: number; limit?: number; type?: string }>({
      query: ({ page = 1, limit = 20, type }) => {
        const params = new URLSearchParams();
        params.append('page', page.toString());
        params.append('limit', limit.toString());
        if (type) params.append('type', type);
        return {
          url: `/transactions?${params}`,
          baseUrl: SERVICES.PAYMENT
        };
      },
      providesTags: ['Transaction'],
    })
  }),
});

// =============================================================================
// Notification Service Endpoints
// =============================================================================

export const notificationApi = api.injectEndpoints({
  endpoints: (builder) => ({
    /**
     * Get user notifications
     */
    getUserNotifications: builder.query<Notification[], { unreadOnly?: boolean; limit?: number }>({
      query: ({ unreadOnly, limit = 50 }) => {
        const params = new URLSearchParams();
        if (unreadOnly) params.append('unread', 'true');
        params.append('limit', limit.toString());
        return {
          url: `/notifications?${params}`,
          baseUrl: SERVICES.NOTIFICATION
        };
      },
      providesTags: ['Notification'],
      transformResponse: (response: ApiResponse<Notification[]>) => response.data,
    }),

    /**
     * Mark notification as read
     */
    markNotificationRead: builder.mutation<void, string>({
      query: (notificationId) => ({
        url: `/notifications/${notificationId}/read`,
        method: 'PATCH',
        baseUrl: SERVICES.NOTIFICATION
      }),
      invalidatesTags: (result, error, id) => [{ type: 'Notification', id }],
    }),

    /**
     * Mark all notifications as read
     */
    markAllNotificationsRead: builder.mutation<void, void>({
      query: () => ({
        url: '/notifications/mark-all-read',
        method: 'PATCH',
        baseUrl: SERVICES.NOTIFICATION
      }),
      invalidatesTags: ['Notification'],
    }),

    /**
     * Subscribe to push notifications
     */
    subscribeToPushNotifications: builder.mutation<void, { subscription: any }>({
      query: (data) => ({
        url: '/push-subscriptions',
        method: 'POST',
        body: data,
        baseUrl: SERVICES.NOTIFICATION
      }),
    })
  }),
});

// =============================================================================
// Social Service Endpoints
// =============================================================================
export const socialApi = api.injectEndpoints({
  endpoints: (builder) => ({
    /**
     * Get social feed posts
     */
    getSocialFeed: builder.query<any[], { type?: string; limit?: number; page?: number }>({
      query: ({ type = 'all', limit = 20, page = 1 }) => {
        const params = new URLSearchParams();
        params.append('type', type);
        params.append('limit', limit.toString());
        params.append('page', page.toString());
        return {
          url: `/social/feed?${params}`,
          baseUrl: SERVICES.CONTENT
        };
      },
      providesTags: ['SocialPost'],
      transformResponse: (response: ApiResponse<any[]>) => response.data,
    }),

    /**
     * Get social stories/moments
     */
    getSocialStories: builder.query<any[], { active?: boolean; limit?: number }>({
      query: ({ active = true, limit = 20 }) => {
        const params = new URLSearchParams();
        if (active) params.append('active', 'true');
        params.append('limit', limit.toString());
        return {
          url: `/social/stories?${params}`,
          baseUrl: SERVICES.CONTENT
        };
      },
      providesTags: ['SocialStory'],
      transformResponse: (response: ApiResponse<any[]>) => response.data,
    }),

    /**
     * Get user friends/connections
     */
    getUserFriends: builder.query<any[], { status?: string; limit?: number }>({
      query: ({ status = 'active', limit = 50 }) => {
        const params = new URLSearchParams();
        params.append('status', status);
        params.append('limit', limit.toString());
        return {
          url: `/friends?${params}`,
          baseUrl: SERVICES.MESSAGING
        };
      },
      providesTags: ['Friend'],
      transformResponse: (response: ApiResponse<any[]>) => response.data,
    }),

    /**
     * Like a social post
     */
    likeSocialPost: builder.mutation<void, { postId: string; isLiked: boolean }>({
      query: ({ postId, isLiked }) => ({
        url: `/social/posts/${postId}/like`,
        method: isLiked ? 'POST' : 'DELETE',
        baseUrl: SERVICES.CONTENT
      }),
      invalidatesTags: ['SocialPost'],
    }),

    /**
     * Create a social post
     */
    createSocialPost: builder.mutation<any, { content: string; images?: string[]; location?: string; tags?: string[]; type?: string }>({
      query: (postData) => ({
        url: '/social/posts',
        method: 'POST',
        body: postData,
        baseUrl: SERVICES.CONTENT
      }),
      invalidatesTags: ['SocialPost'],
    }),

    /**
     * Follow/unfollow a user
     */
    followUser: builder.mutation<void, { userId: string; action: 'follow' | 'unfollow' }>({
      query: ({ userId, action }) => ({
        url: `/friends/${userId}/${action}`,
        method: 'POST',
        baseUrl: SERVICES.MESSAGING
      }),
      invalidatesTags: ['Friend'],
    }),
  }),
});

// =============================================================================
// Export Generated Hooks
// =============================================================================

export const {
  // Messaging hooks
  useGetUserChatsQuery,
  useGetChatMessagesQuery,
  useSendChatMessageMutation,
  useCreateChatMutation,

  // Commerce hooks - Orders only (products/shops exported from commerceApi above)
  useGetProductQuery,
  useCreateOrderMutation,
  useGetUserOrdersQuery,
  useUpdateOrderStatusMutation,

  // Payment hooks
  useGetUserWalletQuery,
  useAddMoneyToWalletMutation,
  useSendMoneyMutation,
  useProcessPaymentMutation,
  useAddPaymentMethodMutation,
  useGetTransactionHistoryQuery,

  // Notification hooks
  useGetUserNotificationsQuery,
  useMarkNotificationReadMutation,
  useMarkAllNotificationsReadMutation,
  useSubscribeToPushNotificationsMutation,

  // Social hooks
  useGetSocialFeedQuery,
  useGetSocialStoriesQuery,
  useGetUserFriendsQuery,
  useLikeSocialPostMutation,
  useCreateSocialPostMutation,
  useFollowUserMutation,

  // Workspace hooks
  useGetUpcomingFeaturesQuery,
  useGetWorkspaceStatsQuery,
  useRequestBetaAccessMutation,
  useSubscribeToFeatureMutation,

  // Discover hooks
  useGetTrendingPostsQuery,
  useGetTrendingTopicsQuery,
  useGetSuggestedUsersQuery,
  useSearchDiscoverQuery,

  // Events hooks
  useGetUpcomingEventsQuery,
  useGetHistoricalEventsQuery,
  useGetEventDetailsQuery,
  useRegisterEventInterestMutation,
  useRsvpToEventMutation,
  useBookEventTicketsMutation
} = api;

// =============================================================================
// Workspace/Work Service Endpoints
// =============================================================================
export const workspaceApi = api.injectEndpoints({
  endpoints: (builder) => ({
    /**
     * Get upcoming work features
     */
    getUpcomingFeatures: builder.query<any[], { status?: string; limit?: number }>({
      query: ({ status, limit = 20 }) => {
        const params = new URLSearchParams();
        if (status) params.append('status', status);
        params.append('limit', limit.toString());
        return {
          url: `/workspace/features?${params}`,
          baseUrl: SERVICES.CONTENT
        };
      },
      providesTags: ['WorkFeature'],
      transformResponse: (response: ApiResponse<any[]>) => response.data,
    }),

    /**
     * Get workspace statistics
     */
    getWorkspaceStats: builder.query<any, void>({
      query: () => ({
        url: `/workspace/stats`,
        baseUrl: SERVICES.CONTENT
      }),
      providesTags: ['WorkStats'],
      transformResponse: (response: ApiResponse<any>) => response.data,
    }),

    /**
     * Request beta access for a feature
     */
    requestBetaAccess: builder.mutation<any, { featureId: string; userInfo: any }>({
      query: ({ featureId, userInfo }) => ({
        url: `/workspace/features/${featureId}/beta-request`,
        method: 'POST',
        body: userInfo,
        baseUrl: SERVICES.CONTENT
      }),
      invalidatesTags: ['WorkFeature'],
    }),

    /**
     * Subscribe to feature notifications
     */
    subscribeToFeature: builder.mutation<any, { featureId: string }>({
      query: ({ featureId }) => ({
        url: `/workspace/features/${featureId}/subscribe`,
        method: 'POST',
        baseUrl: SERVICES.CONTENT
      }),
      invalidatesTags: ['WorkFeature'],
    })
  }),
});

// =============================================================================
// Discover/Trending Service Endpoints
// =============================================================================
export const discoverApi = api.injectEndpoints({
  endpoints: (builder) => ({
    /**
     * Get trending posts
     */
    getTrendingPosts: builder.query<any[], { limit?: number; timeframe?: string; category?: string }>({
      query: ({ limit = 20, timeframe = '24h', category }) => {
        const params = new URLSearchParams();
        params.append('limit', limit.toString());
        params.append('timeframe', timeframe);
        if (category) params.append('category', category);
        return {
          url: `/discover/trending/posts?${params}`,
          baseUrl: SERVICES.CONTENT
        };
      },
      providesTags: ['TrendingPost'],
      transformResponse: (response: ApiResponse<any[]>) => response.data,
    }),

    /**
     * Get trending topics and hashtags
     */
    getTrendingTopics: builder.query<any[], { category?: string; timeframe?: string; limit?: number }>({
      query: ({ category, timeframe = '24h', limit = 20 }) => {
        const params = new URLSearchParams();
        if (category) params.append('category', category);
        params.append('timeframe', timeframe);
        params.append('limit', limit.toString());
        return {
          url: `/discover/trending/topics?${params}`,
          baseUrl: SERVICES.CONTENT
        };
      },
      providesTags: ['TrendingTopic'],
      transformResponse: (response: ApiResponse<any[]>) => response.data,
    }),

    /**
     * Get suggested users to follow
     */
    getSuggestedUsers: builder.query<any[], { category?: string; limit?: number }>({
      query: ({ category, limit = 10 }) => {
        const params = new URLSearchParams();
        if (category) params.append('category', category);
        params.append('limit', limit.toString());
        return {
          url: `/discover/users/suggested?${params}`,
          baseUrl: SERVICES.CONTENT
        };
      },
      providesTags: ['SuggestedUser'],
      transformResponse: (response: ApiResponse<any[]>) => response.data,
    }),

    /**
     * Search discover content
     */
    searchDiscover: builder.query<any, { query: string; type?: string; limit?: number }>({
      query: ({ query, type = 'all', limit = 20 }) => {
        const params = new URLSearchParams();
        params.append('q', query);
        params.append('type', type);
        params.append('limit', limit.toString());
        return {
          url: `/discover/search?${params}`,
          baseUrl: SERVICES.CONTENT
        };
      },
      providesTags: ['SearchResult'],
      transformResponse: (response: ApiResponse<any>) => response.data,
    })
  }),
});

// =============================================================================
// Events Service Endpoints
// =============================================================================
export const eventsApi = api.injectEndpoints({
  endpoints: (builder) => ({
    /**
     * Get upcoming events
     */
    getUpcomingEvents: builder.query<any[], { category?: string; location?: string; limit?: number; page?: number }>({
      query: ({ category, location, limit = 20, page = 1 }) => {
        const params = new URLSearchParams();
        if (category) params.append('category', category);
        if (location) params.append('location', location);
        params.append('limit', limit.toString());
        params.append('page', page.toString());
        return {
          url: `/events/upcoming?${params}`,
          baseUrl: SERVICES.CONTENT
        };
      },
      providesTags: ['Event'],
      transformResponse: (response: ApiResponse<any[]>) => response.data,
    }),

    /**
     * Get historical/past events
     */
    getHistoricalEvents: builder.query<any[], { category?: string; limit?: number; page?: number }>({
      query: ({ category, limit = 10, page = 1 }) => {
        const params = new URLSearchParams();
        if (category) params.append('category', category);
        params.append('limit', limit.toString());
        params.append('page', page.toString());
        return {
          url: `/events/historical?${params}`,
          baseUrl: SERVICES.CONTENT
        };
      },
      providesTags: ['HistoricalEvent'],
      transformResponse: (response: ApiResponse<any[]>) => response.data,
    }),

    /**
     * Get event details
     */
    getEventDetails: builder.query<any, { eventId: string }>({
      query: ({ eventId }) => ({
        url: `/events/${eventId}`,
        baseUrl: SERVICES.CONTENT
      }),
      providesTags: ['Event'],
      transformResponse: (response: ApiResponse<any>) => response.data,
    }),

    /**
     * Register interest in an event
     */
    registerEventInterest: builder.mutation<any, { eventId: string; interested: boolean }>({
      query: ({ eventId, interested }) => ({
        url: `/events/${eventId}/interest`,
        method: 'POST',
        body: { interested },
        baseUrl: SERVICES.CONTENT
      }),
      invalidatesTags: ['Event'],
    }),

    /**
     * RSVP to an event (attending)
     */
    rsvpToEvent: builder.mutation<any, { eventId: string; attending: boolean }>({
      query: ({ eventId, attending }) => ({
        url: `/events/${eventId}/rsvp`,
        method: 'POST',
        body: { attending },
        baseUrl: SERVICES.CONTENT
      }),
      invalidatesTags: ['Event'],
    }),

    /**
     * Book event tickets
     */
    bookEventTickets: builder.mutation<any, { eventId: string; ticketType: string; quantity: number; userInfo?: any }>({
      query: ({ eventId, ticketType, quantity, userInfo }) => ({
        url: `/events/${eventId}/book`,
        method: 'POST',
        body: { ticketType, quantity, userInfo },
        baseUrl: SERVICES.CONTENT
      }),
      invalidatesTags: ['Event'],
    })
  }),
});

// =============================================================================
// Authentication API - Critical for routes.tsx mockUser replacement
// =============================================================================

export const authApi = api.injectEndpoints({
  endpoints: (builder) => ({
    getCurrentUser: builder.query<User, void>({
      query: () => ({
        url: '/users/profile',
        baseUrl: SERVICES.AUTH
      }),
      providesTags: ['User'],
      transformResponse: (response: ApiResponse<User>) => response.data,
    }),

    updateUserProfile: builder.mutation<User, Partial<User>>({
      query: (userData) => ({
        url: '/users/profile',
        method: 'PATCH',
        body: userData,
        baseUrl: SERVICES.AUTH
      }),
      invalidatesTags: ['User'],
      transformResponse: (response: ApiResponse<User>) => response.data,
    }),

    getUserSession: builder.query<any, void>({
      query: () => ({
        url: '/auth/session',
        baseUrl: SERVICES.AUTH
      }),
      providesTags: ['Session'],
      transformResponse: (response: ApiResponse<any>) => response.data,
    })
  }),
});

// Export the generated hooks from authApi
export const {
  useGetCurrentUserQuery,
  useUpdateUserProfileMutation,
  useGetUserSessionQuery,
} = authApi;

// =============================================================================
// Product Reviews Service Endpoints (E-commerce Reviews)
// =============================================================================

export const reviewsApi = api.injectEndpoints({
  endpoints: (builder) => ({
    /**
     * Get product reviews with pagination
     */
    getProductReviews: builder.query<any[], {
      productId: string;
      page?: number;
      limit?: number;
      sortBy?: 'newest' | 'oldest' | 'rating_high' | 'rating_low' | 'helpful';
    }>({
      query: ({ productId, page = 1, limit = 10, sortBy = 'newest' }) => {
        const params = new URLSearchParams();
        params.append('page', page.toString());
        params.append('limit', limit.toString());
        params.append('sort', sortBy);
        return {
          url: `/products/${productId}/reviews?${params}`,
          baseUrl: SERVICES.COMMERCE
        };
      },
      providesTags: (result, error, { productId }) => [
        { type: 'ProductReview', id: productId },
        { type: 'ProductReview', id: 'LIST' }
      ],
      transformResponse: (response: ApiResponse<any[]>) => response.data,
    }),

    /**
     * Add a new product review
     */
    addProductReview: builder.mutation<any, {
      productId: string;
      rating: number;
      comment: string;
      images?: string[];
    }>({
      query: ({ productId, ...reviewData }) => ({
        url: `/products/${productId}/reviews`,
        method: 'POST',
        body: reviewData,
        baseUrl: SERVICES.COMMERCE
      }),
      invalidatesTags: (result, error, { productId }) => [
        { type: 'ProductReview', id: productId },
        { type: 'ProductReview', id: 'LIST' },
        { type: 'Product', id: productId }
      ],
      transformResponse: (response: ApiResponse<any>) => response.data,
    }),

    /**
     * Mark review as helpful
     */
    markReviewHelpful: builder.mutation<any, { reviewId: string }>({
      query: ({ reviewId }) => ({
        url: `/reviews/${reviewId}/helpful`,
        method: 'POST',
        baseUrl: SERVICES.COMMERCE
      }),
      invalidatesTags: ['ProductReview'],
      transformResponse: (response: ApiResponse<any>) => response.data,
    }),

    /**
     * Get review statistics for a product
     */
    getProductReviewStats: builder.query<any, { productId: string }>({
      query: ({ productId }) => ({
        url: `/products/${productId}/reviews/stats`,
        baseUrl: SERVICES.COMMERCE
      }),
      providesTags: (result, error, { productId }) => [
        { type: 'ProductReview', id: `${productId}-stats` }
      ],
      transformResponse: (response: ApiResponse<any>) => response.data,
    })
  }),
});

// Export the generated hooks
export const {
  useGetProductReviewsQuery,
  useAddProductReviewMutation,
  useMarkReviewHelpfulMutation,
  useGetProductReviewStatsQuery
} = reviewsApi;

// =============================================================================
// Service Health Check Utilities
// =============================================================================

/**
 * Check health of all microservices
 */
export const checkServiceHealth = async (): Promise<Record<string, boolean>> => {
  const healthChecks: Record<string, boolean> = {};

  for (const [serviceName, serviceUrl] of Object.entries(SERVICES)) {
    try {
      const response = await fetch(`${serviceUrl}/health`, {
        method: 'GET',
        timeout: 5000
      } as RequestInit);
      healthChecks[serviceName] = response.ok;
    } catch (error) {
      healthChecks[serviceName] = false;
      console.warn(`Service ${serviceName} health check failed:`, error);
    }
  }

  return healthChecks;
};

/**
 * Currency conversion utilities for Southeast Asian markets
 */
export const CURRENCY_CONFIG = {
  THB: { symbol: '฿', name: 'Thai Baht', country: 'Thailand' },
  IDR: { symbol: 'Rp', name: 'Indonesian Rupiah', country: 'Indonesia' },
  SGD: { symbol: 'S$', name: 'Singapore Dollar', country: 'Singapore' },
  PHP: { symbol: '₱', name: 'Philippine Peso', country: 'Philippines' },
  MYR: { symbol: 'RM', name: 'Malaysian Ringgit', country: 'Malaysia' },
  VND: { symbol: '₫', name: 'Vietnamese Dong', country: 'Vietnam' }
};

export const formatCurrency = (amount: number, currency: string): string => {
  const config = CURRENCY_CONFIG[currency as keyof typeof CURRENCY_CONFIG];
  if (!config) return `${amount} ${currency}`;

  return `${config.symbol}${amount.toLocaleString()}`;
};