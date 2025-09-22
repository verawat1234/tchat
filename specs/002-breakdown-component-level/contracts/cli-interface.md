# CLI Interface Contract

**Tool**: component-analyzer
**Version**: 1.0.0

## Command Structure

```bash
component-analyzer [command] [options]
```

## Commands

### analyze
Analyze components in a directory and categorize them.

```bash
component-analyzer analyze [path] [options]

Arguments:
  path              Path to analyze (default: ./src/components)

Options:
  --recursive       Include subdirectories (default: true)
  --output, -o      Output format: json|markdown|both (default: both)
  --save, -s        Save results to registry (default: true)
  --verbose, -v     Verbose output
  --quiet, -q       Suppress non-error output
  --max-depth       Maximum directory depth (default: 10)
  --exclude         Glob patterns to exclude
  --include-tests   Include test files (default: false)

Examples:
  component-analyzer analyze
  component-analyzer analyze apps/web/src/components -o json
  component-analyzer analyze --exclude "*.test.tsx" --verbose
```

### list
List components from the registry.

```bash
component-analyzer list [type] [options]

Arguments:
  type              Component type: all|atoms|molecules|organisms (default: all)

Options:
  --category, -c    Filter by category
  --sort, -s        Sort by: name|usage|created (default: name)
  --limit, -l       Maximum results (default: 100)
  --format, -f      Output format: table|json|list (default: table)
  --deprecated      Include deprecated components

Examples:
  component-analyzer list molecules
  component-analyzer list --category form --sort usage
  component-analyzer list atoms -f json -l 10
```

### duplicates
Find duplicate or similar components.

```bash
component-analyzer duplicates [options]

Options:
  --threshold, -t   Similarity threshold % (default: 75)
  --type            Component type to check (default: all)
  --auto-merge      Automatically suggest merges
  --output, -o      Output format: table|json|report (default: report)

Examples:
  component-analyzer duplicates
  component-analyzer duplicates --threshold 90 --type molecule
  component-analyzer duplicates --auto-merge -o json
```

### validate
Validate components against consistency rules.

```bash
component-analyzer validate [component-id] [options]

Arguments:
  component-id      Specific component to validate (optional)

Options:
  --rules, -r       Specific rules to apply (comma-separated)
  --fix             Attempt to auto-fix issues
  --severity, -s    Minimum severity: error|warning|info (default: warning)
  --output, -o      Output format: summary|detailed|json (default: summary)

Examples:
  component-analyzer validate
  component-analyzer validate btn-primary --rules naming,styling
  component-analyzer validate --fix --severity error
```

### generate
Generate documentation for components.

```bash
component-analyzer generate [type] [options]

Arguments:
  type              Documentation type: docs|storybook|types (default: docs)

Options:
  --components, -c  Specific component IDs (comma-separated)
  --output, -o      Output directory (default: ./docs/components)
  --format, -f      Format: markdown|html|json (default: markdown)
  --examples        Include usage examples (default: true)
  --visuals         Include visual references (default: true)
  --force           Overwrite existing files

Examples:
  component-analyzer generate docs
  component-analyzer generate storybook --components btn-primary,input-text
  component-analyzer generate types -o ./types -f typescript
```

### watch
Watch for component changes and update registry.

```bash
component-analyzer watch [path] [options]

Arguments:
  path              Path to watch (default: ./src/components)

Options:
  --interval, -i    Check interval in seconds (default: 5)
  --auto-fix        Automatically fix consistency issues
  --notify, -n      Send notifications on changes
  --webhook, -w     Webhook URL for notifications

Examples:
  component-analyzer watch
  component-analyzer watch apps/web/src --interval 10
  component-analyzer watch --auto-fix --notify
```

### stats
Display statistics about components.

```bash
component-analyzer stats [options]

Options:
  --period, -p      Time period: today|week|month|all (default: all)
  --format, -f      Output format: summary|detailed|json (default: summary)
  --chart           Generate visual charts

Examples:
  component-analyzer stats
  component-analyzer stats --period week --chart
  component-analyzer stats -f json
```

### config
Configure analyzer settings.

```bash
component-analyzer config [action] [key] [value]

Arguments:
  action            get|set|list|reset
  key               Configuration key
  value             Configuration value

Examples:
  component-analyzer config list
  component-analyzer config get output.format
  component-analyzer config set output.format json
  component-analyzer config reset
```

## Global Options

```bash
Options:
  --config, -c      Config file path (default: .component-analyzer.json)
  --no-color        Disable colored output
  --json            Output raw JSON
  --help, -h        Show help
  --version, -V     Show version
  --debug           Enable debug logging
```

## Output Formats

### Table (default for terminal)
```
┌─────────────┬──────────┬──────────┬───────┬────────┐
│ ID          │ Name     │ Type     │ Usage │ Status │
├─────────────┼──────────┼──────────┼───────┼────────┤
│ btn-primary │ Button   │ atom     │ 23    │ active │
│ search-bar  │ SearchBar│ molecule │ 5     │ active │
└─────────────┴──────────┴──────────┴───────┴────────┘
```

### JSON
```json
{
  "components": [
    {
      "id": "btn-primary",
      "name": "Button",
      "type": "atom",
      "usage": 23,
      "status": "active"
    }
  ],
  "total": 1
}
```

### Markdown
```markdown
## Components

### btn-primary (Button)
- **Type**: atom
- **Usage**: 23 times
- **Status**: active
```

## Exit Codes

- `0`: Success
- `1`: General error
- `2`: Invalid arguments
- `3`: Component not found
- `4`: Validation errors
- `5`: File system error
- `6`: Configuration error

## Configuration File

`.component-analyzer.json`
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
    "includeTests": false
  },
  "validation": {
    "rules": ["naming", "structure", "styling", "accessibility"],
    "autoFix": false,
    "severity": "warning"
  },
  "output": {
    "format": "both",
    "colors": true,
    "verbose": false
  },
  "duplicates": {
    "threshold": 75,
    "autoMerge": false
  }
}
```

## Environment Variables

- `ANALYZER_CONFIG`: Path to config file
- `ANALYZER_REGISTRY`: Path to registry file
- `ANALYZER_OUTPUT_DIR`: Default output directory
- `ANALYZER_LOG_LEVEL`: Log level (debug|info|warn|error)
- `NO_COLOR`: Disable colored output