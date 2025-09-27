// Social Post Types - Comprehensive post system following MessageData pattern
/**
 * Core social post entity with type discrimination and engagement features
 * Supports rich content types optimized for social media interactions
 */

import { ReactNode } from 'react';

// Core Post Entity
export interface PostData {
  readonly id: string;
  readonly authorId: string;
  readonly authorName: string;
  readonly authorAvatar?: string;
  readonly timestamp: Date;
  readonly type: PostType;
  readonly content: PostContent;
  readonly privacy: PostPrivacy;
  readonly audience?: PostAudience;
  readonly engagement?: PostEngagement;
  readonly metadata?: PostMetadata;
  readonly renderState?: PostRenderState;
}

// Post Type Enum - All supported social post types
export enum PostType {
  // Core Content Types (8)
  TEXT = 'text',                    // Simple status updates
  IMAGE = 'image',                  // Single/multiple photos
  VIDEO = 'video',                  // Video posts
  AUDIO = 'audio',                  // Voice notes, music
  LINK_SHARE = 'link_share',        // Shared articles/websites
  POST_MESSAGE = 'post_message',    // Message posted to someone's timeline/wall
  REVIEW = 'review',                // Reviews of places, products, services
  ALBUM = 'album',                  // Photo collections

  // Rich Media Types (6)
  STORY = 'story',                  // Ephemeral 24h content
  REEL = 'reel',                    // Short-form vertical video
  LIVE_STREAM = 'live_stream',      // Live video broadcasts
  PLAYLIST = 'playlist',            // Music/video collections
  MOOD_BOARD = 'mood_board',        // Visual inspiration
  TUTORIAL = 'tutorial',            // How-to content

  // Interactive Content (6)
  POLL = 'poll',                    // Voting posts
  QUIZ = 'quiz',                    // Trivia, personality tests
  SURVEY = 'survey',                // Feedback collection
  Q_AND_A = 'q_and_a',            // Ask me anything
  CHALLENGE = 'challenge',          // Viral challenges/trends
  PETITION = 'petition',            // Social causes

  // Social & Location (8)
  CHECK_IN = 'check_in',            // Location-based posts
  TRAVEL_LOG = 'travel_log',        // Trip updates/itinerary
  LIFE_EVENT = 'life_event',        // Major life moments
  MILESTONE = 'milestone',          // Personal achievements
  MEMORY = 'memory',                // Throwback/flashback posts
  ANNIVERSARY = 'anniversary',      // Yearly memories
  RECOMMENDATION = 'recommendation', // Place/product suggestions
  GROUP_ACTIVITY = 'group_activity', // Group-specific content

  // Commercial & Business (6)
  PRODUCT_SHOWCASE = 'product_showcase', // Selling items
  SERVICE_PROMOTION = 'service_promotion', // Business services
  EVENT_PROMOTION = 'event_promotion',   // Events/meetups
  JOB_POSTING = 'job_posting',          // Hiring/career opportunities
  FUNDRAISER = 'fundraiser',            // Charity/personal causes
  COLLABORATION = 'collaboration',      // Creative projects

  // Specialized Content (8)
  RECIPE = 'recipe',                    // Cooking/food content
  WORKOUT = 'workout',                  // Fitness routines
  BOOK_REVIEW = 'book_review',          // Reading updates
  MOOD_UPDATE = 'mood_update',          // Emotional status
  ACHIEVEMENT = 'achievement',          // Gaming/app achievements
  QUOTE = 'quote',                      // Inspirational quotes
  MUSIC = 'music',                      // Music sharing/streaming
  VENUE = 'venue',                      // Venue information/reviews
}

// Post Privacy Settings
export enum PostPrivacy {
  PUBLIC = 'public',                   // Visible to everyone
  FRIENDS = 'friends',                 // Friends only
  CLOSE_FRIENDS = 'close_friends',     // Close friends list
  FOLLOWERS = 'followers',             // Followers only
  MUTUAL_FRIENDS = 'mutual_friends',   // Mutual connections
  CUSTOM = 'custom',                   // Custom audience
  UNLISTED = 'unlisted',              // Hidden from feeds but shareable
  PRIVATE = 'private',                 // Only author can see
}

// Post Audience Targeting
export interface PostAudience {
  readonly includedUsers?: string[];      // Specific users who can see
  readonly excludedUsers?: string[];      // Users who cannot see
  readonly includedGroups?: string[];     // Specific groups/circles
  readonly excludedGroups?: string[];     // Excluded groups
  readonly locationRestriction?: LocationRestriction;
  readonly ageRange?: AgeRange;
  readonly interests?: string[];          // Interest-based targeting
}

// Post Engagement Metrics
export interface PostEngagement {
  readonly reactions: PostReaction[];
  readonly comments: PostComment[];
  readonly shares: PostShare[];
  readonly saves: string[];               // User IDs who saved
  readonly views: number;
  readonly reach: number;                 // Unique users reached
  readonly impressions: number;           // Total views
  readonly clickThroughs: number;         // Link/action clicks
  readonly engagementRate: number;        // Calculated engagement %
}

// Post Reactions (enhanced reaction system)
export interface PostReaction {
  readonly type: ReactionType;
  readonly userId: string;
  readonly timestamp: Date;
  readonly userName?: string;
}

export enum ReactionType {
  LIKE = 'like',
  LOVE = 'love',
  HAHA = 'haha',
  WOW = 'wow',
  SAD = 'sad',
  ANGRY = 'angry',
  CARE = 'care',
  FIRE = 'fire',                         // Trending/hot content
  CLAP = 'clap',                         // Appreciation
  CELEBRATE = 'celebrate',               // Achievements
}

// Post Comments
export interface PostComment {
  readonly id: string;
  readonly userId: string;
  readonly userName: string;
  readonly userAvatar?: string;
  readonly content: string;
  readonly timestamp: Date;
  readonly replies?: PostComment[];       // Nested replies
  readonly reactions?: PostReaction[];
  readonly isEdited: boolean;
  readonly isDeleted: boolean;
  readonly mentionedUsers?: string[];
}

// Post Shares
export interface PostShare {
  readonly id: string;
  readonly userId: string;
  readonly userName: string;
  readonly timestamp: Date;
  readonly shareType: ShareType;
  readonly addedComment?: string;
  readonly sharedToGroups?: string[];
}

export enum ShareType {
  DIRECT_SHARE = 'direct_share',         // Simple reshare
  QUOTE_SHARE = 'quote_share',          // Share with comment
  STORY_SHARE = 'story_share',          // Share to story
  MESSAGE_SHARE = 'message_share',      // Share via DM
  EXTERNAL_SHARE = 'external_share',    // Share outside platform
}

// Post Metadata
export interface PostMetadata {
  readonly editedAt?: Date;
  readonly originalPostId?: string;      // If this is a shared post
  readonly isPromoted?: boolean;         // Sponsored/boosted content
  readonly promotionDetails?: PromotionDetails;
  readonly tags?: string[];              // Hashtags
  readonly mentionedUsers?: string[];    // @mentions
  readonly location?: PostLocation;
  readonly mood?: MoodType;
  readonly feeling?: FeelingType;
  readonly activity?: ActivityType;
  readonly contentWarning?: ContentWarning;
  readonly language?: string;
  readonly isArchived: boolean;
  readonly archivedAt?: Date;
  readonly expiresAt?: Date;             // For stories/temporary content
  readonly allowComments: boolean;
  readonly allowShares: boolean;
  readonly allowSaves: boolean;
  readonly isPinned: boolean;
  readonly isSticky: boolean;            // Stays at top of profile/group
}

// Content Warning System
export interface ContentWarning {
  readonly type: ContentWarningType;
  readonly reason: string;
  readonly reviewStatus: ReviewStatus;
  readonly reviewedBy?: string;
  readonly reviewedAt?: Date;
}

export enum ContentWarningType {
  NONE = 'none',
  SENSITIVE_CONTENT = 'sensitive_content',
  GRAPHIC_VIOLENCE = 'graphic_violence',
  ADULT_CONTENT = 'adult_content',
  DISTURBING_CONTENT = 'disturbing_content',
  SPOILER = 'spoiler',
  FLASHING_LIGHTS = 'flashing_lights',
  POLITICAL_CONTENT = 'political_content',
}

// Discriminated union for post content based on type
export type PostContent =
  | string                              // Simple text posts
  | ImagePostContent
  | VideoPostContent
  | AudioPostContent
  | LinkShareContent
  | PostMessageContent
  | ReviewContent
  | AlbumContent
  | StoryContent
  | ReelContent
  | LiveStreamContent
  | PlaylistContent
  | MoodBoardContent
  | TutorialContent
  | PollContent
  | QuizContent
  | SurveyContent
  | QAndAContent
  | ChallengeContent
  | PetitionContent
  | CheckInContent
  | TravelLogContent
  | LifeEventContent
  | MilestoneContent
  | MemoryContent
  | AnniversaryContent
  | RecommendationContent
  | GroupActivityContent
  | ProductShowcaseContent
  | ServicePromotionContent
  | EventPromotionContent
  | JobPostingContent
  | FundraiserContent
  | CollaborationContent
  | RecipeContent
  | WorkoutContent
  | BookReviewContent
  | MoodUpdateContent
  | AchievementContent
  | QuoteContent
  | MusicContent
  | VenueContent;

// Content Type Interfaces
export interface ImagePostContent {
  readonly images: MediaFile[];
  readonly caption?: string;
  readonly altText?: string[];           // Accessibility descriptions
  readonly filters?: ImageFilter[];
  readonly location?: PostLocation;
}

export interface VideoPostContent {
  readonly videoUrl: string;
  readonly thumbnailUrl?: string;
  readonly duration: number;
  readonly quality: VideoQuality;
  readonly caption?: string;
  readonly subtitles?: SubtitleTrack[];
  readonly chapters?: VideoChapter[];
  readonly isLooping: boolean;
  readonly autoPlay: boolean;
}

export interface AudioPostContent {
  readonly audioUrl: string;
  readonly duration: number;
  readonly title?: string;
  readonly artist?: string;
  readonly album?: string;
  readonly genre?: string;
  readonly waveformData?: number[];      // For audio visualization
  readonly isMusic: boolean;
  readonly transcript?: string;          // For voice notes
}

export interface LinkShareContent {
  readonly url: string;
  readonly title?: string;
  readonly description?: string;
  readonly imageUrl?: string;
  readonly siteName?: string;
  readonly domain: string;
  readonly comment?: string;             // User's comment about the link
  readonly previewGenerated: boolean;
}

export interface PostMessageContent {
  readonly recipientId: string;          // Who this message is posted to
  readonly recipientName: string;        // Recipient's display name
  readonly message: string;              // The message content
  readonly messageType: PostMessageType;
  readonly isReply?: boolean;            // Reply to another wall post
  readonly replyToId?: string;           // Original post ID if replying
  readonly visibility: PostMessageVisibility;
  readonly allowReplies: boolean;
  readonly taggedFriends?: string[];     // Friends tagged in the message
}

export interface ReviewContent {
  readonly subjectType: ReviewSubjectType;
  readonly subjectId: string;            // ID of what's being reviewed
  readonly subjectName: string;          // Name of business/product/service
  readonly rating: number;               // 1-5 star rating
  readonly reviewTitle?: string;         // Brief title for the review
  readonly reviewText: string;           // Detailed review text
  readonly photos?: string[];            // Photos related to the review
  readonly visitDate?: Date;             // When the experience happened
  readonly isVerifiedPurchase?: boolean; // For product reviews
  readonly pros?: string[];              // List of positive aspects
  readonly cons?: string[];              // List of negative aspects
  readonly wouldRecommend?: boolean;     // Recommendation status
  readonly helpfulCount: number;         // How many found it helpful
  readonly location?: PostLocation;      // For location-based reviews
  readonly priceRange?: PriceRange;      // Price assessment
  readonly serviceRatings?: ServiceRatings; // Detailed service ratings
}

export interface AlbumContent {
  readonly title: string;
  readonly description?: string;
  readonly images: MediaFile[];
  readonly coverImage?: number;          // Index of cover image
  readonly allowDownload: boolean;
  readonly isCollaborative: boolean;     // Allow others to add photos
  readonly contributors?: string[];      // User IDs who can contribute
}

export interface StoryContent {
  readonly mediaUrl: string;
  readonly mediaType: 'image' | 'video';
  readonly duration?: number;            // For videos
  readonly stickers?: StorySticker[];
  readonly filters?: StoryFilter[];
  readonly music?: StoryMusic;
  readonly polls?: StoryPoll[];
  readonly questions?: StoryQuestion[];
  readonly links?: StoryLink[];
  readonly mentions?: StoryMention[];
  readonly viewersList?: string[];       // Who viewed the story
  readonly allowReplies: boolean;
  readonly allowSharing: boolean;
  readonly highlightCategory?: string;   // For story highlights
}

export interface ReelContent {
  readonly videoUrl: string;
  readonly thumbnailUrl?: string;
  readonly duration: number;
  readonly caption?: string;
  readonly music?: ReelMusic;
  readonly effects?: ReelEffect[];
  readonly hashtags?: string[];
  readonly trendingSound?: boolean;
  readonly originalAudio?: boolean;
  readonly isRemix?: boolean;           // Remix of another reel
  readonly remixOriginalId?: string;
}

export interface LiveStreamContent {
  readonly streamUrl: string;
  readonly title: string;
  readonly description?: string;
  readonly thumbnailUrl?: string;
  readonly startTime: Date;
  readonly endTime?: Date;
  readonly maxViewers?: number;
  readonly currentViewers: number;
  readonly totalViewers: number;
  readonly isRecorded: boolean;
  readonly recordingUrl?: string;
  readonly chatEnabled: boolean;
  readonly guestsEnabled: boolean;
  readonly guests?: StreamGuest[];
  readonly donations?: StreamDonation[];
}

export interface PlaylistContent {
  readonly title: string;
  readonly description?: string;
  readonly coverImage?: string;
  readonly tracks: PlaylistTrack[];
  readonly isPublic: boolean;
  readonly isCollaborative: boolean;
  readonly totalDuration: number;
  readonly genre?: string;
  readonly mood?: string;
}

export interface MoodBoardContent {
  readonly title: string;
  readonly description?: string;
  readonly images: MediaFile[];
  readonly colorPalette?: string[];
  readonly theme?: string;
  readonly inspiration?: string;
  readonly tags?: string[];
  readonly isPublicTemplate: boolean;
}

export interface TutorialContent {
  readonly title: string;
  readonly description?: string;
  readonly difficulty: DifficultyLevel;
  readonly duration?: number;           // Estimated completion time
  readonly steps: TutorialStep[];
  readonly materials?: string[];
  readonly tools?: string[];
  readonly category: TutorialCategory;
  readonly videoUrl?: string;
  readonly tags?: string[];
}

export interface PollContent {
  readonly question: string;
  readonly options: PollOption[];
  readonly allowMultipleChoices: boolean;
  readonly showResults: PollResultsVisibility;
  readonly endTime?: Date;
  readonly totalVotes: number;
  readonly voterIds: string[];          // Track who voted
}

export interface QuizContent {
  readonly title: string;
  readonly description?: string;
  readonly questions: QuizQuestion[];
  readonly timeLimit?: number;
  readonly passingScore?: number;
  readonly allowRetakes: boolean;
  readonly showAnswersAfter: 'never' | 'submission' | 'completion';
  readonly category?: string;
  readonly difficulty: DifficultyLevel;
}

export interface SurveyContent {
  readonly title: string;
  readonly description?: string;
  readonly questions: SurveyQuestion[];
  readonly allowAnonymous: boolean;
  readonly showProgress: boolean;
  readonly autoSave: boolean;
  readonly maxResponses?: number;
  readonly endDate?: Date;
}

export interface QAndAContent {
  readonly question: string;
  readonly description?: string;
  readonly category?: string;
  readonly answers: QAAnswer[];
  readonly isAnswered: boolean;
  readonly bestAnswerId?: string;
  readonly tags?: string[];
}

export interface ChallengeContent {
  readonly title: string;
  readonly description: string;
  readonly rules: string[];
  readonly hashtag: string;
  readonly duration?: number;           // In days
  readonly participants: number;
  readonly submissions: ChallengeSubmission[];
  readonly prizes?: string[];
  readonly deadline?: Date;
}

export interface PetitionContent {
  readonly title: string;
  readonly description: string;
  readonly targetAudience: string;      // Who the petition is directed to
  readonly goal: number;               // Target signatures
  readonly currentSignatures: number;
  readonly signers: PetitionSigner[];
  readonly category: PetitionCategory;
  readonly deadline?: Date;
  readonly updates?: PetitionUpdate[];
}

export interface CheckInContent {
  readonly locationId: string;
  readonly locationName: string;
  readonly locationAddress?: string;
  readonly coordinates?: Coordinates;
  readonly category: LocationCategory;
  readonly rating?: number;             // User's rating 1-5
  readonly review?: string;
  readonly photos?: string[];
  readonly companions?: string[];       // Tagged friends
  readonly activity?: CheckInActivity;
}

export interface TravelLogContent {
  readonly title: string;
  readonly destination: string;
  readonly startDate: Date;
  readonly endDate?: Date;
  readonly itinerary: TravelDay[];
  readonly budget?: number;
  readonly currency?: string;
  readonly travelType: TravelType;
  readonly companions?: string[];
  readonly highlights?: string[];
  readonly tips?: string[];
  readonly photos?: string[];
}

export interface LifeEventContent {
  readonly eventType: LifeEventType;
  readonly title: string;
  readonly description?: string;
  readonly date?: Date;                 // When the event happened
  readonly location?: string;
  readonly participants?: string[];     // Tagged people
  readonly privacy: LifeEventPrivacy;
  readonly milestone?: boolean;         // Show as major milestone
  readonly photos?: string[];
  readonly customType?: string;         // For custom life events
}

export interface MilestoneContent {
  readonly title: string;
  readonly description: string;
  readonly category: MilestoneCategory;
  readonly achievementDate: Date;
  readonly metrics?: MilestoneMetric[];
  readonly photos?: string[];
  readonly certificateUrl?: string;
  readonly shareWithEmployer?: boolean;
  readonly tags?: string[];
}

export interface MemoryContent {
  readonly title: string;
  readonly originalDate: Date;
  readonly yearsAgo: number;
  readonly originalPostId?: string;
  readonly memories: MemoryItem[];
  readonly location?: string;
  readonly people?: string[];
  readonly customMessage?: string;
}

export interface AnniversaryContent {
  readonly title: string;
  readonly anniversaryType: AnniversaryType;
  readonly originalDate: Date;
  readonly yearsCount: number;
  readonly description?: string;
  readonly milestonePhotos?: string[];
  readonly timeline?: AnniversaryMilestone[];
  readonly celebrationPlan?: string;
}

export interface RecommendationContent {
  readonly title: string;
  readonly itemType: RecommendationType;
  readonly itemId: string;
  readonly itemName: string;
  readonly rating: number;
  readonly reason: string;
  readonly pros?: string[];
  readonly cons?: string[];
  readonly photos?: string[];
  readonly price?: number;
  readonly currency?: string;
  readonly availableAt?: string[];      // Where to find/buy
  readonly tags?: string[];
}

export interface GroupActivityContent {
  readonly groupId: string;
  readonly groupName: string;
  readonly activityType: GroupActivityType;
  readonly title: string;
  readonly description: string;
  readonly participants: string[];
  readonly maxParticipants?: number;
  readonly dateTime?: Date;
  readonly location?: string;
  readonly requirements?: string[];
  readonly photos?: string[];
}

export interface ProductShowcaseContent {
  readonly productId: string;
  readonly name: string;
  readonly description: string;
  readonly price?: number;
  readonly currency?: string;
  readonly images: string[];
  readonly category: ProductCategory;
  readonly condition: ProductCondition;
  readonly availability: ProductAvailability;
  readonly location?: string;
  readonly shippingInfo?: ShippingInfo;
  readonly tags?: string[];
  readonly contactInfo?: ContactInfo;
  readonly isNegotiable: boolean;
}

export interface ServicePromotionContent {
  readonly serviceId: string;
  readonly name: string;
  readonly description: string;
  readonly category: ServiceCategory;
  readonly pricing: ServicePricing;
  readonly availability: ServiceAvailability;
  readonly location?: string;
  readonly serviceArea?: string[];
  readonly portfolio?: string[];        // Portfolio images
  readonly testimonials?: Testimonial[];
  readonly contactInfo: ContactInfo;
  readonly tags?: string[];
}

export interface EventPromotionContent {
  readonly eventId: string;
  readonly title: string;
  readonly description: string;
  readonly eventType: EventType;
  readonly startDate: Date;
  readonly endDate?: Date;
  readonly location: EventLocation;
  readonly ticketing: EventTicketing;
  readonly organizer: EventOrganizer;
  readonly maxAttendees?: number;
  readonly currentAttendees: number;
  readonly agenda?: EventAgendaItem[];
  readonly sponsors?: EventSponsor[];
  readonly photos?: string[];
}

export interface JobPostingContent {
  readonly jobId: string;
  readonly title: string;
  readonly company: string;
  readonly description: string;
  readonly requirements: string[];
  readonly qualifications: string[];
  readonly jobType: JobType;
  readonly workLocation: WorkLocation;
  readonly salaryRange?: SalaryRange;
  readonly benefits?: string[];
  readonly applicationDeadline?: Date;
  readonly contactEmail?: string;
  readonly applyUrl?: string;
  readonly tags?: string[];
}

export interface FundraiserContent {
  readonly fundraiserId: string;
  readonly title: string;
  readonly description: string;
  readonly cause: FundraiserCause;
  readonly goalAmount: number;
  readonly currentAmount: number;
  readonly currency: string;
  readonly deadline?: Date;
  readonly beneficiary?: string;
  readonly organizer: FundraiserOrganizer;
  readonly updates?: FundraiserUpdate[];
  readonly donors?: FundraiserDonor[];
  readonly photos?: string[];
  readonly verificationStatus: VerificationStatus;
}

export interface CollaborationContent {
  readonly title: string;
  readonly description: string;
  readonly collaborationType: CollaborationType;
  readonly skillsNeeded: string[];
  readonly timeline?: string;
  readonly commitment: CommitmentLevel;
  readonly compensation?: string;
  readonly portfolio?: string[];
  readonly collaborators?: Collaborator[];
  readonly applicationProcess?: string;
  readonly contactInfo?: ContactInfo;
  readonly tags?: string[];
}

export interface RecipeContent {
  readonly title: string;
  readonly description?: string;
  readonly images: string[];
  readonly servings: number;
  readonly prepTime: number;           // minutes
  readonly cookTime: number;           // minutes
  readonly difficulty: DifficultyLevel;
  readonly cuisine?: string;
  readonly dietaryTags?: DietaryTag[];
  readonly ingredients: RecipeIngredient[];
  readonly instructions: RecipeStep[];
  readonly nutrition?: NutritionInfo;
  readonly tips?: string[];
  readonly videoUrl?: string;
}

export interface WorkoutContent {
  readonly title: string;
  readonly description?: string;
  readonly workoutType: WorkoutType;
  readonly difficulty: DifficultyLevel;
  readonly duration: number;           // minutes
  readonly equipment?: string[];
  readonly exercises: WorkoutExercise[];
  readonly targetMuscles?: string[];
  readonly caloriesBurned?: number;
  readonly restPeriods?: number;       // seconds between sets
  readonly videoUrl?: string;
  readonly photos?: string[];
}

export interface BookReviewContent {
  readonly bookId?: string;
  readonly title: string;
  readonly author: string;
  readonly isbn?: string;
  readonly genre: string;
  readonly rating: number;             // 1-5 stars
  readonly reviewText: string;
  readonly readingStatus: ReadingStatus;
  readonly startDate?: Date;
  readonly finishDate?: Date;
  readonly pageCount?: number;
  readonly quotes?: BookQuote[];
  readonly tags?: string[];
  readonly wouldRecommend: boolean;
  readonly coverImage?: string;
}

export interface MoodUpdateContent {
  readonly mood: MoodType;
  readonly intensity: number;          // 1-10 scale
  readonly description?: string;
  readonly triggers?: string[];        // What caused the mood
  readonly location?: string;
  readonly activity?: string;
  readonly companions?: string[];
  readonly coping?: string[];          // Coping mechanisms used
  readonly isAnonymous?: boolean;
  readonly seekingSupport?: boolean;
}

export interface AchievementContent {
  readonly title: string;
  readonly description: string;
  readonly category: AchievementCategory;
  readonly platform?: string;          // Game/app platform
  readonly dateEarned: Date;
  readonly difficulty: DifficultyLevel;
  readonly rarity?: AchievementRarity;
  readonly points?: number;
  readonly badgeUrl?: string;
  readonly screenshot?: string;
  readonly shareStats?: boolean;
}

export interface QuoteContent {
  readonly text: string;
  readonly author?: string;
  readonly source?: string;            // Book, movie, etc.
  readonly category: QuoteCategory;
  readonly language?: string;
  readonly translation?: string;
  readonly personalNote?: string;      // User's thoughts on the quote
  readonly backgroundImage?: string;
  readonly textStyle?: QuoteTextStyle;
  readonly tags?: string[];
}

export interface MusicContent {
  readonly trackId?: string;
  readonly title: string;
  readonly artist: string;
  readonly album?: string;
  readonly genre?: string;
  readonly duration?: number;
  readonly spotifyUrl?: string;
  readonly appleUrl?: string;
  readonly youtubeUrl?: string;
  readonly lyrics?: string;
  readonly personalNote?: string;
  readonly mood?: string;
  readonly rating?: number;           // Personal rating
  readonly coverArt?: string;
}

export interface VenueContent {
  readonly venueId: string;
  readonly name: string;
  readonly address: string;
  readonly coordinates?: Coordinates;
  readonly category: LocationCategory;
  readonly rating?: number;
  readonly priceRange?: PriceRange;
  readonly hours?: BusinessHours[];
  readonly contact?: ContactInfo;
  readonly images?: string[];
  readonly amenities?: string[];
  readonly accessibility?: string[];
  readonly reviews?: VenueReview[];
}

// Supporting Types and Enums
export interface MediaFile {
  readonly url: string;
  readonly thumbnailUrl?: string;
  readonly width?: number;
  readonly height?: number;
  readonly fileSize?: number;
  readonly mimeType?: string;
  readonly altText?: string;
}

export interface PostLocation {
  readonly name: string;
  readonly address?: string;
  readonly coordinates?: Coordinates;
  readonly city?: string;
  readonly country?: string;
  readonly category?: LocationCategory;
}

export interface Coordinates {
  readonly latitude: number;
  readonly longitude: number;
}

export enum LocationCategory {
  RESTAURANT = 'restaurant',
  HOTEL = 'hotel',
  ATTRACTION = 'attraction',
  SHOPPING = 'shopping',
  ENTERTAINMENT = 'entertainment',
  OUTDOORS = 'outdoors',
  TRANSPORTATION = 'transportation',
  HOME = 'home',
  WORK = 'work',
  EDUCATION = 'education',
  HEALTHCARE = 'healthcare',
  GOVERNMENT = 'government',
  RELIGIOUS = 'religious',
  SPORTS = 'sports',
  OTHER = 'other'
}

export enum MoodType {
  HAPPY = 'happy',
  EXCITED = 'excited',
  GRATEFUL = 'grateful',
  LOVED = 'loved',
  BLESSED = 'blessed',
  RELAXED = 'relaxed',
  CONTENT = 'content',
  MOTIVATED = 'motivated',
  PROUD = 'proud',
  ACCOMPLISHED = 'accomplished',
  TIRED = 'tired',
  STRESSED = 'stressed',
  SAD = 'sad',
  ANXIOUS = 'anxious',
  CONFUSED = 'confused',
  FRUSTRATED = 'frustrated',
  ANGRY = 'angry',
  LONELY = 'lonely',
  NOSTALGIC = 'nostalgic',
  CONTEMPLATIVE = 'contemplative'
}

export enum FeelingType {
  AMAZING = 'amazing',
  FANTASTIC = 'fantastic',
  GOOD = 'good',
  OKAY = 'okay',
  MEH = 'meh',
  NOT_GREAT = 'not_great',
  TERRIBLE = 'terrible'
}

export enum ActivityType {
  EATING = 'eating',
  DRINKING = 'drinking',
  TRAVELING = 'traveling',
  EXERCISING = 'exercising',
  WORKING = 'working',
  STUDYING = 'studying',
  READING = 'reading',
  WATCHING = 'watching',
  LISTENING = 'listening',
  PLAYING = 'playing',
  COOKING = 'cooking',
  SHOPPING = 'shopping',
  CELEBRATING = 'celebrating',
  RELAXING = 'relaxing',
  SLEEPING = 'sleeping'
}

// Supporting types for PostMessage and Review content
export enum PostMessageType {
  WALL_POST = 'wall_post',              // Regular post on someone's wall
  BIRTHDAY_WISH = 'birthday_wish',      // Birthday message
  CONGRATULATIONS = 'congratulations',  // Congrats message
  APPRECIATION = 'appreciation',        // Thank you/appreciation
  SUPPORT = 'support',                  // Support/encouragement
  SHOUTOUT = 'shoutout',               // Public shoutout/recognition
  MEMORY = 'memory',                   // Shared memory
  JOKE = 'joke',                       // Funny message
  QUESTION = 'question',               // Ask a question
  GENERAL = 'general'                  // General message
}

export enum PostMessageVisibility {
  PUBLIC = 'public',                   // Visible to everyone
  FRIENDS = 'friends',                 // Friends of both parties
  MUTUAL_FRIENDS = 'mutual_friends',   // Only mutual friends
  RECIPIENT_ONLY = 'recipient_only'    // Only the recipient can see
}

export enum ReviewSubjectType {
  BUSINESS = 'business',               // Local business/restaurant
  PRODUCT = 'product',                 // Physical/digital product
  SERVICE = 'service',                 // Professional service
  APP = 'app',                        // Mobile/web application
  MOVIE = 'movie',                    // Film/movie
  BOOK = 'book',                      // Book/publication
  MUSIC = 'music',                    // Album/song/artist
  TV_SHOW = 'tv_show',                // TV series/show
  PODCAST = 'podcast',                // Podcast series
  EVENT = 'event',                    // Event/conference
  ACCOMMODATION = 'accommodation',     // Hotel/Airbnb
  TRANSPORTATION = 'transportation',  // Uber/taxi/airline
  HEALTHCARE = 'healthcare',          // Doctor/hospital
  EDUCATION = 'education',            // School/course
  OTHER = 'other'                     // Other subjects
}

export enum PriceRange {
  VERY_LOW = 'very_low',              // $
  LOW = 'low',                        // $$
  MEDIUM = 'medium',                  // $$$
  HIGH = 'high',                      // $$$$
  VERY_HIGH = 'very_high'            // $$$$$
}

export interface ServiceRatings {
  readonly quality?: number;           // Quality of service (1-5)
  readonly speed?: number;             // Speed of service (1-5)
  readonly friendliness?: number;      // Staff friendliness (1-5)
  readonly cleanliness?: number;       // Cleanliness (1-5)
  readonly value?: number;             // Value for money (1-5)
  readonly atmosphere?: number;        // Atmosphere/ambiance (1-5)
  readonly accessibility?: number;     // Accessibility features (1-5)
}

// Post Render State (for UI components)
export interface PostRenderState {
  readonly isVisible: boolean;
  readonly isExpanded: boolean;
  readonly loadingState: LoadingState;
  readonly interactionState?: PostInteractionState;
  readonly validationErrors?: ValidationError[];
  readonly componentProps?: Record<string, unknown>;
}

export enum LoadingState {
  IDLE = 'idle',
  LOADING = 'loading',
  SUCCESS = 'success',
  ERROR = 'error'
}

export interface PostInteractionState {
  readonly hasLiked: boolean;
  readonly hasSaved: boolean;
  readonly hasShared: boolean;
  readonly hasCommented: boolean;
  readonly currentReaction?: ReactionType;
  readonly isFollowingAuthor: boolean;
}

// Additional enums and interfaces for comprehensive post system
export enum VideoQuality {
  LOW_240P = '240p',
  MEDIUM_480P = '480p',
  HIGH_720P = '720p',
  FULL_HD_1080P = '1080p',
  ULTRA_HD_4K = '4k'
}

export enum DifficultyLevel {
  BEGINNER = 'beginner',
  INTERMEDIATE = 'intermediate',
  ADVANCED = 'advanced',
  EXPERT = 'expert'
}

export enum LifeEventType {
  RELATIONSHIP_STATUS = 'relationship_status',
  NEW_JOB = 'new_job',
  GRADUATION = 'graduation',
  MOVED = 'moved',
  TRAVEL = 'travel',
  ACHIEVEMENT = 'achievement',
  HEALTH = 'health',
  FAMILY = 'family',
  EDUCATION = 'education',
  WORK = 'work',
  HOME = 'home',
  HOBBY = 'hobby',
  VOLUNTEER = 'volunteer',
  CUSTOM = 'custom'
}

export enum LifeEventPrivacy {
  PUBLIC = 'public',
  FRIENDS = 'friends',
  FAMILY = 'family',
  PRIVATE = 'private'
}

export enum ProductCategory {
  ELECTRONICS = 'electronics',
  CLOTHING = 'clothing',
  HOME = 'home',
  BOOKS = 'books',
  SPORTS = 'sports',
  BEAUTY = 'beauty',
  AUTOMOTIVE = 'automotive',
  FOOD = 'food',
  HEALTH = 'health',
  TOYS = 'toys',
  ART = 'art',
  MUSIC = 'music',
  OTHER = 'other'
}

export enum ProductCondition {
  NEW = 'new',
  LIKE_NEW = 'like_new',
  GOOD = 'good',
  FAIR = 'fair',
  POOR = 'poor'
}

export enum ProductAvailability {
  AVAILABLE = 'available',
  PENDING = 'pending',
  SOLD = 'sold',
  RESERVED = 'reserved'
}

// Utility types for components
export type PostComponentProps<T extends PostContent = PostContent> = {
  post: PostData & { content: T };
  onReaction?: (reactionType: ReactionType) => void;
  onComment?: (comment: string) => void;
  onShare?: (shareType: ShareType, comment?: string) => void;
  onSave?: () => void;
  className?: string;
  children?: ReactNode;
};

// Type guards for post content
export const isImagePost = (content: PostContent): content is ImagePostContent => {
  return typeof content === 'object' && content !== null && 'images' in content;
};

export const isVideoPost = (content: PostContent): content is VideoPostContent => {
  return typeof content === 'object' && content !== null && 'videoUrl' in content && 'duration' in content;
};

export const isPollPost = (content: PostContent): content is PollContent => {
  return typeof content === 'object' && content !== null && 'question' in content && 'options' in content;
};

export const isCheckInPost = (content: PostContent): content is CheckInContent => {
  return typeof content === 'object' && content !== null && 'locationId' in content;
};

export const isLifeEventPost = (content: PostContent): content is LifeEventContent => {
  return typeof content === 'object' && content !== null && 'eventType' in content;
};

export const isProductPost = (content: PostContent): content is ProductShowcaseContent => {
  return typeof content === 'object' && content !== null && 'productId' in content;
};

export const isPostMessage = (content: PostContent): content is PostMessageContent => {
  return typeof content === 'object' && content !== null && 'recipientId' in content && 'message' in content;
};

export const isReviewPost = (content: PostContent): content is ReviewContent => {
  return typeof content === 'object' && content !== null && 'rating' in content && 'reviewText' in content;
};

export const isStoryPost = (content: PostContent): content is StoryContent => {
  return typeof content === 'object' && content !== null && 'mediaUrl' in content && 'mediaType' in content;
};

export const isReelPost = (content: PostContent): content is ReelContent => {
  return typeof content === 'object' && content !== null && 'videoUrl' in content && 'duration' in content && 'caption' in content;
};

export const isRecipePost = (content: PostContent): content is RecipeContent => {
  return typeof content === 'object' && content !== null && 'ingredients' in content && 'instructions' in content;
};

// Post validation utilities
export const validatePostData = (post: Partial<PostData>): ValidationError[] => {
  const errors: ValidationError[] = [];

  if (!post.id) {
    errors.push({
      field: 'id',
      message: 'Post ID is required',
      code: 'REQUIRED_FIELD',
      severity: ValidationSeverity.ERROR
    });
  }

  if (!post.authorId) {
    errors.push({
      field: 'authorId',
      message: 'Author ID is required',
      code: 'REQUIRED_FIELD',
      severity: ValidationSeverity.ERROR
    });
  }

  if (!post.type || !Object.values(PostType).includes(post.type)) {
    errors.push({
      field: 'type',
      message: 'Valid post type is required',
      code: 'INVALID_TYPE',
      severity: ValidationSeverity.ERROR
    });
  }

  if (!post.privacy || !Object.values(PostPrivacy).includes(post.privacy)) {
    errors.push({
      field: 'privacy',
      message: 'Valid privacy setting is required',
      code: 'INVALID_PRIVACY',
      severity: ValidationSeverity.ERROR
    });
  }

  return errors;
};

// Additional interfaces needed for complete system
interface ValidationError {
  readonly field: string;
  readonly message: string;
  readonly code: string;
  readonly severity: ValidationSeverity;
}

enum ValidationSeverity {
  INFO = 'info',
  WARNING = 'warning',
  ERROR = 'error'
}

interface LocationRestriction {
  readonly countries?: string[];
  readonly regions?: string[];
  readonly cities?: string[];
  readonly radius?: { center: Coordinates; distance: number };
}

interface AgeRange {
  readonly min?: number;
  readonly max?: number;
}

interface PromotionDetails {
  readonly campaignId: string;
  readonly budget: number;
  readonly targetAudience: PostAudience;
  readonly startDate: Date;
  readonly endDate: Date;
  readonly objective: PromotionObjective;
}

enum PromotionObjective {
  AWARENESS = 'awareness',
  ENGAGEMENT = 'engagement',
  TRAFFIC = 'traffic',
  CONVERSIONS = 'conversions',
  APP_INSTALLS = 'app_installs'
}

enum ReviewStatus {
  PENDING = 'pending',
  APPROVED = 'approved',
  REJECTED = 'rejected',
  APPEALED = 'appealed'
}

// Additional detailed interfaces for specialized content
interface StorySticker {
  readonly type: StickerType;
  readonly position: { x: number; y: number };
  readonly size: number;
  readonly rotation: number;
  readonly data?: Record<string, unknown>;
}

enum StickerType {
  EMOJI = 'emoji',
  GIF = 'gif',
  LOCATION = 'location',
  MENTION = 'mention',
  HASHTAG = 'hashtag',
  TIME = 'time',
  TEMPERATURE = 'temperature',
  MUSIC = 'music',
  POLL = 'poll',
  QUESTION = 'question',
  LINK = 'link'
}

interface StoryFilter {
  readonly name: string;
  readonly intensity: number;
}

interface StoryMusic {
  readonly songId: string;
  readonly title: string;
  readonly artist: string;
  readonly startTime: number;
  readonly duration: number;
}

interface StoryPoll {
  readonly question: string;
  readonly options: [string, string];
  readonly votes: [number, number];
}

interface StoryQuestion {
  readonly question: string;
  readonly responses: StoryQuestionResponse[];
}

interface StoryQuestionResponse {
  readonly userId: string;
  readonly userName: string;
  readonly response: string;
  readonly timestamp: Date;
}

interface StoryLink {
  readonly url: string;
  readonly title?: string;
}

interface StoryMention {
  readonly userId: string;
  readonly userName: string;
  readonly position: { x: number; y: number };
}

interface SubtitleTrack {
  readonly language: string;
  readonly label: string;
  readonly url: string;
  readonly isDefault: boolean;
}

interface VideoChapter {
  readonly title: string;
  readonly startTime: number;
  readonly endTime: number;
  readonly thumbnailUrl?: string;
}

interface ImageFilter {
  readonly name: string;
  readonly intensity: number;
}

interface ReelMusic {
  readonly songId: string;
  readonly title: string;
  readonly artist: string;
  readonly duration: number;
  readonly isOriginal: boolean;
}

interface ReelEffect {
  readonly name: string;
  readonly category: string;
  readonly parameters?: Record<string, unknown>;
}

interface StreamGuest {
  readonly userId: string;
  readonly userName: string;
  readonly joinTime: Date;
  readonly leaveTime?: Date;
}

interface StreamDonation {
  readonly userId: string;
  readonly userName: string;
  readonly amount: number;
  readonly currency: string;
  readonly message?: string;
  readonly timestamp: Date;
}

interface PlaylistTrack {
  readonly trackId: string;
  readonly title: string;
  readonly artist: string;
  readonly duration: number;
  readonly url?: string;
  readonly order: number;
}

interface TutorialStep {
  readonly stepNumber: number;
  readonly title: string;
  readonly description: string;
  readonly duration?: number;
  readonly imageUrl?: string;
  readonly videoUrl?: string;
  readonly materials?: string[];
  readonly tips?: string[];
}

enum TutorialCategory {
  COOKING = 'cooking',
  DIY = 'diy',
  TECHNOLOGY = 'technology',
  ART = 'art',
  MUSIC = 'music',
  SPORTS = 'sports',
  EDUCATION = 'education',
  BUSINESS = 'business',
  LIFESTYLE = 'lifestyle',
  OTHER = 'other'
}

interface PollOption {
  readonly id: string;
  readonly text: string;
  readonly imageUrl?: string;
  readonly votes: number;
  readonly percentage: number;
}

enum PollResultsVisibility {
  ALWAYS = 'always',
  AFTER_VOTE = 'after_vote',
  AFTER_END = 'after_end',
  NEVER = 'never'
}

interface QuizQuestion {
  readonly id: string;
  readonly question: string;
  readonly type: QuestionType;
  readonly options?: string[];
  readonly correctAnswer: string | string[];
  readonly explanation?: string;
  readonly points: number;
}

interface SurveyQuestion {
  readonly id: string;
  readonly question: string;
  readonly type: QuestionType;
  readonly options?: string[];
  readonly isRequired: boolean;
  readonly branching?: BranchingRule[];
}

enum QuestionType {
  MULTIPLE_CHOICE = 'multiple_choice',
  MULTIPLE_SELECT = 'multiple_select',
  TRUE_FALSE = 'true_false',
  SHORT_ANSWER = 'short_answer',
  LONG_ANSWER = 'long_answer',
  SCALE = 'scale',
  RATING = 'rating'
}

interface BranchingRule {
  readonly condition: BranchingCondition;
  readonly targetQuestion: string;
}

interface BranchingCondition {
  readonly operator: 'equals' | 'not_equals' | 'greater_than' | 'less_than' | 'contains';
  readonly value: unknown;
}

interface QAAnswer {
  readonly id: string;
  readonly userId: string;
  readonly userName: string;
  readonly answer: string;
  readonly timestamp: Date;
  readonly likes: number;
  readonly isAccepted: boolean;
  readonly isBestAnswer: boolean;
}

interface ChallengeSubmission {
  readonly id: string;
  readonly userId: string;
  readonly userName: string;
  readonly mediaUrl: string;
  readonly caption?: string;
  readonly timestamp: Date;
  readonly likes: number;
  readonly isWinner?: boolean;
}

interface PetitionSigner {
  readonly userId: string;
  readonly userName: string;
  readonly timestamp: Date;
  readonly comment?: string;
  readonly isAnonymous: boolean;
}

enum PetitionCategory {
  SOCIAL_JUSTICE = 'social_justice',
  ENVIRONMENT = 'environment',
  POLITICS = 'politics',
  HEALTHCARE = 'healthcare',
  EDUCATION = 'education',
  ANIMAL_RIGHTS = 'animal_rights',
  HUMAN_RIGHTS = 'human_rights',
  LOCAL_ISSUES = 'local_issues',
  OTHER = 'other'
}

interface PetitionUpdate {
  readonly id: string;
  readonly title: string;
  readonly description: string;
  readonly timestamp: Date;
  readonly imageUrl?: string;
}

enum CheckInActivity {
  DINING = 'dining',
  SHOPPING = 'shopping',
  SIGHTSEEING = 'sightseeing',
  BUSINESS = 'business',
  SOCIAL = 'social',
  EXERCISE = 'exercise',
  ENTERTAINMENT = 'entertainment',
  TRAVEL = 'travel'
}

interface TravelDay {
  readonly date: Date;
  readonly activities: TravelActivity[];
  readonly accommodations?: string;
  readonly meals?: string[];
  readonly transportation?: string;
  readonly cost?: number;
  readonly notes?: string;
}

interface TravelActivity {
  readonly name: string;
  readonly location: string;
  readonly time?: string;
  readonly duration?: number;
  readonly cost?: number;
  readonly rating?: number;
  readonly photos?: string[];
}

enum TravelType {
  SOLO = 'solo',
  COUPLE = 'couple',
  FAMILY = 'family',
  FRIENDS = 'friends',
  BUSINESS = 'business',
  GROUP = 'group'
}

enum MilestoneCategory {
  CAREER = 'career',
  EDUCATION = 'education',
  PERSONAL = 'personal',
  HEALTH = 'health',
  FINANCIAL = 'financial',
  RELATIONSHIP = 'relationship',
  CREATIVE = 'creative',
  SPIRITUAL = 'spiritual',
  SPORTS = 'sports',
  TRAVEL = 'travel'
}

interface MilestoneMetric {
  readonly name: string;
  readonly value: number;
  readonly unit: string;
  readonly previousValue?: number;
  readonly improvement?: number;
}

interface MemoryItem {
  readonly type: 'photo' | 'video' | 'post' | 'event';
  readonly url?: string;
  readonly description: string;
  readonly timestamp: Date;
  readonly location?: string;
}

enum AnniversaryType {
  RELATIONSHIP = 'relationship',
  WEDDING = 'wedding',
  WORK = 'work',
  FRIENDSHIP = 'friendship',
  ACHIEVEMENT = 'achievement',
  MEMORIAL = 'memorial',
  BUSINESS = 'business',
  PERSONAL = 'personal'
}

interface AnniversaryMilestone {
  readonly year: number;
  readonly description: string;
  readonly photos?: string[];
  readonly significance: string;
}

enum RecommendationType {
  PRODUCT = 'product',
  SERVICE = 'service',
  PLACE = 'place',
  EXPERIENCE = 'experience',
  CONTENT = 'content',
  PERSON = 'person'
}

enum GroupActivityType {
  MEETUP = 'meetup',
  EVENT = 'event',
  CHALLENGE = 'challenge',
  DISCUSSION = 'discussion',
  PROJECT = 'project',
  LEARNING = 'learning',
  SOCIAL = 'social',
  VOLUNTEER = 'volunteer'
}

interface ShippingInfo {
  readonly shipsFrom: string;
  readonly shipsTo: string[];
  readonly shippingCost?: number;
  readonly freeShipping?: boolean;
  readonly estimatedDays?: number;
}

interface ContactInfo {
  readonly phone?: string;
  readonly email?: string;
  readonly website?: string;
  readonly socialMedia?: Record<string, string>;
}

enum ServiceCategory {
  PROFESSIONAL = 'professional',
  CREATIVE = 'creative',
  TECHNICAL = 'technical',
  HEALTH = 'health',
  EDUCATION = 'education',
  CONSULTING = 'consulting',
  REPAIR = 'repair',
  PERSONAL = 'personal'
}

interface ServicePricing {
  readonly type: 'fixed' | 'hourly' | 'project' | 'subscription';
  readonly amount: number;
  readonly currency: string;
  readonly billingPeriod?: string;
  readonly minimumOrder?: number;
}

interface ServiceAvailability {
  readonly schedule: string;
  readonly timezone: string;
  readonly leadTime?: number;        // Days notice needed
  readonly maxCapacity?: number;
}

interface Testimonial {
  readonly id: string;
  readonly clientName: string;
  readonly rating: number;
  readonly text: string;
  readonly timestamp: Date;
  readonly projectType?: string;
}

enum EventType {
  CONFERENCE = 'conference',
  WORKSHOP = 'workshop',
  SEMINAR = 'seminar',
  NETWORKING = 'networking',
  SOCIAL = 'social',
  CULTURAL = 'cultural',
  SPORTS = 'sports',
  ENTERTAINMENT = 'entertainment',
  CHARITY = 'charity',
  BUSINESS = 'business'
}

interface EventLocation {
  readonly name: string;
  readonly address: string;
  readonly coordinates?: Coordinates;
  readonly capacity?: number;
  readonly accessibility?: string[];
}

interface EventTicketing {
  readonly isFree: boolean;
  readonly prices?: EventPrice[];
  readonly salesStart?: Date;
  readonly salesEnd?: Date;
  readonly refundPolicy?: string;
}

interface EventPrice {
  readonly tier: string;
  readonly price: number;
  readonly currency: string;
  readonly benefits?: string[];
  readonly quantity?: number;
  readonly earlyBird?: boolean;
}

interface EventOrganizer {
  readonly name: string;
  readonly type: 'individual' | 'company' | 'organization';
  readonly description?: string;
  readonly website?: string;
  readonly logo?: string;
  readonly contactInfo?: ContactInfo;
}

interface EventAgendaItem {
  readonly startTime: Date;
  readonly endTime: Date;
  readonly title: string;
  readonly description?: string;
  readonly speaker?: string;
  readonly location?: string;
  readonly type: 'session' | 'break' | 'networking' | 'meal';
}

interface EventSponsor {
  readonly name: string;
  readonly tier: 'platinum' | 'gold' | 'silver' | 'bronze' | 'partner';
  readonly logo: string;
  readonly website?: string;
  readonly description?: string;
}

enum JobType {
  FULL_TIME = 'full_time',
  PART_TIME = 'part_time',
  CONTRACT = 'contract',
  FREELANCE = 'freelance',
  INTERNSHIP = 'internship',
  TEMPORARY = 'temporary',
  REMOTE = 'remote',
  HYBRID = 'hybrid'
}

enum WorkLocation {
  REMOTE = 'remote',
  ON_SITE = 'on_site',
  HYBRID = 'hybrid'
}

interface SalaryRange {
  readonly min: number;
  readonly max: number;
  readonly currency: string;
  readonly period: 'hourly' | 'monthly' | 'yearly';
  readonly negotiable: boolean;
}

enum FundraiserCause {
  MEDICAL = 'medical',
  EDUCATION = 'education',
  DISASTER_RELIEF = 'disaster_relief',
  ANIMAL_WELFARE = 'animal_welfare',
  ENVIRONMENT = 'environment',
  COMMUNITY = 'community',
  MEMORIAL = 'memorial',
  SPORTS = 'sports',
  ARTS = 'arts',
  RELIGIOUS = 'religious',
  PERSONAL = 'personal',
  CHARITY = 'charity'
}

interface FundraiserOrganizer {
  readonly name: string;
  readonly relationship?: string;       // Relationship to beneficiary
  readonly isVerified: boolean;
  readonly previousFundraisers?: number;
  readonly location?: string;
}

interface FundraiserUpdate {
  readonly id: string;
  readonly title: string;
  readonly description: string;
  readonly timestamp: Date;
  readonly photos?: string[];
}

interface FundraiserDonor {
  readonly name?: string;              // Anonymous if not provided
  readonly amount: number;
  readonly timestamp: Date;
  readonly message?: string;
  readonly isAnonymous: boolean;
}

enum VerificationStatus {
  UNVERIFIED = 'unverified',
  PENDING = 'pending',
  VERIFIED = 'verified',
  REJECTED = 'rejected'
}

enum CollaborationType {
  CREATIVE = 'creative',
  TECHNICAL = 'technical',
  BUSINESS = 'business',
  RESEARCH = 'research',
  ARTISTIC = 'artistic',
  VOLUNTEER = 'volunteer',
  ACADEMIC = 'academic',
  SOCIAL = 'social'
}

enum CommitmentLevel {
  LOW = 'low',                        // 1-5 hours/week
  MEDIUM = 'medium',                  // 6-15 hours/week
  HIGH = 'high',                      // 16-30 hours/week
  FULL_TIME = 'full_time'             // 30+ hours/week
}

interface Collaborator {
  readonly userId: string;
  readonly name: string;
  readonly role: string;
  readonly skills: string[];
  readonly joinedAt: Date;
  readonly avatar?: string;
}

enum DietaryTag {
  VEGETARIAN = 'vegetarian',
  VEGAN = 'vegan',
  GLUTEN_FREE = 'gluten_free',
  DAIRY_FREE = 'dairy_free',
  NUT_FREE = 'nut_free',
  LOW_CARB = 'low_carb',
  KETO = 'keto',
  PALEO = 'paleo',
  HALAL = 'halal',
  KOSHER = 'kosher'
}

interface RecipeIngredient {
  readonly name: string;
  readonly amount: number;
  readonly unit: string;
  readonly notes?: string;
  readonly isOptional: boolean;
}

interface RecipeStep {
  readonly stepNumber: number;
  readonly instruction: string;
  readonly duration?: number;
  readonly temperature?: number;
  readonly images?: string[];
  readonly tips?: string;
}

interface NutritionInfo {
  readonly calories: number;
  readonly protein: number;
  readonly carbs: number;
  readonly fat: number;
  readonly fiber?: number;
  readonly sugar?: number;
  readonly sodium?: number;
}

enum WorkoutType {
  CARDIO = 'cardio',
  STRENGTH = 'strength',
  FLEXIBILITY = 'flexibility',
  BALANCE = 'balance',
  HIIT = 'hiit',
  YOGA = 'yoga',
  PILATES = 'pilates',
  SPORTS = 'sports',
  DANCE = 'dance',
  MARTIAL_ARTS = 'martial_arts'
}

interface WorkoutExercise {
  readonly name: string;
  readonly sets?: number;
  readonly reps?: number;
  readonly duration?: number;           // seconds
  readonly weight?: number;
  readonly restBetweenSets?: number;    // seconds
  readonly instructions?: string;
  readonly targetMuscles?: string[];
  readonly demonstration?: string;      // Video URL
}

enum ReadingStatus {
  WANT_TO_READ = 'want_to_read',
  CURRENTLY_READING = 'currently_reading',
  FINISHED = 'finished',
  DNF = 'dnf',                         // Did not finish
  REREADING = 'rereading'
}

interface BookQuote {
  readonly text: string;
  readonly page?: number;
  readonly chapter?: string;
  readonly notes?: string;
}

enum AchievementCategory {
  GAMING = 'gaming',
  FITNESS = 'fitness',
  LEARNING = 'learning',
  CAREER = 'career',
  CREATIVE = 'creative',
  SOCIAL = 'social',
  FINANCIAL = 'financial',
  TRAVEL = 'travel',
  HOBBY = 'hobby',
  PERSONAL = 'personal'
}

enum AchievementRarity {
  COMMON = 'common',
  UNCOMMON = 'uncommon',
  RARE = 'rare',
  EPIC = 'epic',
  LEGENDARY = 'legendary'
}

enum QuoteCategory {
  MOTIVATIONAL = 'motivational',
  INSPIRATIONAL = 'inspirational',
  WISDOM = 'wisdom',
  LOVE = 'love',
  FRIENDSHIP = 'friendship',
  SUCCESS = 'success',
  LIFE = 'life',
  HAPPINESS = 'happiness',
  FUNNY = 'funny',
  PHILOSOPHICAL = 'philosophical',
  SPIRITUAL = 'spiritual',
  LITERARY = 'literary'
}

interface QuoteTextStyle {
  readonly fontFamily?: string;
  readonly fontSize?: number;
  readonly color?: string;
  readonly alignment?: 'left' | 'center' | 'right';
  readonly style?: 'normal' | 'italic' | 'bold';
}

interface BusinessHours {
  readonly dayOfWeek: number;          // 0-6, Sunday-Saturday
  readonly openTime: string;           // HH:mm format
  readonly closeTime: string;          // HH:mm format
  readonly closed: boolean;
}

interface VenueReview {
  readonly id: string;
  readonly userId: string;
  readonly userName: string;
  readonly rating: number;
  readonly text: string;
  readonly photos?: string[];
  readonly timestamp: Date;
  readonly helpful: number;
}

export default PostType;