#!/usr/bin/env node

/**
 * Component Analyzer CLI
 * Command-line interface for the component analyzer
 */

import { ComponentAnalyzer } from './ComponentAnalyzer';
import * as path from 'path';
import * as fs from 'fs';

// ANSI color codes for terminal output
const colors = {
  reset: '\x1b[0m',
  bright: '\x1b[1m',
  dim: '\x1b[2m',
  red: '\x1b[31m',
  green: '\x1b[32m',
  yellow: '\x1b[33m',
  blue: '\x1b[34m',
  cyan: '\x1b[36m'
};

function printHeader() {
  console.log(`${colors.cyan}${colors.bright}`);
  console.log('╔════════════════════════════════════╗');
  console.log('║    Component Analyzer CLI v1.0     ║');
  console.log('╚════════════════════════════════════╝');
  console.log(`${colors.reset}\n`);
}

async function main() {
  const command = process.argv[2];
  const targetPath = process.argv[3] || 'src/components';

  printHeader();

  if (!command || command === 'help') {
    showHelp();
    process.exit(0);
  }

  const analyzer = new ComponentAnalyzer({
    paths: {
      components: targetPath,
      output: 'docs/components',
      registry: 'docs/components/registry.json'
    },
    analysis: {
      recursive: true,
      maxDepth: 10,
      exclude: ['*.test.tsx', '*.stories.tsx', '*.spec.tsx'],
      includeTests: false
    },
    duplicates: {
      threshold: 75,
      autoMerge: false
    },
    output: {
      format: 'both',
      includeVisuals: false,
      includeExamples: true,
      colors: true,
      verbose: true
    }
  });

  // Add progress listeners
  analyzer.on('analysis:start', (data) => {
    console.log(`${colors.blue}🔍 Starting analysis of ${data.totalFiles} files...${colors.reset}`);
  });

  analyzer.on('file:processed', (data) => {
    process.stdout.write(`${colors.dim}  Processing: ${data.file} (${data.index}/${data.total})${colors.reset}\r`);
  });

  analyzer.on('file:error', (data) => {
    console.log(`${colors.red}  ❌ Error in ${data.file}: ${data.error}${colors.reset}`);
  });

  try {
    switch (command) {
      case 'analyze':
        await runAnalysis(analyzer, targetPath);
        break;
      case 'validate':
        await runValidation(analyzer, targetPath);
        break;
      case 'duplicates':
        await findDuplicates(analyzer, targetPath);
        break;
      case 'docs':
        await generateDocs(analyzer, targetPath);
        break;
      default:
        console.error(`${colors.red}Unknown command: ${command}${colors.reset}`);
        showHelp();
        process.exit(1);
    }
  } catch (error: any) {
    console.error(`${colors.red}${colors.bright}❌ Error: ${error.message}${colors.reset}`);
    if (error.stack) {
      console.error(`${colors.dim}${error.stack}${colors.reset}`);
    }
    process.exit(1);
  }
}

function showHelp() {
  console.log(`${colors.bright}Usage:${colors.reset}`);
  console.log('  tsx cli.ts <command> [path]\n');

  console.log(`${colors.bright}Commands:${colors.reset}`);
  console.log(`  ${colors.cyan}analyze${colors.reset}    Analyze all components in the target directory`);
  console.log(`  ${colors.cyan}validate${colors.reset}   Validate components against consistency rules`);
  console.log(`  ${colors.cyan}duplicates${colors.reset} Find duplicate components`);
  console.log(`  ${colors.cyan}docs${colors.reset}       Generate component documentation\n`);

  console.log(`${colors.bright}Examples:${colors.reset}`);
  console.log('  tsx cli.ts analyze src/components');
  console.log('  tsx cli.ts validate src/components');
  console.log('  tsx cli.ts duplicates');
  console.log('  tsx cli.ts docs\n');

  console.log(`${colors.bright}Options:${colors.reset}`);
  console.log('  path       Target directory (default: src/components)');
}

async function runAnalysis(analyzer: ComponentAnalyzer, targetPath: string) {
  console.log(`${colors.bright}📊 Analyzing components in ${targetPath}...${colors.reset}\n`);

  const absolutePath = path.resolve(targetPath);

  if (!fs.existsSync(absolutePath)) {
    throw new Error(`Directory not found: ${absolutePath}`);
  }

  const result = await analyzer.analyze(absolutePath, { verbose: true });

  console.log('\n');
  console.log(`${colors.green}${colors.bright}✅ Analysis Complete!${colors.reset}\n`);

  console.log(`${colors.bright}📈 Summary:${colors.reset}`);
  console.log(`  Total components found: ${colors.cyan}${result.componentsFound}${colors.reset}`);
  console.log(`  ├─ Atoms:     ${colors.green}${result.categorized.atoms}${colors.reset}`);
  console.log(`  ├─ Molecules: ${colors.blue}${result.categorized.molecules}${colors.reset}`);
  console.log(`  └─ Organisms: ${colors.yellow}${result.categorized.organisms}${colors.reset}`);

  if (result.uncategorized > 0) {
    console.log(`  ⚠️  Uncategorized: ${colors.yellow}${result.uncategorized}${colors.reset}`);
  }

  console.log(`\n  ⏱️  Duration: ${result.duration}ms`);

  if (result.errors.length > 0) {
    console.log(`\n${colors.red}Errors encountered:${colors.reset}`);
    result.errors.forEach(err => console.log(`  - ${err}`));
  }

  // Save registry
  await analyzer.saveRegistry();
  console.log(`\n${colors.dim}Registry saved to: docs/components/registry.json${colors.reset}`);
}

async function runValidation(analyzer: ComponentAnalyzer, targetPath: string) {
  console.log(`${colors.bright}✓ Validating components...${colors.reset}\n`);

  const absolutePath = path.resolve(targetPath);

  // First analyze to populate registry
  await analyzer.analyze(absolutePath, { verbose: false });

  const validation = analyzer.validateComponents();
  const { summary } = validation;

  console.log(`${colors.bright}📋 Validation Results:${colors.reset}`);
  console.log(`  Total checked: ${summary.totalChecked}`);
  console.log(`  ✅ Passed: ${colors.green}${summary.passed}${colors.reset}`);
  console.log(`  ❌ Failed: ${colors.red}${summary.failed}${colors.reset}`);

  if (summary.errorCount > 0 || summary.warningCount > 0) {
    console.log(`\n${colors.bright}Issues by severity:${colors.reset}`);
    if (summary.errorCount > 0) {
      console.log(`  🔴 Errors:   ${colors.red}${summary.errorCount}${colors.reset}`);
    }
    if (summary.warningCount > 0) {
      console.log(`  🟡 Warnings: ${colors.yellow}${summary.warningCount}${colors.reset}`);
    }
    if (summary.infoCount > 0) {
      console.log(`  🔵 Info:     ${colors.blue}${summary.infoCount}${colors.reset}`);
    }
  }

  // Show details for failed components
  const failedResults = validation.results.filter((r: any) => !r.valid);
  if (failedResults.length > 0) {
    console.log(`\n${colors.bright}Failed Components:${colors.reset}`);
    failedResults.slice(0, 5).forEach((result: any) => {
      console.log(`  ${colors.yellow}${result.componentId}:${colors.reset}`);
      result.violations.forEach((v: any) => {
        console.log(`    - ${v.message}`);
      });
    });

    if (failedResults.length > 5) {
      console.log(`  ${colors.dim}... and ${failedResults.length - 5} more${colors.reset}`);
    }
  }

  process.exit(summary.failed > 0 ? 1 : 0);
}

async function findDuplicates(analyzer: ComponentAnalyzer, targetPath: string) {
  console.log(`${colors.bright}🔎 Finding duplicate components...${colors.reset}\n`);

  const absolutePath = path.resolve(targetPath);

  // First analyze to populate registry
  await analyzer.analyze(absolutePath, { verbose: false });

  const duplicates = await analyzer.findDuplicates({ threshold: 70 });

  if (duplicates.length === 0) {
    console.log(`${colors.green}✅ No duplicate components found${colors.reset}`);
    return;
  }

  console.log(`${colors.yellow}⚠️  Found ${duplicates.length} potential duplicate groups:${colors.reset}\n`);

  duplicates.forEach((group: any, index: number) => {
    console.log(`${colors.bright}Group ${index + 1}${colors.reset} (${group.similarity}% similarity):`);

    group.components.forEach((comp: any) => {
      console.log(`  - ${colors.cyan}${comp.name}${colors.reset} (${comp.filePath})`);
    });

    console.log(`  ${colors.dim}Suggested merge: Keep ${group.suggestedMerge}${colors.reset}`);

    if (group.reasoning.length > 0) {
      console.log(`  ${colors.dim}Reasoning:${colors.reset}`);
      group.reasoning.forEach((r: any) => console.log(`    - ${r}`));
    }

    console.log('');
  });
}

async function generateDocs(analyzer: ComponentAnalyzer, targetPath: string) {
  console.log(`${colors.bright}📝 Generating component documentation...${colors.reset}\n`);

  const absolutePath = path.resolve(targetPath);

  // First analyze to populate registry
  console.log(`${colors.dim}Analyzing components...${colors.reset}`);
  await analyzer.analyze(absolutePath, { verbose: false });

  console.log(`${colors.dim}Generating documentation...${colors.reset}`);
  const docs = await analyzer.generateDocumentation({
    format: 'markdown',
    includeExamples: true,
    outputPath: 'docs/components'
  });

  console.log(`\n${colors.green}✅ Documentation generated successfully!${colors.reset}`);
  console.log(`\n${colors.bright}📄 Files created:${colors.reset}`);

  docs.filesGenerated.forEach((file: any) => {
    const filename = path.basename(file);
    console.log(`  - ${colors.cyan}${filename}${colors.reset}`);
  });

  console.log(`\n  Total components documented: ${colors.bright}${docs.totalComponents}${colors.reset}`);
  console.log(`  Output directory: ${colors.dim}${docs.outputPath}${colors.reset}`);
}

// Run the CLI
main().catch(error => {
  console.error(`${colors.red}${colors.bright}Fatal error:${colors.reset}`, error);
  process.exit(1);
});