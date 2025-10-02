import React, { useEffect, useState } from 'react';

interface Product {
  feature_id: string;
  product_id: string;
  display_order: number;
  start_time: string;
  end_time: string | null;
  is_active: boolean;
  view_count: number;
  click_count: number;
  // Product details from commerce service
  name?: string;
  price?: number;
  image_url?: string;
  availability?: 'in_stock' | 'low_stock' | 'out_of_stock';
}

interface ProductOverlayProps {
  streamId: string;
  products: Product[];
  onProductClick?: (productId: string) => void;
}

export const ProductOverlay: React.FC<ProductOverlayProps> = ({
  streamId,
  products,
  onProductClick,
}) => {
  const [activeProducts, setActiveProducts] = useState<Product[]>([]);
  const [currentProductIndex, setCurrentProductIndex] = useState(0);

  // Filter active products
  useEffect(() => {
    const active = products.filter((p) => p.is_active);
    setActiveProducts(active);
    setCurrentProductIndex(0);
  }, [products]);

  // Cycle through products every 10 seconds
  useEffect(() => {
    if (activeProducts.length <= 1) return;

    const interval = setInterval(() => {
      setCurrentProductIndex((prev) => (prev + 1) % activeProducts.length);
    }, 10000);

    return () => clearInterval(interval);
  }, [activeProducts.length]);

  // Track view event
  useEffect(() => {
    if (activeProducts.length === 0) return;

    const currentProduct = activeProducts[currentProductIndex];
    if (!currentProduct) return;

    // Track view
    trackProductView(streamId, currentProduct.feature_id);
  }, [currentProductIndex, activeProducts, streamId]);

  // Track product view
  const trackProductView = async (streamId: string, featureId: string) => {
    try {
      // Analytics tracking would go here
      console.log(`[ProductOverlay] Viewed product ${featureId}`);
    } catch (error) {
      console.error('[ProductOverlay] Failed to track view:', error);
    }
  };

  // Handle product click
  const handleProductClick = async (product: Product) => {
    // Track click event
    try {
      console.log(`[ProductOverlay] Clicked product ${product.feature_id}`);
      // Analytics tracking would go here
    } catch (error) {
      console.error('[ProductOverlay] Failed to track click:', error);
    }

    // Callback for parent component
    if (onProductClick) {
      onProductClick(product.product_id);
    }
  };

  // Handle add to cart
  const handleAddToCart = async (product: Product, e: React.MouseEvent) => {
    e.stopPropagation();

    try {
      // Commerce service integration would go here
      const response = await fetch('/api/v1/cart/items', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${localStorage.getItem('auth_token')}`,
        },
        body: JSON.stringify({
          product_id: product.product_id,
          quantity: 1,
        }),
      });

      if (response.ok) {
        // Show success notification
        console.log(`[ProductOverlay] Added ${product.name} to cart`);
      }
    } catch (error) {
      console.error('[ProductOverlay] Failed to add to cart:', error);
    }
  };

  if (activeProducts.length === 0) {
    return null;
  }

  const currentProduct = activeProducts[currentProductIndex];

  return (
    <div className="product-overlay">
      {/* Product Card */}
      <div
        onClick={() => handleProductClick(currentProduct)}
        className="bg-white rounded-lg shadow-lg overflow-hidden cursor-pointer transform transition-all hover:scale-105"
      >
        {/* Product Image */}
        {currentProduct.image_url && (
          <div className="relative h-48 bg-gray-200">
            <img
              src={currentProduct.image_url}
              alt={currentProduct.name || 'Product'}
              className="w-full h-full object-cover"
            />
            {/* Availability Badge */}
            {currentProduct.availability === 'low_stock' && (
              <div className="absolute top-2 right-2 bg-amber-500 text-white text-xs font-medium px-2 py-1 rounded">
                Low Stock
              </div>
            )}
            {currentProduct.availability === 'out_of_stock' && (
              <div className="absolute top-2 right-2 bg-red-500 text-white text-xs font-medium px-2 py-1 rounded">
                Out of Stock
              </div>
            )}
          </div>
        )}

        {/* Product Details */}
        <div className="p-4">
          <h4 className="text-lg font-semibold text-gray-900 line-clamp-2">
            {currentProduct.name || 'Featured Product'}
          </h4>

          {/* Price */}
          {currentProduct.price && (
            <div className="mt-2 text-2xl font-bold text-blue-600">
              ${currentProduct.price.toFixed(2)}
            </div>
          )}

          {/* Add to Cart Button */}
          <button
            onClick={(e) => handleAddToCart(currentProduct, e)}
            disabled={currentProduct.availability === 'out_of_stock'}
            className="mt-3 w-full px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors font-medium"
          >
            {currentProduct.availability === 'out_of_stock' ? 'Out of Stock' : 'Add to Cart'}
          </button>

          {/* Analytics Info */}
          <div className="mt-3 flex items-center justify-between text-xs text-gray-500">
            <span>{currentProduct.view_count} views</span>
            <span>{currentProduct.click_count} clicks</span>
          </div>
        </div>

        {/* Product Indicator (if multiple) */}
        {activeProducts.length > 1 && (
          <div className="px-4 pb-3 flex items-center justify-center space-x-1">
            {activeProducts.map((_, index) => (
              <div
                key={index}
                className={`w-2 h-2 rounded-full transition-colors ${
                  index === currentProductIndex ? 'bg-blue-600' : 'bg-gray-300'
                }`}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
};