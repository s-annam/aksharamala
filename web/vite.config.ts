import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// Load environment variables
const API_BASE_URL = process.env.VITE_API_BASE_URL || 'http://localhost:8081';

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': {
        target: API_BASE_URL,
        changeOrigin: true
      }
    }
  },
  define: {
    'import.meta.env.VITE_API_BASE_URL': JSON.stringify(API_BASE_URL)
  }
})
