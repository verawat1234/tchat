#!/usr/bin/env node

/**
 * Component Analyzer CLI
 * Analyze and categorize React components following atomic design principles
 */

import { Command } from 'commander';
import chalk from 'chalk';
import { version } from '../package.json';
import { AnalyzeCommand } from './commands/analyze';
import { ListCommand } from './commands/list';
import { DuplicatesCommand } from './commands/duplicates';
import { ValidateCommand } from './commands/validate';
import { GenerateCommand } from './commands/generate';
import { WatchCommand } from './commands/watch';
import { StatsCommand } from './commands/stats';
import { ConfigCommand } from './commands/config';

const program = new Command();

// Configure the CLI
program
  .name('component-analyzer')
  .description('Analyze and categorize React components following atomic design principles')
  .version(version)
  .option('--config <path>', 'Path to configuration file', '.component-analyzer.json')
  .option('--no-color', 'Disable colored output')
  .option('--json', 'Output raw JSON')
  .option('--debug', 'Enable debug logging');

// Analyze command
program
  .command('analyze [path]')
  .description('Analyze components in a directory and categorize them')
  .option('-r, --recursive', 'Include subdirectories (default: true)')
  .option('-o, --output <format>', 'Output format: json|markdown|both (default: both)')
  .option('-s, --save', 'Save results to registry (default: true)')
  .option('-v, --verbose', 'Verbose output')
  .option('-q, --quiet', 'Suppress non-error output')
  .option('--max-depth <n>', 'Maximum directory depth', parseInt)
  .option('--exclude <patterns>', 'Glob patterns to exclude')
  .option('--include-tests', 'Include test files')
  .action(async (path, options) => {
    try {
      const command = new AnalyzeCommand();
      await command.execute(path, options);
    } catch (error) {
      console.error(chalk.red('Error:'), error instanceof Error ? error.message : error);
      process.exit(1);
    }
  });

// List command
program
  .command('list [type]')
  .description('List components from the registry')
  .option('-c, --category <category>', 'Filter by category')
  .option('-s, --sort <field>', 'Sort by: name|usage|created (default: name)')
  .option('-l, --limit <n>', 'Maximum results (default: 100)', parseInt)
  .option('-f, --format <format>', 'Output format: table|json|list (default: table)')
  .option('--deprecated', 'Include deprecated components')
  .action(async (type, options) => {
    try {
      const command = new ListCommand();
      await command.execute(type, options);
    } catch (error) {
      console.error(chalk.red('Error:'), error instanceof Error ? error.message : error);
      process.exit(1);
    }
  });

// Duplicates command
program
  .command('duplicates')
  .description('Find duplicate or similar components')
  .option('-t, --threshold <n>', 'Similarity threshold % (default: 75)', parseInt)
  .option('--type <type>', 'Component type to check (default: all)')
  .option('--auto-merge', 'Automatically suggest merges')
  .option('-o, --output <format>', 'Output format: table|json|report (default: report)')
  .action(async (options) => {
    try {
      const command = new DuplicatesCommand();
      await command.execute(options);
    } catch (error) {
      console.error(chalk.red('Error:'), error instanceof Error ? error.message : error);
      process.exit(1);
    }
  });

// Validate command
program
  .command('validate [component-id]')
  .description('Validate components against consistency rules')
  .option('-r, --rules <rules>', 'Specific rules to apply (comma-separated)')
  .option('--fix', 'Attempt to auto-fix issues')
  .option('-s, --severity <level>', 'Minimum severity: error|warning|info (default: warning)')
  .option('-o, --output <format>', 'Output format: summary|detailed|json (default: summary)')
  .action(async (componentId, options) => {
    try {
      const command = new ValidateCommand();
      await command.execute(componentId, options);
    } catch (error) {
      console.error(chalk.red('Error:'), error instanceof Error ? error.message : error);
      process.exit(1);
    }
  });

// Generate command
program
  .command('generate [type]')
  .description('Generate documentation for components')
  .option('-c, --components <ids>', 'Specific component IDs (comma-separated)')
  .option('-o, --output <dir>', 'Output directory (default: ./docs/components)')
  .option('-f, --format <format>', 'Format: markdown|html|json (default: markdown)')
  .option('--examples', 'Include usage examples (default: true)')
  .option('--visuals', 'Include visual references (default: true)')
  .option('--force', 'Overwrite existing files')
  .action(async (type, options) => {
    try {
      const command = new GenerateCommand();
      await command.execute(type, options);
    } catch (error) {
      console.error(chalk.red('Error:'), error instanceof Error ? error.message : error);
      process.exit(1);
    }
  });

// Watch command
program
  .command('watch [path]')
  .description('Watch for component changes and update registry')
  .option('-i, --interval <n>', 'Check interval in seconds (default: 5)', parseInt)
  .option('--auto-fix', 'Automatically fix consistency issues')
  .option('-n, --notify', 'Send notifications on changes')
  .option('-w, --webhook <url>', 'Webhook URL for notifications')
  .action(async (path, options) => {
    try {
      const command = new WatchCommand();
      await command.execute(path, options);
    } catch (error) {
      console.error(chalk.red('Error:'), error instanceof Error ? error.message : error);
      process.exit(1);
    }
  });

// Stats command
program
  .command('stats')
  .description('Display statistics about components')
  .option('-p, --period <period>', 'Time period: today|week|month|all (default: all)')
  .option('-f, --format <format>', 'Output format: summary|detailed|json (default: summary)')
  .option('--chart', 'Generate visual charts')
  .action(async (options) => {
    try {
      const command = new StatsCommand();
      await command.execute(options);
    } catch (error) {
      console.error(chalk.red('Error:'), error instanceof Error ? error.message : error);
      process.exit(1);
    }
  });

// Config command
program
  .command('config <action> [key] [value]')
  .description('Configure analyzer settings')
  .action(async (action, key, value, options) => {
    try {
      const command = new ConfigCommand();
      await command.execute(action, key, value, options);
    } catch (error) {
      console.error(chalk.red('Error:'), error instanceof Error ? error.message : error);
      process.exit(1);
    }
  });

// Handle unknown commands
program.on('command:*', function () {
  console.error(chalk.red('Invalid command: %s'), program.args.join(' '));
  console.log('See --help for a list of available commands.');
  process.exit(1);
});

// Parse command line arguments
program.parse(process.argv);

// Show help if no command provided
if (!process.argv.slice(2).length) {
  program.outputHelp();
}