/**
 * Config command - manage configuration
 */

import chalk from 'chalk';
import * as fs from 'fs';
import * as path from 'path';

export class ConfigCommand {
  async execute(action: string, key?: string, value?: string, options: any = {}): Promise<void> {
    switch (action) {
      case 'list':
        await this.listConfig();
        break;
      case 'get':
        if (!key) {
          console.error(chalk.red('Error: Key is required for get action'));
          process.exit(1);
        }
        await this.getConfig(key);
        break;
      case 'set':
        if (!key || value === undefined) {
          console.error(chalk.red('Error: Key and value are required for set action'));
          process.exit(1);
        }
        await this.setConfig(key, value);
        break;
      case 'reset':
        await this.resetConfig();
        break;
      default:
        console.error(chalk.red(`Error: Unknown action: ${action}`));
        console.log('Available actions: list, get, set, reset');
        process.exit(1);
    }
  }

  private async listConfig(): Promise<void> {
    const config = await this.loadConfig();
    console.log(chalk.cyan('Current Configuration:'));
    console.log(JSON.stringify(config, null, 2));
  }

  private async getConfig(key: string): Promise<void> {
    const config = await this.loadConfig();
    const value = this.getNestedValue(config, key);
    
    if (value !== undefined) {
      console.log(value);
    } else {
      console.error(chalk.red(`Error: Key not found: ${key}`));
      process.exit(1);
    }
  }

  private async setConfig(key: string, value: string): Promise<void> {
    const config = await this.loadConfig();
    this.setNestedValue(config, key, value);
    await this.saveConfig(config);
    console.log(chalk.green(`✓ Configuration updated: ${key} = ${value}`));
  }

  private async resetConfig(): Promise<void> {
    const defaultConfig = {
      paths: {
        components: 'src/components',
        output: 'docs/components',
        registry: 'docs/components/registry.json'
      },
      analysis: {
        recursive: true,
        maxDepth: 10,
        exclude: ['*.test.tsx', '*.stories.tsx'],
        includeTests: false
      },
      validation: {
        autoFix: false,
        severity: 'warning'
      },
      output: {
        format: 'both',
        colors: true,
        verbose: false
      }
    };

    await this.saveConfig(defaultConfig);
    console.log(chalk.green('✓ Configuration reset to defaults'));
  }

  private async loadConfig(): Promise<any> {
    const configPath = '.component-analyzer.json';
    
    if (fs.existsSync(configPath)) {
      const content = fs.readFileSync(configPath, 'utf8');
      return JSON.parse(content);
    }

    return {};
  }

  private async saveConfig(config: any): Promise<void> {
    const configPath = '.component-analyzer.json';
    fs.writeFileSync(configPath, JSON.stringify(config, null, 2));
  }

  private getNestedValue(obj: any, key: string): any {
    const keys = key.split('.');
    let value = obj;

    for (const k of keys) {
      if (value && typeof value === 'object' && k in value) {
        value = value[k];
      } else {
        return undefined;
      }
    }

    return value;
  }

  private setNestedValue(obj: any, key: string, value: string): void {
    const keys = key.split('.');
    let current = obj;

    for (let i = 0; i < keys.length - 1; i++) {
      const k = keys[i];
      if (!(k in current) || typeof current[k] !== 'object') {
        current[k] = {};
      }
      current = current[k];
    }

    const lastKey = keys[keys.length - 1];
    
    // Try to parse value as JSON first
    try {
      current[lastKey] = JSON.parse(value);
    } catch {
      // If not JSON, treat as string
      current[lastKey] = value;
    }
  }
}