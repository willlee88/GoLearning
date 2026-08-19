/**
 * Compare inline ```go blocks in lessons vs referenced examples/*.go
 * Usage: node scripts/audit-example-sync.mjs
 */
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');

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

function parseFrontmatter(raw) {
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
    data[m[1]] = m[2].trim().replace(/^["']|["']$/g, '');
  }
  return { data, body };
}

function goBlocks(body) {
  const out = [];
  const re = /```go\n([\s\S]*?)```/g;
  let m;
  while ((m = re.exec(body))) out.push(m[1].trim());
  return out;
}

function collectGo(dir) {
  if (!fs.existsSync(dir)) return [];
  const files = [];
  const w = (d) => {
    for (const e of fs.readdirSync(d, { withFileTypes: true })) {
      const f = path.join(d, e.name);
      if (e.isDirectory()) w(f);
      else if (e.name.endsWith('.go') && !e.name.endsWith('_test.go')) files.push(f);
    }
  };
  w(dir);
  return files;
}

const rows = [];
for (const file of walk(path.join(root, 'content'))) {
  const { data, body } = parseFrontmatter(fs.readFileSync(file, 'utf8'));
  if (!data.example || !String(data.example).startsWith('examples/')) continue;
  const exDir = path.join(root, data.example);
  const files = collectGo(exDir);
  const exCode = files.map((f) => fs.readFileSync(f, 'utf8')).join('\n');
  const blocks = goBlocks(body);
  let overlap = 0;
  let total = 0;
  const miss = [];
  for (const b of blocks.slice(0, 3)) {
    const lines = b
      .split(/\n/)
      .map((l) => l.trim())
      .filter(
        (l) =>
          l &&
          !l.startsWith('//') &&
          l.length > 10 &&
          !l.startsWith('package ') &&
          !l.startsWith('import') &&
          l !== ')' &&
          l !== '}' &&
          l !== '{',
      );
    for (const l of lines) {
      total++;
      if (exCode.includes(l)) overlap++;
      else if (miss.length < 4) miss.push(l.slice(0, 90));
    }
  }
  const ratio = total ? overlap / total : 1;
  rows.push({
    id: data.lessonId,
    file: path.relative(root, file).replaceAll('\\', '/'),
    ex: data.example,
    blocks: blocks.length,
    goFiles: files.length,
    total,
    overlap,
    ratio: Number(ratio.toFixed(2)),
    miss,
  });
}

rows.sort((a, b) => a.ratio - b.ratio);
console.log('=== Likely mismatches (overlap < 0.5, enough lines) ===');
for (const r of rows.filter((r) => r.total >= 3 && r.ratio < 0.5)) {
  console.log(`\n${r.id}  ${r.file}`);
  console.log(`  example=${r.ex}  overlap=${r.overlap}/${r.total} (${r.ratio})`);
  for (const m of r.miss) console.log(`  miss: ${m}`);
}
console.log('\n=== Summary ===');
const scored = rows.filter((r) => r.total >= 3);
const bad = scored.filter((r) => r.ratio < 0.5);
console.log(`pairs=${rows.length} scored=${scored.length} mismatched=${bad.length}`);
console.log(
  'avg ratio',
  (scored.reduce((s, r) => s + r.ratio, 0) / Math.max(scored.length, 1)).toFixed(2),
);
