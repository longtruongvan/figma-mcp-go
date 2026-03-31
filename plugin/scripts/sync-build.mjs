import { cp, mkdir, stat, watch } from "node:fs/promises";
import path from "node:path";
import process from "node:process";

const rootDir = path.resolve(import.meta.dirname, "..");
const distDir = path.join(rootDir, "dist");
const watchMode = process.argv.includes("--watch");

const syncTargets = [
  ["index.html", "index.html"],
  ["code.js", "code.js"],
];

const syncFiles = async () => {
  await mkdir(distDir, { recursive: true });

  for (const [fromName, toName] of syncTargets) {
    const from = path.join(distDir, fromName);
    const to = path.join(rootDir, toName);

    try {
      await stat(from);
      await cp(from, to, { force: true });
    } catch (error) {
      if (error && typeof error === "object" && "code" in error && error.code === "ENOENT") {
        continue;
      }
      throw error;
    }
  }
};

await syncFiles();

if (!watchMode) {
  process.exit(0);
}

console.log(`Watching ${distDir} for build output changes...`);

let pending = null;
for await (const _event of watch(distDir)) {
  if (pending) clearTimeout(pending);
  pending = setTimeout(async () => {
    try {
      await syncFiles();
    } catch (error) {
      console.error("Failed to sync plugin build output:", error);
    }
  }, 50);
}
