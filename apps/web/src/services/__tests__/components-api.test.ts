/**
 * Contract Test: GET /components API endpoint
 * This test MUST FAIL until the API service is implemented (TDD requirement)
 */
import { describe, it, expect } from 'vitest';
import { componentsApi } from '../componentsApi';

describe('Components API Contract Tests', () => {
  describe('GET /components', () => {
    it('should retrieve all component definitions with status', async () => {
      // Contract: API must return component list with platform status
      const result = await componentsApi.getComponents();

      expect(result).toHaveProperty('components');
      expect(result).toHaveProperty('total');
      expect(result).toHaveProperty('platforms');
      expect(Array.isArray(result.components)).toBe(true);

      // Platform status structure validation
      expect(result.platforms).toHaveProperty('web');
      expect(result.platforms).toHaveProperty('ios');
      expect(result.platforms).toHaveProperty('android');

      expect(result.platforms.web).toHaveProperty('total');
      expect(result.platforms.web).toHaveProperty('complete');
      expect(result.platforms.web).toHaveProperty('partial');
      expect(result.platforms.web).toHaveProperty('missing');

      // Component structure validation
      if (result.components.length > 0) {
        const component = result.components[0];
        expect(component).toHaveProperty('id');
        expect(component).toHaveProperty('name');
        expect(component).toHaveProperty('type');
        expect(component).toHaveProperty('platforms');
        expect(component).toHaveProperty('status');
        expect(['atom', 'molecule', 'organism']).toContain(component.type);
        expect(['missing', 'partial', 'complete']).toContain(component.status);
      }
    });

    it('should filter components by platform', async () => {
      const result = await componentsApi.getComponents({ platform: 'web' });

      expect(result.components).toBeDefined();
      // Each component should include 'web' in its platforms array
      result.components.forEach(component => {
        expect(component.platforms).toContain('web');
      });
    });

    it('should filter components by status', async () => {
      const result = await componentsApi.getComponents({ status: 'missing' });

      expect(result.components).toBeDefined();
      result.components.forEach(component => {
        expect(component.status).toBe('missing');
      });
    });
  });
});