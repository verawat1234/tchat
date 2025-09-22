/**
 * @file content-category-demo.ts
 * @description Demonstration of enhanced ContentCategory implementation
 */

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
  updateCategoryHierarchies,
  type EnhancedContentCategory,
  type CategoryAccessContext,
} from '../content-category';

// Example 1: Creating a basic category hierarchy
console.log('=== Example 1: Creating Category Hierarchy ===');

const rootCategory = createContentCategory({
  id: 'documentation',
  name: 'Documentation',
  description: 'All documentation content',
  createdBy: 'admin@example.com',
});

const apiDocsCategory = createContentCategory({
  id: 'api-docs',
  name: 'API Documentation',
  description: 'API reference documentation',
  parentId: 'documentation',
  createdBy: 'admin@example.com',
});

const userGuideCategory = createContentCategory({
  id: 'user-guide',
  name: 'User Guide',
  description: 'End-user documentation',
  parentId: 'documentation',
  createdBy: 'admin@example.com',
});

const quickstartCategory = createContentCategory({
  id: 'quickstart',
  name: 'Quickstart Guide',
  description: 'Getting started documentation',
  parentId: 'user-guide',
  createdBy: 'admin@example.com',
});

console.log('Created categories:');
console.log('- Root:', rootCategory.name);
console.log('- API Docs:', apiDocsCategory.name);
console.log('- User Guide:', userGuideCategory.name);
console.log('- Quickstart:', quickstartCategory.name);

// Example 2: Computing hierarchy metadata
console.log('\n=== Example 2: Computing Hierarchy Metadata ===');

const allCategories = [rootCategory, apiDocsCategory, userGuideCategory, quickstartCategory];
const updatedCategories = updateCategoryHierarchies(allCategories);

const quickstartWithHierarchy = updatedCategories.find(c => c.id === 'quickstart')!;
console.log('Quickstart category hierarchy:', {
  level: quickstartWithHierarchy.hierarchy.level,
  path: quickstartWithHierarchy.hierarchy.path,
  isLeaf: quickstartWithHierarchy.hierarchy.isLeaf,
  isRoot: quickstartWithHierarchy.hierarchy.isRoot,
});

// Example 3: Working with hierarchy relationships
console.log('\n=== Example 3: Hierarchy Relationships ===');

const parents = getCategoryParents(quickstartWithHierarchy, updatedCategories);
console.log('Quickstart parents:', parents.map(p => p.name));

const rootChildren = getCategoryChildren(updatedCategories[0], updatedCategories, true);
console.log('All descendants of root:', rootChildren.map(c => c.name));

// Example 4: Advanced permission system
console.log('\n=== Example 4: Advanced Permission System ===');

const restrictedCategory = createContentCategory({
  id: 'internal-docs',
  name: 'Internal Documentation',
  description: 'Internal company documentation',
  parentId: 'documentation',
  createdBy: 'admin@example.com',
  permissions: createCategoryPermissions({
    readRoles: ['employee', 'manager', 'admin'],
    writeRoles: ['manager', 'admin'],
    publishRoles: ['admin'],
    deleteRoles: ['admin'],
    adminRoles: ['admin'],
    inheritFromParent: false, // Don't inherit from parent
    overrides: {
      [PermissionAction.READ]: {
        inherit: false,
        additionalRoles: ['contractor'], // Contractors can read but not from inheritance
        excludeRoles: [], // No roles to exclude
      },
    },
  }),
});

// Test permissions for different user roles
const users = [
  { id: 'user1', roles: ['user'], name: 'Regular User' },
  { id: 'user2', roles: ['employee'], name: 'Employee' },
  { id: 'user3', roles: ['contractor'], name: 'Contractor' },
  { id: 'user4', roles: ['manager'], name: 'Manager' },
  { id: 'user5', roles: ['admin'], name: 'Admin' },
];

console.log('Permission check for internal docs:');
users.forEach(user => {
  const context: CategoryAccessContext = {
    userId: user.id,
    userRoles: user.roles,
    action: PermissionAction.READ,
  };

  const canRead = hasPermission(context, restrictedCategory);
  console.log(`- ${user.name}: ${canRead ? 'CAN' : 'CANNOT'} read`);
});

// Example 5: Permission inheritance
console.log('\n=== Example 5: Permission Inheritance ===');

const childWithInheritance = createContentCategory({
  id: 'child-inherit',
  name: 'Child with Inheritance',
  description: 'Child category that inherits permissions',
  parentId: 'internal-docs',
  createdBy: 'admin@example.com',
  permissions: createCategoryPermissions({
    readRoles: ['guest'], // Add guest role in addition to inherited roles
    writeRoles: [],
    publishRoles: [],
    deleteRoles: [],
    adminRoles: [],
    inheritFromParent: true,
  }),
});

const allCategoriesWithInheritance = [
  ...updatedCategories,
  restrictedCategory,
  childWithInheritance,
];

const effectivePermissions = getEffectivePermissions(
  childWithInheritance,
  allCategoriesWithInheritance
);

console.log('Effective read permissions for child:', effectivePermissions.read);

// Example 6: Validation
console.log('\n=== Example 6: Hierarchy Validation ===');

const validationResult = validateCategoryHierarchy(allCategoriesWithInheritance);
console.log('Hierarchy validation:', {
  valid: validationResult.valid,
  errors: validationResult.errors.length,
  warnings: validationResult.warnings.length,
});

if (validationResult.errors.length > 0) {
  console.log('Validation errors:');
  validationResult.errors.forEach(error => {
    console.log(`- ${error.field}: ${error.message} (${error.code})`);
  });
}

// Example 7: Complex permission scenarios
console.log('\n=== Example 7: Complex Permission Override ===');

const complexCategory = createContentCategory({
  id: 'complex-perms',
  name: 'Complex Permissions',
  description: 'Category with complex permission overrides',
  parentId: 'documentation',
  createdBy: 'admin@example.com',
  permissions: createCategoryPermissions({
    readRoles: ['user'],
    writeRoles: ['editor'],
    publishRoles: ['admin'],
    deleteRoles: ['admin'],
    adminRoles: ['admin'],
    inheritFromParent: true,
    overrides: {
      [PermissionAction.WRITE]: {
        inherit: true, // Inherit from parent
        additionalRoles: ['contributor'], // Add contributors
        excludeRoles: ['junior-editor'], // Exclude junior editors
      },
      [PermissionAction.READ]: {
        inherit: false, // Don't inherit read permissions
        additionalRoles: ['guest', 'user', 'editor', 'admin'], // Explicit roles only
      },
    },
  }),
});

console.log('Complex category permissions created successfully');

console.log('\n=== Demo Complete ===');
console.log('Enhanced ContentCategory implementation provides:');
console.log('1. ✅ Complete ContentCategory interface with hierarchy support');
console.log('2. ✅ Category permission system implementation');
console.log('3. ✅ Type guards and validation functions');
console.log('4. ✅ Factory functions for creating categories');
console.log('5. ✅ Category hierarchy utility functions');
console.log('6. ✅ Permission checking helpers');

export {
  rootCategory,
  apiDocsCategory,
  userGuideCategory,
  quickstartCategory,
  restrictedCategory,
  childWithInheritance,
  complexCategory,
  allCategoriesWithInheritance,
};