/// <reference types="vitest" />
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react-swc';
import path from 'path';

/**
 * Vitest Configuration for Contract Tests
 *
 * This configuration is specifically for Pact consumer contract tests.
 * It isolates contract testing from regular unit tests and prevents
 * MSW interference.
 */
export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'node', // Use node environment to avoid browser APIs
    setupFiles: ['./src/contracts/test-setup.ts'],
    globals: true,
    include: ['src/contracts/**/*.test.ts'],
    exclude: [
      '**/node_modules/**',
      '**/dist/**',
      'src/!(contracts)/**/*.test.ts', // Exclude non-contract tests
    ],
    // Disable watch mode for contract tests
    watch: false,
    // Increase timeout for contract tests
    testTimeout: 30000,
    hookTimeout: 30000,
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
    preserveSymlinks: false,
  },
  optimizeDeps: {
    include: ['@pact-foundation/pact'],
  },
});