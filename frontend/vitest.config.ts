import { defaultExclude, defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";
import path from "path";

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  test: {
    environment: "jsdom",
    exclude: [...defaultExclude, "**/*.spec.{js,mjs,cjs,ts,mts,cts,jsx,tsx}"],
    globals: true,
    setupFiles: [path.resolve(__dirname, "./src/utils/vitest/setup.ts")],
    coverage: {
      enabled: true,
      reporter: ["html"],
    },
  },
});
