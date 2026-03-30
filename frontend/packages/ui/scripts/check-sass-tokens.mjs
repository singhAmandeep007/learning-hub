import { readdir, readFile } from "node:fs/promises";
import path from "node:path";
import process from "node:process";

const ROOT = process.cwd();
const COMPONENT_STYLES_DIR = path.join(ROOT, "src", "styles", "components");
const BANNED_PATTERN = /var\(--lhui-/;

async function listScssFiles(dir) {
  const entries = await readdir(dir, { withFileTypes: true });
  const files = [];

  for (const entry of entries) {
    const fullPath = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      files.push(...(await listScssFiles(fullPath)));
      continue;
    }

    if (entry.isFile() && entry.name.endsWith(".scss")) {
      files.push(fullPath);
    }
  }

  return files;
}

function findViolations(content) {
  const lines = content.split(/\r?\n/);
  const violations = [];

  for (let i = 0; i < lines.length; i += 1) {
    if (BANNED_PATTERN.test(lines[i])) {
      violations.push({ line: i + 1, text: lines[i].trim() });
    }
  }

  return violations;
}

async function main() {
  const files = await listScssFiles(COMPONENT_STYLES_DIR);
  const report = [];

  for (const file of files) {
    const content = await readFile(file, "utf8");
    const violations = findViolations(content);
    if (violations.length > 0) {
      report.push({ file, violations });
    }
  }

  if (report.length === 0) {
    console.log("Sass token check passed: no var(--lhui-...) usage in component styles.");
    return;
  }

  console.error("Sass token check failed: use Sass tokens from scss/tokens.scss in component styles.");
  for (const item of report) {
    const relativeFile = path.relative(ROOT, item.file);
    for (const violation of item.violations) {
      console.error(`- ${relativeFile}:${violation.line} -> ${violation.text}`);
    }
  }

  process.exit(1);
}

main().catch((error) => {
  console.error("Sass token check crashed:", error);
  process.exit(1);
});
