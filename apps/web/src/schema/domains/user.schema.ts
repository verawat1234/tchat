/**
 * User & Authentication Domain Schema
 *
 * Handles user profiles, authentication, preferences, KYC, and friend relationships
 * Designed for Southeast Asian markets with multi-locale and KYC tier support
 */

import { UUID, Timestamp, Locale, CountryCode } from '../schema';

// =============================================================================
// USER CORE
// =============================================================================

export interface User {
  id: UUID;
  phone?: string;
  email?: string;
  name: string;
  avatar?: string;
  country: CountryCode;
  locale: Locale;
  kycTier: 0 | 1 | 2 | 3; // Aligned with backend KYCTier (0=unverified, 1=basic, 2=standard, 3=premium)
  verificationTier: 0 | 1 | 2 | 3 | 4; // Added to match backend VerificationTier (0=none, 1=phone, 2=email, 3=kyc, 4=full)
  status: UserStatus;
  lastSeen?: Timestamp;
  phoneVerified: boolean; // Added to match backend phone_verified
  phoneVerifiedAt?: Timestamp; // Added to match backend phone_verified_at
  emailVerified: boolean; // Added to match backend email_verified
  emailVerifiedAt?: Timestamp; // Added to match backend email_verified_at
  isVerified: boolean; // Computed field for backward compatibility
  settings: UserSettings;
  profile: UserProfile;
  preferences: UserPreferences;
  createdAt: Timestamp;
  updatedAt: Timestamp;
}

export type UserStatus = 'active' | 'suspended' | 'deleted'; // Aligned with backend UserStatus

// =============================================================================
// USER PROFILE
// =============================================================================

export interface UserProfile {
  displayName?: string;
  bio?: string;
  birthday?: string;
  gender?: 'male' | 'female' | 'other' | 'prefer_not_to_say';
  location?: string;
  website?: string;
  occupation?: string;
  interests: string[];
  languages: Locale[];
  timezone: string;
}

// =============================================================================
// USER SETTINGS
// =============================================================================

export interface UserSettings {
  privacy: PrivacySettings;
  notifications: NotificationSettings;
  appearance: AppearanceSettings;
  security: SecuritySettings;
}

export interface PrivacySettings {
  profileVisibility: 'public' | 'friends' | 'private';
  phoneVisibility: 'public' | 'friends' | 'private';
  lastSeenVisibility: 'everyone' | 'friends' | 'nobody';
  readReceiptsEnabled: boolean;
  onlineStatusVisible: boolean;
}

export interface NotificationSettings {
  pushEnabled: boolean;
  emailEnabled: boolean;
  messageNotifications: boolean;
  postNotifications: boolean;
  eventNotifications: boolean;
  paymentNotifications: boolean;
  mutedChats: UUID[];
  mutedUsers: UUID[];
}

export interface AppearanceSettings {
  theme: 'light' | 'dark' | 'auto';
  language: Locale;
  fontSize: 'small' | 'medium' | 'large';
  chatWallpaper?: string;
}

export interface SecuritySettings {
  twoFactorEnabled: boolean;
  biometricEnabled: boolean;
  autoLockTimeout: number; // minutes
  trustedDevices: TrustedDevice[];
}

export interface TrustedDevice {
  id: UUID;
  deviceName: string;
  deviceType: 'mobile' | 'desktop' | 'tablet';
  lastUsed: Timestamp;
  ipAddress?: string;
  location?: string;
}

// =============================================================================
// USER PREFERENCES
// =============================================================================

export interface UserPreferences {
  defaultCurrency: Currency;
  defaultPaymentMethod?: UUID;
  eventCategories: string[];
  productCategories: string[];
  contentLanguages: Locale[];
  contentFilters: string[];
}

// =============================================================================
// KYC (Know Your Customer)
// =============================================================================

export interface KYCInfo {
  tier: 0 | 1 | 2 | 3; // Aligned with backend KYCTier (0=unverified, 1=basic, 2=standard, 3=premium)
  status: 'pending' | 'approved' | 'rejected' | 'incomplete';
  documents: KYCDocument[];
  verifiedAt?: Timestamp;
  expiresAt?: Timestamp;
  dailyLimit: number;
  monthlyLimit: number;
  usedThisMonth: number;
}

export interface KYCDocument {
  id: UUID;
  type: KYCDocumentType;
  fileUrl: string;
  status: 'pending' | 'approved' | 'rejected';
  uploadedAt: Timestamp;
  verifiedAt?: Timestamp;
  rejectionReason?: string;
  metadata?: Record<string, any>;
}

export type KYCDocumentType =
  | 'id_card'
  | 'passport'
  | 'driving_license'
  | 'utility_bill'
  | 'bank_statement';

// =============================================================================
// FRIEND RELATIONSHIPS
// =============================================================================

export interface Friend {
  id: UUID;
  userId: UUID;
  friendId: UUID;
  status: FriendStatus;
  mutualFriends: number;
  commonInterests: string[];
  createdAt: Timestamp;
  acceptedAt?: Timestamp;
}

export type FriendStatus = 'pending' | 'accepted' | 'blocked';

// =============================================================================
// AUTHENTICATION
// =============================================================================

export interface AuthSession {
  id: UUID;
  userId: UUID;
  deviceId: string;
  deviceInfo: DeviceInfo;
  accessTokenHash: string;
  refreshTokenHash: string;
  expiresAt: Timestamp;
  isActive: boolean;
  ipAddress?: string;
  location?: SessionLocation;
  userAgent?: string;
  createdAt: Timestamp;
  lastUsed: Timestamp;
}

export interface DeviceInfo {
  type: 'mobile' | 'desktop' | 'tablet' | 'web';
  os: string;
  browser?: string;
  appVersion?: string;
  pushToken?: string;
  notificationEnabled: boolean;
}

export interface SessionLocation {
  country?: CountryCode;
  city?: string;
  region?: string;
  latitude?: number;
  longitude?: number;
  timezone?: string;
  ipLookupProvider?: string;
}

export interface OTPVerification {
  id: UUID;
  userId?: UUID;
  phone?: string;
  email?: string;
  code: string;
  type: 'registration' | 'login' | 'password_reset' | 'phone_verification';
  attempts: number;
  maxAttempts: number;
  isUsed: boolean;
  expiresAt: Timestamp;
  createdAt: Timestamp;
  usedAt?: Timestamp;
}

// =============================================================================
// USER ACTIVITY & ANALYTICS
// =============================================================================

export interface UserActivity {
  id: UUID;
  userId: UUID;
  action: UserActionType;
  targetType?: string;
  targetId?: UUID;
  metadata?: Record<string, any>;
  ipAddress?: string;
  userAgent?: string;
  location?: string;
  createdAt: Timestamp;
}

export type UserActionType =
  | 'login'
  | 'logout'
  | 'profile_update'
  | 'settings_change'
  | 'kyc_submission'
  | 'password_change'
  | 'device_added'
  | 'friend_request'
  | 'friend_accept'
  | 'privacy_change';

export interface UserStats {
  userId: UUID;
  totalMessages: number;
  totalPosts: number;
  totalFriends: number;
  totalOrders: number;
  totalSpent: number;
  eventsAttended: number;
  videosWatched: number;
  achievementsEarned: number;
  joinedAt: Timestamp;
  lastActive: Timestamp;
  updatedAt: Timestamp;
}

// =============================================================================
// BUSINESS LOGIC HELPERS
// =============================================================================

/**
 * KYC Tier Limits (Thailand specific) - Aligned with backend GetMaxTransactionLimit()
 */
export const KYC_LIMITS = {
  0: { daily: 100, monthly: 100 },        // THB - Unverified (backend: 100.0)
  1: { daily: 1000, monthly: 1000 },      // THB - Basic (backend: 1000.0)
  2: { daily: 10000, monthly: 10000 },    // THB - Standard (backend: 10000.0)
  3: { daily: 100000, monthly: 100000 }   // THB - Premium (backend: 100000.0)
} as const;

/**
 * KYC Tier Names - Aligned with backend KYCTier constants
 */
export const KYC_TIER_NAMES = {
  0: 'Unverified',
  1: 'Basic',
  2: 'Standard',
  3: 'Premium'
} as const;

/**
 * Verification Tier Names - Aligned with backend VerificationTier constants
 */
export const VERIFICATION_TIER_NAMES = {
  0: 'None',
  1: 'Phone',
  2: 'Email',
  3: 'KYC',
  4: 'Full'
} as const;

/**
 * Supported Southeast Asian Countries
 */
export const SEA_COUNTRIES: CountryCode[] = ['TH', 'ID', 'MY', 'VN', 'SG', 'PH'];

/**
 * Supported Locales with Country Mapping
 */
export const LOCALE_COUNTRY_MAP: Record<Locale, CountryCode> = {
  'th-TH': 'TH',
  'id-ID': 'ID',
  'ms-MY': 'MY',
  'vi-VN': 'VN',
  'en-US': 'SG' // Default for Singapore and Philippines
};

/**
 * Default user preferences by country
 */
export const DEFAULT_PREFERENCES: Record<CountryCode, Partial<UserPreferences>> = {
  TH: {
    defaultCurrency: 'THB',
    contentLanguages: ['th-TH', 'en-US'],
    eventCategories: ['music', 'food', 'cultural', 'temple']
  },
  ID: {
    defaultCurrency: 'IDR',
    contentLanguages: ['id-ID', 'en-US'],
    eventCategories: ['music', 'food', 'cultural']
  },
  MY: {
    defaultCurrency: 'MYR',
    contentLanguages: ['ms-MY', 'en-US'],
    eventCategories: ['music', 'food', 'cultural']
  },
  VN: {
    defaultCurrency: 'VND',
    contentLanguages: ['vi-VN', 'en-US'],
    eventCategories: ['music', 'food', 'cultural']
  },
  SG: {
    defaultCurrency: 'SGD',
    contentLanguages: ['en-US'],
    eventCategories: ['music', 'food', 'business']
  },
  PH: {
    defaultCurrency: 'PHP',
    contentLanguages: ['en-US'],
    eventCategories: ['music', 'food', 'cultural']
  }
};
