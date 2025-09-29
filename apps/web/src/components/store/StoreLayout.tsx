import React, { useState } from 'react';
import { Store, Grid3X3, Eye, Play } from 'lucide-react';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '../ui/tabs';
import StreamTab from '../StreamTab';

interface StoreLayoutProps {
  user: any;
  // Store-related props
  onShopClick?: (shopId: string) => void;
  onProductClick?: (productId: string) => void;
  onAddToCart?: (productId: string, quantity?: number) => void;
  onLiveStreamClick?: (streamId: string) => void;
  onProductShare?: (productId: string, productData: any) => void;
  onShopShare?: (shopId: string, shopData: any) => void;
  // Stream-related props
  onStreamContentClick?: (contentId: string) => void;
  onStreamContentShare?: (contentId: string, contentData: any) => void;
  cartItems?: string[];
  children?: React.ReactNode;
  defaultTab?: 'shops' | 'products' | 'stream' | 'live';
}

export function StoreLayout({
  user,
  onShopClick,
  onProductClick,
  onAddToCart,
  onLiveStreamClick,
  onProductShare,
  onShopShare,
  onStreamContentClick,
  onStreamContentShare,
  cartItems = [],
  children,
  defaultTab = 'shops'
}: StoreLayoutProps) {
  const [activeTab, setActiveTab] = useState(defaultTab);

  // Handle unified cart functionality
  const handleAddToCart = (itemId: string, quantity: number = 1, itemType: 'product' | 'stream' = 'product') => {
    if (onAddToCart) {
      onAddToCart(itemId, quantity);
    }
    // Could add specific handling for different item types here
  };

  return (
    <div className="flex-1 overflow-hidden">
      <Tabs value={activeTab} onValueChange={setActiveTab} className="h-full flex flex-col">
        {/* Unified Store Navigation */}
        <div className="mx-3 sm:mx-4 mt-3 sm:mt-4 flex justify-center">
          <TabsList className="inline-flex h-auto p-0.5 sm:p-1 bg-muted rounded-full">
            {/* Shops Tab */}
            <TabsTrigger
              value="shops"
              className="h-9 sm:h-10 px-3 sm:px-4 py-2 rounded-full data-[state=active]:bg-primary data-[state=active]:text-primary-foreground hover:bg-accent transition-colors touch-manipulation flex items-center gap-2"
            >
              <Store className="w-4 h-4 sm:w-5 sm:h-5" />
              <span className="text-sm">Shops</span>
            </TabsTrigger>

            {/* Products Tab */}
            <TabsTrigger
              value="products"
              className="h-9 sm:h-10 px-3 sm:px-4 py-2 rounded-full data-[state=active]:bg-primary data-[state=active]:text-primary-foreground hover:bg-accent transition-colors touch-manipulation flex items-center gap-2"
            >
              <Grid3X3 className="w-4 h-4 sm:w-5 sm:h-5" />
              <span className="text-sm">Products</span>
            </TabsTrigger>

            {/* Stream Tab - NEW: Positioned before Live tab */}
            <TabsTrigger
              value="stream"
              className="h-9 sm:h-10 px-3 sm:px-4 py-2 rounded-full data-[state=active]:bg-primary data-[state=active]:text-primary-foreground hover:bg-accent transition-colors touch-manipulation flex items-center gap-2"
            >
              <Play className="w-4 h-4 sm:w-5 sm:h-5" />
              <span className="text-sm">Stream</span>
            </TabsTrigger>

            {/* Live Tab */}
            <TabsTrigger
              value="live"
              className="h-9 sm:h-10 px-3 sm:px-4 py-2 rounded-full data-[state=active]:bg-primary data-[state=active]:text-primary-foreground hover:bg-accent transition-colors touch-manipulation flex items-center gap-2"
            >
              <Eye className="w-4 h-4 sm:w-5 sm:h-5" />
              <span className="text-sm">Live</span>
            </TabsTrigger>
          </TabsList>
        </div>

        {/* Tab Content */}
        {/* Shops Content */}
        <TabsContent value="shops" className="flex-1 overflow-hidden mt-4">
          <div className="h-full">
            {children ? (
              <div className="h-full">{children}</div>
            ) : (
              <div className="flex items-center justify-center h-full text-muted-foreground">
                <div className="text-center">
                  <Store className="w-16 h-16 mx-auto mb-4 opacity-50" />
                  <p>Shops content will be loaded here</p>
                  <p className="text-sm mt-2">Integrate with existing StoreTab shops content</p>
                </div>
              </div>
            )}
          </div>
        </TabsContent>

        {/* Products Content */}
        <TabsContent value="products" className="flex-1 overflow-hidden mt-4">
          <div className="h-full">
            {children ? (
              <div className="h-full">{children}</div>
            ) : (
              <div className="flex items-center justify-center h-full text-muted-foreground">
                <div className="text-center">
                  <Grid3X3 className="w-16 h-16 mx-auto mb-4 opacity-50" />
                  <p>Products content will be loaded here</p>
                  <p className="text-sm mt-2">Integrate with existing StoreTab products content</p>
                </div>
              </div>
            )}
          </div>
        </TabsContent>

        {/* Stream Content - NEW: Integrated Stream functionality */}
        <TabsContent value="stream" className="flex-1 overflow-hidden mt-4">
          <div className="h-full">
            <StreamTab
              user={user}
              onContentClick={onStreamContentClick}
              onAddToCart={handleAddToCart}
              onContentShare={onStreamContentShare}
              cartItems={cartItems}
            />
          </div>
        </TabsContent>

        {/* Live Content */}
        <TabsContent value="live" className="flex-1 overflow-hidden mt-4">
          <div className="h-full">
            {children ? (
              <div className="h-full">{children}</div>
            ) : (
              <div className="flex items-center justify-center h-full text-muted-foreground">
                <div className="text-center">
                  <Eye className="w-16 h-16 mx-auto mb-4 opacity-50" />
                  <p>Live streams content will be loaded here</p>
                  <p className="text-sm mt-2">Integrate with existing StoreTab live content</p>
                </div>
              </div>
            )}
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}

export default StoreLayout;