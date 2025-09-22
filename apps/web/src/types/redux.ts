import type { Action, ThunkAction } from '@reduxjs/toolkit';
import type { TypedUseSelectorHook } from 'react-redux';

/**
 * Redux Store Types
 *
 * This file defines the core TypeScript types for Redux configuration
 * using Redux Toolkit and React Redux v9.2.0.
 */

/**
 * Root state interface - to be extended when store is configured
 * This will be overridden with the actual RootState type from the store
 */
export interface RootState {
  // This will be populated by the actual store configuration
}

/**
 * App dispatch type - to be extended when store is configured
 * This will be overridden with the actual AppDispatch type from the store
 */
export type AppDispatch = any;

/**
 * Typed hooks for use throughout the application
 * These provide type safety when using Redux hooks
 */
export type AppThunk<ReturnType = void> = ThunkAction<
  ReturnType,
  RootState,
  unknown,
  Action<string>
>;

/**
 * Typed useSelector hook
 * Use this instead of the plain useSelector from react-redux
 */
export type TypedUseSelector = TypedUseSelectorHook<RootState>;

/**
 * Common Redux action payload types
 */
export interface BaseAction<T = any> {
  type: string;
  payload?: T;
  error?: boolean;
  meta?: any;
}

/**
 * Generic async action states for handling loading states
 */
export interface AsyncState<T = any> {
  data: T | null;
  loading: boolean;
  error: string | null;
}

/**
 * Generic collection state for handling lists/arrays
 */
export interface CollectionState<T = any> extends AsyncState<T[]> {
  selectedItems: T[];
  filters: Record<string, any>;
  pagination: {
    page: number;
    limit: number;
    total: number;
  };
}

/**
 * Entity state for normalized data structures
 */
export interface EntityState<T = any> {
  entities: Record<string | number, T>;
  ids: (string | number)[];
  loading: boolean;
  error: string | null;
}

/**
 * Standard API response wrapper
 */
export interface ApiResponse<T = any> {
  data: T;
  message?: string;
  success: boolean;
  errors?: string[];
}

/**
 * Standard error state
 */
export interface ErrorState {
  message: string;
  code?: string | number;
  details?: any;
}