import { describe, it, expect, beforeEach } from 'vitest';
import { setupServer } from 'msw/node';
import { http, HttpResponse } from 'msw';
import type { DesignToken, DesignTokenCreateRequest, DesignTokenResponse } from '../../tools/design-tokens/src/types';

/**
 * Contract Test: Design Token API POST /design-tokens
 *
 * CRITICAL TDD: This test MUST FAIL initially because:
 * 1. DesignToken types don't exist yet
 * 2. API endpoint doesn't exist yet
 * 3. Color conversion service doesn't exist yet
 *
 * These tests drive the implementation of the design token system.
 */

const server = setupServer();

describe('Design Token API - POST /design-tokens (Contract)', () => {
  beforeEach(() => {
    server.listen({ onUnhandledRequest: 'error' });
  });

  afterEach(() => {
    server.resetHandlers();
  });

  afterAll(() => {
    server.close();
  });

  it('should create a new color design token with oklch values', async () => {
    // EXPECTED TO FAIL: DesignTokenCreateRequest type doesn't exist
    const tokenRequest: DesignTokenCreateRequest = {
      name: 'primary-500',
      category: 'color',
      value: {
        oklch: {
          l: 0.6,
          c: 0.2,
          h: 250
        }
      },
      platforms: ['web', 'ios'],
      metadata: {
        description: 'Primary brand color at 500 level',
        usage: 'buttons, links, primary actions'
      }
    };

    // EXPECTED TO FAIL: API endpoint doesn't exist
    server.use(
      http.post('/api/design-tokens', async ({ request }) => {
        const body = await request.json() as DesignTokenCreateRequest;

        // Contract validation
        expect(body.name).toBe('primary-500');
        expect(body.category).toBe('color');
        expect(body.value.oklch).toBeDefined();
        expect(body.platforms).toContain('web');
        expect(body.platforms).toContain('ios');

        const response: DesignTokenResponse = {
          id: '550e8400-e29b-41d4-a716-446655440000',
          ...body,
          generatedValues: {
            web: {
              css: '--color-primary-500: oklch(60% 0.2 250);',
              hex: '#6366f1'
            },
            ios: {
              swift: 'static let primary500 = Color(hex: "#6366f1")',
              uiColor: 'UIColor(red: 0.388, green: 0.4, blue: 0.945, alpha: 1.0)'
            }
          },
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z'
        };

        return HttpResponse.json(response, { status: 201 });
      })
    );

    // EXPECTED TO FAIL: Design token API client doesn't exist
    const response = await fetch('/api/design-tokens', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(tokenRequest)
    });

    expect(response.status).toBe(201);

    const result: DesignTokenResponse = await response.json();
    expect(result.id).toBeDefined();
    expect(result.name).toBe('primary-500');
    expect(result.generatedValues.web.hex).toBe('#6366f1');
    expect(result.generatedValues.ios.swift).toContain('Color(hex: "#6366f1")');
  });

  it('should create a spacing design token', async () => {
    // EXPECTED TO FAIL: Spacing token types don't exist
    const spacingToken: DesignTokenCreateRequest = {
      name: 'spacing-md',
      category: 'spacing',
      value: {
        base: 16,
        scale: {
          xs: 8,
          sm: 12,
          md: 16,
          lg: 24,
          xl: 32
        }
      },
      platforms: ['web', 'ios'],
      metadata: {
        description: 'Medium spacing value',
        usage: 'standard component padding'
      }
    };

    server.use(
      http.post('/api/design-tokens', async ({ request }) => {
        const body = await request.json() as DesignTokenCreateRequest;

        const response: DesignTokenResponse = {
          id: '550e8400-e29b-41d4-a716-446655440001',
          ...body,
          generatedValues: {
            web: {
              css: '--spacing-md: 16px;'
            },
            ios: {
              swift: 'static let spacingMd: CGFloat = 16'
            }
          },
          createdAt: '2024-01-01T00:00:00Z',
          updatedAt: '2024-01-01T00:00:00Z'
        };

        return HttpResponse.json(response, { status: 201 });
      })
    );

    const response = await fetch('/api/design-tokens', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(spacingToken)
    });

    expect(response.status).toBe(201);

    const result: DesignTokenResponse = await response.json();
    expect(result.generatedValues.web.css).toContain('--spacing-md: 16px;');
    expect(result.generatedValues.ios.swift).toContain('static let spacingMd: CGFloat = 16');
  });

  it('should validate required fields', async () => {
    // EXPECTED TO FAIL: Validation doesn't exist yet
    const invalidToken = {
      // Missing required name field
      category: 'color',
      value: { oklch: { l: 0.6, c: 0.2, h: 250 } }
    };

    server.use(
      http.post('/api/design-tokens', async () => {
        return HttpResponse.json(
          { error: 'Validation failed', details: ['name is required'] },
          { status: 400 }
        );
      })
    );

    const response = await fetch('/api/design-tokens', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(invalidToken)
    });

    expect(response.status).toBe(400);

    const error = await response.json();
    expect(error.error).toBe('Validation failed');
    expect(error.details).toContain('name is required');
  });

  it('should handle color conversion errors gracefully', async () => {
    // EXPECTED TO FAIL: Color conversion service doesn't exist
    const invalidColorToken: DesignTokenCreateRequest = {
      name: 'invalid-color',
      category: 'color',
      value: {
        oklch: {
          l: -0.1, // Invalid lightness (negative)
          c: 0.2,
          h: 250
        }
      },
      platforms: ['web', 'ios']
    };

    server.use(
      http.post('/api/design-tokens', async () => {
        return HttpResponse.json(
          {
            error: 'Color conversion failed',
            details: ['Invalid OKLCH lightness value: must be between 0 and 1']
          },
          { status: 422 }
        );
      })
    );

    const response = await fetch('/api/design-tokens', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(invalidColorToken)
    });

    expect(response.status).toBe(422);

    const error = await response.json();
    expect(error.error).toBe('Color conversion failed');
    expect(error.details[0]).toContain('Invalid OKLCH lightness value');
  });
});