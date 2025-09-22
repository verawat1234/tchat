/**
 * Tests for Organism entity
 */

import { describe, it, expect, beforeEach } from 'vitest';
import { Organism } from '../Organism';
import { ComponentType } from '../Component';

describe('Organism', () => {
  let organism: Organism;

  beforeEach(() => {
    organism = new Organism({
      id: 'organism-header',
      name: 'Header',
      type: ComponentType.ORGANISM,
      filePath: '/src/components/organisms/Header.tsx',
      contains: ['Logo', 'Navigation', 'SearchBar', 'UserProfile'],
      standalone: true,
      pageSection: true,
      dataSource: 'user-api'
    });
  });

  describe('constructor', () => {
    it('should create an organism with required properties', () => {
      expect(organism).toBeInstanceOf(Organism);
      expect(organism.type).toBe(ComponentType.ORGANISM);
      expect(organism.name).toBe('Header');
      expect(organism.contains).toHaveLength(4);
    });

    it('should initialize with default values when not provided', () => {
      const minimalOrganism = new Organism({
        id: 'minimal-organism',
        name: 'MinimalOrganism',
        type: ComponentType.ORGANISM,
        filePath: '/src/components/organisms/Minimal.tsx'
      });

      expect(minimalOrganism.contains).toEqual([]);
      expect(minimalOrganism.standalone).toBe(false);
      expect(minimalOrganism.pageSection).toBe(false);
      expect(minimalOrganism.dataSource).toBeNull();
    });

    it('should validate component type is ORGANISM', () => {
      expect(() => {
        new Organism({
          id: 'wrong-type',
          name: 'WrongType',
          type: ComponentType.ATOM as any,
          filePath: '/src/components/wrong.tsx'
        });
      }).toThrow();
    });
  });

  describe('contains property', () => {
    it('should manage contained components', () => {
      expect(organism.contains).toEqual(['Logo', 'Navigation', 'SearchBar', 'UserProfile']);

      organism.contains.push('NotificationBell');
      expect(organism.contains).toContain('NotificationBell');
      expect(organism.contains.length).toBe(5);
    });

    it('should allow empty contains array', () => {
      organism.contains = [];
      expect(organism.contains).toEqual([]);
    });

    it('should track unique components', () => {
      organism.contains = ['Button', 'Card', 'Button', 'Form'];
      // Note: Currently allows duplicates, but could be enhanced
      expect(organism.contains.length).toBe(4);
    });
  });

  describe('standalone property', () => {
    it('should indicate if organism can function independently', () => {
      expect(organism.standalone).toBe(true);

      organism.standalone = false;
      expect(organism.standalone).toBe(false);
    });

    it('should differentiate standalone vs dependent organisms', () => {
      const dependentOrganism = new Organism({
        id: 'dependent-organism',
        name: 'TabContent',
        type: ComponentType.ORGANISM,
        filePath: '/src/components/organisms/TabContent.tsx',
        standalone: false
      });

      expect(dependentOrganism.standalone).toBe(false);
    });
  });

  describe('pageSection property', () => {
    it('should indicate if organism represents a page section', () => {
      expect(organism.pageSection).toBe(true);

      organism.pageSection = false;
      expect(organism.pageSection).toBe(false);
    });

    it('should identify major page sections', () => {
      const footerOrganism = new Organism({
        id: 'footer',
        name: 'Footer',
        type: ComponentType.ORGANISM,
        filePath: '/src/components/organisms/Footer.tsx',
        pageSection: true,
        standalone: true
      });

      expect(footerOrganism.pageSection).toBe(true);
      expect(footerOrganism.standalone).toBe(true);
    });
  });

  describe('dataSource property', () => {
    it('should track data source for organism', () => {
      expect(organism.dataSource).toBe('user-api');

      organism.dataSource = 'graphql-endpoint';
      expect(organism.dataSource).toBe('graphql-endpoint');
    });

    it('should allow null data source', () => {
      organism.dataSource = null;
      expect(organism.dataSource).toBeNull();
    });

    it('should identify data-driven organisms', () => {
      const dashboardOrganism = new Organism({
        id: 'dashboard',
        name: 'Dashboard',
        type: ComponentType.ORGANISM,
        filePath: '/src/components/organisms/Dashboard.tsx',
        dataSource: 'analytics-api',
        contains: ['Chart', 'DataTable', 'MetricCard']
      });

      expect(dashboardOrganism.dataSource).toBe('analytics-api');
      expect(dashboardOrganism.contains).toContain('Chart');
    });
  });

  describe('toJSON', () => {
    it('should serialize organism to JSON', () => {
      const json = organism.toJSON();

      expect(json).toMatchObject({
        id: 'organism-header',
        name: 'Header',
        type: ComponentType.ORGANISM,
        contains: ['Logo', 'Navigation', 'SearchBar', 'UserProfile'],
        standalone: true,
        pageSection: true,
        dataSource: 'user-api'
      });
    });

    it('should include inherited Component properties', () => {
      const json = organism.toJSON();

      expect(json).toHaveProperty('filePath');
      expect(json).toHaveProperty('category');
      expect(json).toHaveProperty('props');
      expect(json).toHaveProperty('dependencies');
      expect(json).toHaveProperty('version');
    });
  });

  describe('fromJSON static method', () => {
    it('should deserialize organism from JSON', () => {
      const json = organism.toJSON();
      const deserializedOrganism = Organism.fromJSON(json);

      expect(deserializedOrganism).toBeInstanceOf(Organism);
      expect(deserializedOrganism.id).toBe(organism.id);
      expect(deserializedOrganism.name).toBe(organism.name);
      expect(deserializedOrganism.contains).toEqual(organism.contains);
      expect(deserializedOrganism.standalone).toBe(organism.standalone);
      expect(deserializedOrganism.pageSection).toBe(organism.pageSection);
      expect(deserializedOrganism.dataSource).toBe(organism.dataSource);
    });

    it('should handle dates correctly in deserialization', () => {
      const json = organism.toJSON();
      const deserializedOrganism = Organism.fromJSON(json);

      expect(deserializedOrganism.createdAt).toBeInstanceOf(Date);
      expect(deserializedOrganism.updatedAt).toBeInstanceOf(Date);
    });
  });

  describe('organism-specific patterns', () => {
    it('should represent complex UI sections', () => {
      const heroOrganism = new Organism({
        id: 'hero-section',
        name: 'HeroSection',
        type: ComponentType.ORGANISM,
        filePath: '/src/components/organisms/HeroSection.tsx',
        contains: ['Heading', 'Subheading', 'CTAButton', 'BackgroundImage'],
        standalone: true,
        pageSection: true,
        dataSource: null
      });

      expect(heroOrganism.contains).toHaveLength(4);
      expect(heroOrganism.standalone).toBe(true);
      expect(heroOrganism.pageSection).toBe(true);
    });

    it('should support form organisms', () => {
      const checkoutForm = new Organism({
        id: 'checkout-form',
        name: 'CheckoutForm',
        type: ComponentType.ORGANISM,
        filePath: '/src/components/organisms/CheckoutForm.tsx',
        contains: ['ShippingForm', 'PaymentForm', 'OrderSummary', 'SubmitButton'],
        standalone: false,
        pageSection: false,
        dataSource: 'checkout-api'
      });

      expect(checkoutForm.contains).toContain('PaymentForm');
      expect(checkoutForm.standalone).toBe(false);
      expect(checkoutForm.dataSource).toBe('checkout-api');
    });

    it('should support layout organisms', () => {
      const layout = new Organism({
        id: 'main-layout',
        name: 'MainLayout',
        type: ComponentType.ORGANISM,
        filePath: '/src/components/organisms/MainLayout.tsx',
        contains: ['Header', 'Sidebar', 'MainContent', 'Footer'],
        standalone: true,
        pageSection: false,
        dataSource: null
      });

      expect(layout.contains).toContain('Header');
      expect(layout.contains).toContain('Footer');
      expect(layout.standalone).toBe(true);
      expect(layout.pageSection).toBe(false);
    });
  });
});