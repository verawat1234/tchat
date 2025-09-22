/**
 * Analyze command - scan and categorize components
 */

import * as fs from 'fs';
import * as path from 'path';
import chalk from 'chalk';
import ora from 'ora';
import { EventEmitter } from 'events';
import { ComponentAnalyzer } from '../core/ComponentAnalyzer';

export interface AnalyzeOptions {
  recursive?: boolean;
  output?: string;
  save?: boolean;
  verbose?: boolean;
  quiet?: boolean;
  maxDepth?: number;
  exclude?: string;
  includeTests?: boolean;
  config?: string;
}

export class AnalyzeCommand extends EventEmitter {
  private analyzer: ComponentAnalyzer;

  constructor() {
    super();
    this.analyzer = new ComponentAnalyzer();
  }

  async execute(targetPath?: string, options: AnalyzeOptions = {}): Promise<void> {
    // Default path
    const analyzePath = targetPath || 'src/components';

    // Validate path exists
    if (!fs.existsSync(analyzePath)) {
      throw new Error(`Path does not exist: ${analyzePath}`);
    }

    // Check if path is directory
    const stat = fs.statSync(analyzePath);
    if (!stat.isDirectory()) {
      throw new Error(`Path is not a directory: ${analyzePath}`);
    }

    // Load config if specified
    const config = await this.loadConfig(options.config);
    const mergedOptions = { ...config, ...options };

    // Parse exclude patterns
    if (typeof mergedOptions.exclude === 'string') {
      mergedOptions.exclude = mergedOptions.exclude.split(',');
    }

    // Show progress spinner if not quiet
    const spinner = !mergedOptions.quiet ? ora('Analyzing components...').start() : null;

    // Set up progress reporting
    if (mergedOptions.verbose) {
      this.analyzer.on('file:processed', ({ file, index, total }) => {
        if (spinner) {
          spinner.text = `Processing ${file} (${index}/${total})`;
        } else {
          console.log(chalk.gray(`Processing ${file} (${index}/${total})`));
        }
      });
    }

    this.emit('progress', { status: 'started' });

    try {
      // Run analysis
      const results = await this.analyzer.analyze(analyzePath, mergedOptions);

      if (spinner) {
        spinner.succeed(`✓ Analyzed ${results.componentsFound} components`);
      }

      // Display results based on output format
      this.displayResults(results, mergedOptions);

      // Save registry if requested
      if (mergedOptions.save !== false) {
        await this.analyzer.saveRegistry();
        if (!mergedOptions.quiet) {
          console.log(chalk.green('✓ Registry saved'));
        }
      }

      // Show errors if any
      if (results.errors && results.errors.length > 0) {
        console.log(chalk.yellow('\n⚠ Errors encountered:'));
        results.errors.forEach(error => {
          console.log(chalk.yellow(`  - ${error}`));
        });
      }

      this.emit('progress', { status: 'completed' });
    } catch (error) {
      if (spinner) {
        spinner.fail('Analysis failed');
      }
      throw error;
    }
  }

  private async loadConfig(configPath?: string): Promise<any> {
    const configFile = configPath || '.component-analyzer.json';

    if (fs.existsSync(configFile)) {
      const content = fs.readFileSync(configFile, 'utf8');
      const config = JSON.parse(content);
      return config.analysis || {};
    }

    return {};
  }

  private displayResults(results: any, options: AnalyzeOptions): void {
    if (options.quiet && options.output === 'json') {
      console.log(JSON.stringify(results));
      return;
    }

    const { componentsFound, categorized, uncategorized, duration } = results;

    switch (options.output) {
      case 'json':
        console.log(JSON.stringify(results, null, 2));
        break;

      case 'markdown':
        console.log('\n## Components Analysis\n');
        console.log(`- **Total Components**: ${componentsFound}`);
        console.log(`- **Atoms**: ${categorized.atoms}`);
        console.log(`- **Molecules**: ${categorized.molecules}`);
        console.log(`- **Organisms**: ${categorized.organisms}`);
        console.log(`- **Uncategorized**: ${uncategorized}`);
        console.log(`- **Duration**: ${duration}ms\n`);

        if (results.components) {
          results.components.forEach((comp: any) => {
            console.log(`### ${comp.name}`);
            console.log(`- **Type**: ${comp.type}`);
            console.log(`- **Path**: ${comp.path}`);
            console.log('');
          });
        }
        break;

      case 'table':
        if (results.components && results.components.length > 0) {
          console.log('\n');
          console.table(
            results.components.map((c: any) => ({
              Name: c.name,
              Type: c.type,
              Usage: c.usage
            }))
          );
        }
        break;

      default:
        // Default output
        if (!options.quiet) {
          console.log(chalk.cyan('\n📊 Analysis Results:'));
          console.log(chalk.white(`  Total Components: ${componentsFound}`));
          console.log(chalk.green(`  ⚛️  Atoms: ${categorized.atoms}`));
          console.log(chalk.blue(`  🧩 Molecules: ${categorized.molecules}`));
          console.log(chalk.magenta(`  🏗️  Organisms: ${categorized.organisms}`));
          if (uncategorized > 0) {
            console.log(chalk.yellow(`  ❓ Uncategorized: ${uncategorized}`));
          }
          console.log(chalk.gray(`  ⏱️  Duration: ${duration}ms`));
        }

        // Show detailed component list in verbose mode
        if (options.verbose && results.details) {
          console.log('\n📋 Component Details:');
          results.details.forEach((comp: any) => {
            console.log(`  ${comp.name} (${comp.type}) - ${comp.path}`);
          });
        }
    }
  }
}

// Stub ComponentAnalyzer for now
class ComponentAnalyzer extends EventEmitter {
  async analyze(path: string, options: any): Promise<any> {
    // This will be implemented with actual analysis logic
    return {
      componentsFound: 10,
      categorized: {
        atoms: 5,
        molecules: 3,
        organisms: 2
      },
      uncategorized: 0,
      errors: [],
      duration: 1234,
      registry: {
        id: 'test-registry',
        components: []
      }
    };
  }

  async saveRegistry(): Promise<void> {
    // Save registry implementation
  }
}