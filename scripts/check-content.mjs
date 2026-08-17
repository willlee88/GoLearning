/**
 * 檢查 content frontmatter 與 example 路徑是否存在。
 * Usage: node scripts/check-content.mjs
 */
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
const contentRoot = path.join(root, 'content');

function walk(dir, acc = []) {
  for (const e of fs.readdirSync(dir, { withFileTypes: true })) {
    const full = path.join(dir, e.name);
    if (e.isDirectory()) {
      if (e.name === '_templates' || e.name === 'paths') continue;
      walk(full, acc);
    } else if (e.name.endsWith('.md')) acc.push(full);
  }
  return acc;
}

/** Minimal YAML-ish frontmatter parse for our flat fields. */
function parseFrontmatter(raw) {
  // Strip UTF-8 BOM (PowerShell Set-Content often adds it)
  if (raw.charCodeAt(0) === 0xfeff) raw = raw.slice(1);
  if (!raw.startsWith('---')) return { data: {}, body: raw };
  const end = raw.indexOf('\n---', 3);
  if (end === -1) return { data: {}, body: raw };
  const block = raw.slice(4, end).trim();
  const body = raw.slice(end + 4);
  const data = {};
  for (const line of block.split('\n')) {
    const m = line.match(/^([A-Za-z0-9_]+):\s*(.*)$/);
    if (!m) continue;
    const key = m[1];
    let val = m[2].trim();
    if (val === 'null') data[key] = null;
    else if (val === 'true') data[key] = true;
    else if (val === 'false') data[key] = false;
    else if (/^-?\d+$/.test(val)) data[key] = Number(val);
    else if (val.startsWith('[') && val.endsWith(']')) {
      data[key] = val
        .slice(1, -1)
        .split(',')
        .map((s) => s.trim().replace(/^["']|["']$/g, ''))
        .filter(Boolean);
    } else {
      data[key] = val.replace(/^["']|["']$/g, '');
    }
  }
  return { data, body };
}

let errors = 0;
const ids = new Map();

for (const file of walk(contentRoot)) {
  const raw = fs.readFileSync(file, 'utf8');
  const { data } = parseFrontmatter(raw);
  const lessonId = data.lessonId ?? data.id;
  if (!lessonId) {
    console.error(`MISSING lessonId: ${file}`);
    errors++;
    continue;
  }
  if (ids.has(lessonId)) {
    console.error(`DUPLICATE lessonId ${lessonId}: ${file} and ${ids.get(lessonId)}`);
    errors++;
  } else {
    ids.set(lessonId, file);
  }
  for (const key of ['title', 'volume', 'status']) {
    if (data[key] == null || data[key] === '') {
      console.error(`MISSING ${key}: ${file}`);
      errors++;
    }
  }
  if (data.example) {
    const ex = path.join(root, String(data.example));
    if (!fs.existsSync(ex)) {
      console.error(`MISSING example path ${data.example}: ${file}`);
      errors++;
    }
  }
}

for (const file of walk(contentRoot)) {
  const { data } = parseFrontmatter(fs.readFileSync(file, 'utf8'));
  for (const key of ['prev', 'next']) {
    const v = data[key];
    if (v && !ids.has(String(v))) {
      console.error(`BROKEN ${key}=${v}: ${file}`);
      errors++;
    }
  }
}

if (errors) {
  console.error(`\nFailed with ${errors} error(s).`);
  process.exit(1);
}
console.log(`OK · ${ids.size} lessons`);
