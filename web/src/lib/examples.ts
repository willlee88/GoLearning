import fs from 'node:fs';
import path from 'node:path';

const repoRoot = path.resolve(process.cwd(), '..');

export type ExampleFile = {
  /** Path relative to the example root, e.g. main.go */
  relPath: string;
  content: string;
  language: string;
};

const SKIP_DIRS = new Set(['node_modules', 'vendor', '.git', 'dist', '.astro']);
const MAX_FILES = 10;
const MAX_BYTES_TOTAL = 120_000;
const MAX_FILE_BYTES = 40_000;

function langFromExt(file: string): string {
  const ext = path.extname(file).toLowerCase();
  switch (ext) {
    case '.go':
      return 'go';
    case '.mod':
      return 'go';
    case '.html':
    case '.htm':
      return 'xml';
    case '.json':
      return 'json';
    case '.md':
      return 'markdown';
    case '.yml':
    case '.yaml':
      return 'yaml';
    case '.ts':
    case '.tsx':
      return 'typescript';
    case '.js':
      return 'javascript';
    case '.sh':
    case '.bash':
      return 'bash';
    case '.ps1':
      return 'powershell';
    default:
      return 'plaintext';
  }
}

function interestingFile(name: string): boolean {
  const lower = name.toLowerCase();
  if (lower.endsWith('.exe') || lower.endsWith('.test') || lower.endsWith('.out')) return false;
  return (
    lower.endsWith('.go') ||
    lower === 'go.mod' ||
    lower === 'readme.md' ||
    lower.endsWith('.html') ||
    lower.endsWith('.json') ||
    lower.endsWith('.yml') ||
    lower.endsWith('.yaml')
  );
}

function scoreFile(rel: string): number {
  const base = path.basename(rel).toLowerCase();
  if (base === 'main.go') return 0;
  if (base.endsWith('_test.go')) return 1;
  if (base.endsWith('.go')) return 2;
  if (base === 'go.mod') return 3;
  if (base === 'readme.md') return 4;
  if (base.endsWith('.html')) return 5;
  return 6;
}

function walkFiles(dir: string, base: string, acc: string[] = []): string[] {
  if (!fs.existsSync(dir)) return acc;
  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    if (entry.name.startsWith('.')) continue;
    const full = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      if (SKIP_DIRS.has(entry.name)) continue;
      walkFiles(full, base, acc);
    } else if (interestingFile(entry.name)) {
      acc.push(path.relative(base, full).replaceAll('\\', '/'));
    }
  }
  return acc;
}

/** Load source files for a lesson `example` field (repo-relative path). */
export function loadExampleFiles(examplePath: string): ExampleFile[] {
  const trimmed = String(examplePath || '').trim();
  if (!trimmed) return [];

  const abs = path.resolve(repoRoot, trimmed);
  const rootResolved = path.resolve(repoRoot);
  if (!abs.startsWith(rootResolved + path.sep) && abs !== rootResolved) {
    return [];
  }
  if (!fs.existsSync(abs)) return [];

  if (fs.statSync(abs).isFile()) {
    const content = fs.readFileSync(abs, 'utf8');
    return [
      {
        relPath: path.basename(abs),
        content: content.length > MAX_FILE_BYTES ? content.slice(0, MAX_FILE_BYTES) + '\n…\n' : content,
        language: langFromExt(abs),
      },
    ];
  }

  const rels = walkFiles(abs, abs).sort((a, b) => scoreFile(a) - scoreFile(b) || a.localeCompare(b));
  const out: ExampleFile[] = [];
  let bytes = 0;
  for (const rel of rels) {
    if (out.length >= MAX_FILES) break;
    const full = path.join(abs, rel);
    let content = fs.readFileSync(full, 'utf8');
    if (content.length > MAX_FILE_BYTES) {
      content = content.slice(0, MAX_FILE_BYTES) + '\n…\n';
    }
    if (bytes + content.length > MAX_BYTES_TOTAL) break;
    bytes += content.length;
    out.push({ relPath: rel, content, language: langFromExt(rel) });
  }
  return out;
}
