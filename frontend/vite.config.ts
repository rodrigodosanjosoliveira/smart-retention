import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  base: "./",
  server: {
    allowedHosts: ['smart-retention-alb-1299321591.us-east-1.elb.amazonaws.com','localhost:8080']
  }
})
