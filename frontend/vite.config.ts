import {defineConfig} from 'vite'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    AutoImport({
      resolvers: [ElementPlusResolver()],
    }),
    Components({
      resolvers: [ElementPlusResolver()],
    }),
  ],
  resolve: {
    alias: {
      // wails3 生成的绑定按 Go module 路径落盘，别名缩短导入路径
      '@backend-bindings': fileURLToPath(
        new URL('./bindings/github.com/yann0917/dedao-gui/backend', import.meta.url),
      ),
    },
  },
  server: {
    // vite 8 默认 host=localhost 只落 IPv6(::1)，
    // 而 wails3 的 ExternalAssetHandler 代理固定拨 IPv4 127.0.0.1，不显式绑定会 dev 白屏
    host: '127.0.0.1',
  },
})
