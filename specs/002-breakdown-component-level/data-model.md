# Data Model: Component Level Molecules Breakdown

**Date**: 2025-09-21
**Version**: 1.0.0

## Entity Relationship Diagram

```
Component (base)
    ↑
    ├── Atom
    ├── Molecule (contains 2+ Atoms)
    └── Organism (contains Molecules/Atoms)

ComponentRegistry ←→ Component (1:many)
Component ←→ UsagePattern (1:many)
Component ←→ ConsistencyRule (many:many)
Molecule ←→ Atom (many:many via Composition)
```

## Core Entities

### 1. Component (Base Entity)
Base class for all UI components in the system.

```typescript
interface Component {
  id: string;                    // Unique identifier (e.g., "btn-primary")
  name: string;                   // Display name (e.g., "Primary Button")
  type: ComponentType;            // "atom" | "molecule" | "organism"
  filePath: string;               // Relative path from project root
  category: string;               // Functional category (e.g., "form", "navigation")
  description: string;            // Purpose and usage description
  props: PropDefinition[];        // Component props/attributes
  dependencies: string[];         // Import dependencies
  usageCount: number;             // Times used in codebase
  createdAt: Date;                // When first detected
  updatedAt: Date;                // Last modification
  deprecated: boolean;            // Deprecation status
  version: string;                // Component version
}

enum ComponentType {
  ATOM = "atom",
  MOLECULE = "molecule",
  ORGANISM = "organism"
}
```

### 2. Atom
The most basic UI building block.

```typescript
interface Atom extends Component {
  type: ComponentType.ATOM;
  htmlElement: string;           // Base HTML element (e.g., "button", "input")
  variants: AtomVariant[];        // Different variations (size, color, etc.)
  accessibility: AccessibilityInfo;
}

interface AtomVariant {
  name: string;                   // Variant name (e.g., "large", "primary")
  props: Record<string, any>;    // Props for this variant
  exampleCode: string;            // Usage example
}

interface AccessibilityInfo {
  ariaLabel: boolean;             // Supports aria-label
  ariaDescribedBy: boolean;       // Supports aria-describedby
  role: string | null;            // ARIA role if applicable
  keyboardNav: boolean;           // Keyboard accessible
  wcagLevel: "A" | "AA" | "AAA"; // WCAG compliance level
}
```

### 3. Molecule
Combination of atoms working together.

```typescript
interface Molecule extends Component {
  type: ComponentType.MOLECULE;
  composition: Composition[];      // Atoms that make up this molecule
  layout: LayoutType;              // How atoms are arranged
  interactions: Interaction[];     // How atoms interact
  slots: SlotDefinition[];         // Customizable slots
}

interface Composition {
  atomId: string;                  // Reference to Atom.id
  quantity: number;                 // How many instances
  required: boolean;                // If this atom is required
  role: string;                     // Role in the molecule (e.g., "trigger", "content")
}

enum LayoutType {
  HORIZONTAL = "horizontal",
  VERTICAL = "vertical",
  GRID = "grid",
  ABSOLUTE = "absolute",
  FLEXIBLE = "flexible"
}

interface Interaction {
  trigger: string;                  // Event trigger (e.g., "onClick")
  source: string;                   // Source atom ID
  target: string;                   // Target atom ID
  action: string;                   // What happens
}
```

### 4. Organism
Complex, self-contained component sections.

```typescript
interface Organism extends Component {
  type: ComponentType.ORGANISM;
  contains: ComponentReference[];   // Molecules and atoms it contains
  standalone: boolean;              // Can function independently
  pageSection: boolean;             // Represents page section
  dataSource: string | null;        // External data dependency
}

interface ComponentReference {
  componentId: string;              // Reference to Component.id
  quantity: number;                 // How many instances
  configuration: Record<string, any>; // Configuration for this instance
}
```

### 5. ComponentRegistry
Central catalog of all components.

```typescript
interface ComponentRegistry {
  id: string;                      // Registry ID
  projectName: string;              // Project this registry belongs to
  components: Map<string, Component>; // All components keyed by ID
  lastUpdated: Date;                // Last scan timestamp
  statistics: RegistryStats;        // Aggregate statistics
  version: string;                  // Registry schema version
}

interface RegistryStats {
  totalComponents: number;
  atomCount: number;
  moleculeCount: number;
  organismCount: number;
  duplicatesFound: number;
  inconsistenciesFound: number;
  averageUsageCount: number;
  mostUsedComponents: string[];    // Top 10 component IDs
}
```

### 6. UsagePattern
Where and how components are used.

```typescript
interface UsagePattern {
  id: string;                      // Usage ID
  componentId: string;             // Reference to Component.id
  filePath: string;                // Where it's used
  lineNumber: number;              // Line in file
  context: UsageContext;           // How it's being used
  props: Record<string, any>;      // Props passed in this usage
  parentComponent: string | null;  // Parent component ID if nested
}

enum UsageContext {
  PAGE = "page",                   // Used in a page
  LAYOUT = "layout",               // Part of layout
  FEATURE = "feature",             // Feature component
  SHARED = "shared",               // Shared/common usage
  TEST = "test"                    // Used in tests
}
```

### 7. ConsistencyRule
Standards that components must follow.

```typescript
interface ConsistencyRule {
  id: string;                      // Rule ID
  name: string;                    // Rule name
  description: string;             // What this rule checks
  category: RuleCategory;          // Type of rule
  severity: RuleSeverity;          // How critical
  validator: string;               // Validation function name
  appliesTo: ComponentType[];      // Which component types
  enabled: boolean;                // If rule is active
}

enum RuleCategory {
  NAMING = "naming",               // Naming conventions
  STRUCTURE = "structure",         // Component structure
  STYLING = "styling",             // Style consistency
  ACCESSIBILITY = "accessibility", // A11y requirements
  PERFORMANCE = "performance",     // Performance standards
  DOCUMENTATION = "documentation"  // Documentation requirements
}

enum RuleSeverity {
  ERROR = "error",                 // Must fix
  WARNING = "warning",             // Should fix
  INFO = "info"                    // Nice to have
}
```

### 8. PropDefinition
Component prop/attribute definition.

```typescript
interface PropDefinition {
  name: string;                    // Prop name
  type: string;                    // TypeScript type
  required: boolean;               // If required
  defaultValue: any;               // Default value if any
  description: string;             // What this prop does
  examples: string[];              // Example values
}
```

### 9. SlotDefinition
Customizable slot in a component.

```typescript
interface SlotDefinition {
  name: string;                    // Slot name
  description: string;             // What goes in this slot
  accepts: ComponentType[];        // What types can go here
  required: boolean;               // If content is required
  defaultContent: string | null;   // Default if not provided
}
```

## State Transitions

### Component Lifecycle States
```
DISCOVERED → ANALYZED → CATEGORIZED → DOCUMENTED → ACTIVE
                ↓            ↓            ↓          ↓
            INVALID     UNCATEGORIZED  OUTDATED  DEPRECATED
```

### Validation States
```
PENDING → VALIDATING → VALID
             ↓
          INVALID → FIXED → VALID
```

## Validation Rules

### Component Validation
1. **Unique ID**: No duplicate component IDs
2. **Valid Type**: Must be atom, molecule, or organism
3. **File Exists**: Referenced file must exist
4. **Props Defined**: All props must have type definitions

### Molecule Validation
1. **Has Atoms**: Must contain at least 2 atoms
2. **Atoms Exist**: Referenced atoms must be in registry
3. **Valid Composition**: Composition must be logically valid

### Registry Validation
1. **No Orphans**: All references must resolve
2. **Version Compatible**: Components compatible with registry version
3. **Statistics Accurate**: Stats match actual component counts

## Indexing Strategy

### Primary Indexes
- `component.id` - Unique component identifier
- `component.name` - For search operations
- `component.type` - For filtering by type
- `component.filePath` - For file-based lookups

### Secondary Indexes
- `component.usageCount` - For popularity sorting
- `component.category` - For category grouping
- `usagePattern.componentId` - For usage lookups
- `composition.atomId` - For dependency tracking

## Data Persistence

### Storage Format
```json
{
  "version": "1.0.0",
  "timestamp": "2025-09-21T10:00:00Z",
  "registry": {
    "components": [...],
    "usagePatterns": [...],
    "consistencyRules": [...]
  },
  "metadata": {
    "projectPath": "/Users/weerawat/Tchat",
    "scanDuration": 4500,
    "filesAnalyzed": 127
  }
}
```

### File Locations
- Registry: `docs/components/registry.json`
- Backups: `docs/components/backups/registry-{timestamp}.json`
- Cache: `.component-analyzer/cache.json`

## Migration Strategy

### From Existing Components
1. Scan all component files
2. Auto-categorize based on patterns
3. Generate initial registry
4. Manual review and corrections
5. Establish baseline

### Version Migrations
- Registry version changes trigger migration scripts
- Backward compatibility for 2 major versions
- Migration logs maintained for audit