import React, { useState, useMemo } from 'react';
import {
  ArrowLeft,
  Share,
  Heart,
  MessageCircle,
  Bookmark,
  MoreHorizontal,
  Star,
  Filter,
  Camera,
  Play,
  Link,
  Check,
  Verified,
  ThumbsUp,
  Store
} from 'lucide-react';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import { Card, CardContent } from './ui/card';
import { ScrollArea } from './ui/scroll-area';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Separator } from './ui/separator';
import { Input } from './ui/input';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from './ui/dropdown-menu';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from './ui/dialog';
import { toast } from "sonner";

/**
 * Review data model - TikTok/Lemon8 style
 */
interface Review {
  id: string;
  userId: string;
  userName: string;
  userAvatar?: string;
  rating: number; // 1-5 stars
  title: string;
  content: string;
  date: string;
  isVerifiedPurchase: boolean;
  likeCount: number;
  commentCount: number;
  bookmarkCount: number;
  isLiked: boolean;
  isBookmarked: boolean;
  images: string[];
  hashtags: string[]; // #beauty #skincare #review
  productId?: string;
  productName?: string;
  shopId?: string;
  shopName?: string;
  response?: ReviewResponse;
  mood?: string; // "love it", "obsessed", "meh"
  skinType?: string; // for beauty reviews
  occasion?: string; // "daily use", "special occasion"
  ageRange?: string; // "20s", "30s", etc.
}

interface ReviewResponse {
  id: string;
  content: string;
  date: string;
  shopName: string;
}

/**
 * Review statistics data
 */
interface ReviewStats {
  totalReviews: number;
  averageRating: number;
  ratingDistribution: Record<number, number>; // star rating -> count
}

/**
 * Review filter options
 */
enum ReviewFilter {
  ALL = "All Reviews",
  FIVE_STAR = "5 Stars",
  FOUR_STAR = "4 Stars",
  THREE_STAR = "3 Stars",
  TWO_STAR = "2 Stars",
  ONE_STAR = "1 Star",
  WITH_PHOTOS = "With Photos",
  VERIFIED = "Verified Purchase"
}

/**
 * Review sort options
 */
enum ReviewSort {
  MOST_RECENT = "Most Recent",
  MOST_HELPFUL = "Most Helpful",
  HIGHEST_RATING = "Highest Rating",
  LOWEST_RATING = "Lowest Rating"
}

interface SharePlatform {
  id: string;
  name: string;
  icon: string;
  color: string;
  isAvailable: boolean;
}

interface ShareContent {
  title: string;
  description?: string;
  url?: string;
  imageUrl?: string;
  type: 'GENERAL' | 'PRODUCT' | 'SHOP' | 'LIVE_STREAM' | 'REVIEW';
}

interface ReviewScreenProps {
  targetId: string; // Product ID, Shop ID, etc.
  targetType: string; // "product", "shop", "user"
  targetName: string;
  onBack: () => void;
  onUserClick?: (userId: string) => void;
  onProductClick?: (productId: string) => void;
  onShopClick?: (shopId: string) => void;
}

// Mock data - TikTok/Lemon8 style
const getMockReviews = (): Review[] => [
  {
    id: "1",
    userId: "user1",
    userName: "Sarah Chen",
    rating: 5,
    title: "OMG this is a game changer! ✨",
    content: "Girl, I am OBSESSED! This product literally changed my life. The quality is chef's kiss and shipping was super fast. Already ordered 2 more for my besties! 💕",
    date: "2 days ago",
    isVerifiedPurchase: true,
    likeCount: 247,
    commentCount: 32,
    bookmarkCount: 89,
    isLiked: true,
    isBookmarked: false,
    images: ["https://via.placeholder.com/300x300?text=Before", "https://via.placeholder.com/300x300?text=After", "https://via.placeholder.com/300x300?text=Result"],
    hashtags: ["#obsessed", "#gamechanger", "#musthave", "#beauty", "#skincare", "#glowup"],
    mood: "obsessed",
    skinType: "combination",
    ageRange: "20s",
    occasion: "daily use",
    productId: "prod1",
    productName: "Glow Serum Pro"
  },
  {
    id: "2",
    userId: "user2",
    userName: "Kimmie Lifestyle",
    rating: 4,
    title: "Pretty good but...",
    content: "Okay so this is actually really nice! Love the texture and how it makes my skin feel. Battery could be better tho. Still would recommend to my followers! 💅",
    date: "1 week ago",
    isVerifiedPurchase: true,
    likeCount: 156,
    commentCount: 18,
    bookmarkCount: 45,
    isLiked: false,
    isBookmarked: true,
    images: ["https://via.placeholder.com/300x300?text=Product1", "https://via.placeholder.com/300x300?text=Product2", "https://via.placeholder.com/300x300?text=Product3", "https://via.placeholder.com/300x300?text=Product4", "https://via.placeholder.com/300x300?text=Product5"],
    hashtags: ["#honest", "#review", "#skincare", "#selfcare", "#nightroutine"],
    mood: "love it",
    skinType: "sensitive",
    ageRange: "30s",
    occasion: "night routine",
    response: {
      id: "resp1",
      content: "Thank you babe! We're working on improving battery life in our next version. Check DM for surprise! 💕",
      date: "5 days ago",
      shopName: "Glow Beauty Co"
    }
  },
  {
    id: "3",
    userId: "user3",
    userName: "Anna Minimalist",
    rating: 3,
    title: "It's okay I guess",
    content: "Not gonna lie, it's just... fine? Does what it says but nothing special. Good packaging though and arrived on time. Maybe I expected too much from the hype? 🤷‍♀️",
    date: "2 weeks ago",
    isVerifiedPurchase: false,
    likeCount: 67,
    commentCount: 12,
    bookmarkCount: 15,
    isLiked: false,
    isBookmarked: false,
    images: ["https://via.placeholder.com/300x300?text=Package"],
    hashtags: ["#honest", "#meh", "#overhyped", "#minimalist"],
    mood: "meh",
    skinType: "normal",
    ageRange: "20s",
    occasion: "testing"
  },
  {
    id: "4",
    userId: "user4",
    userName: "Beauty Guru TH",
    rating: 5,
    title: "Holy grail status! 🙌",
    content: "Y'ALL I've been testing this for 3 months now and WOW. My skin has never looked better! Even my dermatologist asked what I'm using. This is going straight to my holy grail list! 🌟",
    date: "3 days ago",
    isVerifiedPurchase: true,
    likeCount: 892,
    commentCount: 156,
    bookmarkCount: 234,
    isLiked: true,
    isBookmarked: true,
    images: ["https://via.placeholder.com/300x300?text=Before", "https://via.placeholder.com/300x300?text=After", "https://via.placeholder.com/300x300?text=Process", "https://via.placeholder.com/300x300?text=Final"],
    hashtags: ["#holygrail", "#transformation", "#skincare", "#glowup", "#beforeafter", "#3monthsupdate"],
    mood: "obsessed",
    skinType: "acne-prone",
    ageRange: "20s",
    occasion: "daily use"
  }
];

const calculateReviewStats = (reviews: Review[]): ReviewStats => {
  if (reviews.length === 0) {
    return { totalReviews: 0, averageRating: 0, ratingDistribution: {} };
  }

  const totalReviews = reviews.length;
  const averageRating = reviews.reduce((sum, r) => sum + r.rating, 0) / totalReviews;
  const distribution: Record<number, number> = {};

  reviews.forEach(review => {
    distribution[review.rating] = (distribution[review.rating] || 0) + 1;
  });

  return { totalReviews, averageRating, ratingDistribution: distribution };
};

const filterAndSortReviews = (reviews: Review[], filter: ReviewFilter, sort: ReviewSort): Review[] => {
  let filtered = reviews;

  switch (filter) {
    case ReviewFilter.FIVE_STAR:
      filtered = reviews.filter(r => r.rating === 5);
      break;
    case ReviewFilter.FOUR_STAR:
      filtered = reviews.filter(r => r.rating === 4);
      break;
    case ReviewFilter.THREE_STAR:
      filtered = reviews.filter(r => r.rating === 3);
      break;
    case ReviewFilter.TWO_STAR:
      filtered = reviews.filter(r => r.rating === 2);
      break;
    case ReviewFilter.ONE_STAR:
      filtered = reviews.filter(r => r.rating === 1);
      break;
    case ReviewFilter.WITH_PHOTOS:
      filtered = reviews.filter(r => r.images.length > 0);
      break;
    case ReviewFilter.VERIFIED:
      filtered = reviews.filter(r => r.isVerifiedPurchase);
      break;
    default:
      filtered = reviews;
  }

  switch (sort) {
    case ReviewSort.MOST_HELPFUL:
      return filtered.sort((a, b) => b.likeCount - a.likeCount);
    case ReviewSort.HIGHEST_RATING:
      return filtered.sort((a, b) => b.rating - a.rating);
    case ReviewSort.LOWEST_RATING:
      return filtered.sort((a, b) => a.rating - b.rating);
    default:
      return filtered; // Most recent by default
  }
};

const getMoodColor = (rating: number, mood: string) => {
  switch (rating) {
    case 5:
      return { bg: "bg-pink-100", text: "text-pink-600", border: "border-pink-200" };
    case 4:
      return { bg: "bg-purple-100", text: "text-purple-600", border: "border-purple-200" };
    case 3:
      return { bg: "bg-blue-100", text: "text-blue-600", border: "border-blue-200" };
    default:
      return { bg: "bg-gray-100", text: "text-gray-600", border: "border-gray-200" };
  }
};

const ReviewScreen: React.FC<ReviewScreenProps> = ({
  targetId,
  targetType,
  targetName,
  onBack,
  onUserClick = () => {},
  onProductClick = () => {},
  onShopClick = () => {}
}) => {
  const [selectedFilter, setSelectedFilter] = useState<ReviewFilter>(ReviewFilter.ALL);
  const [selectedSort, setSelectedSort] = useState<ReviewSort>(ReviewSort.MOST_RECENT);
  const [showShareModal, setShowShareModal] = useState(false);
  const [reviewToShare, setReviewToShare] = useState<Review | null>(null);
  const [showCopiedFeedback, setShowCopiedFeedback] = useState(false);

  const reviews = getMockReviews();
  const reviewStats = calculateReviewStats(reviews);

  const filteredReviews = useMemo(() =>
    filterAndSortReviews(reviews, selectedFilter, selectedSort),
    [reviews, selectedFilter, selectedSort]
  );

  const handleShareReview = (review: Review) => {
    setReviewToShare(review);
    setShowShareModal(true);
  };

  const handleShareAll = () => {
    setReviewToShare(null);
    setShowShareModal(true);
  };

  const handleLike = (reviewId: string) => {
    toast.success("Liked! ❤️");
  };

  const handleBookmark = (reviewId: string) => {
    toast.success("Saved! 🔖");
  };

  const handleComment = (reviewId: string) => {
    toast.info("Opening comments...");
  };

  const handleCopyLink = (url: string) => {
    navigator.clipboard.writeText(url);
    setShowCopiedFeedback(true);
    toast.success("Link copied! 📋");
    setTimeout(() => setShowCopiedFeedback(false), 2000);
  };

  const shareContent: ShareContent = reviewToShare ? {
    title: `${reviewToShare.userName}'s Review of ${targetName}`,
    description: reviewToShare.content,
    url: `https://tchat.app/reviews/${reviewToShare.id}`,
    type: 'REVIEW'
  } : {
    title: `${targetName} - Reviews (${reviewStats.averageRating.toFixed(1)}⭐)`,
    description: `Check out what people are saying about ${targetName}`,
    url: `https://tchat.app/reviews/${targetId}`,
    type: 'REVIEW'
  };

  const sharePlatforms: SharePlatform[] = [
    { id: "whatsapp", name: "WhatsApp", icon: "💬", color: "#25D366", isAvailable: true },
    { id: "facebook", name: "Facebook", icon: "📘", color: "#1877F2", isAvailable: true },
    { id: "twitter", name: "Twitter", icon: "🐦", color: "#1DA1F2", isAvailable: true },
    { id: "instagram", name: "Instagram", icon: "📷", color: "#E4405F", isAvailable: true },
    { id: "linkedin", name: "LinkedIn", icon: "💼", color: "#0A66C2", isAvailable: true },
    { id: "telegram", name: "Telegram", icon: "✈️", color: "#0088CC", isAvailable: true },
  ];

  return (
    <div className="flex flex-col h-screen bg-gray-50">
      {/* Header */}
      <div className="sticky top-0 z-10 bg-white border-b border-gray-200 px-4 py-3">
        <div className="flex items-center justify-between max-w-6xl mx-auto">
          <div className="flex items-center space-x-3">
            <Button variant="ghost" size="icon" onClick={onBack}>
              <ArrowLeft className="w-5 h-5" />
            </Button>
            <div>
              <h1 className="text-lg font-bold text-gray-900">Reviews</h1>
              <p className="text-sm text-gray-500">{targetName}</p>
            </div>
          </div>
          <Button variant="ghost" size="icon" onClick={handleShareAll}>
            <Share className="w-5 h-5" />
          </Button>
        </div>
      </div>

      <div className="flex-1 overflow-auto">
        <div className="max-w-4xl mx-auto p-4">
          {reviews.length === 0 ? (
            <div className="text-center py-20">
              <div className="text-6xl mb-4">⭐</div>
              <h2 className="text-xl font-semibold text-gray-900 mb-2">No Reviews Yet</h2>
              <p className="text-gray-500">Be the first to share your experience!</p>
            </div>
          ) : (
            <div className="space-y-6">
              {/* Review Statistics */}
              <Card className="bg-white">
                <CardContent className="p-6">
                  <div className="flex items-center justify-between">
                    <div className="flex-1">
                      <div className="flex items-center space-x-3 mb-2">
                        <span className="text-3xl font-bold text-gray-900">
                          {reviewStats.averageRating.toFixed(1)}
                        </span>
                        <div className="flex items-center space-x-1">
                          {[1, 2, 3, 4, 5].map(star => (
                            <Star
                              key={star}
                              className={`w-5 h-5 ${star <= reviewStats.averageRating ? 'text-yellow-500 fill-yellow-500' : 'text-gray-300'}`}
                            />
                          ))}
                        </div>
                      </div>
                      <p className="text-gray-600">
                        {reviewStats.totalReviews} {reviewStats.totalReviews === 1 ? 'review' : 'reviews'}
                      </p>
                    </div>

                    {/* Rating Distribution */}
                    <div className="flex-1 max-w-sm">
                      {[5, 4, 3, 2, 1].map(stars => {
                        const count = reviewStats.ratingDistribution[stars] || 0;
                        const percentage = reviewStats.totalReviews > 0 ? (count / reviewStats.totalReviews) * 100 : 0;

                        return (
                          <div key={stars} className="flex items-center space-x-2 mb-1">
                            <span className="text-sm text-gray-600 w-2">{stars}</span>
                            <Star className="w-3 h-3 text-yellow-500 fill-yellow-500" />
                            <div className="flex-1 h-2 bg-gray-200 rounded-full overflow-hidden">
                              <div
                                className="h-full bg-yellow-500 transition-all duration-300"
                                style={{ width: `${percentage}%` }}
                              />
                            </div>
                            <span className="text-sm text-gray-500 w-6 text-right">{count}</span>
                          </div>
                        );
                      })}
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Filter and Sort Controls */}
              <div className="bg-white rounded-lg border border-gray-200 p-4">
                <div className="space-y-4">
                  {/* Filters */}
                  <div>
                    <h3 className="text-sm font-medium text-gray-900 mb-2">Filter</h3>
                    <div className="flex flex-wrap gap-2">
                      {Object.values(ReviewFilter).map(filter => (
                        <Badge
                          key={filter}
                          variant={selectedFilter === filter ? "default" : "outline"}
                          className="cursor-pointer hover:bg-gray-100"
                          onClick={() => setSelectedFilter(filter)}
                        >
                          {filter}
                        </Badge>
                      ))}
                    </div>
                  </div>

                  {/* Sort */}
                  <div>
                    <h3 className="text-sm font-medium text-gray-900 mb-2">Sort by</h3>
                    <div className="flex flex-wrap gap-2">
                      {Object.values(ReviewSort).map(sort => (
                        <Badge
                          key={sort}
                          variant={selectedSort === sort ? "secondary" : "outline"}
                          className="cursor-pointer hover:bg-gray-100"
                          onClick={() => setSelectedSort(sort)}
                        >
                          {sort}
                        </Badge>
                      ))}
                    </div>
                  </div>
                </div>
              </div>

              {/* Reviews List */}
              <div className="space-y-6">
                {filteredReviews.map(review => (
                  <Card key={review.id} className="bg-white">
                    <CardContent className="p-6">
                      {/* Header - TikTok style */}
                      <div className="flex items-start justify-between mb-4">
                        <div className="flex items-start space-x-3">
                          <Avatar
                            className="w-12 h-12 cursor-pointer"
                            onClick={() => onUserClick(review.userId)}
                          >
                            <AvatarImage src={review.userAvatar} />
                            <AvatarFallback className="bg-gradient-to-r from-pink-500 to-purple-500 text-white font-bold">
                              {review.userName.charAt(0)}
                            </AvatarFallback>
                          </Avatar>

                          <div className="flex-1">
                            <div className="flex items-center space-x-2 mb-1">
                              <span
                                className="font-bold text-gray-900 cursor-pointer hover:underline"
                                onClick={() => onUserClick(review.userId)}
                              >
                                @{review.userName.toLowerCase().replace(/\s+/g, '')}
                              </span>
                              {review.isVerifiedPurchase && (
                                <Badge variant="secondary" className="text-xs">
                                  <Verified className="w-3 h-3 mr-1" />
                                  Verified
                                </Badge>
                              )}
                            </div>

                            <div className="flex items-center space-x-2">
                              {review.mood && (
                                <Badge
                                  className={`text-xs ${getMoodColor(review.rating, review.mood).bg} ${getMoodColor(review.rating, review.mood).text} ${getMoodColor(review.rating, review.mood).border}`}
                                  variant="outline"
                                >
                                  {review.mood}
                                </Badge>
                              )}
                              <span className="text-sm text-gray-500">{review.date}</span>
                            </div>
                          </div>
                        </div>

                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="icon">
                              <MoreHorizontal className="w-4 h-4" />
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent>
                            <DropdownMenuItem onClick={() => handleShareReview(review)}>
                              <Share className="w-4 h-4 mr-2" />
                              Share Review
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </div>

                      {/* Images - Instagram/Lemon8 style gallery */}
                      {review.images.length > 0 && (
                        <div className="mb-4">
                          <div className="flex space-x-2 overflow-x-auto pb-2">
                            {review.images.slice(0, 3).map((image, idx) => (
                              <div
                                key={idx}
                                className="flex-shrink-0 w-32 h-32 rounded-2xl bg-gradient-to-br from-pink-100 to-purple-100 flex items-center justify-center overflow-hidden"
                              >
                                <Camera className="w-8 h-8 text-gray-400" />
                              </div>
                            ))}
                            {review.images.length > 3 && (
                              <div className="flex-shrink-0 w-32 h-32 rounded-2xl bg-black bg-opacity-70 flex items-center justify-center text-white font-bold">
                                +{review.images.length - 3}
                              </div>
                            )}
                          </div>
                        </div>
                      )}

                      {/* Content */}
                      <div className="mb-4">
                        {review.title && (
                          <h3 className="text-lg font-bold text-gray-900 mb-2">
                            {review.title}
                          </h3>
                        )}
                        <p className="text-gray-800 leading-relaxed">
                          {review.content}
                        </p>
                      </div>

                      {/* Hashtags */}
                      {review.hashtags.length > 0 && (
                        <div className="mb-4">
                          <div className="flex flex-wrap gap-1">
                            {review.hashtags.map(hashtag => (
                              <span
                                key={hashtag}
                                className="text-blue-600 font-medium cursor-pointer hover:underline"
                              >
                                {hashtag}
                              </span>
                            ))}
                          </div>
                        </div>
                      )}

                      {/* Product Context */}
                      {review.productName && (
                        <div className="mb-4">
                          <Card
                            className="bg-blue-50 border-blue-200 cursor-pointer hover:bg-blue-100 transition-colors"
                            onClick={() => review.productId && onProductClick(review.productId)}
                          >
                            <CardContent className="p-3">
                              <div className="flex items-center space-x-2">
                                <Store className="w-4 h-4 text-blue-600" />
                                <span className="text-blue-800 font-medium">{review.productName}</span>
                              </div>
                            </CardContent>
                          </Card>
                        </div>
                      )}

                      {/* Info Chips - Lemon8 style */}
                      {(review.skinType || review.ageRange || review.occasion) && (
                        <div className="mb-4">
                          <div className="flex flex-wrap gap-2">
                            {review.skinType && (
                              <Badge variant="secondary" className="text-xs">
                                Skin: {review.skinType}
                              </Badge>
                            )}
                            {review.ageRange && (
                              <Badge variant="secondary" className="text-xs">
                                Age: {review.ageRange}
                              </Badge>
                            )}
                            {review.occasion && (
                              <Badge variant="secondary" className="text-xs">
                                Use: {review.occasion}
                              </Badge>
                            )}
                          </div>
                        </div>
                      )}

                      {/* Star Rating */}
                      <div className="flex items-center space-x-1 mb-4">
                        {[1, 2, 3, 4, 5].map(star => (
                          <Star
                            key={star}
                            className={`w-5 h-5 ${star <= review.rating ? 'text-yellow-500 fill-yellow-500' : 'text-gray-300'}`}
                          />
                        ))}
                      </div>

                      {/* Social Actions - TikTok style */}
                      <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-6">
                          {/* Like */}
                          <button
                            className="flex items-center space-x-2 text-gray-600 hover:text-pink-600 transition-colors"
                            onClick={() => handleLike(review.id)}
                          >
                            <Heart className={`w-5 h-5 ${review.isLiked ? 'fill-pink-600 text-pink-600' : ''}`} />
                            <span className="text-sm">{review.likeCount}</span>
                          </button>

                          {/* Comment */}
                          <button
                            className="flex items-center space-x-2 text-gray-600 hover:text-blue-600 transition-colors"
                            onClick={() => handleComment(review.id)}
                          >
                            <MessageCircle className="w-5 h-5" />
                            <span className="text-sm">{review.commentCount}</span>
                          </button>

                          {/* Share */}
                          <button
                            className="flex items-center space-x-2 text-gray-600 hover:text-green-600 transition-colors"
                            onClick={() => handleShareReview(review)}
                          >
                            <Share className="w-5 h-5" />
                          </button>
                        </div>

                        {/* Bookmark */}
                        <button
                          className="text-gray-600 hover:text-blue-600 transition-colors"
                          onClick={() => handleBookmark(review.id)}
                        >
                          <Bookmark className={`w-5 h-5 ${review.isBookmarked ? 'fill-blue-600 text-blue-600' : ''}`} />
                        </button>
                      </div>

                      {/* Shop Response */}
                      {review.response && (
                        <div className="mt-4">
                          <Card className="bg-blue-50 border-blue-200">
                            <CardContent className="p-4">
                              <div className="flex items-center justify-between mb-2">
                                <div className="flex items-center space-x-2">
                                  <Store className="w-4 h-4 text-blue-600" />
                                  <span className="text-blue-800 font-bold text-sm">
                                    {review.response.shopName}
                                  </span>
                                </div>
                                <span className="text-xs text-blue-600">
                                  {review.response.date}
                                </span>
                              </div>
                              <p className="text-blue-800 text-sm">
                                {review.response.content}
                              </p>
                            </CardContent>
                          </Card>
                        </div>
                      )}
                    </CardContent>
                  </Card>
                ))}

                {/* Load More */}
                {filteredReviews.length >= 10 && (
                  <div className="text-center py-6">
                    <Button variant="outline" size="lg">
                      <Plus className="w-4 h-4 mr-2" />
                      Load More Reviews
                    </Button>
                  </div>
                )}
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Share Modal */}
      <Dialog open={showShareModal} onOpenChange={setShowShareModal}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Share {shareContent.type === 'REVIEW' && reviewToShare ? 'Review' : 'Reviews'}</DialogTitle>
          </DialogHeader>

          <div className="space-y-4">
            {/* Content Preview */}
            <Card className="bg-gray-50">
              <CardContent className="p-4">
                <div className="flex items-start space-x-3">
                  <div className="w-10 h-10 bg-blue-600 rounded-full flex items-center justify-center">
                    <Star className="w-5 h-5 text-white" />
                  </div>
                  <div className="flex-1">
                    <h3 className="font-medium text-gray-900">{shareContent.title}</h3>
                    {shareContent.description && (
                      <p className="text-sm text-gray-600 mt-1 line-clamp-2">{shareContent.description}</p>
                    )}
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Share Platforms */}
            <div>
              <h4 className="text-sm font-medium text-gray-900 mb-3">Share via</h4>
              <div className="grid grid-cols-3 gap-3">
                {sharePlatforms.map(platform => (
                  <button
                    key={platform.id}
                    className="flex flex-col items-center p-3 rounded-lg hover:bg-gray-100 transition-colors"
                    onClick={() => {
                      toast.success(`Shared to ${platform.name}!`);
                      setShowShareModal(false);
                    }}
                  >
                    <div
                      className="w-12 h-12 rounded-full flex items-center justify-center text-2xl mb-2"
                      style={{ backgroundColor: `${platform.color}20` }}
                    >
                      {platform.icon}
                    </div>
                    <span className="text-xs text-gray-600 text-center">{platform.name}</span>
                  </button>
                ))}
              </div>
            </div>

            {/* Copy Link */}
            {shareContent.url && (
              <div>
                <h4 className="text-sm font-medium text-gray-900 mb-2">Or copy link</h4>
                <div
                  className="flex items-center p-3 bg-gray-100 rounded-lg cursor-pointer hover:bg-gray-200 transition-colors"
                  onClick={() => handleCopyLink(shareContent.url!)}
                >
                  {showCopiedFeedback ? (
                    <>
                      <Check className="w-5 h-5 text-green-600 mr-3" />
                      <span className="text-green-600 flex-1">Link copied to clipboard!</span>
                    </>
                  ) : (
                    <>
                      <Link className="w-5 h-5 text-gray-600 mr-3" />
                      <span className="text-gray-800 flex-1 truncate">{shareContent.url}</span>
                      <span className="text-blue-600 text-sm font-medium ml-2">Copy</span>
                    </>
                  )}
                </div>
              </div>
            )}
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default ReviewScreen;