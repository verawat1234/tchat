/**
 * DocumentationGenerator - Generates documentation for components
 */

import { Component, ComponentType, PropDefinition } from '../models/Component';
import { Molecule } from '../models/Molecule';
import { ComponentRegistry } from '../models/ComponentRegistry';
import * as fs from 'fs/promises';
import * as path from 'path';

export interface DocumentationOptions {
  format: 'markdown' | 'json' | 'html';
  includeExamples: boolean;
  includeVisuals: boolean;
  outputPath: string;
}

export interface GenerationResult {
  outputPath: string;
  filesGenerated: string[];
  totalComponents: number;
}

export class DocumentationGenerator {
  private defaultOptions: DocumentationOptions = {
    format: 'markdown',
    includeExamples: true,
    includeVisuals: false,
    outputPath: 'docs/components'
  };

  /**
   * Generate documentation for all components in registry
   */
  async generateFromRegistry(
    registry: ComponentRegistry,
    options?: Partial<DocumentationOptions>
  ): Promise<GenerationResult> {
    const opts = { ...this.defaultOptions, ...options };
    const filesGenerated: string[] = [];

    // Ensure output directory exists
    await fs.mkdir(opts.outputPath, { recursive: true });

    // Group components by type
    const atoms = registry.getComponentsByType(ComponentType.ATOM);
    const molecules = registry.getComponentsByType(ComponentType.MOLECULE);
    const organisms = registry.getComponentsByType(ComponentType.ORGANISM);

    // Generate documentation based on format
    switch (opts.format) {
      case 'markdown':
        filesGenerated.push(...await this.generateMarkdown(atoms, molecules, organisms, opts));
        break;

      case 'json':
        filesGenerated.push(await this.generateJSON(registry, opts));
        break;

      case 'html':
        filesGenerated.push(...await this.generateHTML(atoms, molecules, organisms, opts));
        break;
    }

    return {
      outputPath: opts.outputPath,
      filesGenerated,
      totalComponents: registry.getComponentCount()
    };
  }

  /**
   * Generate markdown documentation
   */
  private async generateMarkdown(
    atoms: Component[],
    molecules: Component[],
    organisms: Component[],
    options: DocumentationOptions
  ): Promise<string[]> {
    const files: string[] = [];

    // Generate index file
    const indexPath = path.join(options.outputPath, 'index.md');
    await fs.writeFile(indexPath, this.generateMarkdownIndex(atoms, molecules, organisms));
    files.push(indexPath);

    // Generate atoms documentation
    if (atoms.length > 0) {
      const atomsPath = path.join(options.outputPath, 'atoms.md');
      await fs.writeFile(atomsPath, this.generateMarkdownForComponents(atoms, 'Atoms', options));
      files.push(atomsPath);
    }

    // Generate molecules documentation
    if (molecules.length > 0) {
      const moleculesPath = path.join(options.outputPath, 'molecules.md');
      await fs.writeFile(moleculesPath, this.generateMarkdownForComponents(molecules, 'Molecules', options));
      files.push(moleculesPath);
    }

    // Generate organisms documentation
    if (organisms.length > 0) {
      const organismsPath = path.join(options.outputPath, 'organisms.md');
      await fs.writeFile(organismsPath, this.generateMarkdownForComponents(organisms, 'Organisms', options));
      files.push(organismsPath);
    }

    return files;
  }

  /**
   * Generate markdown index
   */
  private generateMarkdownIndex(
    atoms: Component[],
    molecules: Component[],
    organisms: Component[]
  ): string {
    const lines: string[] = [
      '# Component Library Documentation',
      '',
      '## Overview',
      '',
      `This documentation covers all UI components following atomic design principles.`,
      '',
      '## Statistics',
      '',
      `- **Total Components**: ${atoms.length + molecules.length + organisms.length}`,
      `- **Atoms**: ${atoms.length}`,
      `- **Molecules**: ${molecules.length}`,
      `- **Organisms**: ${organisms.length}`,
      '',
      '## Component Categories',
      '',
      '### [Atoms](./atoms.md)',
      'Basic building blocks of the UI. These are the smallest functional units.',
      ''
    ];

    // List top atoms
    if (atoms.length > 0) {
      lines.push('**Featured Atoms:**');
      atoms.slice(0, 5).forEach(atom => {
        lines.push(`- ${atom.name} - ${atom.description || 'No description'}`);
      });
      lines.push('');
    }

    lines.push(
      '### [Molecules](./molecules.md)',
      'Combinations of atoms working together as a unit.',
      ''
    );

    // List top molecules
    if (molecules.length > 0) {
      lines.push('**Featured Molecules:**');
      molecules.slice(0, 5).forEach(molecule => {
        lines.push(`- ${molecule.name} - ${molecule.description || 'No description'}`);
      });
      lines.push('');
    }

    lines.push(
      '### [Organisms](./organisms.md)',
      'Complex, self-contained sections of the interface.',
      ''
    );

    // List top organisms
    if (organisms.length > 0) {
      lines.push('**Featured Organisms:**');
      organisms.slice(0, 5).forEach(organism => {
        lines.push(`- ${organism.name} - ${organism.description || 'No description'}`);
      });
      lines.push('');
    }

    lines.push(
      '## Usage Guidelines',
      '',
      '1. **Atoms** should be used as basic building blocks',
      '2. **Molecules** combine atoms for specific functionality',
      '3. **Organisms** are complete sections ready for page composition',
      '',
      '## Contributing',
      '',
      'To add or modify components, please follow the atomic design principles and update this documentation.',
      '',
      '---',
      `*Generated on ${new Date().toISOString()}*`
    );

    return lines.join('\n');
  }

  /**
   * Generate markdown for a set of components
   */
  private generateMarkdownForComponents(
    components: Component[],
    title: string,
    options: DocumentationOptions
  ): string {
    const lines: string[] = [
      `# ${title}`,
      '',
      `Total: ${components.length} components`,
      '',
      '---',
      ''
    ];

    // Sort components alphabetically
    const sorted = components.sort((a, b) => a.name.localeCompare(b.name));

    for (const component of sorted) {
      lines.push(...this.generateComponentMarkdown(component, options));
      lines.push('', '---', '');
    }

    lines.push(`*Generated on ${new Date().toISOString()}*`);

    return lines.join('\n');
  }

  /**
   * Generate markdown for a single component
   */
  private generateComponentMarkdown(component: Component, options: DocumentationOptions): string[] {
    const lines: string[] = [];

    // Component header
    lines.push(`## ${component.name}`);

    if (component.deprecated) {
      lines.push('', '> ⚠️ **DEPRECATED**' + (component.deprecationMessage ? `: ${component.deprecationMessage}` : ''));
    }

    lines.push('', component.description || '*No description available*', '');

    // Metadata
    lines.push('### Metadata');
    lines.push('');
    lines.push(`- **Type**: ${component.type}`);
    lines.push(`- **Category**: ${component.category}`);
    lines.push(`- **File**: \`${component.filePath}\``);
    lines.push(`- **Usage Count**: ${component.usageCount}`);
    lines.push(`- **Version**: ${component.version}`);
    lines.push('');

    // Props
    if (component.props.length > 0) {
      lines.push('### Props');
      lines.push('');
      lines.push('| Name | Type | Required | Default | Description |');
      lines.push('|------|------|----------|---------|-------------|');

      for (const prop of component.props) {
        lines.push(this.generatePropRow(prop));
      }

      lines.push('');
    }

    // Molecule-specific information
    if (component.type === ComponentType.MOLECULE && (component as any).composition) {
      const molecule = component as unknown as Molecule;
      if (molecule.composition && molecule.composition.length > 0) {
        lines.push('### Composition');
        lines.push('');
        lines.push('This molecule is composed of the following atoms:');
        lines.push('');

        for (const comp of molecule.composition) {
          lines.push(`- **${comp.atomId}** (${comp.quantity}x) - ${comp.role}`);
        }

        lines.push('');
      }
    }

    // Usage examples
    if (options.includeExamples) {
      lines.push('### Usage Example');
      lines.push('');
      lines.push('```jsx');
      lines.push(this.generateUsageExample(component));
      lines.push('```');
      lines.push('');
    }

    // Dependencies
    if (component.dependencies.length > 0) {
      lines.push('### Dependencies');
      lines.push('');
      component.dependencies.forEach(dep => {
        lines.push(`- ${dep}`);
      });
      lines.push('');
    }

    return lines;
  }

  /**
   * Generate prop table row
   */
  private generatePropRow(prop: PropDefinition): string {
    const name = prop.name;
    const type = `\`${prop.type}\``;
    const required = prop.required ? '✓' : '';
    const defaultValue = prop.defaultValue !== undefined ? `\`${JSON.stringify(prop.defaultValue)}\`` : '-';
    const description = prop.description || '-';

    return `| ${name} | ${type} | ${required} | ${defaultValue} | ${description} |`;
  }

  /**
   * Generate usage example for a component
   */
  private generateUsageExample(component: Component): string {
    const props = component.props
      .filter(p => p.required)
      .map(p => {
        if (p.type === 'string') return `${p.name}="${p.defaultValue || 'value'}"`;
        if (p.type === 'boolean') return p.name;
        if (p.type.includes('=>')) return `${p.name}={() => console.log('${p.name}')}`;
        return `${p.name}={${p.defaultValue || 'value'}}`;
      })
      .join('\n    ');

    if (props) {
      return `<${component.name}\n    ${props}\n/>`;
    } else {
      return `<${component.name} />`;
    }
  }

  /**
   * Generate JSON documentation
   */
  private async generateJSON(
    registry: ComponentRegistry,
    options: DocumentationOptions
  ): Promise<string> {
    const outputPath = path.join(options.outputPath, 'components.json');
    const json = registry.toJSON();

    await fs.writeFile(outputPath, JSON.stringify(json, null, 2));

    return outputPath;
  }

  /**
   * Generate HTML documentation
   */
  private async generateHTML(
    atoms: Component[],
    molecules: Component[],
    organisms: Component[],
    options: DocumentationOptions
  ): Promise<string[]> {
    const files: string[] = [];

    // Generate main HTML file
    const indexPath = path.join(options.outputPath, 'index.html');
    await fs.writeFile(indexPath, this.generateHTMLPage(atoms, molecules, organisms));
    files.push(indexPath);

    // Generate CSS file
    const cssPath = path.join(options.outputPath, 'styles.css');
    await fs.writeFile(cssPath, this.generateCSS());
    files.push(cssPath);

    return files;
  }

  /**
   * Generate HTML page
   */
  private generateHTMLPage(
    atoms: Component[],
    molecules: Component[],
    organisms: Component[]
  ): string {
    return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Component Library Documentation</title>
    <link rel="stylesheet" href="styles.css">
</head>
<body>
    <header>
        <h1>Component Library Documentation</h1>
        <p>Generated on ${new Date().toLocaleDateString()}</p>
    </header>

    <nav>
        <ul>
            <li><a href="#atoms">Atoms (${atoms.length})</a></li>
            <li><a href="#molecules">Molecules (${molecules.length})</a></li>
            <li><a href="#organisms">Organisms (${organisms.length})</a></li>
        </ul>
    </nav>

    <main>
        <section id="atoms">
            <h2>Atoms</h2>
            ${this.generateHTMLComponents(atoms)}
        </section>

        <section id="molecules">
            <h2>Molecules</h2>
            ${this.generateHTMLComponents(molecules)}
        </section>

        <section id="organisms">
            <h2>Organisms</h2>
            ${this.generateHTMLComponents(organisms)}
        </section>
    </main>

    <footer>
        <p>Total Components: ${atoms.length + molecules.length + organisms.length}</p>
    </footer>
</body>
</html>`;
  }

  /**
   * Generate HTML for components
   */
  private generateHTMLComponents(components: Component[]): string {
    if (components.length === 0) {
      return '<p>No components in this category</p>';
    }

    return components.map(component => `
        <article class="component">
            <h3>${component.name}</h3>
            ${component.deprecated ? '<span class="deprecated">Deprecated</span>' : ''}
            <p>${component.description || 'No description'}</p>
            <dl>
                <dt>Category:</dt><dd>${component.category}</dd>
                <dt>Usage:</dt><dd>${component.usageCount} times</dd>
                <dt>File:</dt><dd><code>${component.filePath}</code></dd>
            </dl>
        </article>
    `).join('\n');
  }

  /**
   * Generate CSS styles
   */
  private generateCSS(): string {
    return `
body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    line-height: 1.6;
    color: #333;
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
}

header {
    background: #f4f4f4;
    padding: 20px;
    border-radius: 5px;
    margin-bottom: 20px;
}

nav ul {
    list-style: none;
    padding: 0;
    display: flex;
    gap: 20px;
}

nav a {
    color: #0066cc;
    text-decoration: none;
}

nav a:hover {
    text-decoration: underline;
}

.component {
    background: #fff;
    border: 1px solid #ddd;
    border-radius: 5px;
    padding: 15px;
    margin-bottom: 15px;
}

.component h3 {
    margin-top: 0;
    color: #0066cc;
}

.deprecated {
    background: #ff6b6b;
    color: white;
    padding: 2px 8px;
    border-radius: 3px;
    font-size: 12px;
}

code {
    background: #f4f4f4;
    padding: 2px 5px;
    border-radius: 3px;
    font-family: 'Consolas', 'Monaco', monospace;
}

dl {
    display: grid;
    grid-template-columns: 100px 1fr;
    gap: 10px;
}

dt {
    font-weight: bold;
}

dd {
    margin: 0;
}

footer {
    margin-top: 40px;
    padding-top: 20px;
    border-top: 1px solid #ddd;
    text-align: center;
    color: #666;
}
    `;
  }
}