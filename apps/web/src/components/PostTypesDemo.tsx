import React, { useState } from 'react';
import {
  PostType,
  PostData,
  PostPrivacy,
  ReactionType,
  ShareType,
  PostMessageContent,
  ReviewContent,
  PollContent,
  ImagePostContent,
  CheckInContent,
  LocationCategory,
  ReviewSubjectType,
  isImagePost,
  isReviewPost,
  isPollPost,
  isPostMessage,
  isCheckInPost
} from '../types/PostData';

/**
 * Demo Component showcasing the unified 42 Post Types System
 * This demonstrates the comprehensive post type architecture with real examples
 */
export const PostTypesDemo: React.FC = () => {
  const [selectedPostType, setSelectedPostType] = useState<PostType>(PostType.TEXT);

  // Sample post data for each major category
  const samplePosts: Record<string, PostData> = {
    // Core Content Types
    [PostType.TEXT]: {
      id: '1',
      authorId: 'user1',
      authorName: 'Alice Johnson',
      authorAvatar: 'https://images.unsplash.com/photo-1494790108755-2616b612b1b7?w=150',
      timestamp: new Date('2024-01-15T10:30:00Z'),
      type: PostType.TEXT,
      content: "Just finished reading 'The Psychology of Social Media' - fascinating insights on how we connect online! 🧠✨",
      privacy: PostPrivacy.PUBLIC,
      engagement: {
        reactions: [
          { type: ReactionType.LIKE, userId: 'u1', timestamp: new Date(), userName: 'Bob' },
          { type: ReactionType.LOVE, userId: 'u2', timestamp: new Date(), userName: 'Carol' },
          { type: ReactionType.FIRE, userId: 'u3', timestamp: new Date(), userName: 'David' }
        ],
        comments: [
          {
            id: 'c1',
            userId: 'u4',
            userName: 'Emma',
            userAvatar: 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=50',
            content: 'Thanks for the recommendation! Adding it to my reading list 📚',
            timestamp: new Date(),
            replies: [],
            isEdited: false,
            isDeleted: false,
            mentionedUsers: []
          }
        ],
        shares: [],
        saves: ['u5', 'u6'],
        views: 127,
        reach: 98,
        impressions: 145,
        clickThroughs: 12,
        engagementRate: 0.078
      },
      metadata: {
        tags: ['#psychology', '#socialmedia', '#reading'],
        mentionedUsers: [],
        mood: 'contemplative' as const,
        feeling: 'good' as const,
        activity: 'reading' as const,
        isArchived: false,
        allowComments: true,
        allowShares: true,
        allowSaves: true,
        isPinned: false,
        isSticky: false
      }
    },

    [PostType.POST_MESSAGE]: {
      id: '2',
      authorId: 'user2',
      authorName: 'Mike Chen',
      authorAvatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=150',
      timestamp: new Date('2024-01-15T14:22:00Z'),
      type: PostType.POST_MESSAGE,
      content: {
        recipientId: 'user3',
        recipientName: 'Sarah Williams',
        message: "Happy Birthday Sarah! 🎉🎂 Hope you have an amazing day filled with joy and laughter. Thanks for being such an incredible friend!",
        messageType: 'birthday_wish' as const,
        isReply: false,
        visibility: 'friends' as const,
        allowReplies: true,
        taggedFriends: ['user4', 'user5']
      } as PostMessageContent,
      privacy: PostPrivacy.FRIENDS,
      engagement: {
        reactions: [
          { type: ReactionType.LOVE, userId: 'user3', timestamp: new Date(), userName: 'Sarah' },
          { type: ReactionType.CELEBRATE, userId: 'user4', timestamp: new Date(), userName: 'Tom' }
        ],
        comments: [
          {
            id: 'c2',
            userId: 'user3',
            userName: 'Sarah Williams',
            content: 'Thank you so much Mike! You made my day! 💕',
            timestamp: new Date(),
            replies: [],
            isEdited: false,
            isDeleted: false,
            mentionedUsers: ['user2']
          }
        ],
        shares: [],
        saves: [],
        views: 43,
        reach: 35,
        impressions: 47,
        clickThroughs: 2,
        engagementRate: 0.093
      }
    },

    [PostType.REVIEW]: {
      id: '3',
      authorId: 'user3',
      authorName: 'David Park',
      authorAvatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150',
      timestamp: new Date('2024-01-15T19:45:00Z'),
      type: PostType.REVIEW,
      content: {
        subjectType: ReviewSubjectType.BUSINESS,
        subjectId: 'rest_123',
        subjectName: 'Sakura Ramen House',
        rating: 5,
        reviewTitle: 'Outstanding Authentic Ramen Experience!',
        reviewText: 'Incredible authentic ramen with the richest, most flavorful broth I\'ve had outside of Japan. The chashu pork melts in your mouth and the noodles have perfect texture. Service was attentive and the atmosphere is cozy. Definitely my new favorite spot!',
        photos: [
          'https://images.unsplash.com/photo-1569718212165-3a8278d5f624?w=400',
          'https://images.unsplash.com/photo-1617093727343-374698b1b08d?w=400'
        ],
        visitDate: new Date('2024-01-15T18:30:00Z'),
        isVerifiedPurchase: true,
        pros: ['Amazing broth flavor', 'Perfect noodle texture', 'Generous portions', 'Authentic ingredients'],
        cons: ['Can get crowded', 'Limited parking'],
        wouldRecommend: true,
        helpfulCount: 12,
        location: {
          name: 'Sakura Ramen House',
          address: '123 Main St, Downtown',
          city: 'San Francisco',
          country: 'USA',
          category: LocationCategory.RESTAURANT
        },
        priceRange: 'medium' as const,
        serviceRatings: {
          quality: 5,
          speed: 4,
          friendliness: 5,
          cleanliness: 5,
          value: 4,
          atmosphere: 5
        }
      } as ReviewContent,
      privacy: PostPrivacy.PUBLIC,
      engagement: {
        reactions: [
          { type: ReactionType.FIRE, userId: 'u1', timestamp: new Date() },
          { type: ReactionType.LIKE, userId: 'u2', timestamp: new Date() },
          { type: ReactionType.WOW, userId: 'u3', timestamp: new Date() }
        ],
        comments: [
          {
            id: 'c3',
            userId: 'u4',
            userName: 'Food Lover',
            content: 'Thanks for the detailed review! Definitely going to try this place 🍜',
            timestamp: new Date(),
            replies: [],
            isEdited: false,
            isDeleted: false
          }
        ],
        shares: [
          {
            id: 's1',
            userId: 'u5',
            userName: 'Local Foodie',
            timestamp: new Date(),
            shareType: ShareType.QUOTE_SHARE,
            addedComment: 'Adding this to our restaurant guide!'
          }
        ],
        saves: ['u6', 'u7', 'u8'],
        views: 89,
        reach: 67,
        impressions: 102,
        clickThroughs: 8,
        engagementRate: 0.112
      }
    },

    [PostType.POLL]: {
      id: '4',
      authorId: 'user4',
      authorName: 'Emma Thompson',
      timestamp: new Date('2024-01-15T16:15:00Z'),
      type: PostType.POLL,
      content: {
        question: "What's your favorite way to spend a weekend? 🌟",
        options: [
          { id: 'opt1', text: 'Outdoor adventures & hiking', imageUrl: 'https://images.unsplash.com/photo-1551632811-561732d1e306?w=200', votes: 23, percentage: 31.5 },
          { id: 'opt2', text: 'Cozy reading at home', imageUrl: 'https://images.unsplash.com/photo-1481627834876-b7833e8f5570?w=200', votes: 18, percentage: 24.7 },
          { id: 'opt3', text: 'Exploring new restaurants', imageUrl: 'https://images.unsplash.com/photo-1517248135467-4c7edcad34c4?w=200', votes: 20, percentage: 27.4 },
          { id: 'opt4', text: 'Creative projects & hobbies', imageUrl: 'https://images.unsplash.com/photo-1513475382585-d06e58bcb0e0?w=200', votes: 12, percentage: 16.4 }
        ],
        allowMultipleChoices: false,
        showResults: 'after_vote' as const,
        endTime: new Date('2024-01-22T16:15:00Z'),
        totalVotes: 73,
        voterIds: ['u1', 'u2', 'u3', 'u4', 'u5']
      } as PollContent,
      privacy: PostPrivacy.PUBLIC,
      engagement: {
        reactions: [
          { type: ReactionType.LIKE, userId: 'u1', timestamp: new Date() },
          { type: ReactionType.WOW, userId: 'u2', timestamp: new Date() }
        ],
        comments: [
          {
            id: 'c4',
            userId: 'u3',
            userName: 'Adventure Seeker',
            content: 'Love this poll! I\'m definitely team outdoor adventures 🏔️',
            timestamp: new Date(),
            replies: [],
            isEdited: false,
            isDeleted: false
          }
        ],
        shares: [],
        saves: ['u4'],
        views: 156,
        reach: 134,
        impressions: 178,
        clickThroughs: 73, // Poll votes count as clicks
        engagementRate: 0.468
      }
    },

    [PostType.CHECK_IN]: {
      id: '5',
      authorId: 'user5',
      authorName: 'Jessica Wu',
      timestamp: new Date('2024-01-15T11:30:00Z'),
      type: PostType.CHECK_IN,
      content: {
        locationId: 'loc_456',
        locationName: 'Golden Gate Bridge',
        locationAddress: 'Golden Gate Bridge, San Francisco, CA',
        coordinates: { latitude: 37.8199, longitude: -122.4783 },
        category: LocationCategory.ATTRACTION,
        rating: 5,
        review: 'Perfect morning for a bridge walk! The fog cleared just in time for amazing photos. Such an iconic and breathtaking view of the bay! 🌉',
        photos: [
          'https://images.unsplash.com/photo-1449824913935-59a10b8d2000?w=400',
          'https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=400'
        ],
        companions: ['user6', 'user7'],
        activity: 'sightseeing' as const
      } as CheckInContent,
      privacy: PostPrivacy.PUBLIC,
      engagement: {
        reactions: [
          { type: ReactionType.LOVE, userId: 'u1', timestamp: new Date() },
          { type: ReactionType.WOW, userId: 'u2', timestamp: new Date() },
          { type: ReactionType.FIRE, userId: 'u3', timestamp: new Date() }
        ],
        comments: [
          {
            id: 'c5',
            userId: 'u4',
            userName: 'Travel Enthusiast',
            content: 'Gorgeous photos! The lighting is perfect 📸',
            timestamp: new Date(),
            replies: [
              {
                id: 'c5r1',
                userId: 'user5',
                userName: 'Jessica Wu',
                content: 'Thank you! We got so lucky with the weather ☀️',
                timestamp: new Date(),
                replies: [],
                isEdited: false,
                isDeleted: false
              }
            ],
            isEdited: false,
            isDeleted: false
          }
        ],
        shares: [
          {
            id: 's2',
            userId: 'u5',
            userName: 'SF Explorer',
            timestamp: new Date(),
            shareType: ShareType.STORY_SHARE
          }
        ],
        saves: ['u6', 'u7'],
        views: 203,
        reach: 178,
        impressions: 234,
        clickThroughs: 15,
        engagementRate: 0.083
      }
    }
  };

  const postTypeCategories = [
    {
      title: 'Core Content Types (8)',
      types: [PostType.TEXT, PostType.IMAGE, PostType.VIDEO, PostType.AUDIO, PostType.LINK_SHARE, PostType.POST_MESSAGE, PostType.REVIEW, PostType.ALBUM]
    },
    {
      title: 'Rich Media Types (6)',
      types: [PostType.STORY, PostType.REEL, PostType.LIVE_STREAM, PostType.PLAYLIST, PostType.MOOD_BOARD, PostType.TUTORIAL]
    },
    {
      title: 'Interactive Content (6)',
      types: [PostType.POLL, PostType.QUIZ, PostType.SURVEY, PostType.Q_AND_A, PostType.CHALLENGE, PostType.PETITION]
    },
    {
      title: 'Social & Location (8)',
      types: [PostType.CHECK_IN, PostType.TRAVEL_LOG, PostType.LIFE_EVENT, PostType.MILESTONE, PostType.MEMORY, PostType.ANNIVERSARY, PostType.RECOMMENDATION, PostType.GROUP_ACTIVITY]
    },
    {
      title: 'Commercial & Business (6)',
      types: [PostType.PRODUCT_SHOWCASE, PostType.SERVICE_PROMOTION, PostType.EVENT_PROMOTION, PostType.JOB_POSTING, PostType.FUNDRAISER, PostType.COLLABORATION]
    },
    {
      title: 'Specialized Content (8)',
      types: [PostType.RECIPE, PostType.WORKOUT, PostType.BOOK_REVIEW, PostType.MOOD_UPDATE, PostType.ACHIEVEMENT, PostType.QUOTE, PostType.MUSIC, PostType.VENUE]
    }
  ];

  const currentPost = samplePosts[selectedPostType];

  const renderPostContent = (post: PostData) => {
    if (typeof post.content === 'string') {
      return (
        <div className="bg-gray-50 p-4 rounded-lg">
          <p className="text-gray-800">{post.content}</p>
        </div>
      );
    }

    if (isPostMessage(post.content)) {
      return (
        <div className="bg-blue-50 p-4 rounded-lg border-l-4 border-blue-400">
          <div className="flex items-center gap-2 mb-2">
            <span className="text-sm text-blue-600">Message to {post.content.recipientName}</span>
            <span className="bg-blue-100 text-blue-700 px-2 py-1 rounded text-xs">{post.content.messageType}</span>
          </div>
          <p className="text-gray-800">{post.content.message}</p>
        </div>
      );
    }

    if (isReviewPost(post.content)) {
      return (
        <div className="space-y-4">
          <div className="bg-yellow-50 p-4 rounded-lg">
            <div className="flex items-center justify-between mb-2">
              <h3 className="font-semibold text-lg">{post.content.reviewTitle}</h3>
              <div className="flex items-center">
                {'★'.repeat(post.content.rating)}{'☆'.repeat(5 - post.content.rating)}
                <span className="ml-2 text-sm text-gray-600">({post.content.rating}/5)</span>
              </div>
            </div>
            <p className="text-gray-700 mb-3">{post.content.reviewText}</p>
            {post.content.photos && (
              <div className="grid grid-cols-2 gap-2">
                {post.content.photos.slice(0, 2).map((photo, index) => (
                  <img key={index} src={photo} alt="Review photo" className="w-full h-24 object-cover rounded" />
                ))}
              </div>
            )}
            {post.content.pros && post.content.pros.length > 0 && (
              <div className="mt-3">
                <h4 className="font-medium text-green-700 mb-1">Pros:</h4>
                <ul className="text-sm text-green-600">
                  {post.content.pros.map((pro, index) => (
                    <li key={index}>• {pro}</li>
                  ))}
                </ul>
              </div>
            )}
          </div>
        </div>
      );
    }

    if (isPollPost(post.content)) {
      return (
        <div className="bg-purple-50 p-4 rounded-lg">
          <h3 className="font-semibold mb-3">{post.content.question}</h3>
          <div className="space-y-3">
            {post.content.options.map((option) => (
              <div key={option.id} className="relative">
                <div className="flex items-center justify-between p-3 bg-white rounded border">
                  <div className="flex items-center gap-3">
                    {option.imageUrl && (
                      <img src={option.imageUrl} alt="" className="w-12 h-12 object-cover rounded" />
                    )}
                    <span>{option.text}</span>
                  </div>
                  <div className="text-right">
                    <div className="font-semibold">{option.percentage}%</div>
                    <div className="text-sm text-gray-500">{option.votes} votes</div>
                  </div>
                </div>
                <div
                  className="absolute bottom-0 left-0 h-1 bg-purple-400 rounded-full transition-all"
                  style={{ width: `${option.percentage}%` }}
                />
              </div>
            ))}
          </div>
          <p className="text-sm text-gray-500 mt-3">Total votes: {post.content.totalVotes}</p>
        </div>
      );
    }

    if (isCheckInPost(post.content)) {
      return (
        <div className="bg-green-50 p-4 rounded-lg">
          <div className="flex items-center gap-2 mb-2">
            <span className="text-green-600">📍</span>
            <span className="font-semibold">{post.content.locationName}</span>
            {post.content.rating && (
              <div className="flex items-center ml-auto">
                {'★'.repeat(post.content.rating)}{'☆'.repeat(5 - post.content.rating)}
              </div>
            )}
          </div>
          {post.content.review && (
            <p className="text-gray-700 mb-3">{post.content.review}</p>
          )}
          {post.content.photos && (
            <div className="grid grid-cols-2 gap-2">
              {post.content.photos.slice(0, 2).map((photo, index) => (
                <img key={index} src={photo} alt="Check-in photo" className="w-full h-24 object-cover rounded" />
              ))}
            </div>
          )}
        </div>
      );
    }

    return (
      <div className="bg-gray-50 p-4 rounded-lg">
        <p className="text-gray-500 italic">Content preview for {selectedPostType} posts</p>
      </div>
    );
  };

  const renderEngagementStats = (engagement: any) => {
    const reactionCounts = engagement.reactions?.reduce((acc: Record<string, number>, reaction: any) => {
      acc[reaction.type] = (acc[reaction.type] || 0) + 1;
      return acc;
    }, {}) || {};

    return (
      <div className="bg-white p-4 rounded-lg border">
        <h3 className="font-semibold mb-3">Engagement Analytics</h3>
        <div className="grid grid-cols-2 gap-4 text-sm">
          <div>
            <span className="text-gray-600">Views:</span> <span className="font-medium">{engagement.views || 0}</span>
          </div>
          <div>
            <span className="text-gray-600">Reach:</span> <span className="font-medium">{engagement.reach || 0}</span>
          </div>
          <div>
            <span className="text-gray-600">Comments:</span> <span className="font-medium">{engagement.comments?.length || 0}</span>
          </div>
          <div>
            <span className="text-gray-600">Shares:</span> <span className="font-medium">{engagement.shares?.length || 0}</span>
          </div>
          <div>
            <span className="text-gray-600">Saves:</span> <span className="font-medium">{engagement.saves?.length || 0}</span>
          </div>
          <div>
            <span className="text-gray-600">Engagement Rate:</span> <span className="font-medium">{((engagement.engagementRate || 0) * 100).toFixed(1)}%</span>
          </div>
        </div>
        {Object.keys(reactionCounts).length > 0 && (
          <div className="mt-3">
            <span className="text-gray-600 text-sm">Reactions:</span>
            <div className="flex gap-2 mt-1">
              {Object.entries(reactionCounts).map(([type, count]) => (
                <span key={type} className="bg-gray-100 px-2 py-1 rounded text-xs">
                  {type}: {count as number}
                </span>
              ))}
            </div>
          </div>
        )}
      </div>
    );
  };

  return (
    <div className="max-w-6xl mx-auto p-6 space-y-8">
      <div className="text-center">
        <h1 className="text-3xl font-bold mb-4">📱 Unified Post Types System Demo</h1>
        <p className="text-lg text-gray-600 mb-2">
          Comprehensive 42 Post Types Architecture
        </p>
        <p className="text-sm text-gray-500">
          Cross-platform consistency between web TypeScript and mobile Kotlin
        </p>
      </div>

      {/* Post Type Selector */}
      <div className="space-y-6">
        {postTypeCategories.map((category) => (
          <div key={category.title} className="bg-white p-4 rounded-lg border">
            <h2 className="font-semibold text-lg mb-3">{category.title}</h2>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-2">
              {category.types.map((type) => (
                <button
                  key={type}
                  onClick={() => setSelectedPostType(type)}
                  className={`p-2 text-sm rounded transition-colors ${
                    selectedPostType === type
                      ? 'bg-blue-500 text-white'
                      : 'bg-gray-100 hover:bg-gray-200 text-gray-700'
                  }`}
                >
                  {type.replace(/_/g, ' ')}
                </button>
              ))}
            </div>
          </div>
        ))}
      </div>

      {/* Selected Post Preview */}
      {currentPost ? (
        <div className="grid md:grid-cols-2 gap-6">
          {/* Post Content */}
          <div className="bg-white p-6 rounded-lg border">
            <div className="flex items-center gap-3 mb-4">
              <img
                src={currentPost.authorAvatar}
                alt={currentPost.authorName}
                className="w-12 h-12 rounded-full object-cover"
              />
              <div className="flex-1">
                <h3 className="font-semibold">{currentPost.authorName}</h3>
                <div className="flex items-center gap-2 text-sm text-gray-500">
                  <span>{currentPost.timestamp.toLocaleDateString()}</span>
                  <span>•</span>
                  <span className="bg-blue-100 text-blue-700 px-2 py-1 rounded">
                    {selectedPostType.replace(/_/g, ' ')}
                  </span>
                  <span>•</span>
                  <span className="bg-green-100 text-green-700 px-2 py-1 rounded">
                    {currentPost.privacy}
                  </span>
                </div>
              </div>
            </div>

            {renderPostContent(currentPost)}

            {currentPost.metadata?.tags && (
              <div className="flex flex-wrap gap-1 mt-3">
                {currentPost.metadata.tags.map((tag, index) => (
                  <span key={index} className="text-blue-500 text-sm">{tag}</span>
                ))}
              </div>
            )}
          </div>

          {/* Engagement Analytics */}
          {currentPost.engagement && renderEngagementStats(currentPost.engagement)}
        </div>
      ) : (
        <div className="text-center p-12 bg-gray-50 rounded-lg">
          <p className="text-gray-500">Select a post type to see preview</p>
        </div>
      )}

      {/* System Features Overview */}
      <div className="bg-gradient-to-r from-blue-50 to-purple-50 p-6 rounded-lg">
        <h2 className="text-2xl font-bold mb-4">🚀 System Features</h2>
        <div className="grid md:grid-cols-3 gap-6 text-sm">
          <div>
            <h3 className="font-semibold text-lg mb-2">✨ Enhanced Engagement</h3>
            <ul className="space-y-1 text-gray-600">
              <li>• 10 reaction types (Like, Love, Fire, etc.)</li>
              <li>• Nested comment threads</li>
              <li>• 5 share types with custom messages</li>
              <li>• Advanced analytics & metrics</li>
            </ul>
          </div>
          <div>
            <h3 className="font-semibold text-lg mb-2">🔒 Privacy & Targeting</h3>
            <ul className="space-y-1 text-gray-600">
              <li>• 8 privacy levels</li>
              <li>• Custom audience targeting</li>
              <li>• Location-based restrictions</li>
              <li>• Age and interest targeting</li>
            </ul>
          </div>
          <div>
            <h3 className="font-semibold text-lg mb-2">🔧 Developer Experience</h3>
            <ul className="space-y-1 text-gray-600">
              <li>• Type-safe discriminated unions</li>
              <li>• Comprehensive validation</li>
              <li>• Cross-platform consistency</li>
              <li>• Easy extensibility</li>
            </ul>
          </div>
        </div>
      </div>

      {/* Statistics */}
      <div className="bg-white p-6 rounded-lg border text-center">
        <h2 className="text-xl font-bold mb-4">📊 Architecture Statistics</h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div>
            <div className="text-3xl font-bold text-blue-600">42</div>
            <div className="text-sm text-gray-600">Post Types</div>
          </div>
          <div>
            <div className="text-3xl font-bold text-purple-600">10</div>
            <div className="text-sm text-gray-600">Reaction Types</div>
          </div>
          <div>
            <div className="text-3xl font-bold text-green-600">8</div>
            <div className="text-sm text-gray-600">Privacy Levels</div>
          </div>
          <div>
            <div className="text-3xl font-bold text-orange-600">100%</div>
            <div className="text-sm text-gray-600">Type Safety</div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default PostTypesDemo;