import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import jotaiDebugLabel from "jotai/babel/plugin-debug-label";
import jotaiReactRefresh from "jotai/babel/plugin-react-refresh";
import { resolve } from "path";
import EnvironmentPlugin from "vite-plugin-environment";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    react({ babel: { plugins: [jotaiDebugLabel, jotaiReactRefresh] } }),
    EnvironmentPlugin(["SERVER_URL"]),
  ],
  resolve: {
    alias: {
      src: resolve(__dirname, "src")
      // lib: resolve(__dirname, "src/lib"),
      // routes: resolve(__dirname, "src/routes"),
    },
  },
});
