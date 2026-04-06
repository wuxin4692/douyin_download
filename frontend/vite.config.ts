import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    host: '0.0.0.0',
    port: 5173,
  },
  build: {
    // Wails 需要的输出配置
    outDir: 'dist',
    emptyOutDir: true,
  },
})
