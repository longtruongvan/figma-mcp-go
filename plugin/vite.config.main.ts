import { defineConfig } from "vite";

export default defineConfig({
  build: {
    target: "es2015",
    lib: {
      entry: "src/main.ts",
      formats: ["iife"],
      name: "code",
      fileName: () => "code.js",
    },
    outDir: ".",
    emptyOutDir: false,
    minify: false,
  },
});
