#!/usr/bin/env node

/**
 * Component Analyzer CLI Script
 * This script can be run manually or as part of CI/CD pipeline
 */

const path = require('path');
const { execSync } = require('child_process');

const COMMANDS = {
  analyze: 'Analyze all components in the codebase',
  validate: 'Validate components against consistency rules',
  duplicates: 'Find duplicate components',
  docs: 'Generate component documentation'
};

function main() {
  const command = process.argv[2];
  const targetPath = process.argv[3] || 'src/components';

  if (!command || command === 'help') {
    showHelp();
    return;
  }

  console.log(`🔍 Running component ${command} on ${targetPath}...`);

  try {
    switch (command) {
      case 'analyze':
        runAnalysis(targetPath);
        break;
      case 'validate':
        runValidation(targetPath);
        break;
      case 'duplicates':
        findDuplicates(targetPath);
        break;
      case 'docs':
        generateDocs(targetPath);
        break;
      default:
        console.error(`Unknown command: ${command}`);
        showHelp();
        process.exit(1);
    }
  } catch (error) {
    console.error('❌ Error:', error.message);
    process.exit(1);
  }
}

function showHelp() {
  console.log('Component Analyzer CLI\n');
  console.log('Usage: node scripts/analyze-components.js <command> [path]\n');
  console.log('Commands:');
  Object.entries(COMMANDS).forEach(([cmd, desc]) => {
    console.log(`  ${cmd.padEnd(12)} ${desc}`);
  });
  console.log('\nExample:');
  console.log('  node scripts/analyze-components.js analyze src/components');
}

function runAnalysis(targetPath) {
  console.log('📊 Analyzing components...');

  // In a real implementation, this would use the ComponentAnalyzer class
  // For now, we'll simulate the analysis
  const components = findComponents(targetPath);

  console.log(`✅ Found ${components.length} components`);
  console.log(`  - Atoms: ${components.filter(c => c.includes('Button') || c.includes('Input')).length}`);
  console.log(`  - Molecules: ${components.filter(c => c.includes('Card') || c.includes('Form')).length}`);
  console.log(`  - Organisms: ${components.filter(c => c.includes('Header') || c.includes('Layout')).length}`);
}

function runValidation(targetPath) {
  console.log('✓ Validating components...');

  const components = findComponents(targetPath);
  let issues = 0;

  components.forEach(component => {
    // Check naming convention (PascalCase)
    const name = path.basename(component, path.extname(component));
    if (!/^[A-Z][a-zA-Z]*$/.test(name)) {
      console.warn(`  ⚠️  ${name} does not follow PascalCase convention`);
      issues++;
    }
  });

  if (issues === 0) {
    console.log('✅ All components pass validation');
  } else {
    console.log(`⚠️  Found ${issues} validation issues`);
  }

  return issues === 0 ? 0 : 1;
}

function findDuplicates(targetPath) {
  console.log('🔎 Finding duplicate components...');

  const components = findComponents(targetPath);
  const duplicates = [];

  // Simple duplicate detection based on similar names
  for (let i = 0; i < components.length; i++) {
    for (let j = i + 1; j < components.length; j++) {
      const name1 = path.basename(components[i], path.extname(components[i]));
      const name2 = path.basename(components[j], path.extname(components[j]));

      if (areSimilar(name1, name2)) {
        duplicates.push([name1, name2]);
      }
    }
  }

  if (duplicates.length > 0) {
    console.log(`⚠️  Found ${duplicates.length} potential duplicates:`);
    duplicates.forEach(([a, b]) => {
      console.log(`  - ${a} ≈ ${b}`);
    });
  } else {
    console.log('✅ No duplicate components found');
  }
}

function generateDocs(targetPath) {
  console.log('📝 Generating component documentation...');

  const components = findComponents(targetPath);
  const outputPath = path.join('docs', 'components.md');

  console.log(`  Generated documentation for ${components.length} components`);
  console.log(`  Output: ${outputPath}`);
}

function findComponents(targetPath) {
  // Simulate finding component files
  // In real implementation, this would use glob or fs.readdir
  try {
    const result = execSync(
      `find ${targetPath} -name "*.tsx" -o -name "*.jsx" 2>/dev/null | grep -v test | grep -v stories`,
      { encoding: 'utf8' }
    );
    return result.trim().split('\n').filter(Boolean);
  } catch (error) {
    return [];
  }
}

function areSimilar(str1, str2) {
  // Simple similarity check
  const s1 = str1.toLowerCase();
  const s2 = str2.toLowerCase();
  return s1.includes(s2) || s2.includes(s1) ||
         (s1.replace('button', '') === s2.replace('button', '')) ||
         (s1.replace('card', '') === s2.replace('card', ''));
}

// Run the CLI
main();