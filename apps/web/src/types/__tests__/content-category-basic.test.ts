/**
 * @file content-category-basic.test.ts
 * @description Basic test for enhanced ContentCategory implementation
 */

import { describe, expect, it } from 'vitest';
import {
  PermissionAction,
  createContentCategory,
  createCategoryPermissions,
  hasPermission,
  isEnhancedContentCategory,
} from '../content-category';

describe('ContentCategory Basic Tests', () => {
  it('should create a content category', () => {
    const category = createContentCategory({
      id: 'test',
      name: 'Test Category',
      description: 'A test category',
      createdBy: 'user123',
    });

    expect(category.id).toBe('test');
    expect(category.name).toBe('Test Category');
    expect(isEnhancedContentCategory(category)).toBe(true);
  });

  it('should create permissions', () => {
    const permissions = createCategoryPermissions({
      readRoles: ['user'],
      writeRoles: ['admin'],
    });

    expect(permissions.read).toEqual(['user']);
    expect(permissions.write).toEqual(['admin']);
  });

  it('should check permissions', () => {
    const category = createContentCategory({
      id: 'test',
      name: 'Test',
      description: 'Test',
      createdBy: 'user123',
    });

    const hasReadPermission = hasPermission(
      {
        userId: 'user1',
        userRoles: ['user'],
        action: PermissionAction.READ,
      },
      category
    );

    expect(hasReadPermission).toBe(true);
  });
});