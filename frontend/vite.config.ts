import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import wails from "@wailsio/runtime/plugins/vite";

const devPort = Number(process.env.WAILS_VITE_PORT || 9245);

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue(), wails("./bindings")],
  server: {
    port: devPort,
    strictPort: true,
  },
});
