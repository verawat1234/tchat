# Component Analyzer Quickstart Guide

**Time to first result**: ~5 minutes
**Prerequisites**: Node.js 18+, npm/yarn, existing React/TypeScript project

## Installation

```bash
# Install globally
npm install -g @tchat/component-analyzer

# Or add to project
npm install --save-dev @tchat/component-analyzer
```

## Quick Start (3 Steps)

### 1. Run Initial Analysis
```bash
# From your project root
component-analyzer analyze

# Output:
# ✓ Analyzed 87 components
# ✓ Categorized: 23 atoms, 45 molecules, 19 organisms
# ✓ Found 12 potential duplicates
# ✓ Registry saved to docs/components/registry.json
```

### 2. Review Results
```bash
# List all molecules
component-analyzer list molecules

# Check for duplicates
component-analyzer duplicates

# Validate consistency
component-analyzer validate
```

### 3. Generate Documentation
```bash
# Generate markdown docs
component-analyzer generate docs

# Files created in docs/components/:
# - atoms.md
# - molecules.md
# - organisms.md
# - index.md
```

## Common Use Cases

### Use Case 1: Developer Creating New Feature
**Goal**: Find existing components to reuse

```bash
# Search for form-related molecules
component-analyzer list molecules --category form

# Get detailed info about a component
component-analyzer show input-field

# Generate usage examples
component-analyzer generate examples --components input-field
```

### Use Case 2: Designer Reviewing UI Consistency
**Goal**: Identify inconsistent components

```bash
# Find duplicates with high similarity
component-analyzer duplicates --threshold 80

# Validate against design standards
component-analyzer validate --rules styling,accessibility

# Generate visual reference
component-analyzer generate visuals
```

### Use Case 3: QA Testing Component Usage
**Goal**: Verify consistent component usage

```bash
# Check component usage patterns
component-analyzer stats --detailed

# Find components with low usage (potential removal candidates)
component-analyzer list --sort usage --limit 10

# Validate all components
component-analyzer validate
```

### Use Case 4: Team Lead Auditing Component Library
**Goal**: Optimize component library

```bash
# Get comprehensive statistics
component-analyzer stats --chart

# Find and merge duplicates
component-analyzer duplicates --auto-merge

# Generate full documentation
component-analyzer generate docs --force
```

## Automated Workflow Setup

### 1. Add to Git Hooks
```bash
# .git/hooks/pre-commit
#!/bin/sh
component-analyzer validate --severity error || exit 1
```

### 2. Add to CI/CD Pipeline
```yaml
# .github/workflows/components.yml
name: Component Analysis
on: [push, pull_request]

jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - run: npm install -g @tchat/component-analyzer
      - run: component-analyzer analyze
      - run: component-analyzer validate
      - run: component-analyzer duplicates --threshold 90
```

### 3. Watch Mode for Development
```bash
# Start watching for changes
component-analyzer watch --auto-fix --notify

# Running in background...
# ✓ Watching apps/web/src/components
# → New molecule detected: SearchBar
# → Auto-categorized and added to registry
# → Documentation updated
```

## Configuration

### Basic Configuration
Create `.component-analyzer.json` in project root:

```json
{
  "paths": {
    "components": "apps/web/src/components"
  },
  "analysis": {
    "exclude": ["*.test.tsx", "*.stories.tsx"]
  },
  "validation": {
    "autoFix": true,
    "severity": "warning"
  }
}
```

### Advanced Configuration
```json
{
  "paths": {
    "components": "apps/web/src/components",
    "output": "docs/components",
    "registry": "docs/components/registry.json"
  },
  "analysis": {
    "recursive": true,
    "maxDepth": 10,
    "exclude": ["*.test.tsx", "*.stories.tsx", "__mocks__/*"],
    "categorization": {
      "atomPatterns": ["Button*.tsx", "Input*.tsx", "Icon*.tsx"],
      "moleculePatterns": ["*Form.tsx", "*Card.tsx", "*List.tsx"],
      "organismPatterns": ["*Section.tsx", "*Layout.tsx", "*Page.tsx"]
    }
  },
  "validation": {
    "rules": {
      "naming": {
        "enabled": true,
        "pattern": "^[A-Z][a-zA-Z]*$"
      },
      "accessibility": {
        "enabled": true,
        "wcagLevel": "AA"
      },
      "documentation": {
        "enabled": true,
        "requireDescription": true,
        "requireExamples": false
      }
    },
    "autoFix": false,
    "severity": "error"
  },
  "duplicates": {
    "threshold": 75,
    "factors": {
      "structural": 0.4,
      "visual": 0.3,
      "functional": 0.3
    }
  },
  "output": {
    "format": "both",
    "includeVisuals": true,
    "includeExamples": true
  }
}
```

## Validation Test Scenarios

### Test 1: Initial Analysis
```bash
# Run analysis
component-analyzer analyze apps/web/src/components

# Expected output:
# - All components detected and categorized
# - Registry file created
# - No crashes or errors
```

### Test 2: Duplicate Detection
```bash
# Create duplicate component for testing
cp src/components/Button.tsx src/components/ButtonCopy.tsx

# Run duplicate detection
component-analyzer duplicates

# Expected output:
# - Button and ButtonCopy identified as duplicates
# - Similarity score > 95%
# - Merge suggestion provided
```

### Test 3: Consistency Validation
```bash
# Run validation
component-analyzer validate

# Expected output:
# - All consistency rules applied
# - Violations reported with severity
# - Suggestions for fixes provided
```

### Test 4: Documentation Generation
```bash
# Generate docs
component-analyzer generate docs

# Expected output:
# - Markdown files created for each component type
# - All components documented
# - Usage examples included
```

### Test 5: Watch Mode
```bash
# Start watch mode
component-analyzer watch &

# Modify a component
echo "// Modified" >> src/components/Button.tsx

# Expected output:
# - Change detected within 5 seconds
# - Component re-analyzed
# - Registry updated
```

## Troubleshooting

### Issue: Components not detected
**Solution**: Check path configuration and file patterns
```bash
component-analyzer analyze --verbose
component-analyzer config get paths.components
```

### Issue: Wrong categorization
**Solution**: Review and adjust categorization patterns
```bash
component-analyzer list --format json | jq '.uncategorized'
component-analyzer config set analysis.categorization.atomPatterns '["*.atom.tsx"]'
```

### Issue: High number of duplicates
**Solution**: Adjust similarity threshold
```bash
component-analyzer duplicates --threshold 85
component-analyzer config set duplicates.threshold 85
```

### Issue: Validation too strict
**Solution**: Adjust severity levels
```bash
component-analyzer validate --severity warning
component-analyzer config set validation.severity warning
```

## Best Practices

1. **Run analysis regularly**: Weekly or bi-weekly
2. **Review duplicates before merging**: Ensure functional equivalence
3. **Document as you go**: Update component descriptions immediately
4. **Use watch mode in development**: Catch issues early
5. **Integrate with CI/CD**: Enforce standards automatically
6. **Start with warnings**: Gradually increase to errors
7. **Involve designers**: Review categorization together

## Next Steps

1. **Customize rules**: Add project-specific validation rules
2. **Create component templates**: Based on identified patterns
3. **Build component library**: Using categorized molecules
4. **Set up Storybook**: Auto-generate stories from analysis
5. **Track metrics**: Monitor component usage over time

## Support

- **Documentation**: `/docs/components/analyzer`
- **Issues**: GitHub Issues
- **Community**: Tchat Discord #components channel