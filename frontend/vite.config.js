import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
    server: {
        // eslint-disable-next-line
        host: process.env.PROXY_API ? true : false,
        port: 8080,
        proxy: {
            "/api": {
                // eslint-disable-next-line
                target: process.env.PROXY_API || "http://localhost:5000/",
                changeOrigin: true,
                secure: false,
                ws:true,
                rewrite: path => path.replace(/^\/api/, '')
            },
            "/static": {
                // eslint-disable-next-line
                target: process.env.PROXY_API || "http://localhost:5000/",
                changeOrigin: true,
                secure: false
            }
        },
    },
    plugins: [
        vue(),
    ],
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url))
        }
    }
})
