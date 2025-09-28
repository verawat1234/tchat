import React, { useState, useEffect, useMemo } from 'react';
import { Heart, MessageCircle, Share, MoreVertical, MapPin, Clock, Users, Calendar, Bell, ArrowRight, Navigation, Music, Flame, Sparkles, Star, Play, Camera, ShoppingCart, Ticket, TrendingUp, Eye, Bookmark, UserPlus, ChevronRight, Volume2, CheckCircle, Award, Image as ImageIcon, Video, AlertTriangle, Loader2 } from 'lucide-react';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import { Card, CardContent, CardHeader } from './ui/card';
import { ScrollArea } from './ui/scroll-area';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from './ui/dropdown-menu';
import { ImageWithFallback } from './figma/ImageWithFallback';
import { Progress } from './ui/progress';
import { Separator } from './ui/separator';
import { toast } from "sonner";
import {
  useGetUpcomingEventsQuery,
  useGetHistoricalEventsQuery,
  useGetEventDetailsQuery,
  useRegisterEventInterestMutation,
  useRsvpToEventMutation,
  useBookEventTicketsMutation
} from '../services/microservicesApi';

interface EventsTabProps {
  user: any;
  onBack?: () => void;
}

interface HistoricalEvent {
  id: string;
  title: string;
  date: string;
  location: string;
  category: 'music' | 'food' | 'cultural' | 'festival' | 'temple' | 'market';
  mainImage: string;
  galleryImages: string[];
  highlights: string[];
  attendeeCount: number;
  rating: number;
  reviews: {
    id: string;
    user: {
      name: string;
      avatar?: string;
    };
    rating: number;
    comment: string;
    photos?: string[];
    timestamp: string;
  }[];
  organizer: {
    name: string;
    avatar?: string;
    verified: boolean;
  };
  tags: string[];
  socialStats: {
    photos: number;
    videos: number;
    mentions: number;
  };
}

interface UpcomingEvent {
  id: string;
  title: string;
  description: string;
  startDate: string;
  endDate: string;
  location: string;
  venue: {
    name: string;
    address: string;
    capacity: number;
    facilities: string[];
  };
  category: 'music' | 'food' | 'cultural' | 'festival' | 'temple' | 'market';
  mainImage: string;
  galleryImages: string[];
  organizer: {
    name: string;
    avatar?: string;
    verified: boolean;
    pastEvents: number;
    rating: number;
  };
  ticketing: {
    available: boolean;
    prices: {
      type: string;
      price: number;
      currency: string;
      perks: string[];
      available: number;
      total: number;
    }[];
    salesEnd: string;
  };
  lineup?: {
    headliners: {
      name: string;
      image: string;
      genre: string;
      popularity: number;
    }[];
    supporting: {
      name: string;
      genre: string;
    }[];
  };
  schedule?: {
    day: string;
    stages: {
      name: string;
      acts: {
        time: string;
        artist: string;
        duration: string;
      }[];
    }[];
  }[];
  popularity: {
    trending: boolean;
    rank?: number;
    interest: number;
    attending: number;
    views: number;
  };
  socialProof: {
    friendsGoing: {
      name: string;
      avatar?: string;
    }[];
    influencersGoing: {
      name: string;
      avatar?: string;
      verified: boolean;
      followers: number;
    }[];
    totalFriendsGoing: number;
  };
  tags: string[];
  amenities: string[];
  ageRestriction?: string;
  weather?: {
    forecast: string;
    temperature: string;
    recommendation: string;
  };
}

export function EventsTab({ user, onBack }: EventsTabProps) {
  const [selectedEvent, setSelectedEvent] = useState<string | null>(null);
  const [selectedTicketType, setSelectedTicketType] = useState<string>('');
  const [ticketQuantity, setTicketQuantity] = useState(1);
  const [showGallery, setShowGallery] = useState(false);
  const [selectedGalleryIndex, setSelectedGalleryIndex] = useState(0);
  const [interestedEvents, setInterestedEvents] = useState<string[]>(['upcoming-1', 'upcoming-3']);
  const [attendingEvents, setAttendingEvents] = useState<string[]>(['upcoming-2']);
  const [followedOrganizers, setFollowedOrganizers] = useState<string[]>(['bangkok-festivals', 'thai-music-co']);
  const [isMounted, setIsMounted] = useState(false);

  useEffect(() => {
    setIsMounted(true);
  }, []);

  // RTK Query hooks for events data
  const { data: historicalEventsData, isLoading: historicalLoading, error: historicalError } = useGetHistoricalEventsQuery({
    limit: 10,
    category: 'all'
  }, {
    skip: !isMounted
  });

  const { data: upcomingEventsData, isLoading: upcomingLoading, error: upcomingError } = useGetUpcomingEventsQuery({
    limit: 10,
    timeframe: 'upcoming'
  }, {
    skip: !isMounted
  });

  // RTK Query mutations for event interactions
  const [registerInterest] = useRegisterEventInterestMutation();
  const [rsvpToEvent] = useRsvpToEventMutation();
  const [bookTickets] = useBookEventTicketsMutation();

  // Fallback historical events data
  const fallbackHistoricalData: HistoricalEvent[] = [
    {
      id: 'hist-1',
      title: 'Songkran Music Festival 2024',
      date: '2024-04-13',
      location: 'Lumpini Park, Bangkok',
      category: 'music',
      mainImage: 'https://images.unsplash.com/photo-1523050854058-8df90110c9d1?w=800&h=600&fit=crop',
      galleryImages: [
        'https://images.unsplash.com/photo-1523050854058-8df90110c9d1?w=400&h=300&fit=crop',
        'https://images.unsplash.com/photo-1540039155733-5bb30b53aa14?w=400&h=300&fit=crop',
        'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=400&h=300&fit=crop',
        'https://images.unsplash.com/photo-1542909168-82c3e7fdca5c?w=400&h=300&fit=crop'
      ],
      highlights: [
        'Epic water fights with traditional Thai music',
        'International DJs mixed with local artists',
        'Traditional Thai food vendors throughout the park',
        'Spectacular fireworks finale over the lake'
      ],
      attendeeCount: 85000,
      rating: 4.8,
      reviews: [
        {
          id: 'rev-1',
          user: { name: 'Sarah Chen', avatar: 'https://images.unsplash.com/photo-1494790108755-2616b612b820?w=50&h=50&fit=crop&crop=face' },
          rating: 5,
          comment: 'Absolutely incredible! The combination of traditional Songkran celebrations with modern music was perfect. The water stages were genius! 🎉💦',
          photos: ['https://images.unsplash.com/photo-1540039155733-5bb30b53aa14?w=300&h=200&fit=crop'],
          timestamp: '2024-04-14'
        },
        {
          id: 'rev-2',
          user: { name: 'Mike Thompson', avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=50&h=50&fit=crop&crop=face' },
          rating: 5,
          comment: 'Best festival experience in Thailand! The lineup was amazing and the cultural elements made it so authentic. Already planning for next year!',
          timestamp: '2024-04-15'
        }
      ],
      organizer: {
        name: 'Bangkok Music Festivals',
        avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=50&h=50&fit=crop',
        verified: true
      },
      tags: ['#Songkran2024', '#WaterFestival', '#ThaiMusic', '#Bangkok'],
      socialStats: {
        photos: 12450,
        videos: 2300,
        mentions: 45600
      }
    },
    {
      id: 'hist-2',
      title: 'Floating Market Food Festival 2024',
      date: '2024-02-10',
      location: 'Damnoen Saduak',
      category: 'food',
      mainImage: 'https://images.unsplash.com/photo-1513475382585-d06e58bcb0e0?w=800&h=600&fit=crop',
      galleryImages: [
        'https://images.unsplash.com/photo-1513475382585-d06e58bcb0e0?w=400&h=300&fit=crop',
        'https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?w=400&h=300&fit=crop',
        'https://images.unsplash.com/photo-1743485753872-3b24372fcd24?w=400&h=300&fit=crop'
      ],
      highlights: [
        'Over 200 traditional Thai vendors on boats',
        'Cooking demonstrations by renowned chefs',
        'Traditional longtail boat tours',
        'Live traditional Thai music performances'
      ],
      attendeeCount: 25000,
      rating: 4.6,
      reviews: [
        {
          id: 'rev-3',
          user: { name: 'Anna Liu', avatar: 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=50&h=50&fit=crop&crop=face' },
          rating: 4,
          comment: 'Such an authentic experience! The boat rides through the market while tasting fresh food was unforgettable. A bit crowded but totally worth it!',
          photos: ['https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?w=300&h=200&fit=crop'],
          timestamp: '2024-02-11'
        }
      ],
      organizer: {
        name: 'Thai Heritage Foundation',
        avatar: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=50&h=50&fit=crop',
        verified: true
      },
      tags: ['#FloatingMarket', '#ThaiFood', '#Heritage', '#Traditional'],
      socialStats: {
        photos: 8200,
        videos: 1100,
        mentions: 18500
      }
    }
  ];

  // Transform RTK Query historical events data to local format
  const historicalEvents = useMemo(() => {
    if (!historicalEventsData || historicalLoading || historicalError) {
      return fallbackHistoricalData;
    }
    return historicalEventsData.map(event => ({
      id: event.id,
      title: event.title || 'Unknown Event',
      date: event.date || 'TBD',
      location: event.location || 'Unknown Location',
      category: event.category || 'cultural',
      mainImage: event.mainImage || 'https://images.unsplash.com/photo-1523050854058-8df90110c9d1?w=800&h=600&fit=crop',
      galleryImages: event.galleryImages || [],
      highlights: event.highlights || [],
      attendeeCount: event.attendeeCount || 0,
      rating: event.rating || 0,
      reviews: event.reviews || [],
      organizer: event.organizer || { name: 'Unknown Organizer', verified: false },
      tags: event.tags || [],
      socialStats: event.socialStats || { photos: 0, videos: 0, mentions: 0 }
    }));
  }, [historicalEventsData, historicalLoading, historicalError]);

  // Fallback upcoming events data
  const fallbackUpcomingData: UpcomingEvent[] = [
    {
      id: 'upcoming-1',
      title: 'Bangkok Electronic Music Festival 2025',
      description: 'Southeast Asia\'s biggest electronic music festival returns with world-class DJs, state-of-the-art sound systems, and immersive visual experiences across multiple stages.',
      startDate: '2025-03-15',
      endDate: '2025-03-17',
      location: 'BITEC Bangna, Bangkok',
      venue: {
        name: 'BITEC Convention Center',
        address: '88 Bangna-Trad Road, Bangna, Bangkok',
        capacity: 50000,
        facilities: ['Air Conditioning', 'Multiple Food Courts', 'VIP Lounges', 'Parking', 'ATM', 'Medical Station']
      },
      category: 'music',
      mainImage: 'https://images.unsplash.com/photo-1540039155733-5bb30b53aa14?w=800&h=600&fit=crop',
      galleryImages: [
        'https://images.unsplash.com/photo-1540039155733-5bb30b53aa14?w=400&h=300&fit=crop',
        'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=400&h=300&fit=crop',
        'https://images.unsplash.com/photo-1542909168-82c3e7fdca5c?w=400&h=300&fit=crop'
      ],
      organizer: {
        name: 'Bangkok Music Co.',
        avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=50&h=50&fit=crop',
        verified: true,
        pastEvents: 15,
        rating: 4.7
      },
      ticketing: {
        available: true,
        prices: [
          {
            type: 'General Admission',
            price: 2500,
            currency: 'THB',
            perks: ['3-Day Access', 'Basic Amenities'],
            available: 15000,
            total: 20000
          },
          {
            type: 'VIP Experience',
            price: 8500,
            currency: 'THB',
            perks: ['3-Day Access', 'VIP Lounge', 'Premium Food & Drinks', 'Artist Meet & Greet', 'Exclusive Merchandise'],
            available: 450,
            total: 500
          },
          {
            type: 'Artist Circle',
            price: 15000,
            currency: 'THB',
            perks: ['3-Day Access', 'Backstage Tours', 'Private Artist Performances', 'Luxury Accommodation', 'Personal Concierge'],
            available: 8,
            total: 50
          }
        ],
        salesEnd: '2025-03-10'
      },
      lineup: {
        headliners: [
          {
            name: 'Calvin Harris',
            image: 'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=100&h=100&fit=crop&crop=face',
            genre: 'Electronic',
            popularity: 95
          },
          {
            name: 'Armin van Buuren',
            image: 'https://images.unsplash.com/photo-1542909168-82c3e7fdca5c?w=100&h=100&fit=crop&crop=face',
            genre: 'Trance',
            popularity: 92
          }
        ],
        supporting: [
          { name: 'Local Thai DJ Collective', genre: 'House' },
          { name: 'Bangkok Bass Club', genre: 'Dubstep' },
          { name: 'Siam Electronic', genre: 'Techno' }
        ]
      },
      schedule: [
        {
          day: 'Day 1 - March 15',
          stages: [
            {
              name: 'Main Stage',
              acts: [
                { time: '18:00', artist: 'Opening Ceremony', duration: '30min' },
                { time: '19:00', artist: 'Siam Electronic', duration: '90min' },
                { time: '21:00', artist: 'Calvin Harris', duration: '2hr' }
              ]
            }
          ]
        }
      ],
      popularity: {
        trending: true,
        rank: 1,
        interest: 125000,
        attending: 18500,
        views: 890000
      },
      socialProof: {
        friendsGoing: [
          { name: 'Sarah', avatar: 'https://images.unsplash.com/photo-1494790108755-2616b612b820?w=50&h=50&fit=crop&crop=face' },
          { name: 'Mike', avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=50&h=50&fit=crop&crop=face' }
        ],
        influencersGoing: [
          {
            name: 'DJ ThaiBeats',
            avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=50&h=50&fit=crop&crop=face',
            verified: true,
            followers: 850000
          }
        ],
        totalFriendsGoing: 12
      },
      tags: ['#BangkokEMF2025', '#Electronic', '#Calvin Harris', '#Armin'],
      amenities: ['Food Courts', 'Bars', 'Rest Areas', 'Merchandise', 'Photo Booths', 'Charging Stations'],
      ageRestriction: '18+',
      weather: {
        forecast: 'Partly Cloudy',
        temperature: '28-32°C',
        recommendation: 'Bring sunscreen and stay hydrated!'
      }
    },
    {
      id: 'upcoming-2',
      title: 'Thai Street Food Championship',
      description: 'The ultimate celebration of Thai street food culture featuring competitions, tastings, and cooking workshops with master chefs.',
      startDate: '2025-02-28',
      endDate: '2025-03-01',
      location: 'Chatuchak Weekend Market',
      venue: {
        name: 'Chatuchak Market Plaza',
        address: 'Kamphaeng Phet 2 Road, Chatuchak, Bangkok',
        capacity: 30000,
        facilities: ['Covered Areas', 'Traditional Seating', 'Cooking Stations', 'Food Courts']
      },
      category: 'food',
      mainImage: 'https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?w=800&h=600&fit=crop',
      galleryImages: [
        'https://images.unsplash.com/photo-1628432021231-4bbd431e6a04?w=400&h=300&fit=crop',
        'https://images.unsplash.com/photo-1513475382585-d06e58bcb0e0?w=400&h=300&fit=crop'
      ],
      organizer: {
        name: 'Thai Culinary Institute',
        avatar: 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=50&h=50&fit=crop',
        verified: true,
        pastEvents: 8,
        rating: 4.9
      },
      ticketing: {
        available: true,
        prices: [
          {
            type: 'Tasting Pass',
            price: 500,
            currency: 'THB',
            perks: ['All-Day Access', '10 Food Credits', 'Event Program'],
            available: 8000,
            total: 10000
          },
          {
            type: 'Chef Experience',
            price: 1500,
            currency: 'THB',
            perks: ['All-Day Access', 'Unlimited Tastings', 'Cooking Workshop', 'Meet the Chefs', 'Recipe Book'],
            available: 150,
            total: 200
          }
        ],
        salesEnd: '2025-02-25'
      },
      popularity: {
        trending: false,
        interest: 45000,
        attending: 12000,
        views: 230000
      },
      socialProof: {
        friendsGoing: [
          { name: 'Emma', avatar: 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=50&h=50&fit=crop&crop=face' }
        ],
        influencersGoing: [],
        totalFriendsGoing: 5
      },
      tags: ['#ThaiStreetFood', '#Chatuchak', '#FoodFestival', '#Championship'],
      amenities: ['Cooking Stations', 'Seating Areas', 'Washrooms', 'First Aid', 'Information Desk'],
      weather: {
        forecast: 'Sunny',
        temperature: '26-30°C',
        recommendation: 'Perfect weather for outdoor food tasting!'
      }
    }
  ];

  // Transform RTK Query upcoming events data to local format
  const upcomingEvents = useMemo(() => {
    if (!upcomingEventsData || upcomingLoading || upcomingError) {
      return fallbackUpcomingData;
    }
    return upcomingEventsData.map(event => ({
      id: event.id,
      title: event.title || 'Unknown Event',
      description: event.description || 'No description available',
      startDate: event.startDate || 'TBD',
      endDate: event.endDate || 'TBD',
      location: event.location || 'Unknown Location',
      category: event.category || 'cultural',
      mainImage: event.mainImage || 'https://images.unsplash.com/photo-1523050854058-8df90110c9d1?w=800&h=600&fit=crop',
      organizer: event.organizer || { name: 'Unknown Organizer', verified: false },
      lineup: event.lineup || [],
      ticketTypes: event.ticketTypes || [],
      amenities: event.amenities || [],
      tags: event.tags || [],
      socialProof: event.socialProof || { attending: 0, interested: 0 },
      ageRestriction: event.ageRestriction,
      weather: event.weather
    }));
  }, [upcomingEventsData, upcomingLoading, upcomingError]);

  const handleEventInterest = async (eventId: string) => {
    try {
      if (interestedEvents.includes(eventId)) {
        setInterestedEvents(prev => prev.filter(id => id !== eventId));
        await registerInterest({ eventId, interested: false });
        toast.success('Removed from interested');
      } else {
        setInterestedEvents(prev => [...prev, eventId]);
        await registerInterest({ eventId, interested: true });
        toast.success('Added to interested events');
      }
    } catch (error) {
      toast.error('Failed to update interest status');
    }
  };

  const handleEventRsvp = async (eventId: string) => {
    try {
      if (attendingEvents.includes(eventId)) {
        setAttendingEvents(prev => prev.filter(id => id !== eventId));
        await rsvpToEvent({ eventId, attending: false });
        toast.success('No longer attending');
      } else {
        setAttendingEvents(prev => [...prev, eventId]);
        await rsvpToEvent({ eventId, attending: true });
        toast.success('You\'re attending this event!');
      }
    } catch (error) {
      toast.error('Failed to update RSVP status');
    }
  };

  const handleEventAttend = async (eventId: string) => {
    // This is a duplicate of handleEventRsvp - calling the same function
    await handleEventRsvp(eventId);
  };

  const handleFollowOrganizer = (organizerId: string, organizerName: string) => {
    if (followedOrganizers.includes(organizerId)) {
      setFollowedOrganizers(prev => prev.filter(id => id !== organizerId));
      toast.success(`Unfollowed ${organizerName}`);
    } else {
      setFollowedOrganizers(prev => [...prev, organizerId]);
      toast.success(`Following ${organizerName}`);
    }
  };

  const handleBookTicket = (eventId: string, ticketType: string, quantity: number) => {
    const event = upcomingEvents.find(e => e.id === eventId);
    const ticket = event?.ticketing.prices.find(p => p.type === ticketType);
    if (ticket) {
      const total = ticket.price * quantity;
      toast.success(`Booking ${quantity}x ${ticketType} tickets for ฿${total.toLocaleString()}`);
    }
  };

  const getCategoryIcon = (category: string) => {
    switch (category) {
      case 'music':
        return <Music className="w-4 h-4" />;
      case 'food':
        return <ShoppingCart className="w-4 h-4" />;
      case 'cultural':
        return <Sparkles className="w-4 h-4" />;
      case 'festival':
        return <Calendar className="w-4 h-4" />;
      default:
        return <Calendar className="w-4 h-4" />;
    }
  };

  const renderHistoricalEvent = (event: HistoricalEvent) => (
    <Card key={event.id} className="mb-6 overflow-hidden">
      <div className="relative">
        <ImageWithFallback
          src={event.mainImage}
          alt={event.title}
          className="w-full h-48 object-cover"
        />
        <div className="absolute top-4 left-4">
          <Badge className="bg-black/70 text-white">
            Past Event
          </Badge>
        </div>
        <div className="absolute top-4 right-4 flex gap-2">
          <Button size="sm" variant="secondary" className="bg-white/90 hover:bg-white">
            <ImageIcon className="w-4 h-4 mr-1" />
            {event.socialStats.photos}
          </Button>
          <Button size="sm" variant="secondary" className="bg-white/90 hover:bg-white">
            <Video className="w-4 h-4 mr-1" />
            {event.socialStats.videos}
          </Button>
        </div>
      </div>
      
      <CardContent className="p-3 sm:p-4 lg:p-6 space-y-4 sm:space-y-6">
        {/* Event Header - Mobile First Design */}
        <div className="space-y-4">
          {/* Title and Basic Info */}
          <div className="space-y-3">
            <h3 className="text-lg sm:text-xl lg:text-2xl font-medium leading-tight pr-2">{event.title}</h3>
            
            {/* Event Details - Stacked on Mobile, Grid on Desktop */}
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 text-sm sm:text-base">
              <div className="flex items-center gap-2 p-2 bg-muted/30 rounded-lg">
                <Calendar className="w-4 h-4 text-muted-foreground flex-shrink-0" />
                <span className="font-medium">
                  {new Date(event.date).toLocaleDateString('th-TH', { 
                    year: 'numeric', 
                    month: 'short', 
                    day: 'numeric' 
                  })}
                </span>
              </div>
              
              <div className="flex items-center gap-2 p-2 bg-muted/30 rounded-lg">
                <MapPin className="w-4 h-4 text-muted-foreground flex-shrink-0" />
                <span className="truncate">{event.location}</span>
              </div>
              
              <div className="flex items-center gap-2 p-2 bg-muted/30 rounded-lg sm:col-span-2 lg:col-span-1">
                <Users className="w-4 h-4 text-muted-foreground flex-shrink-0" />
                <span className="font-medium text-chart-1">
                  {event.attendeeCount.toLocaleString()} attended
                </span>
              </div>
            </div>

            {/* Rating and Organizer - Responsive Layout */}
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
              <div className="flex items-center gap-3">
                <div className="flex items-center gap-1">
                  {[...Array(5)].map((_, i) => (
                    <Star
                      key={i}
                      className={`w-4 h-4 ${i < Math.floor(event.rating) ? 'text-yellow-500 fill-current' : 'text-gray-300'}`}
                    />
                  ))}
                  <span className="text-sm sm:text-base font-medium ml-1">{event.rating}</span>
                </div>
                <span className="text-xs sm:text-sm text-muted-foreground">
                  ({event.reviews.length} reviews)
                </span>
              </div>
              
              {/* Organizer Info */}
              <div className="flex items-center gap-3 p-2 bg-muted/20 rounded-lg">
                <Avatar className="w-8 h-8 sm:w-10 sm:h-10 cursor-pointer ring-1 ring-border">
                  <AvatarImage src={event.organizer.avatar} />
                  <AvatarFallback className="bg-gradient-to-br from-chart-1 to-chart-2 text-white">
                    {event.organizer.name.charAt(0)}
                  </AvatarFallback>
                </Avatar>
                <div className="min-w-0">
                  <div className="flex items-center gap-1">
                    <span className="text-sm font-medium truncate">{event.organizer.name}</span>
                    {event.organizer.verified && (
                      <CheckCircle className="w-4 h-4 text-blue-500 flex-shrink-0" />
                    )}
                  </div>
                  <div className="text-xs text-muted-foreground">Event Organizer</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Event Highlights - Enhanced Responsive Grid */}
        <div className="space-y-3">
          <h4 className="text-base sm:text-lg font-medium flex items-center gap-2">
            <Sparkles className="w-5 h-5 text-yellow-500" />
            Event Highlights
          </h4>
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-3">
            {event.highlights.map((highlight, index) => (
              <div key={index} className="flex items-start gap-3 p-3 bg-gradient-to-r from-muted/20 to-muted/30 rounded-lg border border-border/30">
                <div className="w-2 h-2 bg-chart-1 rounded-full mt-2 flex-shrink-0"></div>
                <span className="text-sm sm:text-base leading-relaxed">{highlight}</span>
              </div>
            ))}
          </div>
        </div>

        {/* Photo Gallery - Mobile Optimized */}
        <div className="space-y-3">
          <h4 className="text-base sm:text-lg font-medium flex items-center gap-2">
            <ImageIcon className="w-5 h-5 text-chart-2" />
            Event Gallery
          </h4>
          <div className="flex gap-2 sm:gap-3 overflow-x-auto pb-2 sm:pb-0">
            {event.galleryImages.slice(0, 4).map((image, index) => (
              <div key={index} className="relative flex-shrink-0">
                <ImageWithFallback
                  src={image}
                  alt={`Gallery ${index + 1}`}
                  className="w-20 h-20 sm:w-24 sm:h-24 lg:w-28 lg:h-28 object-cover rounded-xl cursor-pointer ring-2 ring-border/20 hover:ring-primary/30 transition-all"
                  onClick={() => {
                    setSelectedEvent(event.id);
                    setShowGallery(true);
                    setSelectedGalleryIndex(index);
                  }}
                />
                {index === 3 && event.galleryImages.length > 4 && (
                  <div className="absolute inset-0 bg-black/60 rounded-xl flex items-center justify-center">
                    <span className="text-white text-xs sm:text-sm font-medium">
                      +{event.galleryImages.length - 4}
                    </span>
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>

        {/* Recent Reviews - Enhanced Mobile Layout */}
        <div className="space-y-4">
          <h4 className="text-base sm:text-lg font-medium flex items-center gap-2">
            <MessageCircle className="w-5 h-5 text-chart-3" />
            What People Said
          </h4>
          <div className="space-y-4">
            {event.reviews.slice(0, 2).map((review) => (
              <div key={review.id} className="bg-gradient-to-r from-muted/30 to-muted/50 rounded-xl p-4 border border-border/30">
                <div className="flex items-start gap-3">
                  <Avatar className="w-10 h-10 sm:w-12 sm:h-12 ring-2 ring-border/20">
                    <AvatarImage src={review.user.avatar} />
                    <AvatarFallback className="bg-gradient-to-br from-chart-1 to-chart-2 text-white">
                      {review.user.name.charAt(0)}
                    </AvatarFallback>
                  </Avatar>
                  <div className="flex-1 min-w-0 space-y-2">
                    <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2">
                      <span className="font-medium text-sm sm:text-base truncate">{review.user.name}</span>
                      <div className="flex items-center gap-1">
                        {[...Array(5)].map((_, i) => (
                          <Star
                            key={i}
                            className={`w-3 h-3 sm:w-4 sm:h-4 ${i < review.rating ? 'text-yellow-500 fill-current' : 'text-gray-300'}`}
                          />
                        ))}
                      </div>
                    </div>
                    <p className="text-sm sm:text-base text-muted-foreground leading-relaxed">{review.comment}</p>
                    {review.photos && (
                      <div className="flex gap-2 mt-3">
                        {review.photos.map((photo, photoIndex) => (
                          <ImageWithFallback
                            key={photoIndex}
                            src={photo}
                            alt={`Review photo ${photoIndex + 1}`}
                            className="w-12 h-12 sm:w-16 sm:h-16 object-cover rounded-lg ring-1 ring-border/20"
                          />
                        ))}
                      </div>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Social Stats - Mobile Optimized Layout */}
        <div className="pt-4 border-t border-border/50">
          <div className="flex flex-col sm:flex-row items-start sm:items-center gap-3">
            <h5 className="text-sm font-medium text-muted-foreground">Event Impact:</h5>
            <div className="flex flex-wrap items-center gap-4 text-sm">
              <div className="flex items-center gap-2 px-3 py-1.5 bg-muted/30 rounded-full">
                <ImageIcon className="w-4 h-4 text-chart-1" />
                <span className="font-medium">{event.socialStats.photos.toLocaleString()}</span>
                <span className="text-muted-foreground">photos</span>
              </div>
              <div className="flex items-center gap-2 px-3 py-1.5 bg-muted/30 rounded-full">
                <Video className="w-4 h-4 text-chart-2" />
                <span className="font-medium">{event.socialStats.videos.toLocaleString()}</span>
                <span className="text-muted-foreground">videos</span>
              </div>
              <div className="flex items-center gap-2 px-3 py-1.5 bg-muted/30 rounded-full">
                <MessageCircle className="w-4 h-4 text-chart-3" />
                <span className="font-medium">{event.socialStats.mentions.toLocaleString()}</span>
                <span className="text-muted-foreground">mentions</span>
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );

  const renderUpcomingEvent = (event: UpcomingEvent) => (
    <Card key={event.id} className="mb-6 overflow-hidden">
      <div className="relative">
        <ImageWithFallback
          src={event.mainImage}
          alt={event.title}
          className="w-full h-64 object-cover"
        />
        {event.popularity.trending && (
          <div className="absolute top-4 left-4">
            <Badge className="bg-red-500 text-white flex items-center gap-1">
              <Flame className="w-3 h-3" />
              Trending #{event.popularity.rank}
            </Badge>
          </div>
        )}
        <div className="absolute top-4 right-4 flex flex-col gap-2">
          <Badge className="bg-white/90 text-gray-900">
            <Eye className="w-3 h-3 mr-1" />
            {(event.popularity.views / 1000).toFixed(0)}K views
          </Badge>
          <Badge className="bg-white/90 text-gray-900">
            <Users className="w-3 h-3 mr-1" />
            {event.popularity.attending.toLocaleString()} going
          </Badge>
        </div>
      </div>
      
      <CardContent className="p-3 sm:p-4 lg:p-6">
        {/* Header Section - Responsive Layout */}
        <div className="space-y-4 mb-6">
          {/* Title & Description */}
          <div className="space-y-2">
            <h3 className="text-lg sm:text-xl lg:text-2xl font-medium leading-tight">{event.title}</h3>
            <p className="text-sm sm:text-base text-muted-foreground leading-relaxed">{event.description}</p>
          </div>

          {/* Organizer Info - Mobile First */}
          <div className="flex items-center gap-3 p-3 bg-muted/20 rounded-xl">
            <Avatar className="w-10 h-10 sm:w-12 sm:h-12 cursor-pointer ring-2 ring-primary/10">
              <AvatarImage src={event.organizer.avatar} />
              <AvatarFallback className="bg-gradient-to-br from-chart-1 to-chart-2 text-white">
                {event.organizer.name.charAt(0)}
              </AvatarFallback>
            </Avatar>
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2 flex-wrap">
                <span className="text-sm sm:text-base font-medium truncate">{event.organizer.name}</span>
                {event.organizer.verified && (
                  <CheckCircle className="w-4 h-4 text-blue-500 flex-shrink-0" />
                )}
              </div>
              <div className="text-xs sm:text-sm text-muted-foreground">
                {event.organizer.pastEvents} events • ⭐ {event.organizer.rating}
              </div>
            </div>
          </div>
            
          {/* Event Details Grid - Responsive */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            <div className="space-y-3">
              <div className="flex items-start gap-3 p-3 bg-muted/10 rounded-lg">
                <Calendar className="w-4 h-4 text-muted-foreground mt-0.5 flex-shrink-0" />
                <div className="text-sm sm:text-base leading-relaxed">
                  <div className="font-medium">
                    {new Date(event.startDate).toLocaleDateString('th-TH', { 
                      year: 'numeric', 
                      month: 'long', 
                      day: 'numeric' 
                    })}
                  </div>
                  {event.endDate !== event.startDate && (
                    <div className="text-muted-foreground text-xs sm:text-sm">
                      until {new Date(event.endDate).toLocaleDateString('th-TH', { 
                        month: 'long', 
                        day: 'numeric' 
                      })}
                    </div>
                  )}
                </div>
              </div>
              
              <div className="flex items-start gap-3 p-3 bg-muted/10 rounded-lg">
                <MapPin className="w-4 h-4 text-muted-foreground mt-0.5 flex-shrink-0" />
                <div className="text-sm sm:text-base">
                  <div className="font-medium">{event.venue.name}</div>
                  <div className="text-muted-foreground text-xs sm:text-sm">
                    Capacity: {event.venue.capacity.toLocaleString()}
                  </div>
                </div>
              </div>

              {event.ageRestriction && (
                <div className="flex items-center gap-3 p-3 bg-muted/10 rounded-lg">
                  <Award className="w-4 h-4 text-muted-foreground flex-shrink-0" />
                  <span className="text-sm sm:text-base">{event.ageRestriction}</span>
                </div>
              )}
            </div>
            
            {/* Weather Card - Responsive */}
            {event.weather && (
              <div className="bg-gradient-to-br from-muted/30 to-muted/50 rounded-xl p-4 border border-border/50">
                <h5 className="font-medium text-sm sm:text-base mb-3 flex items-center gap-2">
                  <div className="w-2 h-2 bg-chart-1 rounded-full"></div>
                  Weather Forecast
                </h5>
                <div className="space-y-2 text-sm sm:text-base">
                  <div className="font-medium text-chart-1">{event.weather.forecast}</div>
                  <div className="text-xl sm:text-2xl font-medium">{event.weather.temperature}</div>
                  <div className="text-xs sm:text-sm text-muted-foreground italic">{event.weather.recommendation}</div>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Social Proof - Enhanced Mobile Layout */}
        {(event.socialProof.friendsGoing.length > 0 || event.socialProof.influencersGoing.length > 0) && (
          <div className="mb-6 p-4 bg-gradient-to-r from-muted/20 to-muted/30 rounded-xl border border-border/30">
            <div className="flex items-center gap-2 mb-3">
              <Users className="w-4 h-4 text-chart-1" />
              <span className="text-sm sm:text-base font-medium">People you might know</span>
            </div>
            
            <div className="space-y-3">
              <div className="flex items-center gap-3 flex-wrap">
                <div className="flex -space-x-2">
                  {event.socialProof.friendsGoing.slice(0, 3).map((friend, index) => (
                    <Avatar key={index} className="w-8 h-8 border-2 border-background ring-1 ring-border/20">
                      <AvatarImage src={friend.avatar} />
                      <AvatarFallback className="text-xs bg-gradient-to-br from-chart-1 to-chart-2 text-white">
                        {friend.name.charAt(0)}
                      </AvatarFallback>
                    </Avatar>
                  ))}
                </div>
                <div className="text-xs sm:text-sm text-muted-foreground leading-relaxed">
                  <span className="font-medium text-foreground">
                    {event.socialProof.friendsGoing.slice(0, 2).map(f => f.name).join(', ')}
                  </span>
                  {event.socialProof.totalFriendsGoing > 2 && (
                    <span> and {event.socialProof.totalFriendsGoing - 2} other friends</span>
                  )} are going
                </div>
              </div>
              
              {event.socialProof.influencersGoing.length > 0 && (
                <div className="flex items-center gap-3 flex-wrap">
                  <div className="flex -space-x-2">
                    {event.socialProof.influencersGoing.map((influencer, index) => (
                      <Avatar key={index} className="w-8 h-8 border-2 border-background ring-1 ring-border/20">
                        <AvatarImage src={influencer.avatar} />
                        <AvatarFallback className="text-xs bg-gradient-to-br from-chart-2 to-chart-3 text-white">
                          {influencer.name.charAt(0)}
                        </AvatarFallback>
                      </Avatar>
                    ))}
                  </div>
                  <div className="text-xs sm:text-sm text-muted-foreground leading-relaxed">
                    <span className="font-medium text-foreground">
                      {event.socialProof.influencersGoing.map(i => i.name).join(', ')}
                    </span> and other creators are going
                  </div>
                </div>
              )}
            </div>
          </div>
        )}

        {/* Lineup - Enhanced Responsive Design */}
        {event.lineup && (
          <div className="mb-6">
            <h4 className="font-medium text-base sm:text-lg mb-4 flex items-center gap-2">
              <Music className="w-5 h-5 text-chart-1" />
              Lineup
            </h4>
            <div className="space-y-6">
              <div>
                <h5 className="text-xs sm:text-sm font-medium text-muted-foreground mb-3 tracking-wide">HEADLINERS</h5>
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  {event.lineup.headliners.map((artist, index) => (
                    <div key={index} className="flex items-center gap-4 p-4 bg-gradient-to-r from-muted/20 to-muted/30 rounded-xl border border-border/30 hover:bg-muted/40 transition-colors">
                      <ImageWithFallback
                        src={artist.image}
                        alt={artist.name}
                        className="w-12 h-12 sm:w-14 sm:h-14 object-cover rounded-full ring-2 ring-primary/10"
                      />
                      <div className="flex-1 min-w-0">
                        <div className="font-medium text-sm sm:text-base truncate">{artist.name}</div>
                        <div className="text-xs sm:text-sm text-muted-foreground">{artist.genre}</div>
                      </div>
                      <div className="flex items-center gap-1.5 bg-orange-500/10 rounded-full px-2 py-1">
                        <Flame className="w-3 h-3 sm:w-4 sm:h-4 text-orange-500" />
                        <span className="text-xs sm:text-sm font-medium text-orange-700">{artist.popularity}%</span>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
              
              <div>
                <h5 className="text-xs sm:text-sm font-medium text-muted-foreground mb-3 tracking-wide">SUPPORTING ACTS</h5>
                <div className="flex flex-wrap gap-2">
                  {event.lineup.supporting.map((artist, index) => (
                    <Badge key={index} variant="secondary" className="flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs sm:text-sm">
                      {getCategoryIcon('music')}
                      {artist.name}
                    </Badge>
                  ))}
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Ticketing - Mobile Optimized */}
        {event.ticketing.available && (
          <div className="mb-6">
            <h4 className="font-medium text-base sm:text-lg mb-4 flex items-center gap-2">
              <Ticket className="w-5 h-5 text-chart-2" />
              Tickets
            </h4>
            <div className="space-y-4">
              {event.ticketing.prices.map((ticket, index) => (
                <div key={index} className="border border-border/50 rounded-xl p-4 sm:p-5 bg-gradient-to-r from-card to-muted/10">
                  <div className="flex flex-col sm:flex-row sm:items-start justify-between gap-4 mb-4">
                    <div className="flex-1">
                      <h5 className="font-medium text-base sm:text-lg mb-1">{ticket.type}</h5>
                      <div className="text-2xl sm:text-3xl font-bold text-green-600 mb-2">
                        ฿{ticket.price.toLocaleString()}
                      </div>
                    </div>
                    <div className="text-center sm:text-right">
                      <div className="text-sm text-muted-foreground mb-2">
                        <span className="font-medium text-foreground">{ticket.available}</span> / {ticket.total} available
                      </div>
                      <Progress 
                        value={(ticket.available / ticket.total) * 100} 
                        className="w-full sm:w-24 h-2"
                      />
                    </div>
                  </div>
                  
                  <div className="mb-4">
                    <h6 className="text-sm font-medium mb-2 text-muted-foreground">INCLUDES:</h6>
                    <div className="flex flex-wrap gap-1.5">
                      {ticket.perks.map((perk, perkIndex) => (
                        <Badge key={perkIndex} variant="outline" className="text-xs px-2 py-1 rounded-full">
                          {perk}
                        </Badge>
                      ))}
                    </div>
                  </div>
                  
                  <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-3">
                    <Button
                      size="sm"
                      onClick={() => handleBookTicket(event.id, ticket.type, 1)}
                      disabled={ticket.available === 0}
                      className="flex-1 h-10 rounded-xl"
                    >
                      <Ticket className="w-4 h-4 mr-2" />
                      {ticket.available === 0 ? 'Sold Out' : 'Book Now'}
                    </Button>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => handleEventInterest(event.id)}
                      className="sm:w-auto rounded-xl"
                    >
                      <Heart className={`w-4 h-4 ${interestedEvents.includes(event.id) ? 'fill-current text-red-500' : ''}`} />
                      <span className="ml-2 sm:hidden">
                        {interestedEvents.includes(event.id) ? 'Interested' : 'Add to Wishlist'}
                      </span>
                    </Button>
                  </div>
                </div>
              ))}
            </div>
            
            <div className="mt-4 p-3 bg-muted/20 rounded-lg border-l-4 border-chart-3">
              <div className="text-sm text-muted-foreground flex items-center gap-2">
                <Clock className="w-4 h-4 flex-shrink-0" />
                <span>
                  Ticket sales end: <span className="font-medium text-foreground">
                    {new Date(event.ticketing.salesEnd).toLocaleDateString('th-TH')}
                  </span>
                </span>
              </div>
            </div>
          </div>
        )}

        {/* Event Actions - Mobile Optimized */}
        <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 pt-4 border-t border-border/50">
          <Button
            variant={attendingEvents.includes(event.id) ? "default" : "outline"}
            onClick={() => handleEventAttend(event.id)}
            className="flex-1 h-11 rounded-xl"
          >
            <CheckCircle className="w-4 h-4 mr-2" />
            {attendingEvents.includes(event.id) ? 'Attending' : 'Join Event'}
          </Button>
          
          <div className="flex items-center gap-2 sm:gap-3">
            <Button
              variant={interestedEvents.includes(event.id) ? "secondary" : "outline"}
              onClick={() => handleEventInterest(event.id)}
              className="rounded-xl px-4"
            >
              <Heart className={`w-4 h-4 ${interestedEvents.includes(event.id) ? 'fill-current' : ''}`} />
              <span className="ml-2 sm:hidden text-sm">
                {interestedEvents.includes(event.id) ? 'Interested' : 'Interest'}
              </span>
            </Button>
            
            <Button variant="outline" className="rounded-xl px-4">
              <Share className="w-4 h-4" />
              <span className="ml-2 sm:hidden text-sm">Share</span>
            </Button>
            
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline" size="icon" className="rounded-xl">
                  <MoreVertical className="w-4 h-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-48">
                <DropdownMenuItem>
                  <Bookmark className="w-4 h-4 mr-2" />
                  Save Event
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Bell className="w-4 h-4 mr-2" />
                  Set Reminder
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <Navigation className="w-4 h-4 mr-2" />
                  Get Directions
                </DropdownMenuItem>
                <DropdownMenuItem
                  onClick={() => handleFollowOrganizer(event.organizer.name.toLowerCase().replace(/\s+/g, '-'), event.organizer.name)}
                >
                  <UserPlus className="w-4 h-4 mr-2" />
                  {followedOrganizers.includes(event.organizer.name.toLowerCase().replace(/\s+/g, '-')) ? 'Unfollow' : 'Follow'} Organizer
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </CardContent>
    </Card>
  );

  return (
    <div className="h-full flex flex-col">
      {onBack && (
        <div className="p-4 border-b">
          <Button variant="ghost" onClick={onBack} className="mb-2">
            ← Back
          </Button>
          <h1 className="text-2xl font-bold">Events</h1>
          <p className="text-muted-foreground">Discover amazing events in Thailand</p>
        </div>
      )}
      
      <Tabs defaultValue="upcoming" className="flex-1 flex flex-col">
        <div className="border-b px-4">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="upcoming" className="flex items-center gap-2">
              <Calendar className="w-4 h-4" />
              Upcoming
            </TabsTrigger>
            <TabsTrigger value="history" className="flex items-center gap-2">
              <Clock className="w-4 h-4" />
              Past Events
            </TabsTrigger>
          </TabsList>
        </div>
        
        <TabsContent value="upcoming" className="flex-1 m-0">
          <ScrollArea className="h-full">
            <div className="p-4 space-y-6">
              {upcomingLoading ? (
                <Card>
                  <CardContent className="p-8 text-center">
                    <Loader2 className="w-8 h-8 text-primary mx-auto mb-4 animate-spin" />
                    <p className="text-muted-foreground">Loading upcoming events...</p>
                  </CardContent>
                </Card>
              ) : upcomingError ? (
                <Card>
                  <CardContent className="p-8 text-center">
                    <AlertTriangle className="w-8 h-8 text-destructive mx-auto mb-4" />
                    <p className="text-destructive mb-2">Failed to load upcoming events</p>
                    <p className="text-muted-foreground text-sm">Using cached content</p>
                  </CardContent>
                </Card>
              ) : upcomingEvents.length === 0 ? (
                <Card>
                  <CardContent className="p-8 text-center">
                    <Calendar className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                    <p className="text-muted-foreground">No upcoming events found</p>
                  </CardContent>
                </Card>
              ) : (
                upcomingEvents.map(renderUpcomingEvent)
              )}
            </div>
          </ScrollArea>
        </TabsContent>
        
        <TabsContent value="history" className="flex-1 m-0">
          <ScrollArea className="h-full">
            <div className="p-3 sm:p-4 lg:p-6 space-y-4 sm:space-y-6 max-w-full w-full min-h-0 flex-1">
              <div className="text-center mb-4 sm:mb-6">
                <h2 className="text-lg sm:text-xl lg:text-2xl font-medium mb-2">Relive Amazing Moments</h2>
                <p className="text-sm sm:text-base text-muted-foreground max-w-md mx-auto leading-relaxed">
                  Explore past events and see what made them special
                </p>
              </div>
              <div className="space-y-4 sm:space-y-6">
                {historicalLoading ? (
                  <Card>
                    <CardContent className="p-8 text-center">
                      <Loader2 className="w-8 h-8 text-primary mx-auto mb-4 animate-spin" />
                      <p className="text-muted-foreground">Loading historical events...</p>
                    </CardContent>
                  </Card>
                ) : historicalError ? (
                  <Card>
                    <CardContent className="p-8 text-center">
                      <AlertTriangle className="w-8 h-8 text-destructive mx-auto mb-4" />
                      <p className="text-destructive mb-2">Failed to load historical events</p>
                      <p className="text-muted-foreground text-sm">Using cached content</p>
                    </CardContent>
                  </Card>
                ) : historicalEvents.length === 0 ? (
                  <Card>
                    <CardContent className="p-8 text-center">
                      <Calendar className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                      <p className="text-muted-foreground">No historical events found</p>
                    </CardContent>
                  </Card>
                ) : (
                  historicalEvents.map(renderHistoricalEvent)
                )}
              </div>
            </div>
          </ScrollArea>
        </TabsContent>
      </Tabs>
    </div>
  );
}