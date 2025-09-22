/**
 * @file content-category.test.ts
 * @description Test suite for enhanced ContentCategory implementation
 */

import { describe, expect, it } from 'vitest';
import {
  PermissionAction,
  createContentCategory,
  createCategoryPermissions,
  computeCategoryHierarchy,
  getCategoryParents,
  getCategoryChildren,
  hasPermission,
  getEffectivePermissions,
  validateCategoryHierarchy,
  validateCategoryMove,
  isEnhancedContentCategory,
  type EnhancedContentCategory,
  type CategoryAccessContext,
} from '../content-category';

describe('Enhanced ContentCategory Implementation', () => {
  // Test data setup
  const mockUser = 'user123';
  const mockAdmin = 'admin123';

  describe('Factory Functions', () => {
    it('should create a basic content category', () => {
      const category = createContentCategory({
        id: 'test-category',
        name: 'Test Category',
        description: 'A test category',
        createdBy: mockUser,
      });

      expect(category.id).toBe('test-category');
      expect(category.name).toBe('Test Category');
      expect(category.description).toBe('A test category');
      expect(category.createdBy).toBe(mockUser);
      expect(category.version).toBe(1);
      expect(category.isActive).toBe(true);
      expect(category.hierarchy.isRoot).toBe(true);
      expect(category.hierarchy.level).toBe(0);
    });

    it('should create a child category', () => {
      const parentCategory = createContentCategory({
        id: 'parent',
        name: 'Parent Category',
        description: 'Parent category',
        createdBy: mockUser,
      });

      const childCategory = createContentCategory({
        id: 'child',
        name: 'Child Category',
        description: 'Child category',
        parentId: 'parent',
        createdBy: mockUser,
      });

      expect(childCategory.parentId).toBe('parent');
      expect(childCategory.hierarchy.isRoot).toBe(false);
    });

    it('should create custom permissions', () => {
      const permissions = createCategoryPermissions({
        readRoles: ['user', 'admin'],
        writeRoles: ['admin'],
        publishRoles: ['admin'],
        deleteRoles: ['admin'],
        adminRoles: ['admin'],
        inheritFromParent: false,
      });

      expect(permissions.read).toEqual(['user', 'admin']);
      expect(permissions.write).toEqual(['admin']);
      expect(permissions.inheritFromParent).toBe(false);
    });
  });

  describe('Type Guards', () => {
    it('should validate enhanced content category', () => {
      const category = createContentCategory({
        id: 'test',
        name: 'Test',
        description: 'Test category',
        createdBy: mockUser,
      });

      expect(isEnhancedContentCategory(category)).toBe(true);
    });

    it('should reject invalid category objects', () => {
      expect(isEnhancedContentCategory(null)).toBe(false);
      expect(isEnhancedContentCategory(undefined)).toBe(false);
      expect(isEnhancedContentCategory({})).toBe(false);
      expect(isEnhancedContentCategory({ id: 'test' })).toBe(false);
    });
  });

  describe('Hierarchy Functions', () => {
    const setupHierarchy = (): EnhancedContentCategory[] => {
      const root = createContentCategory({
        id: 'root',
        name: 'Root',
        description: 'Root category',
        createdBy: mockUser,
      });

      const child1 = createContentCategory({
        id: 'child1',
        name: 'Child 1',
        description: 'First child',
        parentId: 'root',
        createdBy: mockUser,
      });

      const child2 = createContentCategory({
        id: 'child2',
        name: 'Child 2',
        description: 'Second child',
        parentId: 'root',
        createdBy: mockUser,
      });

      const grandchild = createContentCategory({
        id: 'grandchild',
        name: 'Grandchild',
        description: 'Grandchild category',
        parentId: 'child1',
        createdBy: mockUser,
      });

      return [root, child1, child2, grandchild];
    };

    it('should compute hierarchy correctly', () => {
      const categories = setupHierarchy();
      const grandchild = categories.find(c => c.id === 'grandchild')!;

      const hierarchy = computeCategoryHierarchy(grandchild, categories);

      expect(hierarchy.level).toBe(2);
      expect(hierarchy.path).toEqual(['root', 'child1']);
      expect(hierarchy.isLeaf).toBe(true);
      expect(hierarchy.isRoot).toBe(false);
    });

    it('should get category parents', () => {
      const categories = setupHierarchy();
      const grandchild = categories.find(c => c.id === 'grandchild')!;

      const parents = getCategoryParents(grandchild, categories);

      expect(parents).toHaveLength(2);
      expect(parents[0].id).toBe('root');
      expect(parents[1].id).toBe('child1');
    });

    it('should get category children', () => {
      const categories = setupHierarchy();
      const root = categories.find(c => c.id === 'root')!;

      const directChildren = getCategoryChildren(root, categories, false);
      const allDescendants = getCategoryChildren(root, categories, true);

      expect(directChildren).toHaveLength(2);
      expect(allDescendants).toHaveLength(3); // child1, child2, grandchild
    });
  });

  describe('Permission System', () => {
    const setupPermissionTest = () => {
      const parent = createContentCategory({
        id: 'parent',
        name: 'Parent',
        description: 'Parent category',
        createdBy: mockAdmin,
        permissions: createCategoryPermissions({
          readRoles: ['user', 'editor', 'admin'],
          writeRoles: ['editor', 'admin'],
          publishRoles: ['admin'],
          deleteRoles: ['admin'],
          adminRoles: ['admin'],
        }),
      });

      const child = createContentCategory({
        id: 'child',
        name: 'Child',
        description: 'Child category',
        parentId: 'parent',
        createdBy: mockAdmin,
        permissions: createCategoryPermissions({
          readRoles: ['editor', 'admin'],
          writeRoles: ['admin'],
          publishRoles: ['admin'],
          deleteRoles: ['admin'],
          adminRoles: ['admin'],
          inheritFromParent: true,
        }),
      });

      return [parent, child];
    };

    it('should check direct permissions', () => {
      const [parent] = setupPermissionTest();

      const userContext: CategoryAccessContext = {
        userId: mockUser,
        userRoles: ['user'],
        action: PermissionAction.READ,
      };

      const adminContext: CategoryAccessContext = {
        userId: mockAdmin,
        userRoles: ['admin'],
        action: PermissionAction.WRITE,
      };

      expect(hasPermission(userContext, parent)).toBe(true);
      expect(hasPermission(adminContext, parent)).toBe(true);
    });

    it('should handle permission inheritance', () => {
      const [parent, child] = setupPermissionTest();

      const userContext: CategoryAccessContext = {
        userId: mockUser,
        userRoles: ['user'],
        action: PermissionAction.READ,
      };

      // User should have read access to child through inheritance
      expect(hasPermission(userContext, child, [parent, child])).toBe(true);
    });

    it('should get effective permissions', () => {
      const [parent, child] = setupPermissionTest();

      const effectivePermissions = getEffectivePermissions(child, [parent, child]);

      // Should inherit read permissions from parent and merge with child's
      expect(effectivePermissions.read).toContain('user');
      expect(effectivePermissions.read).toContain('editor');
      expect(effectivePermissions.read).toContain('admin');
    });
  });

  describe('Validation Functions', () => {
    it('should validate hierarchy without issues', () => {
      const categories = [
        createContentCategory({
          id: 'root',
          name: 'Root',
          description: 'Root category',
          createdBy: mockUser,
        }),
        createContentCategory({
          id: 'child',
          name: 'Child',
          description: 'Child category',
          parentId: 'root',
          createdBy: mockUser,
        }),
      ];

      const validation = validateCategoryHierarchy(categories);

      expect(validation.valid).toBe(true);
      expect(validation.errors).toHaveLength(0);
    });

    it('should detect circular references', () => {
      const categories = [
        createContentCategory({
          id: 'cat1',
          name: 'Category 1',
          description: 'First category',
          parentId: 'cat2',
          createdBy: mockUser,
        }),
        createContentCategory({
          id: 'cat2',
          name: 'Category 2',
          description: 'Second category',
          parentId: 'cat1',
          createdBy: mockUser,
        }),
      ];

      const validation = validateCategoryHierarchy(categories);

      expect(validation.valid).toBe(false);
      expect(validation.errors.some(e => e.code === 'CIRCULAR_REFERENCE')).toBe(true);
    });

    it('should validate category moves', () => {
      const parent = createContentCategory({
        id: 'parent',
        name: 'Parent',
        description: 'Parent category',
        createdBy: mockUser,
      });

      const child = createContentCategory({
        id: 'child',
        name: 'Child',
        description: 'Child category',
        parentId: 'parent',
        createdBy: mockUser,
      });

      const newParent = createContentCategory({
        id: 'newparent',
        name: 'New Parent',
        description: 'New parent category',
        createdBy: mockUser,
      });

      // Valid move
      const validMove = validateCategoryMove(child, 'newparent', [parent, child, newParent]);
      expect(validMove.valid).toBe(true);

      // Invalid move (circular reference)
      const invalidMove = validateCategoryMove(parent, 'child', [parent, child, newParent]);
      expect(invalidMove.valid).toBe(false);
      expect(invalidMove.errors.some(e => e.code === 'CIRCULAR_REFERENCE')).toBe(true);
    });
  });
});