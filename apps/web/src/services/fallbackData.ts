/**
 * Fallback Data Service for Development
 *
 * Provides mock data when backend API endpoints are not available,
 * ensuring smooth development experience while backend services are being built.
 */

// =============================================================================
// Mock Data Collections
// =============================================================================

export const MOCK_USERS = [
  {
    id: 'user-1',
    name: 'Sarah Chen',
    email: 'sarah.chen@example.com',
    avatar: 'https://images.unsplash.com/photo-1494790108755-2616b2e0d36c?w=100&h=100&fit=crop&crop=face',
    status: 'online',
    country: 'SG'
  },
  {
    id: 'user-2',
    name: 'Arif Rahman',
    email: 'arif.rahman@example.com',
    avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&h=100&fit=crop&crop=face',
    status: 'away',
    country: 'ID'
  },
  {
    id: 'user-3',
    name: 'Thicha Sansern',
    email: 'thicha@example.com',
    avatar: 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=100&h=100&fit=crop&crop=face',
    status: 'online',
    country: 'TH'
  }
];

export const MOCK_PRODUCTS = [
  {
    id: 'prod-1',
    name: 'Traditional Thai Silk Scarf',
    description: 'Handwoven silk scarf from Northern Thailand with traditional patterns',
    price: 1250,
    currency: 'THB',
    images: ['https://images.unsplash.com/photo-1594633312681-425c7b97ccd1?w=400&h=400&fit=crop'],
    category: 'Fashion',
    merchant: {
      id: 'merchant-1',
      name: 'Thai Heritage Crafts',
      country: 'TH'
    },
    rating: 4.8,
    reviewCount: 127,
    shipping: {
      available: true,
      cost: 150,
      estimatedDays: 3
    },
    createdAt: '2024-01-15T10:00:00Z',
    updatedAt: '2024-03-10T14:30:00Z'
  },
  {
    id: 'prod-2',
    name: 'Indonesian Batik Shirt',
    description: 'Authentic Javanese batik shirt with traditional motifs',
    price: 85000,
    currency: 'IDR',
    images: ['https://images.unsplash.com/photo-1621072156002-e2fccdc0b176?w=400&h=400&fit=crop'],
    category: 'Fashion',
    merchant: {
      id: 'merchant-2',
      name: 'Jakarta Batik House',
      country: 'ID'
    },
    rating: 4.6,
    reviewCount: 89,
    shipping: {
      available: true,
      cost: 25000,
      estimatedDays: 2
    },
    createdAt: '2024-02-20T09:15:00Z',
    updatedAt: '2024-03-08T16:45:00Z'
  },
  {
    id: 'prod-3',
    name: 'Singapore Laksa Spice Kit',
    description: 'Authentic laksa spice blend kit with coconut milk and instructions',
    price: 28,
    currency: 'SGD',
    images: ['https://images.unsplash.com/photo-1569718212165-3a8278d5f624?w=400&h=400&fit=crop'],
    category: 'Food',
    merchant: {
      id: 'merchant-3',
      name: 'Singapore Flavors',
      country: 'SG'
    },
    rating: 4.9,
    reviewCount: 203,
    shipping: {
      available: true,
      cost: 5,
      estimatedDays: 1
    },
    createdAt: '2024-01-30T11:20:00Z',
    updatedAt: '2024-03-12T13:10:00Z'
  }
];

export const MOCK_CHATS = [
  {
    id: 'chat-1',
    type: 'direct',
    participants: ['user-1', 'user-2'],
    lastMessage: {
      id: 'msg-1',
      content: 'Hey! How are you doing?',
      senderId: 'user-1',
      timestamp: '2024-03-15T10:30:00Z',
      type: 'text'
    },
    unreadCount: 2,
    updatedAt: '2024-03-15T10:30:00Z'
  },
  {
    id: 'chat-2',
    type: 'group',
    name: 'Southeast Asia Travel',
    participants: ['user-1', 'user-2', 'user-3'],
    lastMessage: {
      id: 'msg-2',
      content: 'Anyone been to Chiang Mai recently?',
      senderId: 'user-3',
      timestamp: '2024-03-15T09:45:00Z',
      type: 'text'
    },
    unreadCount: 0,
    updatedAt: '2024-03-15T09:45:00Z'
  }
];

// Enhanced messages with comprehensive message types aligned with mobile platform
export const MOCK_MESSAGES = [
  {
    id: 'msg-1',
    chatId: 'chat-1',
    content: 'Hey! How are you doing?',
    senderId: 'user-1',
    senderName: 'Sarah Chen',
    senderAvatar: 'https://images.unsplash.com/photo-1494790108755-2616b2e0d36c?w=100&h=100&fit=crop&crop=face',
    timestamp: '2024-03-15T10:30:00Z',
    type: 'text',
    status: 'delivered',
    messageType: 'TEXT',
    isEdited: false,
    isPinned: false,
    isDeleted: false,
    reactions: [],
    attachments: [],
    deliveryStatus: 'DELIVERED',
    readBy: ['user-2']
  },
  {
    id: 'msg-2',
    chatId: 'chat-1',
    content: 'I am doing great! Just got back from Bali. The beaches were amazing! 🏖️',
    senderId: 'user-2',
    senderName: 'Arif Rahman',
    senderAvatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&h=100&fit=crop&crop=face',
    timestamp: '2024-03-15T10:32:00Z',
    type: 'text',
    status: 'read',
    messageType: 'TEXT',
    isEdited: false,
    isPinned: false,
    isDeleted: false,
    reactions: [
      {
        emoji: '🏖️',
        userId: 'user-1',
        userName: 'Sarah Chen',
        timestamp: '2024-03-15T10:33:00Z'
      }
    ],
    attachments: [],
    deliveryStatus: 'READ',
    readBy: ['user-1']
  },
  {
    id: 'msg-3',
    chatId: 'chat-1',
    content: 'That sounds incredible! I am planning a trip to Thailand next month.',
    senderId: 'user-1',
    senderName: 'Sarah Chen',
    senderAvatar: 'https://images.unsplash.com/photo-1494790108755-2616b2e0d36c?w=100&h=100&fit=crop&crop=face',
    timestamp: '2024-03-15T10:35:00Z',
    type: 'text',
    status: 'sent',
    messageType: 'TEXT',
    isEdited: false,
    isPinned: false,
    isDeleted: false,
    reactions: [],
    attachments: [],
    deliveryStatus: 'SENT',
    readBy: []
  },
  {
    id: 'msg-4',
    chatId: 'chat-1',
    content: 'Check out this amazing temple I visited!',
    senderId: 'user-2',
    senderName: 'Arif Rahman',
    senderAvatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&h=100&fit=crop&crop=face',
    timestamp: '2024-03-15T10:37:00Z',
    type: 'image',
    status: 'delivered',
    messageType: 'IMAGE',
    isEdited: false,
    isPinned: false,
    isDeleted: false,
    reactions: [
      {
        emoji: '😍',
        userId: 'user-1',
        userName: 'Sarah Chen',
        timestamp: '2024-03-15T10:38:00Z'
      },
      {
        emoji: '🙏',
        userId: 'user-1',
        userName: 'Sarah Chen',
        timestamp: '2024-03-15T10:38:30Z'
      }
    ],
    attachments: [
      {
        id: 'att-1',
        type: 'IMAGE',
        url: 'https://images.unsplash.com/photo-1563492065037-3b3c1aa9be1d?w=600&h=400&fit=crop',
        thumbnail: 'https://images.unsplash.com/photo-1563492065037-3b3c1aa9be1d?w=200&h=150&fit=crop',
        caption: 'Wat Arun temple in Bangkok - absolutely stunning!',
        width: 600,
        height: 400,
        metadata: {
          location: 'Wat Arun, Bangkok, Thailand'
        }
      }
    ],
    deliveryStatus: 'DELIVERED',
    readBy: ['user-1']
  },
  {
    id: 'msg-5',
    chatId: 'chat-1',
    content: 'Where is this exactly?',
    senderId: 'user-1',
    senderName: 'Sarah Chen',
    senderAvatar: 'https://images.unsplash.com/photo-1494790108755-2616b2e0d36c?w=100&h=100&fit=crop&crop=face',
    timestamp: '2024-03-15T10:40:00Z',
    type: 'text',
    status: 'read',
    messageType: 'TEXT',
    isEdited: false,
    isPinned: false,
    isDeleted: false,
    reactions: [],
    attachments: [],
    deliveryStatus: 'READ',
    readBy: ['user-2']
  },
  {
    id: 'msg-6',
    chatId: 'chat-1',
    content: 'Here\'s the exact location!',
    senderId: 'user-2',
    senderName: 'Arif Rahman',
    senderAvatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&h=100&fit=crop&crop=face',
    timestamp: '2024-03-15T10:42:00Z',
    type: 'location',
    status: 'delivered',
    messageType: 'LOCATION',
    isEdited: false,
    isPinned: false,
    isDeleted: false,
    reactions: [
      {
        emoji: '📍',
        userId: 'user-1',
        userName: 'Sarah Chen',
        timestamp: '2024-03-15T10:43:00Z'
      }
    ],
    attachments: [
      {
        id: 'att-2',
        type: 'LOCATION',
        url: 'https://maps.google.com/?q=13.7443,100.4893',
        metadata: {
          latitude: '13.7443',
          longitude: '100.4893',
          address: 'Wat Arun Ratchawararam Ratchawaramahawihan, Bangkok, Thailand',
          placeName: 'Wat Arun (Temple of Dawn)'
        }
      }
    ],
    deliveryStatus: 'DELIVERED',
    readBy: ['user-1']
  },
  {
    id: 'msg-7',
    chatId: 'chat-1',
    content: 'Let\'s plan a Southeast Asia meetup!',
    senderId: 'user-1',
    senderName: 'Sarah Chen',
    senderAvatar: 'https://images.unsplash.com/photo-1494790108755-2616b2e0d36c?w=100&h=100&fit=crop&crop=face',
    timestamp: '2024-03-15T11:00:00Z',
    type: 'poll',
    status: 'delivered',
    messageType: 'POLL',
    isEdited: false,
    isPinned: false,
    isDeleted: false,
    reactions: [
      {
        emoji: '🎉',
        userId: 'user-2',
        userName: 'Arif Rahman',
        timestamp: '2024-03-15T11:01:00Z'
      }
    ],
    attachments: [],
    poll: {
      question: 'Which city should we meet in?',
      options: ['Bangkok, Thailand', 'Singapore', 'Ho Chi Minh City, Vietnam', 'Kuala Lumpur, Malaysia'],
      votes: { 0: 2, 1: 3, 2: 1, 3: 2 },
      expiresAt: '2024-03-20T11:00:00Z'
    },
    deliveryStatus: 'DELIVERED',
    readBy: ['user-2']
  },
  {
    id: 'msg-8',
    chatId: 'chat-1',
    content: 'Digital Nomad Meetup - Bangkok',
    senderId: 'user-2',
    senderName: 'Arif Rahman',
    senderAvatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&h=100&fit=crop&crop=face',
    timestamp: '2024-03-15T14:30:00Z',
    type: 'event',
    status: 'delivered',
    messageType: 'EVENT',
    isEdited: false,
    isPinned: true,
    isDeleted: false,
    reactions: [
      {
        emoji: '📅',
        userId: 'user-1',
        userName: 'Sarah Chen',
        timestamp: '2024-03-15T14:31:00Z'
      }
    ],
    attachments: [],
    event: {
      title: 'Digital Nomad Meetup - Bangkok',
      description: 'Let\'s meet up for coffee and co-working in Bangkok! Great opportunity to network and share experiences.',
      startDate: '2024-04-15T09:00:00Z',
      endDate: '2024-04-15T17:00:00Z',
      location: {
        name: 'WeWork Ploenchit',
        address: '518 Ploenchit Rd, Lumphini, Pathum Wan, Bangkok 10330',
        latitude: 13.7420,
        longitude: 100.5480
      },
      attendees: ['user-1', 'user-2'],
      maxAttendees: 20,
      rsvpDeadline: '2024-04-10T23:59:59Z'
    },
    deliveryStatus: 'DELIVERED',
    readBy: ['user-1']
  },
  {
    id: 'msg-9',
    chatId: 'chat-1',
    content: 'Check out this great co-working space!',
    senderId: 'user-1',
    senderName: 'Sarah Chen',
    senderAvatar: 'https://images.unsplash.com/photo-1494790108755-2616b2e0d36c?w=100&h=100&fit=crop&crop=face',
    timestamp: '2024-03-15T16:15:00Z',
    type: 'embed',
    status: 'read',
    messageType: 'EMBED',
    isEdited: false,
    isPinned: false,
    isDeleted: false,
    reactions: [
      {
        emoji: '💼',
        userId: 'user-2',
        userName: 'Arif Rahman',
        timestamp: '2024-03-15T16:16:00Z'
      }
    ],
    attachments: [],
    embed: {
      title: 'Hubud - Bali Co-working Space',
      description: 'A beautiful bamboo co-working space in the heart of Ubud, Bali. Perfect for digital nomads!',
      url: 'https://hubud.org',
      image: 'https://images.unsplash.com/photo-1497366216548-37526070297c?w=400&h=300&fit=crop',
      siteName: 'Hubud',
      type: 'website'
    },
    deliveryStatus: 'READ',
    readBy: ['user-2']
  },
  {
    id: 'msg-10',
    chatId: 'chat-1',
    content: 'Split the bill for last night\'s dinner?',
    senderId: 'user-2',
    senderName: 'Arif Rahman',
    senderAvatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&h=100&fit=crop&crop=face',
    timestamp: '2024-03-15T18:20:00Z',
    type: 'payment',
    status: 'delivered',
    messageType: 'PAYMENT_REQUEST',
    isEdited: false,
    isPinned: false,
    isDeleted: false,
    reactions: [],
    attachments: [],
    payment: {
      type: 'request',
      amount: 1850,
      currency: 'THB',
      description: 'Dinner at Som Tam Nua - Thai Restaurant',
      recipient: 'user-1',
      dueDate: '2024-03-20T23:59:59Z',
      status: 'pending'
    },
    deliveryStatus: 'DELIVERED',
    readBy: ['user-1']
  }
];

export const MOCK_WALLET = {
  id: 'wallet-1',
  userId: 'user-1',
  balances: [
    { currency: 'THB', amount: 15420.50 },
    { currency: 'SGD', amount: 875.25 },
    { currency: 'IDR', amount: 2340000 }
  ],
  paymentMethods: [
    {
      id: 'pm-1',
      type: 'PromptPay',
      details: { phoneNumber: '+66812345678' },
      isDefault: true
    },
    {
      id: 'pm-2',
      type: 'GrabPay',
      details: { accountId: 'grabpay-123' },
      isDefault: false
    }
  ],
  transactions: [
    {
      id: 'txn-1',
      type: 'credit',
      amount: 500,
      currency: 'THB',
      description: 'Payment received from Sarah Chen',
      timestamp: '2024-03-15T08:30:00Z'
    },
    {
      id: 'txn-2',
      type: 'debit',
      amount: 75,
      currency: 'SGD',
      description: 'Purchase: Singapore Laksa Spice Kit',
      timestamp: '2024-03-14T16:45:00Z'
    },
    {
      id: 'txn-3',
      type: 'credit',
      amount: 1250000,
      currency: 'IDR',
      description: 'Wallet top-up via GoPay',
      timestamp: '2024-03-13T12:20:00Z'
    }
  ]
};

export const MOCK_NOTIFICATIONS = [
  {
    id: 'notif-1',
    userId: 'user-1',
    type: 'message',
    title: 'New message from Arif Rahman',
    content: 'I am doing great! Just got back from Bali...',
    data: { chatId: 'chat-1', messageId: 'msg-2' },
    read: false,
    createdAt: '2024-03-15T10:32:00Z'
  },
  {
    id: 'notif-2',
    userId: 'user-1',
    type: 'payment',
    title: 'Payment received',
    content: 'You received ฿500.00 from Sarah Chen',
    data: { transactionId: 'txn-1', amount: 500, currency: 'THB' },
    read: false,
    createdAt: '2024-03-15T08:30:00Z'
  },
  {
    id: 'notif-3',
    userId: 'user-1',
    type: 'order',
    title: 'Order shipped',
    content: 'Your Singapore Laksa Spice Kit has been shipped',
    data: { orderId: 'order-1', trackingNumber: 'SG123456789' },
    read: true,
    createdAt: '2024-03-14T14:20:00Z'
  }
];

// Enhanced social feed with comprehensive post types aligned with mobile platform
export const MOCK_SOCIAL_FEED = [
  {
    id: 'post-1',
    userId: 'user-2',
    type: 'IMAGE',
    content: 'Amazing sunset at Tanah Lot, Bali! The colors were absolutely breathtaking. 🌅',
    images: ['https://images.unsplash.com/photo-1537953773345-d172ccf13cf1?w=600&h=400&fit=crop'],
    location: 'Tanah Lot, Bali, Indonesia',
    tags: ['travel', 'bali', 'sunset', 'indonesia'],
    likes: 127,
    comments: 23,
    shares: 8,
    createdAt: '2024-03-14T18:45:00Z',
    postType: 'TRAVEL_LOG'
  },
  {
    id: 'post-2',
    userId: 'user-3',
    type: 'TEXT',
    content: 'Just tried the most incredible Pad Thai at this tiny street stall in Bangkok. Sometimes the best food comes from the most unexpected places! 🍜 #streetfood #bangkok #thailand',
    location: 'Bangkok, Thailand',
    tags: ['food', 'bangkok', 'thailand', 'streetfood'],
    likes: 89,
    comments: 15,
    shares: 4,
    createdAt: '2024-03-15T12:30:00Z',
    postType: 'RECIPE'
  },
  {
    id: 'post-3',
    userId: 'user-1',
    type: 'IMAGE',
    content: 'Marina Bay Sands at night never gets old. Singapore, you are beautiful! ✨',
    images: ['https://images.unsplash.com/photo-1525625293386-3f8f99389edd?w=600&h=400&fit=crop'],
    location: 'Marina Bay, Singapore',
    tags: ['singapore', 'night', 'cityscape', 'marinabay'],
    likes: 156,
    comments: 31,
    shares: 12,
    createdAt: '2024-03-13T20:15:00Z',
    postType: 'CHECK_IN'
  },
  {
    id: 'post-4',
    userId: 'user-1',
    type: 'POLL',
    content: 'Planning my next Southeast Asia adventure! Where should I go next? 🌍',
    poll: {
      question: 'Best Southeast Asian destination for digital nomads?',
      options: ['Chiang Mai, Thailand', 'Canggu, Bali', 'Ho Chi Minh City, Vietnam', 'Penang, Malaysia'],
      votes: { 0: 45, 1: 67, 2: 23, 3: 31 },
      expiresAt: '2024-03-20T18:00:00Z'
    },
    tags: ['travel', 'digitalnomad', 'southeastasia'],
    likes: 89,
    comments: 42,
    shares: 15,
    createdAt: '2024-03-15T14:20:00Z',
    postType: 'POLL'
  },
  {
    id: 'post-5',
    userId: 'user-2',
    type: 'PRODUCT_SHOWCASE',
    content: 'New handwoven silk scarves from local artisans in Luang Prabang! Supporting traditional crafts and fair trade. 🧣✨',
    images: ['https://images.unsplash.com/photo-1594633312681-425c7b97ccd1?w=600&h=400&fit=crop'],
    location: 'Luang Prabang, Laos',
    tags: ['handmade', 'fairtrade', 'laos', 'silk'],
    likes: 234,
    comments: 18,
    shares: 32,
    createdAt: '2024-03-14T09:30:00Z',
    postType: 'PRODUCT_SHOWCASE',
    product: {
      name: 'Traditional Silk Scarf',
      price: 45,
      currency: 'USD'
    }
  },
  {
    id: 'post-6',
    userId: 'user-3',
    type: 'LIVE_STREAM',
    content: 'LIVE: Sunset cooking class from my rooftop in Bangkok! Learning to make authentic Tom Yum 🍲',
    location: 'Bangkok, Thailand',
    tags: ['cooking', 'livestream', 'tomyum', 'thai'],
    likes: 445,
    comments: 127,
    shares: 67,
    createdAt: '2024-03-15T17:30:00Z',
    postType: 'LIVE_STREAM',
    liveData: {
      viewers: 234,
      startTime: '2024-03-15T17:30:00Z',
      isLive: true
    }
  },
  {
    id: 'post-7',
    userId: 'user-1',
    type: 'REVIEW',
    content: 'Just stayed at this incredible eco-resort in Ubud. The infinity pool overlooking rice terraces was magical! Perfect for digital detox 🌾',
    images: ['https://images.unsplash.com/photo-1571003123894-1f0594d2b5d9?w=600&h=400&fit=crop'],
    location: 'Ubud, Bali, Indonesia',
    tags: ['review', 'ubud', 'ecoresort', 'digitaldetox'],
    likes: 312,
    comments: 45,
    shares: 28,
    createdAt: '2024-03-13T11:15:00Z',
    postType: 'REVIEW',
    rating: 4.8,
    targetType: 'venue',
    targetName: 'Kamandalu Ubud Resort'
  },
  {
    id: 'post-8',
    userId: 'user-2',
    type: 'STORY',
    content: 'Quick temple visit before my flight ✈️ Always grateful for these peaceful moments',
    images: ['https://images.unsplash.com/photo-1563492065037-3b3c1aa9be1d?w=600&h=400&fit=crop'],
    location: 'Wat Arun, Bangkok',
    tags: ['temple', 'gratitude', 'travel'],
    likes: 89,
    comments: 7,
    shares: 3,
    createdAt: '2024-03-15T08:45:00Z',
    postType: 'STORY',
    expiresAt: '2024-03-16T08:45:00Z'
  },
  {
    id: 'post-9',
    userId: 'user-3',
    type: 'MILESTONE',
    content: '🎉 One year of living as a digital nomad in Southeast Asia! From beaches to mountains, street food to temple visits, this journey has been incredible. Thank you to everyone who followed along! 🙏',
    images: [
      'https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=300&h=300&fit=crop',
      'https://images.unsplash.com/photo-1552465011-b4e21bf6e79a?w=300&h=300&fit=crop',
      'https://images.unsplash.com/photo-1528181304800-259b08848526?w=300&h=300&fit=crop'
    ],
    tags: ['milestone', 'digitalnomad', 'oneyear', 'grateful'],
    likes: 567,
    comments: 89,
    shares: 156,
    createdAt: '2024-03-12T16:00:00Z',
    postType: 'MILESTONE'
  },
  {
    id: 'post-10',
    userId: 'user-1',
    type: 'RECOMMENDATION',
    content: 'Hidden gem alert! 💎 This local coffee shop in Hoi An serves the most amazing Vietnamese iced coffee and they support local farmers. A must-visit!',
    images: ['https://images.unsplash.com/photo-1501339847302-ac426a4a7cbb?w=600&h=400&fit=crop'],
    location: 'Hoi An, Vietnam',
    tags: ['coffee', 'hoian', 'hidden gem', 'recommendation'],
    likes: 178,
    comments: 34,
    shares: 67,
    createdAt: '2024-03-14T13:20:00Z',
    postType: 'RECOMMENDATION'
  }
];

export const MOCK_PRODUCT_REVIEWS = [
  {
    id: 'review-1',
    productId: 'prod-1',
    userId: 'user-2',
    rating: 5,
    comment: 'Absolutely beautiful scarf! The silk quality is excellent and the traditional patterns are stunning. Fast shipping too!',
    images: ['https://images.unsplash.com/photo-1594633312681-425c7b97ccd1?w=200&h=200&fit=crop'],
    helpful: 12,
    createdAt: '2024-03-10T14:30:00Z'
  },
  {
    id: 'review-2',
    productId: 'prod-1',
    userId: 'user-3',
    rating: 4,
    comment: 'Very good quality silk scarf. The colors are vibrant and it feels luxurious. Took a bit longer to arrive but worth the wait.',
    helpful: 8,
    createdAt: '2024-03-08T16:20:00Z'
  },
  {
    id: 'review-3',
    productId: 'prod-2',
    userId: 'user-1',
    rating: 5,
    comment: 'Authentic Indonesian batik! Perfect fit and the traditional motifs are beautifully crafted. Highly recommend!',
    helpful: 15,
    createdAt: '2024-03-05T11:45:00Z'
  }
];

// =============================================================================
// Fallback Data Service
// =============================================================================

export class FallbackDataService {
  private static instance: FallbackDataService;

  static getInstance(): FallbackDataService {
    if (!FallbackDataService.instance) {
      FallbackDataService.instance = new FallbackDataService();
    }
    return FallbackDataService.instance;
  }

  /**
   * Check if an endpoint should use fallback data
   */
  shouldUseFallback(endpoint: string, error?: any): boolean {
    // Use fallback for 404 errors (endpoint not implemented)
    if (error?.status === 404) return true;

    // Use fallback for network errors
    if (error?.name === 'NetworkError' || error?.message?.includes('fetch')) return true;

    // Use fallback for connection errors
    if (error?.code === 'ECONNREFUSED') return true;

    return false;
  }

  /**
   * Get fallback data for a specific endpoint
   */
  getFallbackData(endpoint: string, params?: any): any {
    console.log(`🔄 Using fallback data for ${endpoint}`, params);

    switch (true) {
      // Products endpoints
      case endpoint.includes('/products') && endpoint.includes('/reviews'):
        const productId = endpoint.match(/\/products\/([^/]+)\/reviews/)?.[1];
        return MOCK_PRODUCT_REVIEWS.filter(r => r.productId === productId);

      case endpoint.includes('/products/featured'):
        return MOCK_PRODUCTS;

      case endpoint.includes('/products/search'):
        return {
          data: MOCK_PRODUCTS,
          total: MOCK_PRODUCTS.length,
          page: 1,
          limit: 20
        };

      case endpoint.includes('/products/'):
        const prodId = endpoint.match(/\/products\/([^/]+)$/)?.[1];
        return MOCK_PRODUCTS.find(p => p.id === prodId) || MOCK_PRODUCTS[0];

      // Chat endpoints
      case endpoint.includes('/chats') && endpoint.includes('/messages'):
        const chatId = endpoint.match(/\/chats\/([^/]+)\/messages/)?.[1];
        return {
          data: MOCK_MESSAGES.filter(m => m.chatId === chatId),
          total: MOCK_MESSAGES.filter(m => m.chatId === chatId).length,
          page: 1,
          limit: 20
        };

      case endpoint.includes('/chats'):
        return MOCK_CHATS;

      // Wallet endpoints
      case endpoint.includes('/wallet'):
        return MOCK_WALLET;

      case endpoint.includes('/transactions'):
        return MOCK_WALLET.transactions;

      // Notification endpoints
      case endpoint.includes('/notifications'):
        return MOCK_NOTIFICATIONS;

      // Social endpoints
      case endpoint.includes('/social/feed'):
        return MOCK_SOCIAL_FEED;

      case endpoint.includes('/friends'):
        return MOCK_USERS;

      // Shop endpoints
      case endpoint.includes('/shops') && endpoint.includes('/products'):
        return MOCK_PRODUCTS;

      case endpoint.includes('/shops/featured'):
        return MOCK_USERS.map(user => ({
          id: `shop-${user.id}`,
          name: `${user.name}'s Store`,
          description: `Authentic products from ${user.country}`,
          country: user.country,
          rating: 4.5 + Math.random() * 0.5,
          productCount: Math.floor(Math.random() * 100) + 10,
          verified: true
        }));

      // Default fallback
      default:
        console.warn(`No fallback data available for ${endpoint}`);
        return [];
    }
  }

  /**
   * Simulate API response with fallback data
   */
  createFallbackResponse(endpoint: string, params?: any): any {
    const data = this.getFallbackData(endpoint, params);

    return {
      data,
      status: 'success',
      message: 'Fallback data (development mode)',
      timestamp: new Date().toISOString(),
      meta: {
        fallback: true,
        endpoint,
        params
      }
    };
  }

  /**
   * Add realistic delay to simulate network requests
   */
  async withDelay<T>(data: T, delay: number = 300): Promise<T> {
    await new Promise(resolve => setTimeout(resolve, delay));
    return data;
  }
}

// Export singleton instance
export const fallbackDataService = FallbackDataService.getInstance();

/**
 * Helper function to handle API calls with fallback
 */
export async function withFallback<T>(
  apiCall: () => Promise<T>,
  endpoint: string,
  params?: any
): Promise<T> {
  try {
    return await apiCall();
  } catch (error) {
    if (fallbackDataService.shouldUseFallback(endpoint, error)) {
      console.warn(`API call failed for ${endpoint}, using fallback data:`, error);
      const fallbackResponse = fallbackDataService.createFallbackResponse(endpoint, params);
      return fallbackDataService.withDelay(fallbackResponse as T);
    }
    throw error;
  }
}