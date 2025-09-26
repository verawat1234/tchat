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

const SERVICES = {
  GATEWAY: import.meta.env.VITE_GATEWAY_URL || 'http://localhost:8080/api/v1',
  AUTH: import.meta.env.VITE_AUTH_URL || 'http://localhost:8081/api/v1',
  CONTENT: import.meta.env.VITE_CONTENT_URL || 'http://localhost:8082/api/v1',
  MESSAGING: import.meta.env.VITE_MESSAGING_URL || 'http://localhost:8083/api/v1',
  COMMERCE: import.meta.env.VITE_COMMERCE_URL || 'http://localhost:8084/api/v1',
  PAYMENT: import.meta.env.VITE_PAYMENT_URL || 'http://localhost:8085/api/v1',
  NOTIFICATION: import.meta.env.VITE_NOTIFICATION_URL || 'http://localhost:8086/api/v1'
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
    })
  }),
});

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
// Export Generated Hooks
// =============================================================================

export const {
  // Messaging hooks
  useGetUserChatsQuery,
  useGetChatMessagesQuery,
  useSendChatMessageMutation,
  useCreateChatMutation,

  // Commerce hooks
  useGetFeaturedProductsQuery,
  useSearchProductsQuery,
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
  useSubscribeToPushNotificationsMutation
} = api;

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