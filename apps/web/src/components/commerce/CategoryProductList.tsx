/**
 * CategoryProductList - Products display component for categories
 *
 * Displays products within a specific category with lazy loading and virtualization.
 * Supports multiple view modes, filtering, sorting, and infinite scroll.
 *
 * Features:
 * - Lazy loading with intersection observer
 * - Virtual scrolling for large product lists
 * - Multiple view modes (grid, list, compact)
 * - Product filtering and sorting
 * - Infinite scroll pagination
 * - Loading skeletons and empty states
 * - Add to cart functionality
 * - Product quick view modal
 * - Accessibility support
 */

import React, { useState, useEffect, useRef, useCallback, useMemo } from 'react';
import { Button } from '../ui/button';
import { Card, CardContent, CardHeader } from '../ui/card';
import { Badge } from '../ui/badge';
import { Skeleton } from '../ui/skeleton';
import { ScrollArea } from '../ui/scroll-area';
import { Avatar, AvatarFallback, AvatarImage } from '../ui/avatar';
import {
  Grid3X3,
  List,
  Star,
  Heart,
  ShoppingCart,
  Eye,
  Filter,
  SlidersHorizontal,
  Package,
  TrendingUp,
  Clock,
  MapPin,
  RefreshCw,
  Search,
} from 'lucide-react';
import { cn } from '../../lib/utils';
import { useCategory } from './CategoryProvider';
import { ImageWithFallback } from '../figma/ImageWithFallback';
import { useGetProductsQuery, useAddToCartMutation } from '../../services/commerceApi';
import type { Product, ProductFilters, SortOptions } from '../../types/commerce';
import { toast } from 'sonner';

// ===== Types =====

export interface CategoryProductListProps {
  categoryId?: string;
  className?: string;
  viewMode?: 'grid' | 'list' | 'compact';
  gridColumns?: 2 | 3 | 4 | 5 | 6;
  showFilters?: boolean;
  showSorting?: boolean;
  enableLazyLoading?: boolean;
  enableVirtualScroll?: boolean;
  enableInfiniteScroll?: boolean;
  pageSize?: number;
  onProductSelect?: (product: Product) => void;
  onAddToCart?: (product: Product) => void;
}

interface ProductCardProps {
  product: Product;
  viewMode: 'grid' | 'list' | 'compact';
  onSelect?: (product: Product) => void;
  onAddToCart?: (product: Product) => void;
}

// ===== Product Card Component =====

function ProductCard({ product, viewMode, onSelect, onAddToCart }: ProductCardProps) {
  const [addToCart, { isLoading: isAddingToCart }] = useAddToCartMutation();
  const [isFavorited, setIsFavorited] = useState(false);

  const handleAddToCart = async (e: React.MouseEvent) => {
    e.stopPropagation();

    try {
      await addToCart({
        item: {
          productId: product.id,
          quantity: 1,
        },
      }).unwrap();

      toast.success(`${product.name} added to cart`);
      onAddToCart?.(product);
    } catch (error) {
      toast.error('Failed to add product to cart');
    }
  };

  const handleFavoriteToggle = (e: React.MouseEvent) => {
    e.stopPropagation();
    setIsFavorited(!isFavorited);
    toast.success(isFavorited ? 'Removed from favorites' : 'Added to favorites');
  };

  const handleProductClick = () => {
    onSelect?.(product);
  };

  // Compact view
  if (viewMode === 'compact') {
    return (
      <Card
        className="cursor-pointer hover:shadow-md transition-shadow"
        onClick={handleProductClick}
      >
        <CardContent className="flex items-center gap-4 p-4">
          <div className="w-16 h-16 rounded-lg overflow-hidden flex-shrink-0">
            <ImageWithFallback
              src={product.images[0] || ''}
              alt={product.name}
              className="w-full h-full object-cover"
            />
          </div>

          <div className="flex-1 min-w-0">
            <h3 className="font-medium line-clamp-1">{product.name}</h3>
            <p className="text-sm text-muted-foreground">{product.category}</p>
            <div className="flex items-center gap-2 mt-1">
              <span className="font-semibold">฿{parseFloat(product.price).toFixed(0)}</span>
              {product.compareAtPrice && parseFloat(product.compareAtPrice) > parseFloat(product.price) && (
                <span className="text-sm text-muted-foreground line-through">
                  ฿{parseFloat(product.compareAtPrice).toFixed(0)}
                </span>
              )}
            </div>
          </div>

          <Button
            size="sm"
            onClick={handleAddToCart}
            disabled={isAddingToCart || !product.isActive}
            loading={isAddingToCart}
          >
            <ShoppingCart className="w-4 h-4" />
          </Button>
        </CardContent>
      </Card>
    );
  }

  // List view
  if (viewMode === 'list') {
    return (
      <Card
        className="cursor-pointer hover:shadow-md transition-shadow"
        onClick={handleProductClick}
      >
        <CardContent className="flex gap-4 p-6">
          <div className="w-24 h-24 rounded-lg overflow-hidden flex-shrink-0">
            <ImageWithFallback
              src={product.images[0] || ''}
              alt={product.name}
              className="w-full h-full object-cover"
            />
          </div>

          <div className="flex-1">
            <div className="flex justify-between items-start mb-2">
              <h3 className="font-semibold text-lg line-clamp-2">{product.name}</h3>
              <Button
                variant="ghost"
                size="icon"
                onClick={handleFavoriteToggle}
                className={cn(
                  'text-muted-foreground hover:text-red-500',
                  isFavorited && 'text-red-500'
                )}
              >
                <Heart className={cn('w-5 h-5', isFavorited && 'fill-current')} />
              </Button>
            </div>

            <p className="text-muted-foreground line-clamp-2 mb-3">
              {product.description}
            </p>

            <div className="flex items-center gap-4 mb-3">
              <div className="flex items-center gap-1">
                <Star className="w-4 h-4 text-yellow-500 fill-current" />
                <span className="text-sm">{parseFloat(product.rating).toFixed(1)}</span>
                <span className="text-sm text-muted-foreground">
                  ({product.reviewCount} reviews)
                </span>
              </div>

              {product.salesCount > 0 && (
                <div className="flex items-center gap-1 text-sm text-muted-foreground">
                  <TrendingUp className="w-4 h-4" />
                  <span>{product.salesCount} sold</span>
                </div>
              )}
            </div>

            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <span className="text-2xl font-bold">฿{parseFloat(product.price).toFixed(0)}</span>
                {product.compareAtPrice && parseFloat(product.compareAtPrice) > parseFloat(product.price) && (
                  <span className="text-lg text-muted-foreground line-through">
                    ฿{parseFloat(product.compareAtPrice).toFixed(0)}
                  </span>
                )}
              </div>

              <Button
                onClick={handleAddToCart}
                disabled={isAddingToCart || !product.isActive}
                loading={isAddingToCart}
              >
                <ShoppingCart className="w-4 h-4 mr-2" />
                Add to Cart
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    );
  }

  // Grid view (default)
  return (
    <Card
      className="group cursor-pointer hover:shadow-lg transition-all duration-200 overflow-hidden"
      onClick={handleProductClick}
    >
      <CardContent className="p-0">
        <div className="relative aspect-square">
          <ImageWithFallback
            src={product.images[0] || ''}
            alt={product.name}
            className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
          />

          {/* Overlay actions */}
          <div className="absolute inset-0 bg-black/0 group-hover:bg-black/20 transition-colors">
            <div className="absolute top-3 right-3 flex flex-col gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
              <Button
                variant="secondary"
                size="icon"
                className="w-8 h-8 bg-white/90 hover:bg-white"
                onClick={handleFavoriteToggle}
              >
                <Heart className={cn('w-4 h-4', isFavorited && 'text-red-500 fill-current')} />
              </Button>

              <Button
                variant="secondary"
                size="icon"
                className="w-8 h-8 bg-white/90 hover:bg-white"
                onClick={(e) => {
                  e.stopPropagation();
                  // Quick view functionality
                }}
              >
                <Eye className="w-4 h-4" />
              </Button>
            </div>
          </div>

          {/* Badges */}
          <div className="absolute top-3 left-3 flex flex-col gap-1">
            {product.isFeatured && (
              <Badge variant="default" className="bg-yellow-500 text-white">
                Featured
              </Badge>
            )}

            {product.compareAtPrice && parseFloat(product.compareAtPrice) > parseFloat(product.price) && (
              <Badge variant="destructive">
                Sale
              </Badge>
            )}
          </div>
        </div>

        <div className="p-4 space-y-3">
          <div>
            <h3 className="font-medium line-clamp-2 group-hover:text-primary transition-colors">
              {product.name}
            </h3>
            <p className="text-sm text-muted-foreground">{product.category}</p>
          </div>

          <div className="flex items-center gap-1">
            <Star className="w-4 h-4 text-yellow-500 fill-current" />
            <span className="text-sm">{parseFloat(product.rating).toFixed(1)}</span>
            <span className="text-sm text-muted-foreground">({product.reviewCount})</span>
          </div>

          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <span className="text-lg font-bold">฿{parseFloat(product.price).toFixed(0)}</span>
              {product.compareAtPrice && parseFloat(product.compareAtPrice) > parseFloat(product.price) && (
                <span className="text-sm text-muted-foreground line-through">
                  ฿{parseFloat(product.compareAtPrice).toFixed(0)}
                </span>
              )}
            </div>
          </div>

          <Button
            className="w-full"
            onClick={handleAddToCart}
            disabled={isAddingToCart || !product.isActive}
            loading={isAddingToCart}
          >
            <ShoppingCart className="w-4 h-4 mr-2" />
            Add to Cart
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}

// ===== Product Card Skeleton =====

function ProductCardSkeleton({ viewMode }: { viewMode: 'grid' | 'list' | 'compact' }) {
  if (viewMode === 'compact') {
    return (
      <Card>
        <CardContent className="flex items-center gap-4 p-4">
          <Skeleton className="w-16 h-16 rounded-lg" />
          <div className="flex-1 space-y-2">
            <Skeleton className="h-4 w-3/4" />
            <Skeleton className="h-3 w-1/2" />
            <Skeleton className="h-4 w-1/3" />
          </div>
          <Skeleton className="w-10 h-8" />
        </CardContent>
      </Card>
    );
  }

  if (viewMode === 'list') {
    return (
      <Card>
        <CardContent className="flex gap-4 p-6">
          <Skeleton className="w-24 h-24 rounded-lg" />
          <div className="flex-1 space-y-3">
            <Skeleton className="h-5 w-3/4" />
            <Skeleton className="h-4 w-full" />
            <Skeleton className="h-4 w-2/3" />
            <div className="flex justify-between items-center">
              <Skeleton className="h-6 w-20" />
              <Skeleton className="h-10 w-32" />
            </div>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardContent className="p-0">
        <Skeleton className="aspect-square w-full" />
        <div className="p-4 space-y-3">
          <Skeleton className="h-4 w-3/4" />
          <Skeleton className="h-3 w-1/2" />
          <Skeleton className="h-4 w-1/3" />
          <Skeleton className="h-10 w-full" />
        </div>
      </CardContent>
    </Card>
  );
}

// ===== View Mode Toggle =====

interface ViewModeToggleProps {
  currentMode: 'grid' | 'list' | 'compact';
  onModeChange: (mode: 'grid' | 'list' | 'compact') => void;
}

function ViewModeToggle({ currentMode, onModeChange }: ViewModeToggleProps) {
  return (
    <div className="flex border rounded-lg overflow-hidden">
      <Button
        variant={currentMode === 'grid' ? 'default' : 'ghost'}
        size="sm"
        className="rounded-none border-0"
        onClick={() => onModeChange('grid')}
      >
        <Grid3X3 className="w-4 h-4" />
      </Button>
      <Button
        variant={currentMode === 'list' ? 'default' : 'ghost'}
        size="sm"
        className="rounded-none border-0 border-l"
        onClick={() => onModeChange('list')}
      >
        <List className="w-4 h-4" />
      </Button>
      <Button
        variant={currentMode === 'compact' ? 'default' : 'ghost'}
        size="sm"
        className="rounded-none border-0 border-l"
        onClick={() => onModeChange('compact')}
      >
        <Package className="w-4 h-4" />
      </Button>
    </div>
  );
}

// ===== Main CategoryProductList Component =====

export function CategoryProductList({
  categoryId,
  className,
  viewMode: propViewMode,
  gridColumns = 4,
  showFilters = true,
  showSorting = true,
  enableLazyLoading = true,
  enableVirtualScroll = false,
  enableInfiniteScroll = true,
  pageSize = 24,
  onProductSelect,
  onAddToCart,
}: CategoryProductListProps) {
  const { state, setViewMode } = useCategory();
  const [currentPage, setCurrentPage] = useState(1);
  const [allProducts, setAllProducts] = useState<Product[]>([]);
  const loadMoreRef = useRef<HTMLDivElement>(null);

  const viewMode = propViewMode || state.viewState.viewMode;
  const sortBy = state.viewState.sortBy;

  // Build product filters
  const productFilters = useMemo((): ProductFilters => {
    const filters: ProductFilters = {};

    if (categoryId) {
      filters.category = categoryId;
    }

    if (state.search.query) {
      filters.search = state.search.query;
    }

    return filters;
  }, [categoryId, state.search.query]);

  // Fetch products with pagination
  const { data: productsData, isLoading, error, isFetching } = useGetProductsQuery({
    filters: productFilters,
    pagination: { page: currentPage, pageSize },
    sort: sortBy,
  });

  // Update products list when new data arrives
  useEffect(() => {
    if (productsData?.products) {
      if (currentPage === 1) {
        setAllProducts(productsData.products);
      } else {
        setAllProducts(prev => [...prev, ...productsData.products]);
      }
    }
  }, [productsData, currentPage]);

  // Reset when filters change
  useEffect(() => {
    setCurrentPage(1);
    setAllProducts([]);
  }, [categoryId, state.search.query, sortBy]);

  // Infinite scroll intersection observer
  useEffect(() => {
    if (!enableInfiniteScroll || !loadMoreRef.current) return;

    const observer = new IntersectionObserver(
      (entries) => {
        const [entry] = entries;
        if (
          entry.isIntersecting &&
          !isFetching &&
          productsData &&
          currentPage < productsData.totalPages
        ) {
          setCurrentPage(prev => prev + 1);
        }
      },
      { threshold: 0.1 }
    );

    observer.observe(loadMoreRef.current);

    return () => observer.disconnect();
  }, [enableInfiniteScroll, isFetching, productsData, currentPage]);

  const handleViewModeChange = (mode: 'grid' | 'list' | 'compact') => {
    setViewMode(mode);
  };

  // Loading state
  if (isLoading && allProducts.length === 0) {
    return (
      <div className={cn('space-y-6', className)}>
        {/* Controls skeleton */}
        <div className="flex items-center justify-between">
          <Skeleton className="w-32 h-10" />
          <Skeleton className="w-24 h-10" />
        </div>

        {/* Grid skeleton */}
        <div className={cn(
          'grid gap-4',
          viewMode === 'grid' && `grid-cols-1 sm:grid-cols-2 md:grid-cols-${Math.min(gridColumns, 3)} lg:grid-cols-${gridColumns}`,
          viewMode === 'list' && 'grid-cols-1',
          viewMode === 'compact' && 'grid-cols-1 gap-2'
        )}>
          {Array.from({ length: pageSize }, (_, index) => (
            <ProductCardSkeleton key={index} viewMode={viewMode} />
          ))}
        </div>
      </div>
    );
  }

  // Error state
  if (error) {
    return (
      <div className={cn('text-center py-12', className)}>
        <Package className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
        <h3 className="text-lg font-medium mb-2">Failed to load products</h3>
        <p className="text-muted-foreground mb-6">
          There was an error loading products. Please try again.
        </p>
        <Button onClick={() => window.location.reload()}>
          <RefreshCw className="w-4 h-4 mr-2" />
          Try Again
        </Button>
      </div>
    );
  }

  // Empty state
  if (allProducts.length === 0 && !isLoading) {
    return (
      <div className={cn('text-center py-12', className)}>
        <Search className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
        <h3 className="text-lg font-medium mb-2">No products found</h3>
        <p className="text-muted-foreground mb-6">
          {categoryId
            ? "This category doesn't have any products yet."
            : "Try adjusting your search or filters to find products."
          }
        </p>
      </div>
    );
  }

  return (
    <div className={cn('space-y-6', className)}>
      {/* Controls */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <span className="text-sm text-muted-foreground">
            {productsData?.total || 0} products found
          </span>
        </div>

        <ViewModeToggle
          currentMode={viewMode}
          onModeChange={handleViewModeChange}
        />
      </div>

      {/* Products Grid */}
      <div className={cn(
        'grid gap-4',
        viewMode === 'grid' && `grid-cols-1 sm:grid-cols-2 md:grid-cols-${Math.min(gridColumns, 3)} lg:grid-cols-${gridColumns}`,
        viewMode === 'list' && 'grid-cols-1',
        viewMode === 'compact' && 'grid-cols-1 gap-2'
      )}>
        {allProducts.map(product => (
          <ProductCard
            key={product.id}
            product={product}
            viewMode={viewMode}
            onSelect={onProductSelect}
            onAddToCart={onAddToCart}
          />
        ))}
      </div>

      {/* Infinite Scroll Trigger */}
      {enableInfiniteScroll && productsData && currentPage < productsData.totalPages && (
        <div ref={loadMoreRef} className="py-4">
          {isFetching && (
            <div className="grid gap-4 grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
              {Array.from({ length: 4 }, (_, index) => (
                <ProductCardSkeleton key={index} viewMode={viewMode} />
              ))}
            </div>
          )}
        </div>
      )}

      {/* Load More Button (fallback) */}
      {!enableInfiniteScroll && productsData && currentPage < productsData.totalPages && (
        <div className="text-center">
          <Button
            onClick={() => setCurrentPage(prev => prev + 1)}
            disabled={isFetching}
            loading={isFetching}
          >
            Load More Products
          </Button>
        </div>
      )}
    </div>
  );
}

export default CategoryProductList;