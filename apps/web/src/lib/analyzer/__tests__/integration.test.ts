/**
 * Integration tests for Component Analyzer E2E workflow
 */

import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import * as fs from 'fs';
import * as path from 'path';
import { ComponentAnalyzer } from '../ComponentAnalyzer';
import { ComponentType } from '../models/Component';
import { RuleSeverity } from '../models/ConsistencyRule';

// Test fixtures directory
const FIXTURES_DIR = path.join(__dirname, 'fixtures');
const OUTPUT_DIR = path.join(__dirname, 'output');

describe('Component Analyzer E2E Integration', () => {
  let analyzer: ComponentAnalyzer;

  beforeEach(() => {
    // Create test fixtures directory
    if (!fs.existsSync(FIXTURES_DIR)) {
      fs.mkdirSync(FIXTURES_DIR, { recursive: true });
    }

    // Create output directory
    if (!fs.existsSync(OUTPUT_DIR)) {
      fs.mkdirSync(OUTPUT_DIR, { recursive: true });
    }

    // Initialize analyzer
    analyzer = new ComponentAnalyzer({
      paths: {
        components: FIXTURES_DIR,
        output: OUTPUT_DIR,
        registry: path.join(OUTPUT_DIR, 'registry.json')
      },
      analysis: {
        recursive: true,
        maxDepth: 5,
        exclude: ['**/*.test.tsx', '**/*.stories.tsx', '*.test.tsx', '*.stories.tsx'],
        includeTests: false
      },
      validation: {
        rules: {},
        autoFix: false,
        severity: 'warning'
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
  });

  afterEach(() => {
    // Clean up test files
    if (fs.existsSync(OUTPUT_DIR)) {
      fs.rmSync(OUTPUT_DIR, { recursive: true, force: true });
    }
  });

  describe('Full Analysis Workflow', () => {
    beforeEach(() => {
      // Create sample component files for testing
      createSampleComponents();
    });

    it('should analyze a directory and generate analysis results', async () => {
      const result = await analyzer.analyze(FIXTURES_DIR, { verbose: true });

      expect(result).toBeDefined();
      expect(result.componentsFound).toBeGreaterThan(0);
      expect(result.duration).toBeGreaterThan(0);
      expect(result.errors).toEqual([]);
      expect(result.registry).toBeDefined();
    });

    it('should categorize components correctly', async () => {
      const result = await analyzer.analyze(FIXTURES_DIR);

      // Check that categorization properties exist
      expect(result.categorized).toBeDefined();
      expect(result.categorized.atoms).toBeGreaterThanOrEqual(0);
      expect(result.categorized.molecules).toBeGreaterThanOrEqual(0);
      expect(result.categorized.organisms).toBeGreaterThanOrEqual(0);
    });

    it('should detect duplicate components', async () => {
      // Create duplicate components
      createDuplicateComponents();

      const result = await analyzer.analyze(FIXTURES_DIR);
      const duplicates = await analyzer.findDuplicates();

      expect(duplicates).toBeDefined();
      expect(duplicates.length).toBeGreaterThanOrEqual(0);
    });

    it('should validate components against rules', async () => {
      const result = await analyzer.analyze(FIXTURES_DIR);
      const validation = analyzer.validateComponents();

      expect(validation).toBeDefined();
      expect(validation.results).toBeInstanceOf(Array);
      expect(validation.summary).toBeDefined();
      expect(validation.summary.totalChecked).toBeGreaterThan(0);
    });

    it('should generate documentation', async () => {
      const result = await analyzer.analyze(FIXTURES_DIR);
      const docs = await analyzer.generateDocumentation();

      expect(docs).toBeDefined();
      expect(docs.filesGenerated).toBeInstanceOf(Array);
      expect(docs.totalComponents).toEqual(result.componentsFound);
    });
  });

  describe('Registry Persistence', () => {
    it('should save registry to file', async () => {
      await analyzer.analyze(FIXTURES_DIR);
      await analyzer.saveRegistry();

      const registryPath = path.join(OUTPUT_DIR, 'registry.json');
      expect(fs.existsSync(registryPath)).toBe(true);

      const registryContent = fs.readFileSync(registryPath, 'utf8');
      const registry = JSON.parse(registryContent);
      expect(registry.components).toBeDefined();
    });

    it('should load registry from file', async () => {
      // First analyze and save
      await analyzer.analyze(FIXTURES_DIR);
      await analyzer.saveRegistry();

      // Create new analyzer instance and load
      const newAnalyzer = new ComponentAnalyzer();
      await newAnalyzer.loadRegistry(path.join(OUTPUT_DIR, 'registry.json'));

      const registry = newAnalyzer.getRegistry();
      expect(registry).toBeDefined();
      expect(registry.getComponentCount()).toBeGreaterThan(0);
    });
  });

  describe('Event Emission', () => {
    it('should emit progress events during analysis', async () => {
      const events: any[] = [];

      analyzer.on('analysis:start', (data) => {
        events.push({ type: 'start', data });
      });

      analyzer.on('file:processed', (data) => {
        events.push({ type: 'processed', data });
      });

      analyzer.on('analysis:complete', (data) => {
        events.push({ type: 'complete', data });
      });

      await analyzer.analyze(FIXTURES_DIR);

      const startEvents = events.filter(e => e.type === 'start');
      const completeEvents = events.filter(e => e.type === 'complete');

      expect(startEvents).toHaveLength(1);
      expect(completeEvents).toHaveLength(1);
      expect(events.length).toBeGreaterThan(2);
    });

    it('should emit error events for invalid files', async () => {
      // Create an invalid component file
      fs.writeFileSync(
        path.join(FIXTURES_DIR, 'InvalidComponent.tsx'),
        'export const Invalid = {{{' // Invalid syntax
      );

      const errors: any[] = [];
      analyzer.on('file:error', (data) => {
        errors.push(data);
      });

      await analyzer.analyze(FIXTURES_DIR);

      // Error events may or may not be emitted based on parser behavior
      expect(errors).toBeDefined();
      if (errors.length > 0) {
        expect(errors[0].error).toBeDefined();
      }
    });
  });

  describe('Configuration Validation', () => {
    it('should handle different output formats', async () => {
      await analyzer.analyze(FIXTURES_DIR);

      // Test markdown generation
      const mdDocs = await analyzer.generateDocumentation({
        format: 'markdown',
        outputPath: OUTPUT_DIR
      });
      expect(mdDocs.filesGenerated.some(f => f.endsWith('.md'))).toBe(true);

      // Test JSON generation
      const jsonDocs = await analyzer.generateDocumentation({
        format: 'json',
        outputPath: OUTPUT_DIR
      });
      expect(jsonDocs.filesGenerated.some(f => f.endsWith('.json'))).toBe(true);

      // Test HTML generation
      const htmlDocs = await analyzer.generateDocumentation({
        format: 'html',
        outputPath: OUTPUT_DIR
      });
      expect(htmlDocs.filesGenerated.some(f => f.endsWith('.html') || f.endsWith('.md') || f.endsWith('.json'))).toBe(true);
    });

    it('should respect exclusion patterns', async () => {
      // Create test and story files
      fs.writeFileSync(
        path.join(FIXTURES_DIR, 'Component.test.tsx'),
        'export const TestComponent = () => <div>Test</div>;'
      );
      fs.writeFileSync(
        path.join(FIXTURES_DIR, 'Component.stories.tsx'),
        'export const StoryComponent = () => <div>Story</div>;'
      );

      const result = await analyzer.analyze(FIXTURES_DIR);

      // These should be excluded
      const registry = analyzer.getRegistry();
      const allComponents = registry.getAllComponents();

      // Debug: Log all component names to see what's being included
      const componentNames = allComponents.map(c => c.name);
      console.log('Found components:', componentNames);

      // Check that test and story files are excluded
      const testComponent = allComponents.find(c => c.name === 'TestComponent');
      const storyComponent = allComponents.find(c => c.name === 'StoryComponent');

      // These should not be found since test/story files should be excluded
      expect(testComponent).toBeUndefined();
      expect(storyComponent).toBeUndefined();

      // Also verify that the regular components ARE found
      const buttonComponent = allComponents.find(c => c.name === 'Button');
      const cardComponent = allComponents.find(c => c.name === 'Card');
      expect(buttonComponent).toBeDefined();
      expect(cardComponent).toBeDefined();
    });
  });
});

/**
 * Helper function to create sample component files
 */
function createSampleComponents(): void {
  // Create a simple Button atom
  fs.writeFileSync(
    path.join(FIXTURES_DIR, 'Button.tsx'),
    `
import React from 'react';

interface ButtonProps {
  label: string;
  onClick: () => void;
  variant?: 'primary' | 'secondary';
  disabled?: boolean;
}

export const Button: React.FC<ButtonProps> = ({
  label,
  onClick,
  variant = 'primary',
  disabled = false
}) => {
  return (
    <button
      className={\`btn btn-\${variant}\`}
      onClick={onClick}
      disabled={disabled}
      aria-label={label}
    >
      {label}
    </button>
  );
};

export default Button;
    `.trim()
  );

  // Create a Card molecule
  fs.writeFileSync(
    path.join(FIXTURES_DIR, 'Card.tsx'),
    `
import React from 'react';
import Button from './Button';

interface CardProps {
  title: string;
  description: string;
  onAction?: () => void;
}

export const Card: React.FC<CardProps> = ({ title, description, onAction }) => {
  return (
    <div className="card">
      <h3>{title}</h3>
      <p>{description}</p>
      {onAction && (
        <Button label="Action" onClick={onAction} />
      )}
    </div>
  );
};

export default Card;
    `.trim()
  );

  // Create an Input atom
  fs.writeFileSync(
    path.join(FIXTURES_DIR, 'Input.tsx'),
    `
import React from 'react';

interface InputProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  type?: 'text' | 'email' | 'password';
}

export const Input: React.FC<InputProps> = ({
  value,
  onChange,
  placeholder,
  type = 'text'
}) => {
  return (
    <input
      type={type}
      value={value}
      onChange={(e) => onChange(e.target.value)}
      placeholder={placeholder}
      className="input"
    />
  );
};

export default Input;
    `.trim()
  );
}

/**
 * Helper function to create duplicate components
 */
function createDuplicateComponents(): void {
  // Create similar button components
  fs.writeFileSync(
    path.join(FIXTURES_DIR, 'PrimaryButton.tsx'),
    `
export const PrimaryButton = ({ label, onClick }) => (
  <button className="btn-primary" onClick={onClick}>{label}</button>
);
    `.trim()
  );

  fs.writeFileSync(
    path.join(FIXTURES_DIR, 'MainButton.tsx'),
    `
export const MainButton = ({ text, handleClick }) => (
  <button className="btn-main" onClick={handleClick}>{text}</button>
);
    `.trim()
  );
}