/**
 * Tests for useNavigationContent hook
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook } from '@testing-library/react';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import { useNavigationContent, NAVIGATION_CONTENT_IDS } from './useNavigationContent';
import { contentApi } from '../services/contentApi';

// Mock the contentApi
vi.mock('../services/contentApi', () => ({
  useGetContentItemQuery: vi.fn(),
  useGetContentByCategoryQuery: vi.fn(),
}));

// Create a mock store
const createMockStore = () => configureStore({
  reducer: {
    [contentApi.reducerPath]: contentApi.reducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(contentApi.middleware),
});

describe('useNavigationContent', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should return fallback text when content is not available', () => {
    const mockUseGetContentItemQuery = vi.mocked(contentApi.useGetContentItemQuery);
    mockUseGetContentItemQuery.mockReturnValue({
      data: undefined,
      isLoading: false,
      isError: false,
    } as any);

    const store = createMockStore();

    const { result } = renderHook(
      () => useNavigationContent(NAVIGATION_CONTENT_IDS.TABS.CHAT),
      {
        wrapper: ({ children }) => <Provider store={store}>{children}</Provider>
      }
    );

    expect(result.current.text).toBe('Chat');
    expect(result.current.isFallback).toBe(true);
    expect(result.current.isLoading).toBe(false);
    expect(result.current.isError).toBe(false);
  });

  it('should return content text when available', () => {
    const mockUseGetContentItemQuery = vi.mocked(contentApi.useGetContentItemQuery);
    mockUseGetContentItemQuery.mockReturnValue({
      data: {
        id: NAVIGATION_CONTENT_IDS.TABS.CHAT,
        status: 'published',
        value: 'Dynamic Chat Text',
        category: 'navigation',
        type: 'text',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        version: 1,
        tags: [],
      },
      isLoading: false,
      isError: false,
    } as any);

    const store = createMockStore();

    const { result } = renderHook(
      () => useNavigationContent(NAVIGATION_CONTENT_IDS.TABS.CHAT),
      {
        wrapper: ({ children }) => <Provider store={store}>{children}</Provider>
      }
    );

    expect(result.current.text).toBe('Dynamic Chat Text');
    expect(result.current.isFallback).toBe(false);
    expect(result.current.isLoading).toBe(false);
    expect(result.current.isError).toBe(false);
  });

  it('should show loading state', () => {
    const mockUseGetContentItemQuery = vi.mocked(contentApi.useGetContentItemQuery);
    mockUseGetContentItemQuery.mockReturnValue({
      data: undefined,
      isLoading: true,
      isError: false,
    } as any);

    const store = createMockStore();

    const { result } = renderHook(
      () => useNavigationContent(NAVIGATION_CONTENT_IDS.TABS.CHAT),
      {
        wrapper: ({ children }) => <Provider store={store}>{children}</Provider>
      }
    );

    expect(result.current.isLoading).toBe(true);
    expect(result.current.text).toBe('Chat'); // Should still show fallback
  });

  it('should show error state', () => {
    const mockUseGetContentItemQuery = vi.mocked(contentApi.useGetContentItemQuery);
    mockUseGetContentItemQuery.mockReturnValue({
      data: undefined,
      isLoading: false,
      isError: true,
    } as any);

    const store = createMockStore();

    const { result } = renderHook(
      () => useNavigationContent(NAVIGATION_CONTENT_IDS.TABS.CHAT),
      {
        wrapper: ({ children }) => <Provider store={store}>{children}</Provider>
      }
    );

    expect(result.current.isError).toBe(true);
    expect(result.current.text).toBe('Chat'); // Should still show fallback
  });
});