# Component Analyzer

A powerful CLI tool to analyze and categorize React components following atomic design principles.

## Features

- 🔍 **Automatic Component Detection** - Scans your codebase to find all React components
- 🏗️ **Atomic Design Categorization** - Categorizes components as Atoms, Molecules, or Organisms
- 🔄 **Duplicate Detection** - Identifies similar components that could be merged
- ✅ **Consistency Validation** - Ensures components follow your team's standards
- 📊 **Usage Statistics** - Tracks how often each component is used
- 📝 **Documentation Generation** - Creates comprehensive component documentation
- 👀 **Watch Mode** - Monitors changes and updates registry in real-time
- 🎨 **Customizable Rules** - Define your own validation rules and categorization patterns

## Installation

```bash
# Install globally
npm install -g @tchat/component-analyzer

# Or add to your project
npm install --save-dev @tchat/component-analyzer
```

## Quick Start

```bash
# Analyze your components
component-analyzer analyze

# List all molecules
component-analyzer list molecules

# Find duplicate components
component-analyzer duplicates

# Generate documentation
component-analyzer generate docs

# Watch for changes
component-analyzer watch
```

## Commands

### `analyze [path]`
Analyze components in a directory and categorize them.

```bash
component-analyzer analyze src/components --recursive --save
```

Options:
- `-r, --recursive` - Include subdirectories (default: true)
- `-o, --output <format>` - Output format: json|markdown|both
- `-s, --save` - Save results to registry
- `--exclude <patterns>` - Glob patterns to exclude

### `list [type]`
List components from the registry.

```bash
component-analyzer list molecules --category form --sort usage
```

Options:
- `-c, --category <category>` - Filter by category
- `-s, --sort <field>` - Sort by: name|usage|created
- `-f, --format <format>` - Output format: table|json|list

### `duplicates`
Find duplicate or similar components.

```bash
component-analyzer duplicates --threshold 80 --auto-merge
```

Options:
- `-t, --threshold <n>` - Similarity threshold % (default: 75)
- `--auto-merge` - Automatically suggest merges

### `validate [component-id]`
Validate components against consistency rules.

```bash
component-analyzer validate --rules naming,styling --fix
```

Options:
- `-r, --rules <rules>` - Specific rules to apply
- `--fix` - Attempt to auto-fix issues
- `-s, --severity <level>` - Minimum severity: error|warning|info

### `generate [type]`
Generate documentation for components.

```bash
component-analyzer generate docs --format markdown --force
```

Options:
- `-f, --format <format>` - Format: markdown|html|json
- `--force` - Overwrite existing files

### `watch [path]`
Watch for component changes and update registry.

```bash
component-analyzer watch --interval 5 --auto-fix
```

Options:
- `-i, --interval <n>` - Check interval in seconds
- `--auto-fix` - Automatically fix consistency issues

### `stats`
Display statistics about components.

```bash
component-analyzer stats --period month --chart
```

Options:
- `-p, --period <period>` - Time period: today|week|month|all
- `--chart` - Generate visual charts

### `config <action> [key] [value]`
Configure analyzer settings.

```bash
component-analyzer config set output.format json
component-analyzer config get paths.components
component-analyzer config list
```

## Configuration

Create a `.component-analyzer.json` file in your project root:

```json
{
  "paths": {
    "components": "src/components",
    "output": "docs/components",
    "registry": "docs/components/registry.json"
  },
  "analysis": {
    "recursive": true,
    "maxDepth": 10,
    "exclude": ["*.test.tsx", "*.stories.tsx"],
    "categorization": {
      "atomPatterns": ["Button*.tsx", "Input*.tsx"],
      "moleculePatterns": ["*Form.tsx", "*Card.tsx"],
      "organismPatterns": ["*Section.tsx", "*Page.tsx"]
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
      }
    },
    "autoFix": false,
    "severity": "warning"
  },
  "duplicates": {
    "threshold": 75,
    "factors": {
      "structural": 0.4,
      "visual": 0.3,
      "functional": 0.3
    }
  }
}
```

## Categorization Rules

Components are categorized based on atomic design principles:

### Atoms
- Basic UI elements that can't be broken down further
- Examples: Button, Input, Label, Icon
- Identified by: Single responsibility, no composition of other components

### Molecules
- Combinations of atoms working together
- Examples: SearchBar (Input + Button), FormField (Label + Input)
- Identified by: Composition of 2+ atoms, clear interaction patterns

### Organisms
- Complex, self-contained sections
- Examples: Header, Footer, ProductCard, NavigationMenu
- Identified by: Multiple molecules/atoms, standalone functionality

## API Usage

You can also use the analyzer programmatically:

```javascript
import { ComponentAnalyzer } from '@tchat/component-analyzer';

const analyzer = new ComponentAnalyzer();

// Analyze components
const results = await analyzer.analyze('src/components', {
  recursive: true,
  exclude: ['*.test.tsx']
});

// Get registry
const registry = analyzer.getRegistry();

// Find duplicates
const duplicates = await analyzer.findDuplicates({
  threshold: 80
});
```

## CI/CD Integration

Add to your GitHub Actions workflow:

```yaml
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

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

MIT © Tchat Team