import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path' // <-- 1. TAMBAHKAN IMPORT INI

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  
  // --- 2. TAMBAHKAN BAGIAN RESOLVE.ALIAS INI ---
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
    },
  },
  // --- BATAS PENAMBAHAN RESOLVE ---

  // --- 3. BIARKAN KONFIGURASI PROXY YANG SUDAH KITA BUAT SEBELUMNYA ---
  server: {
    proxy: {
      '/api/v1': {
        target: 'http://localhost:8080', 
        changeOrigin: true,
      }
    }
  }
})