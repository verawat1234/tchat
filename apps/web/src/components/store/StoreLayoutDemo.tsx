import React from 'react';
import StoreLayout from './StoreLayout';

// Demo component showing how to integrate StoreLayout with existing store functionality
interface StoreLayoutDemoProps {
  user: any;
}

export function StoreLayoutDemo({ user }: StoreLayoutDemoProps) {
  // Example handlers for store functionality
  const handleShopClick = (shopId: string) => {
    console.log('Shop clicked:', shopId);
  };

  const handleProductClick = (productId: string) => {
    console.log('Product clicked:', productId);
  };

  const handleAddToCart = (itemId: string, quantity?: number) => {
    console.log('Add to cart:', itemId, 'quantity:', quantity);
  };

  const handleLiveStreamClick = (streamId: string) => {
    console.log('Live stream clicked:', streamId);
  };

  const handleProductShare = (productId: string, productData: any) => {
    console.log('Product shared:', productId, productData);
  };

  const handleShopShare = (shopId: string, shopData: any) => {
    console.log('Shop shared:', shopId, shopData);
  };

  const handleStreamContentClick = (contentId: string) => {
    console.log('Stream content clicked:', contentId);
  };

  const handleStreamContentShare = (contentId: string, contentData: any) => {
    console.log('Stream content shared:', contentId, contentData);
  };

  return (
    <div className="h-full flex flex-col">
      {/* Header */}
      <div className="p-4 border-b">
        <h1 className="text-2xl font-bold">Store & Stream</h1>
        <p className="text-muted-foreground">Unified shopping and streaming experience</p>
      </div>

      {/* Unified Store Layout with Stream Integration */}
      <StoreLayout
        user={user}
        onShopClick={handleShopClick}
        onProductClick={handleProductClick}
        onAddToCart={handleAddToCart}
        onLiveStreamClick={handleLiveStreamClick}
        onProductShare={handleProductShare}
        onShopShare={handleShopShare}
        onStreamContentClick={handleStreamContentClick}
        onStreamContentShare={handleStreamContentShare}
        cartItems={[]} // Example cart items
        defaultTab="stream" // Start with stream tab for demo
      />
    </div>
  );
}

export default StoreLayoutDemo;