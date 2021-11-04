import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'


// yarn add @types/node
import path from "path";


// https://vitejs.dev/config/
export default defineConfig({
    // base: '/k8sailor/webapp/dist/',
    // resolv 对象中配置
    resolve: {
        // 别名配置
        alias: {
            "@": path.resolve(__dirname, "src"),
        }
    },
    plugins: [vue()]
})
