/**
 * Enhanced ContentCategory Implementation
 *
 * This file provides the enhanced ContentCategory type system with:
 * - Complete ContentCategory interface with hierarchy support
 * - Category permission system implementation
 * - Type guards and validation functions for ContentCategory
 * - Factory functions for creating ContentCategory instances
 * - Category hierarchy utility functions (parent/child relationships)
 * - Permission checking helpers
 */

import { ContentCategory, ValidationResult } from './content';

/**
 * Permission action types for content categories
 */
export enum PermissionAction {
  /** View content in the category */
  READ = 'read',
  /** Edit content in the category */
  WRITE = 'write',
  /** Publish content in the category */
  PUBLISH = 'publish',
  /** Delete content in the category */
  DELETE = 'delete',
  /** Manage category settings and permissions */
  ADMIN = 'admin'
}

/**
 * Enhanced permissions interface for content category access control
 */
export interface CategoryPermissions {
  /** User roles that can view content in this category */
  read: string[];
  /** User roles that can edit content in this category */
  write: string[];
  /** User roles that can publish content in this category */
  publish: string[];
  /** User roles that can delete content in this category */
  delete: string[];
  /** User roles that can manage category settings */
  admin: string[];
  /** Whether permissions inherit from parent category */
  inheritFromParent: boolean;
  /** Override settings for specific permissions */
  overrides?: {
    [key in PermissionAction]?: {
      /** Whether to inherit this specific permission from parent */
      inherit: boolean;
      /** Additional roles for this permission (merged with inherited) */
      additionalRoles?: string[];
      /** Roles to exclude from inherited permissions */
      excludeRoles?: string[];
    };
  };
}

/**
 * Category hierarchy metadata
 */
export interface CategoryHierarchy {
  /** Number of levels deep in the hierarchy (0 = root) */
  level: number;
  /** Full path from root to this category */
  path: string[];
  /** Whether this category has child categories */
  hasChildren: boolean;
  /** Number of direct child categories */
  childCount: number;
  /** Total number of descendant categories */
  descendantCount: number;
  /** Whether this category is a leaf node (no children) */
  isLeaf: boolean;
  /** Whether this category is a root node (no parent) */
  isRoot: boolean;
}

/**
 * Category access control context
 */
export interface CategoryAccessContext {
  /** User ID requesting access */
  userId: string;
  /** User roles */
  userRoles: string[];
  /** Requested action */
  action: PermissionAction;
  /** Whether to check inherited permissions */
  checkInheritance?: boolean;
  /** Additional context data */
  context?: Record<string, any>;
}

/**
 * Enhanced content category for organizing and grouping related content items
 */
export interface EnhancedContentCategory extends Omit<ContentCategory, 'permissions'> {
  /** Enhanced access control permissions */
  permissions: CategoryPermissions;
  /** Hierarchy metadata (computed) */
  hierarchy: CategoryHierarchy;
  /** Category creation timestamp */
  createdAt: string;
  /** User who created the category */
  createdBy: string;
  /** Last update timestamp */
  updatedAt: string;
  /** User who last updated the category */
  updatedBy: string;
  /** Category version for optimistic locking */
  version: number;
  /** Whether the category is active */
  isActive: boolean;
  /** Sort order within parent category */
  sortOrder: number;
  /** Category metadata and tags */
  metadata?: {
    /** Searchable tags */
    tags?: string[];
    /** Category color for UI display */
    color?: string;
    /** Category icon identifier */
    icon?: string;
    /** Additional custom properties */
    custom?: Record<string, any>;
  };
}

// ====================================
// TYPE GUARDS AND VALIDATION
// ====================================

/**
 * Enhanced type guard for ContentCategory with comprehensive validation
 * @param value - The value to check
 * @returns true if the value is a valid ContentCategory
 */
export function isEnhancedContentCategory(value: unknown): value is EnhancedContentCategory {
  if (typeof value !== 'object' || value === null) {
    return false;
  }

  const category = value as Partial<EnhancedContentCategory>;

  return (
    typeof category.id === 'string' &&
    typeof category.name === 'string' &&
    typeof category.description === 'string' &&
    (category.parentId === undefined || typeof category.parentId === 'string') &&
    isValidCategoryPermissions(category.permissions) &&
    isValidCategoryHierarchy(category.hierarchy) &&
    typeof category.createdAt === 'string' &&
    typeof category.createdBy === 'string' &&
    typeof category.updatedAt === 'string' &&
    typeof category.updatedBy === 'string' &&
    typeof category.version === 'number' &&
    category.version >= 1 &&
    typeof category.isActive === 'boolean' &&
    typeof category.sortOrder === 'number' &&
    (category.metadata === undefined || isValidCategoryMetadata(category.metadata))
  );
}

/**
 * Validates CategoryPermissions structure
 * @param permissions - The permissions to validate
 * @returns true if valid
 */
export function isValidCategoryPermissions(permissions: unknown): permissions is CategoryPermissions {
  if (typeof permissions !== 'object' || permissions === null) {
    return false;
  }

  const perms = permissions as Partial<CategoryPermissions>;

  return (
    Array.isArray(perms.read) &&
    Array.isArray(perms.write) &&
    Array.isArray(perms.publish) &&
    Array.isArray(perms.delete) &&
    Array.isArray(perms.admin) &&
    typeof perms.inheritFromParent === 'boolean' &&
    perms.read.every(role => typeof role === 'string') &&
    perms.write.every(role => typeof role === 'string') &&
    perms.publish.every(role => typeof role === 'string') &&
    perms.delete.every(role => typeof role === 'string') &&
    perms.admin.every(role => typeof role === 'string') &&
    (perms.overrides === undefined || isValidPermissionOverrides(perms.overrides))
  );
}

/**
 * Validates permission overrides structure
 * @param overrides - The overrides to validate
 * @returns true if valid
 */
export function isValidPermissionOverrides(overrides: unknown): boolean {
  if (typeof overrides !== 'object' || overrides === null) {
    return false;
  }

  const overrideObj = overrides as Record<string, any>;

  return Object.entries(overrideObj).every(([action, override]) => (
    Object.values(PermissionAction).includes(action as PermissionAction) &&
    typeof override === 'object' &&
    override !== null &&
    typeof override.inherit === 'boolean' &&
    (override.additionalRoles === undefined || Array.isArray(override.additionalRoles)) &&
    (override.excludeRoles === undefined || Array.isArray(override.excludeRoles))
  ));
}

/**
 * Validates CategoryHierarchy structure
 * @param hierarchy - The hierarchy to validate
 * @returns true if valid
 */
export function isValidCategoryHierarchy(hierarchy: unknown): hierarchy is CategoryHierarchy {
  if (typeof hierarchy !== 'object' || hierarchy === null) {
    return false;
  }

  const hier = hierarchy as Partial<CategoryHierarchy>;

  return (
    typeof hier.level === 'number' &&
    hier.level >= 0 &&
    Array.isArray(hier.path) &&
    hier.path.every(p => typeof p === 'string') &&
    typeof hier.hasChildren === 'boolean' &&
    typeof hier.childCount === 'number' &&
    hier.childCount >= 0 &&
    typeof hier.descendantCount === 'number' &&
    hier.descendantCount >= 0 &&
    typeof hier.isLeaf === 'boolean' &&
    typeof hier.isRoot === 'boolean'
  );
}

/**
 * Validates category metadata structure
 * @param metadata - The metadata to validate
 * @returns true if valid
 */
export function isValidCategoryMetadata(metadata: unknown): boolean {
  if (typeof metadata !== 'object' || metadata === null) {
    return false;
  }

  const meta = metadata as Record<string, any>;

  return (
    (meta.tags === undefined || Array.isArray(meta.tags)) &&
    (meta.color === undefined || typeof meta.color === 'string') &&
    (meta.icon === undefined || typeof meta.icon === 'string') &&
    (meta.custom === undefined || (typeof meta.custom === 'object' && meta.custom !== null))
  );
}

// ====================================
// FACTORY FUNCTIONS
// ====================================

/**
 * Factory function to create a new EnhancedContentCategory
 * @param params - Category creation parameters
 * @returns A new EnhancedContentCategory instance
 */
export function createContentCategory(params: {
  id: string;
  name: string;
  description: string;
  parentId?: string;
  permissions?: Partial<CategoryPermissions>;
  createdBy: string;
  sortOrder?: number;
  metadata?: EnhancedContentCategory['metadata'];
}): EnhancedContentCategory {
  const now = new Date().toISOString();

  const defaultPermissions: CategoryPermissions = {
    read: ['user', 'editor', 'admin'],
    write: ['editor', 'admin'],
    publish: ['admin'],
    delete: ['admin'],
    admin: ['admin'],
    inheritFromParent: true,
  };

  const permissions = { ...defaultPermissions, ...params.permissions };

  return {
    id: params.id,
    name: params.name,
    description: params.description,
    parentId: params.parentId,
    permissions,
    hierarchy: {
      level: 0, // Will be computed by hierarchy utility
      path: [], // Will be computed by hierarchy utility
      hasChildren: false,
      childCount: 0,
      descendantCount: 0,
      isLeaf: true,
      isRoot: !params.parentId,
    },
    createdAt: now,
    createdBy: params.createdBy,
    updatedAt: now,
    updatedBy: params.createdBy,
    version: 1,
    isActive: true,
    sortOrder: params.sortOrder ?? 0,
    metadata: params.metadata,
  };
}

/**
 * Factory function to create default permissions
 * @param options - Permission configuration options
 * @returns CategoryPermissions object
 */
export function createCategoryPermissions(options: {
  inheritFromParent?: boolean;
  readRoles?: string[];
  writeRoles?: string[];
  publishRoles?: string[];
  deleteRoles?: string[];
  adminRoles?: string[];
  overrides?: CategoryPermissions['overrides'];
} = {}): CategoryPermissions {
  return {
    read: options.readRoles ?? ['user', 'editor', 'admin'],
    write: options.writeRoles ?? ['editor', 'admin'],
    publish: options.publishRoles ?? ['admin'],
    delete: options.deleteRoles ?? ['admin'],
    admin: options.adminRoles ?? ['admin'],
    inheritFromParent: options.inheritFromParent ?? true,
    overrides: options.overrides,
  };
}

// ====================================
// HIERARCHY UTILITY FUNCTIONS
// ====================================

/**
 * Computes category hierarchy information
 * @param category - The category to compute hierarchy for
 * @param allCategories - All categories in the system
 * @returns Updated CategoryHierarchy
 */
export function computeCategoryHierarchy(
  category: EnhancedContentCategory,
  allCategories: EnhancedContentCategory[]
): CategoryHierarchy {
  const categoryMap = new Map(allCategories.map(cat => [cat.id, cat]));

  // Compute path and level
  const path: string[] = [];
  let level = 0;
  let currentCategory = category;

  while (currentCategory.parentId) {
    const parent = categoryMap.get(currentCategory.parentId);
    if (!parent) break;
    path.unshift(parent.id);
    level++;
    currentCategory = parent;
  }

  // Count children and descendants
  const children = allCategories.filter(cat => cat.parentId === category.id);
  const childCount = children.length;

  const getDescendantCount = (catId: string): number => {
    const directChildren = allCategories.filter(cat => cat.parentId === catId);
    return directChildren.length + directChildren.reduce((sum, child) =>
      sum + getDescendantCount(child.id), 0
    );
  };

  const descendantCount = getDescendantCount(category.id);

  return {
    level,
    path,
    hasChildren: childCount > 0,
    childCount,
    descendantCount,
    isLeaf: childCount === 0,
    isRoot: !category.parentId,
  };
}

/**
 * Gets all parent categories for a given category
 * @param category - The category to get parents for
 * @param allCategories - All categories in the system
 * @returns Array of parent categories from root to immediate parent
 */
export function getCategoryParents(
  category: EnhancedContentCategory,
  allCategories: EnhancedContentCategory[]
): EnhancedContentCategory[] {
  const categoryMap = new Map(allCategories.map(cat => [cat.id, cat]));
  const parents: EnhancedContentCategory[] = [];

  let currentCategory = category;
  while (currentCategory.parentId) {
    const parent = categoryMap.get(currentCategory.parentId);
    if (!parent) break;
    parents.unshift(parent);
    currentCategory = parent;
  }

  return parents;
}

/**
 * Gets all child categories for a given category
 * @param category - The category to get children for
 * @param allCategories - All categories in the system
 * @param recursive - Whether to include all descendants or just direct children
 * @returns Array of child categories
 */
export function getCategoryChildren(
  category: EnhancedContentCategory,
  allCategories: EnhancedContentCategory[],
  recursive = false
): EnhancedContentCategory[] {
  const directChildren = allCategories.filter(cat => cat.parentId === category.id);

  if (!recursive) {
    return directChildren;
  }

  const getAllDescendants = (catId: string): EnhancedContentCategory[] => {
    const children = allCategories.filter(cat => cat.parentId === catId);
    return children.reduce((acc, child) => [
      ...acc,
      child,
      ...getAllDescendants(child.id)
    ], [] as EnhancedContentCategory[]);
  };

  return getAllDescendants(category.id);
}

/**
 * Updates category hierarchy metadata for all categories
 * @param categories - All categories to update
 * @returns Updated categories with correct hierarchy metadata
 */
export function updateCategoryHierarchies(
  categories: EnhancedContentCategory[]
): EnhancedContentCategory[] {
  return categories.map(category => ({
    ...category,
    hierarchy: computeCategoryHierarchy(category, categories)
  }));
}

// ====================================
// PERMISSION CHECKING HELPERS
// ====================================

/**
 * Checks if a user has permission to perform an action on a category
 * @param context - Access control context
 * @param category - The category to check permissions for
 * @param allCategories - All categories (needed for inheritance)
 * @returns true if user has permission
 */
export function hasPermission(
  context: CategoryAccessContext,
  category: EnhancedContentCategory,
  allCategories?: EnhancedContentCategory[]
): boolean {
  const { userRoles, action, checkInheritance = true } = context;
  const { permissions } = category;

  // Check direct permissions
  const actionRoles = permissions[action];
  if (actionRoles.some(role => userRoles.includes(role))) {
    return true;
  }

  // Check permission overrides
  const override = permissions.overrides?.[action];
  if (override) {
    if (!override.inherit) {
      // Override disables inheritance, check only additional roles
      const additionalRoles = override.additionalRoles ?? [];
      return additionalRoles.some(role => userRoles.includes(role));
    }

    // Remove excluded roles from user roles for this check
    const excludeRoles = override.excludeRoles ?? [];
    const effectiveUserRoles = userRoles.filter(role => !excludeRoles.includes(role));

    // Check additional roles
    const additionalRoles = override.additionalRoles ?? [];
    if (additionalRoles.some(role => effectiveUserRoles.includes(role))) {
      return true;
    }
  }

  // Check inherited permissions if enabled and parent exists
  if (checkInheritance && permissions.inheritFromParent && category.parentId && allCategories) {
    const parent = allCategories.find(cat => cat.id === category.parentId);
    if (parent) {
      return hasPermission(context, parent, allCategories);
    }
  }

  return false;
}

/**
 * Gets effective permissions for a category (including inheritance)
 * @param category - The category to get permissions for
 * @param allCategories - All categories (needed for inheritance)
 * @returns Effective permissions with inheritance resolved
 */
export function getEffectivePermissions(
  category: EnhancedContentCategory,
  allCategories: EnhancedContentCategory[]
): CategoryPermissions {
  if (!category.permissions.inheritFromParent || !category.parentId) {
    return category.permissions;
  }

  const parent = allCategories.find(cat => cat.id === category.parentId);
  if (!parent) {
    return category.permissions;
  }

  const parentPerms = getEffectivePermissions(parent, allCategories);
  const categoryPerms = category.permissions;

  // Merge permissions with overrides applied
  const effectivePerms: CategoryPermissions = {
    read: [...parentPerms.read],
    write: [...parentPerms.write],
    publish: [...parentPerms.publish],
    delete: [...parentPerms.delete],
    admin: [...parentPerms.admin],
    inheritFromParent: categoryPerms.inheritFromParent,
    overrides: categoryPerms.overrides,
  };

  // Apply overrides and merge with category-specific permissions
  Object.entries(PermissionAction).forEach(([, action]) => {
    const override = categoryPerms.overrides?.[action];
    const categoryRoles = categoryPerms[action];

    if (override && !override.inherit) {
      // Replace inherited permissions with only additional roles
      effectivePerms[action] = [...(override.additionalRoles ?? [])];
    } else {
      // Merge inherited and category-specific roles
      const inheritedRoles = effectivePerms[action];
      const additionalRoles = override?.additionalRoles ?? [];
      const excludeRoles = override?.excludeRoles ?? [];

      const allRoles = [...inheritedRoles, ...categoryRoles, ...additionalRoles];
      effectivePerms[action] = [...new Set(allRoles)].filter(role => !excludeRoles.includes(role));
    }
  });

  return effectivePerms;
}

/**
 * Checks if a user can access any content in a category tree
 * @param context - Access control context
 * @param category - Root category to check
 * @param allCategories - All categories in the system
 * @returns true if user has any access to the category tree
 */
export function canAccessCategoryTree(
  context: CategoryAccessContext,
  category: EnhancedContentCategory,
  allCategories: EnhancedContentCategory[]
): boolean {
  // Check direct access to this category
  const actions = [PermissionAction.READ, PermissionAction.WRITE, PermissionAction.PUBLISH, PermissionAction.DELETE, PermissionAction.ADMIN];

  for (const action of actions) {
    if (hasPermission({ ...context, action }, category, allCategories)) {
      return true;
    }
  }

  // Check access to any child categories
  const children = getCategoryChildren(category, allCategories, true);
  for (const child of children) {
    for (const action of actions) {
      if (hasPermission({ ...context, action }, child, allCategories)) {
        return true;
      }
    }
  }

  return false;
}

// ====================================
// VALIDATION FUNCTIONS
// ====================================

/**
 * Validates category hierarchy consistency
 * @param categories - All categories to validate
 * @returns Validation result with any hierarchy issues
 */
export function validateCategoryHierarchy(categories: EnhancedContentCategory[]): ValidationResult {
  const errors: ValidationResult['errors'] = [];
  const warnings: ValidationResult['warnings'] = [];
  const categoryMap = new Map(categories.map(cat => [cat.id, cat]));

  categories.forEach(category => {
    // Check for circular references
    const visited = new Set<string>();
    let current = category;

    while (current.parentId) {
      if (visited.has(current.id)) {
        errors.push({
          field: 'hierarchy',
          message: `Circular reference detected in category hierarchy for category "${category.id}"`,
          code: 'CIRCULAR_REFERENCE',
        });
        break;
      }

      visited.add(current.id);
      const parent = categoryMap.get(current.parentId);

      if (!parent) {
        errors.push({
          field: 'parentId',
          message: `Parent category "${current.parentId}" not found for category "${category.id}"`,
          code: 'PARENT_NOT_FOUND',
        });
        break;
      }

      current = parent;
    }

    // Check hierarchy metadata consistency
    const computedHierarchy = computeCategoryHierarchy(category, categories);
    if (category.hierarchy.level !== computedHierarchy.level) {
      warnings.push({
        field: 'hierarchy.level',
        message: `Hierarchy level mismatch for category "${category.id}". Expected ${computedHierarchy.level}, got ${category.hierarchy.level}`,
        code: 'HIERARCHY_LEVEL_MISMATCH',
      });
    }

    if (category.hierarchy.childCount !== computedHierarchy.childCount) {
      warnings.push({
        field: 'hierarchy.childCount',
        message: `Child count mismatch for category "${category.id}". Expected ${computedHierarchy.childCount}, got ${category.hierarchy.childCount}`,
        code: 'CHILD_COUNT_MISMATCH',
      });
    }
  });

  return {
    valid: errors.length === 0,
    errors,
    warnings,
  };
}

/**
 * Validates that a category can be moved to a new parent
 * @param category - Category to move
 * @param newParentId - New parent category ID
 * @param allCategories - All categories in the system
 * @returns Validation result
 */
export function validateCategoryMove(
  category: EnhancedContentCategory,
  newParentId: string | undefined,
  allCategories: EnhancedContentCategory[]
): ValidationResult {
  const errors: ValidationResult['errors'] = [];
  const warnings: ValidationResult['warnings'] = [];

  if (newParentId) {
    // Check if new parent exists
    const newParent = allCategories.find(cat => cat.id === newParentId);
    if (!newParent) {
      errors.push({
        field: 'parentId',
        message: `New parent category "${newParentId}" not found`,
        code: 'PARENT_NOT_FOUND',
      });
      return { valid: false, errors, warnings };
    }

    // Check if moving would create a circular reference
    const descendants = getCategoryChildren(category, allCategories, true);
    if (descendants.some(desc => desc.id === newParentId)) {
      errors.push({
        field: 'parentId',
        message: `Cannot move category "${category.id}" to "${newParentId}" as it would create a circular reference`,
        code: 'CIRCULAR_REFERENCE',
      });
    }

    // Check if new parent is the same as current parent
    if (category.parentId === newParentId) {
      warnings.push({
        field: 'parentId',
        message: `Category "${category.id}" is already a child of "${newParentId}"`,
        code: 'NO_CHANGE',
      });
    }
  }

  return {
    valid: errors.length === 0,
    errors,
    warnings,
  };
}