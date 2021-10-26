import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import * as fs from "fs";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    https: {
      cert: fs.readFileSync('certificates/cert.pem'),
      key: fs.readFileSync('certificates/key.pem'),
    },
  },
})
