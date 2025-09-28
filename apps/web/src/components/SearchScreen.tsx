import React, { useState, useMemo, useEffect } from 'react';
import { ArrowLeft, Search, Filter, Clock, TrendingUp, Hash, User, Store, MessageCircle, X } from 'lucide-react';
import { useSearchShopsQuery, useSearchProductsQuery, useSearchDiscoverQuery } from '../services/microservicesApi';
import { Button } from './ui/button';
import { Input } from './ui/input';
import { Badge } from './ui/badge';
import { Card, CardContent } from './ui/card';
import { ScrollArea } from './ui/scroll-area';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { Separator } from './ui/separator';

interface SearchScreenProps {
  user: any;
  onBack: () => void;
}

interface SearchResult {
  id: string;
  type: 'chat' | 'contact' | 'merchant' | 'product' | 'message' | 'hashtag';
  title: string;
  subtitle?: string;
  avatar?: string;
  timestamp?: string;
  highlight?: string;
  metadata?: any;
}

export function SearchScreen({ user, onBack }: SearchScreenProps) {
  const [searchQuery, setSearchQuery] = useState('');
  const [recentSearches, setRecentSearches] = useState([
    'Pad Thai',
    'Som Tam',
    'PromptPay',
    'Bangkok Street Food',
    'Family Group'
  ]);

  // RTK Query search - only trigger when there's a query
  const shouldSearch = searchQuery.length >= 2;

  const {
    data: shopsData,
    isLoading: shopsLoading,
    error: shopsError
  } = useSearchShopsQuery(
    { query: searchQuery, limit: 5 },
    { skip: !shouldSearch }
  );

  const {
    data: productsData,
    isLoading: productsLoading,
    error: productsError
  } = useSearchProductsQuery(
    { query: searchQuery, limit: 5 },
    { skip: !shouldSearch }
  );

  const {
    data: discoverData,
    isLoading: discoverLoading,
    error: discoverError
  } = useSearchDiscoverQuery(
    { query: searchQuery, limit: 5 },
    { skip: !shouldSearch }
  );

  const searchResults: SearchResult[] = useMemo(() => {
    if (!shouldSearch) {
      // Return empty array when no search query
      return [];
    }

    if (shopsLoading || productsLoading || discoverLoading) {
      // Fallback data while loading
      return [
        {
          id: '1',
          type: 'chat',
          title: 'Family Group',
          subtitle: 'Mom: Dinner at 7pm! 🍽️',
          avatar: '',
          timestamp: '5 min ago'
        },
        {
          id: '2',
          type: 'merchant',
          title: 'Somtam Vendor',
          subtitle: 'Thai Street Food • 0.5 km away',
          avatar: 'https://images.unsplash.com/photo-1743485753872-3b24372fcd24?w=150'
        }
      ];
    }

    const results: SearchResult[] = [];

    // Add shop results
    if (shopsData) {
      shopsData.forEach((shop: any) => {
        results.push({
          id: shop.id || shop.shop_id || `shop-${Math.random()}`,
          type: 'merchant',
          title: shop.name || shop.shop_name || 'Shop',
          subtitle: `${shop.category || 'Store'} • ${shop.distance || '0.5 km away'}`,
          avatar: shop.image || shop.logo || shop.avatar || undefined,
          metadata: shop
        });
      });
    }

    // Add product results
    if (productsData?.data) {
      productsData.data.forEach((product: any) => {
        results.push({
          id: product.id || product.product_id || `product-${Math.random()}`,
          type: 'product',
          title: product.name || product.title || 'Product',
          subtitle: `${product.currency || '฿'}${product.price || 0} • ${product.shop_name || 'Store'}`,
          avatar: product.image || product.image_url || product.thumbnail || undefined,
          metadata: product
        });
      });
    }

    // Add discover results (hashtags, posts, etc.)
    if (discoverData) {
      const discoverResults = Array.isArray(discoverData) ? discoverData : [discoverData];
      discoverResults.forEach((item: any) => {
        results.push({
          id: item.id || item.content_id || `discover-${Math.random()}`,
          type: item.type === 'hashtag' ? 'hashtag' : 'message',
          title: item.title || item.content || item.hashtag || 'Content',
          subtitle: item.subtitle || item.description || `${item.engagement_count || 0} interactions`,
          avatar: item.image || item.thumbnail || undefined,
          highlight: searchQuery,
          timestamp: item.timestamp || item.created_at || undefined,
          metadata: item
        });
      });
    }

    return results;
  }, [searchQuery, shopsData, productsData, discoverData, shopsLoading, productsLoading, discoverLoading, shouldSearch]);

  const trendingSearches = [
    { query: 'Songkran Festival', count: '1.2k' },
    { query: 'PromptPay QR', count: '890' },
    { query: 'Bangkok Street Food', count: '756' },
    { query: 'Thai New Year', count: '645' },
    { query: 'Som Tam Recipe', count: '532' }
  ];

  const popularHashtags = [
    '#ThaiStreetFood',
    '#BangkokEats',
    '#SongkranFestival',
    '#PromptPay',
    '#ThaiCulture',
    '#SEAFood',
    '#LiveCooking',
    '#ThaiMarket'
  ];

  const handleSearch = (query: string) => {
    if (query.trim()) {
      setSearchQuery(query);
      // Add to recent searches if not already there
      if (!recentSearches.includes(query)) {
        setRecentSearches([query, ...recentSearches.slice(0, 4)]);
      }
    }
  };

  const clearRecentSearch = (index: number) => {
    setRecentSearches(recentSearches.filter((_, i) => i !== index));
  };

  const getResultIcon = (type: string) => {
    switch (type) {
      case 'chat':
        return <MessageCircle className="w-5 h-5 text-chart-1" />;
      case 'merchant':
        return <Store className="w-5 h-5 text-chart-2" />;
      case 'contact':
        return <User className="w-5 h-5 text-chart-3" />;
      case 'hashtag':
        return <Hash className="w-5 h-5 text-chart-4" />;
      default:
        return <Search className="w-5 h-5 text-muted-foreground" />;
    }
  };

  const filteredResults = searchQuery 
    ? searchResults.filter(result => 
        result.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
        result.subtitle?.toLowerCase().includes(searchQuery.toLowerCase())
      )
    : [];

  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <header className="border-b border-border bg-card px-4 py-3">
        <div className="flex items-center gap-3">
          <Button variant="ghost" size="icon" onClick={onBack}>
            <ArrowLeft className="w-5 h-5" />
          </Button>
          <div className="flex-1 relative">
            <Search className="absolute left-3 top-3 w-4 h-4 text-muted-foreground" />
            <Input
              placeholder="Search messages, contacts, merchants..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && handleSearch(searchQuery)}
              className="pl-10"
              autoFocus
            />
            {searchQuery && (
              <Button
                variant="ghost"
                size="icon"
                className="absolute right-1 top-1 w-8 h-8"
                onClick={() => setSearchQuery('')}
              >
                <X className="w-4 h-4" />
              </Button>
            )}
          </div>
          <Button variant="ghost" size="icon">
            <Filter className="w-5 h-5" />
          </Button>
        </div>
      </header>

      {/* Content */}
      <div className="flex-1 overflow-hidden">
        {searchQuery ? (
          /* Search Results */
          <div className="h-full">
            <Tabs defaultValue="all" className="h-full flex flex-col">
              <TabsList className="mx-4 mt-4 grid w-full grid-cols-5">
                <TabsTrigger value="all" className="text-xs">All</TabsTrigger>
                <TabsTrigger value="chats" className="text-xs">Chats</TabsTrigger>
                <TabsTrigger value="merchants" className="text-xs">Shops</TabsTrigger>
                <TabsTrigger value="products" className="text-xs">Products</TabsTrigger>
                <TabsTrigger value="messages" className="text-xs">Messages</TabsTrigger>
              </TabsList>

              <TabsContent value="all" className="flex-1 overflow-hidden mt-4">
                <ScrollArea className="h-full px-4">
                  <div className="space-y-2 pb-4">
                    {filteredResults.length > 0 ? (
                      filteredResults.map((result) => (
                        <Card key={result.id} className="cursor-pointer hover:bg-accent/50 transition-colors">
                          <CardContent className="p-4">
                            <div className="flex items-center gap-3">
                              <div className="flex-shrink-0">
                                {result.avatar ? (
                                  <Avatar className="w-10 h-10">
                                    <AvatarImage src={result.avatar} />
                                    <AvatarFallback>
                                      {result.title.charAt(0)}
                                    </AvatarFallback>
                                  </Avatar>
                                ) : (
                                  <div className="w-10 h-10 bg-muted rounded-full flex items-center justify-center">
                                    {getResultIcon(result.type)}
                                  </div>
                                )}
                              </div>
                              
                              <div className="flex-1 min-w-0">
                                <div className="flex items-center gap-2 mb-1">
                                  <p className="font-medium truncate">{result.title}</p>
                                  <Badge variant="outline" className="text-xs">
                                    {result.type}
                                  </Badge>
                                </div>
                                {result.subtitle && (
                                  <p className="text-sm text-muted-foreground truncate">
                                    {result.highlight ? (
                                      <>
                                        {result.subtitle.split(result.highlight)[0]}
                                        <mark className="bg-yellow-200 dark:bg-yellow-800 px-1 rounded">
                                          {result.highlight}
                                        </mark>
                                        {result.subtitle.split(result.highlight)[1]}
                                      </>
                                    ) : (
                                      result.subtitle
                                    )}
                                  </p>
                                )}
                              </div>
                              
                              {result.timestamp && (
                                <span className="text-xs text-muted-foreground">
                                  {result.timestamp}
                                </span>
                              )}
                            </div>
                          </CardContent>
                        </Card>
                      ))
                    ) : (
                      <div className="text-center py-8">
                        <Search className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                        <h3 className="font-medium mb-2">No results found</h3>
                        <p className="text-sm text-muted-foreground">
                          Try different keywords or check the spelling
                        </p>
                      </div>
                    )}
                  </div>
                </ScrollArea>
              </TabsContent>

              <TabsContent value="chats" className="flex-1 overflow-hidden mt-4">
                <ScrollArea className="h-full px-4">
                  <div className="space-y-2 pb-4">
                    {filteredResults.filter(r => r.type === 'chat').map((result) => (
                      <Card key={result.id} className="cursor-pointer hover:bg-accent/50 transition-colors">
                        <CardContent className="p-4">
                          <div className="flex items-center gap-3">
                            <Avatar className="w-10 h-10">
                              <AvatarFallback>{result.title.charAt(0)}</AvatarFallback>
                            </Avatar>
                            <div className="flex-1">
                              <p className="font-medium">{result.title}</p>
                              <p className="text-sm text-muted-foreground">{result.subtitle}</p>
                            </div>
                            <span className="text-xs text-muted-foreground">{result.timestamp}</span>
                          </div>
                        </CardContent>
                      </Card>
                    ))}
                  </div>
                </ScrollArea>
              </TabsContent>

              <TabsContent value="merchants" className="flex-1 overflow-hidden mt-4">
                <ScrollArea className="h-full px-4">
                  <div className="space-y-2 pb-4">
                    {filteredResults.filter(r => r.type === 'merchant').map((result) => (
                      <Card key={result.id} className="cursor-pointer hover:bg-accent/50 transition-colors">
                        <CardContent className="p-4">
                          <div className="flex items-center gap-3">
                            <Avatar className="w-12 h-12">
                              <AvatarImage src={result.avatar} />
                              <AvatarFallback>{result.title.charAt(0)}</AvatarFallback>
                            </Avatar>
                            <div className="flex-1">
                              <p className="font-medium">{result.title}</p>
                              <p className="text-sm text-muted-foreground">{result.subtitle}</p>
                            </div>
                            <Badge variant="secondary">Verified</Badge>
                          </div>
                        </CardContent>
                      </Card>
                    ))}
                  </div>
                </ScrollArea>
              </TabsContent>
            </Tabs>
          </div>
        ) : (
          /* Search Discovery */
          <ScrollArea className="h-full">
            <div className="p-4 space-y-6">
              {/* Recent Searches */}
              {recentSearches.length > 0 && (
                <div>
                  <h3 className="font-medium mb-3 flex items-center gap-2">
                    <Clock className="w-4 h-4" />
                    Recent Searches
                  </h3>
                  <div className="space-y-2">
                    {recentSearches.map((search, index) => (
                      <div key={index} className="flex items-center justify-between p-3 rounded-lg hover:bg-accent/50 cursor-pointer transition-colors">
                        <button
                          onClick={() => handleSearch(search)}
                          className="flex items-center gap-3 flex-1 text-left"
                        >
                          <Clock className="w-4 h-4 text-muted-foreground" />
                          <span>{search}</span>
                        </button>
                        <Button
                          variant="ghost"
                          size="icon"
                          className="w-8 h-8"
                          onClick={() => clearRecentSearch(index)}
                        >
                          <X className="w-4 h-4" />
                        </Button>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              <Separator />

              {/* Trending */}
              <div>
                <h3 className="font-medium mb-3 flex items-center gap-2">
                  <TrendingUp className="w-4 h-4" />
                  Trending in Thailand
                </h3>
                <div className="space-y-2">
                  {trendingSearches.map((trend, index) => (
                    <button
                      key={index}
                      onClick={() => handleSearch(trend.query)}
                      className="flex items-center justify-between w-full p-3 rounded-lg hover:bg-accent/50 transition-colors text-left"
                    >
                      <div className="flex items-center gap-3">
                        <div className="w-6 h-6 bg-chart-1 rounded-full flex items-center justify-center text-xs font-bold text-white">
                          {index + 1}
                        </div>
                        <span>{trend.query}</span>
                      </div>
                      <Badge variant="secondary" className="text-xs">
                        {trend.count} searches
                      </Badge>
                    </button>
                  ))}
                </div>
              </div>

              <Separator />

              {/* Popular Hashtags */}
              <div>
                <h3 className="font-medium mb-3 flex items-center gap-2">
                  <Hash className="w-4 h-4" />
                  Popular Hashtags
                </h3>
                <div className="flex flex-wrap gap-2">
                  {popularHashtags.map((hashtag, index) => (
                    <button
                      key={index}
                      onClick={() => handleSearch(hashtag)}
                      className="px-3 py-2 bg-muted rounded-full text-sm hover:bg-accent transition-colors"
                    >
                      {hashtag}
                    </button>
                  ))}
                </div>
              </div>

              <Separator />

              {/* Quick Actions */}
              <div>
                <h3 className="font-medium mb-3">Quick Actions</h3>
                <div className="grid grid-cols-2 gap-3">
                  <Card className="cursor-pointer hover:bg-accent/50 transition-colors">
                    <CardContent className="p-4 text-center">
                      <Store className="w-8 h-8 text-chart-2 mx-auto mb-2" />
                      <p className="font-medium text-sm">Find Merchants</p>
                      <p className="text-xs text-muted-foreground">Nearby food vendors</p>
                    </CardContent>
                  </Card>
                  
                  <Card className="cursor-pointer hover:bg-accent/50 transition-colors">
                    <CardContent className="p-4 text-center">
                      <Hash className="w-8 h-8 text-chart-4 mx-auto mb-2" />
                      <p className="font-medium text-sm">Explore Tags</p>
                      <p className="text-xs text-muted-foreground">Discover content</p>
                    </CardContent>
                  </Card>
                </div>
              </div>
            </div>
          </ScrollArea>
        )}
      </div>
    </div>
  );
}