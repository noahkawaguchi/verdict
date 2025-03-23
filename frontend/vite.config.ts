import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  test: {
    globals: true, // Use global test APIs
    environment: 'jsdom', // Simulate browser environment
    setupFiles: './src/vitest.setup.ts', // Path to the setup file
  },
});
