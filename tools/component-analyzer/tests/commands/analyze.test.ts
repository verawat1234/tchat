import { describe, it, expect, beforeEach, vi } from 'vitest';
import { AnalyzeCommand } from '../../src/commands/analyze';
import { ComponentAnalyzer } from '../../src/core/ComponentAnalyzer';
import * as fs from 'fs';
import * as path from 'path';

vi.mock('fs');
vi.mock('../../src/core/ComponentAnalyzer');

describe('Analyze Command', () => {
  let command: AnalyzeCommand;
  let mockAnalyzer: ComponentAnalyzer;
  let consoleLogSpy: any;

  beforeEach(() => {
    command = new AnalyzeCommand();
    mockAnalyzer = new ComponentAnalyzer();
    consoleLogSpy = vi.spyOn(console, 'log').mockImplementation(() => {});
    vi.clearAllMocks();
  });

  afterEach(() => {
    consoleLogSpy.mockRestore();
  });

  describe('Command Execution', () => {
    it('should analyze components in specified directory', async () => {
      const testPath = '/test/components';
      const mockResults = {
        componentsFound: 10,
        categorized: {
          atoms: 5,
          molecules: 3,
          organisms: 2
        },
        uncategorized: 0,
        errors: [],
        duration: 1234
      };

      vi.spyOn(mockAnalyzer, 'analyze').mockResolvedValue(mockResults);
      vi.spyOn(fs, 'existsSync').mockReturnValue(true);
      vi.spyOn(fs, 'statSync').mockReturnValue({ isDirectory: () => true } as any);

      await command.execute(testPath);

      expect(mockAnalyzer.analyze).toHaveBeenCalledWith(testPath, expect.any(Object));
      expect(consoleLogSpy).toHaveBeenCalledWith(expect.stringContaining('✓ Analyzed 10 components'));
    });

    it('should use default path if not specified', async () => {
      const defaultPath = 'src/components';
      const mockResults = {
        componentsFound: 5,
        categorized: {
          atoms: 3,
          molecules: 2,
          organisms: 0
        },
        uncategorized: 0,
        errors: [],
        duration: 500
      };

      vi.spyOn(mockAnalyzer, 'analyze').mockResolvedValue(mockResults);
      vi.spyOn(fs, 'existsSync').mockReturnValue(true);
      vi.spyOn(fs, 'statSync').mockReturnValue({ isDirectory: () => true } as any);

      await command.execute();

      expect(mockAnalyzer.analyze).toHaveBeenCalledWith(defaultPath, expect.any(Object));
    });

    it('should handle non-existent directory', async () => {
      const invalidPath = '/non/existent/path';

      vi.spyOn(fs, 'existsSync').mockReturnValue(false);

      await expect(command.execute(invalidPath)).rejects.toThrow(
        `Path does not exist: ${invalidPath}`
      );
    });

    it('should handle file path instead of directory', async () => {
      const filePath = '/test/component.tsx';

      vi.spyOn(fs, 'existsSync').mockReturnValue(true);
      vi.spyOn(fs, 'statSync').mockReturnValue({ isDirectory: () => false } as any);

      await expect(command.execute(filePath)).rejects.toThrow(
        `Path is not a directory: ${filePath}`
      );
    });
  });

  describe('Options Handling', () => {
    it('should handle recursive option', async () => {
      const options = { recursive: false };
      const mockResults = {
        componentsFound: 3,
        categorized: { atoms: 3, molecules: 0, organisms: 0 },
        uncategorized: 0,
        errors: [],
        duration: 200
      };

      vi.spyOn(mockAnalyzer, 'analyze').mockResolvedValue(mockResults);
      vi.spyOn(fs, 'existsSync').mockReturnValue(true);
      vi.spyOn(fs, 'statSync').mockReturnValue({ isDirectory: () => true } as any);

      await command.execute('src/components', options);

      expect(mockAnalyzer.analyze).toHaveBeenCalledWith(
        'src/components',
        expect.objectContaining({ recursive: false })
      );
    });

    it('should handle output format option', async () => {
      const options = { output: 'json' };
      const mockResults = {
        componentsFound: 5,
        categorized: { atoms: 2, molecules: 2, organisms: 1 },
        uncategorized: 0,
        errors: [],
        duration: 300
      };

      vi.spyOn(mockAnalyzer, 'analyze').mockResolvedValue(mockResults);
      vi.spyOn(fs, 'existsSync').mockReturnValue(true);
      vi.spyOn(fs, 'statSync').mockReturnValue({ isDirectory: () => true } as any);

      await command.execute('src/components', options);

      expect(consoleLogSpy).toHaveBeenCalledWith(
        expect.stringMatching(/^{.*\"componentsFound\":5.*}$/)
      );
    });

    it('should handle save option', async () => {
      const options = { save: true };
      const mockResults = {
        componentsFound: 8,
        categorized: { atoms: 4, molecules: 3, organisms: 1 },
        uncategorized: 0,
        errors: [],
        duration: 400,
        registry: {
          id: 'test-registry',
          components: []
        }
      };

      vi.spyOn(mockAnalyzer, 'analyze').mockResolvedValue(mockResults);
      vi.spyOn(mockAnalyzer, 'saveRegistry').mockResolvedValue(undefined);
      vi.spyOn(fs, 'existsSync').mockReturnValue(true);
      vi.spyOn(fs, 'statSync').mockReturnValue({ isDirectory: () => true } as any);

      await command.execute('src/components', options);

      expect(mockAnalyzer.saveRegistry).toHaveBeenCalled();
      expect(consoleLogSpy).toHaveBeenCalledWith(
        expect.stringContaining('✓ Registry saved')
      );
    });

    it('should handle verbose option', async () => {
      const options = { verbose: true };
      const mockResults = {
        componentsFound: 2,
        categorized: { atoms: 1, molecules: 1, organisms: 0 },
        uncategorized: 0,
        errors: [],
        duration: 150,
        details: [
          { name: 'Button', type: 'atom', path: 'src/Button.tsx' },
          { name: 'SearchBar', type: 'molecule', path: 'src/SearchBar.tsx' }
        ]
      };

      vi.spyOn(mockAnalyzer, 'analyze').mockResolvedValue(mockResults);
      vi.spyOn(fs, 'existsSync').mockReturnValue(true);
      vi.spyOn(fs, 'statSync').mockReturnValue({ isDirectory: () => true } as any);

      await command.execute('src/components', options);

      expect(consoleLogSpy).toHaveBeenCalledWith(expect.stringContaining('Button'));
      expect(consoleLogSpy).toHaveBeenCalledWith(expect.stringContaining('SearchBar'));
    });

    it('should handle exclude patterns', async () => {
      const options = { exclude: '*.test.tsx,*.stories.tsx' };
      const mockResults = {
        componentsFound: 6,
        categorized: { atoms: 3, molecules: 2, organisms: 1 },
        uncategorized: 0,
        errors: [],
        duration: 350
      };

      vi.spyOn(mockAnalyzer, 'analyze').mockResolvedValue(mockResults);
      vi.spyOn(fs, 'existsSync').mockReturnValue(true);
      vi.spyOn(fs, 'statSync').mockReturnValue({ isDirectory: () => true } as any);

      await command.execute('src/components', options);

      expect(mockAnalyzer.analyze).toHaveBeenCalledWith(
        'src/components',
        expect.objectContaining({
          exclude: ['*.test.tsx', '*.stories.tsx']
        })
      );
    });

    it('should handle max depth option', async () => {
      const options = { maxDepth: 3 };
      const mockResults = {
        componentsFound: 4,
        categorized: { atoms: 2, molecules: 1, organisms: 1 },
        uncategorized: 0,
        errors: [],
        duration: 250
      };

      vi.spyOn(mockAnalyzer, 'analyze').mockResolvedValue(mockResults);
      vi.spyOn(fs, 'existsSync').mockReturnValue(true);
      vi.spyOn(fs, 'statSync').mockReturnValue({ isDirectory: () => true } as any);

      await command.execute('src/components', options);

      expect(mockAnalyzer.analyze).toHaveBeenCalledWith(
        'src/components',
        expect.objectContaining({ maxDepth: 3 })
      );
    });
  });

  describe('Error Handling', () => {
    it('should handle analysis errors', async () => {
      const mockResults = {
        componentsFound: 5,
        categorized: { atoms: 2, molecules: 1, organisms: 1 },
        uncategorized: 1,
        errors: [
          'Failed to parse src/BadComponent.tsx',
          'Invalid syntax in src/BrokenComponent.tsx'
        ],
        duration: 450
      };

      vi.spyOn(mockAnalyzer, 'analyze').mockResolvedValue(mockResults);
      vi.spyOn(fs, 'existsSync').mockReturnValue(true);
      vi.spyOn(fs, 'statSync').mockReturnValue({ isDirectory: () => true } as any);

      await command.execute('src/components');

      expect(consoleLogSpy).toHaveBeenCalledWith(
        expect.stringContaining('⚠ Errors encountered')
      );
      expect(consoleLogSpy).toHaveBeenCalledWith(
        expect.stringContaining('Failed to parse src/BadComponent.tsx')
      );
    });

    it('should handle analyzer exceptions', async () => {
      vi.spyOn(mockAnalyzer, 'analyze').mockRejectedValue(
        new Error('Analyzer initialization failed')
      );
      vi.spyOn(fs, 'existsSync').mockReturnValue(true);
      vi.spyOn(fs, 'statSync').mockReturnValue({ isDirectory: () => true } as any);

      await expect(command.execute('src/components')).rejects.toThrow(
        'Analyzer initialization failed'
      );
    });

    it('should handle permission errors', async () => {
      vi.spyOn(fs, 'existsSync').mockReturnValue(true);
      vi.spyOn(fs, 'statSync').mockImplementation(() => {
        throw new Error('Permission denied');
      });

      await expect(command.execute('/restricted/path')).rejects.toThrow(
        'Permission denied'
      );
    });
  });

  describe('Output Formatting', () => {
    it('should format table output correctly', async () => {
      const mockResults = {
        componentsFound: 10,
        categorized: { atoms: 5, molecules: 3, organisms: 2 },
        uncategorized: 0,
        errors: [],
        duration: 600,
        components: [
          { id: 'btn-primary', name: 'Button', type: 'atom', usage: 23 },
          { id: 'search-bar', name: 'SearchBar', type: 'molecule', usage: 5 }
        ]
      };

      vi.spyOn(mockAnalyzer, 'analyze').mockResolvedValue(mockResults);
      vi.spyOn(fs, 'existsSync').mockReturnValue(true);
      vi.spyOn(fs, 'statSync').mockReturnValue({ isDirectory: () => true } as any);

      await command.execute('src/components', { output: 'table' });

      expect(consoleLogSpy).toHaveBeenCalledWith(
        expect.stringMatching(/Button.*atom.*23/)
      );
      expect(consoleLogSpy).toHaveBeenCalledWith(
        expect.stringMatching(/SearchBar.*molecule.*5/)
      );
    });

    it('should format JSON output correctly', async () => {
      const mockResults = {
        componentsFound: 3,
        categorized: { atoms: 2, molecules: 1, organisms: 0 },
        uncategorized: 0,
        errors: [],
        duration: 200
      };

      vi.spyOn(mockAnalyzer, 'analyze').mockResolvedValue(mockResults);
      vi.spyOn(fs, 'existsSync').mockReturnValue(true);
      vi.spyOn(fs, 'statSync').mockReturnValue({ isDirectory: () => true } as any);

      await command.execute('src/components', { output: 'json', quiet: true });

      const output = JSON.parse(consoleLogSpy.mock.calls[0][0]);
      expect(output.componentsFound).toBe(3);
      expect(output.categorized.atoms).toBe(2);
    });

    it('should format markdown output correctly', async () => {
      const mockResults = {
        componentsFound: 5,
        categorized: { atoms: 2, molecules: 2, organisms: 1 },
        uncategorized: 0,
        errors: [],
        duration: 300,
        components: [
          { id: 'btn', name: 'Button', type: 'atom', path: 'src/Button.tsx' }
        ]
      };

      vi.spyOn(mockAnalyzer, 'analyze').mockResolvedValue(mockResults);
      vi.spyOn(fs, 'existsSync').mockReturnValue(true);
      vi.spyOn(fs, 'statSync').mockReturnValue({ isDirectory: () => true } as any);

      await command.execute('src/components', { output: 'markdown' });

      expect(consoleLogSpy).toHaveBeenCalledWith(
        expect.stringMatching(/## Components Analysis/)
      );
      expect(consoleLogSpy).toHaveBeenCalledWith(
        expect.stringMatching(/### Button/)
      );
    });
  });

  describe('Progress Reporting', () => {
    it('should show progress spinner during analysis', async () => {
      const mockResults = {
        componentsFound: 15,
        categorized: { atoms: 7, molecules: 5, organisms: 3 },
        uncategorized: 0,
        errors: [],
        duration: 2000
      };

      vi.spyOn(mockAnalyzer, 'analyze').mockImplementation(() => {
        return new Promise(resolve => setTimeout(() => resolve(mockResults), 100));
      });
      vi.spyOn(fs, 'existsSync').mockReturnValue(true);
      vi.spyOn(fs, 'statSync').mockReturnValue({ isDirectory: () => true } as any);

      const progressSpy = vi.fn();
      command.on('progress', progressSpy);

      await command.execute('src/components');

      expect(progressSpy).toHaveBeenCalled();
    });

    it('should report file processing progress in verbose mode', async () => {
      const mockResults = {
        componentsFound: 5,
        categorized: { atoms: 3, molecules: 2, organisms: 0 },
        uncategorized: 0,
        errors: [],
        duration: 500
      };

      const progressCallback = vi.fn();
      vi.spyOn(mockAnalyzer, 'on').mockImplementation((event, callback) => {
        if (event === 'file:processed') {
          progressCallback.mockImplementation(callback);
          callback({ file: 'Button.tsx', index: 1, total: 5 });
          callback({ file: 'Input.tsx', index: 2, total: 5 });
        }
      });

      vi.spyOn(mockAnalyzer, 'analyze').mockResolvedValue(mockResults);
      vi.spyOn(fs, 'existsSync').mockReturnValue(true);
      vi.spyOn(fs, 'statSync').mockReturnValue({ isDirectory: () => true } as any);

      await command.execute('src/components', { verbose: true });

      expect(progressCallback).toHaveBeenCalled();
      expect(consoleLogSpy).toHaveBeenCalledWith(
        expect.stringContaining('Processing Button.tsx (1/5)')
      );
    });
  });

  describe('Configuration', () => {
    it('should load configuration from file', async () => {
      const configPath = '.component-analyzer.json';
      const config = {
        paths: {
          components: 'src/components',
          output: 'docs/components'
        },
        analysis: {
          recursive: true,
          exclude: ['*.test.tsx']
        }
      };

      vi.spyOn(fs, 'existsSync').mockImplementation(path => {
        if (path === configPath) return true;
        if (path === 'src/components') return true;
        return false;
      });
      vi.spyOn(fs, 'readFileSync').mockReturnValue(JSON.stringify(config));
      vi.spyOn(fs, 'statSync').mockReturnValue({ isDirectory: () => true } as any);

      const mockResults = {
        componentsFound: 7,
        categorized: { atoms: 4, molecules: 2, organisms: 1 },
        uncategorized: 0,
        errors: [],
        duration: 400
      };

      vi.spyOn(mockAnalyzer, 'analyze').mockResolvedValue(mockResults);

      await command.execute(undefined, { config: configPath });

      expect(fs.readFileSync).toHaveBeenCalledWith(configPath, 'utf8');
      expect(mockAnalyzer.analyze).toHaveBeenCalledWith(
        'src/components',
        expect.objectContaining({
          recursive: true,
          exclude: ['*.test.tsx']
        })
      );
    });

    it('should merge CLI options with config file', async () => {
      const configPath = '.component-analyzer.json';
      const config = {
        analysis: {
          recursive: true,
          maxDepth: 5
        }
      };

      vi.spyOn(fs, 'existsSync').mockReturnValue(true);
      vi.spyOn(fs, 'readFileSync').mockReturnValue(JSON.stringify(config));
      vi.spyOn(fs, 'statSync').mockReturnValue({ isDirectory: () => true } as any);

      const mockResults = {
        componentsFound: 3,
        categorized: { atoms: 2, molecules: 1, organisms: 0 },
        uncategorized: 0,
        errors: [],
        duration: 150
      };

      vi.spyOn(mockAnalyzer, 'analyze').mockResolvedValue(mockResults);

      await command.execute('src/components', {
        config: configPath,
        recursive: false, // Override config
        verbose: true // Additional option
      });

      expect(mockAnalyzer.analyze).toHaveBeenCalledWith(
        'src/components',
        expect.objectContaining({
          recursive: false, // CLI option takes precedence
          maxDepth: 5, // From config
          verbose: true // Additional CLI option
        })
      );
    });
  });
});