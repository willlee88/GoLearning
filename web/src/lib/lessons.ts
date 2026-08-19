import fs from 'node:fs';
import path from 'node:path';
import matter from 'gray-matter';

const contentRoot = path.resolve(process.cwd(), '../content');

export type Lesson = {
  lessonId: string;
  title: string;
  description: string;
  volume: string;
  order: number;
  level: string;
  status: string;
  path_required: boolean;
  tags: string[];
  example: string;
  prev: string | null;
  next: string | null;
  body: string;
  slug: string;
  filePath: string;
};

const volumeMeta: Record<string, { title: string; blurb: string }> = {
  p0: { title: 'P0 · 從 Python 走進 Go', blurb: '先搞懂「為什麼不一樣」——用你已有的 Python 當橋。' },
  a: { title: 'A · 語言與資料', blurb: '型別、切片、map、介面……寫 Go 天天會碰到的東西。' },
  b: { title: 'B · 錯誤與 API', blurb: '沒有例外怎麼辦？怎樣把錯誤講清楚、好維護。' },
  c: { title: 'C · 併發', blurb: '一次做很多事：goroutine、channel，還有 race 怎麼抓。' },
  d: { title: 'D · 工程化', blurb: 'module、怎麼寫測試、怎麼測效能。' },
  e: { title: 'E · 標準庫', blurb: '時間、讀寫、設定、JSON……官方工具箱怎麼用。' },
  f: { title: 'F · 網路', blurb: 'TCP、HTTP、WebSocket，最後練一個房間 Lab。' },
  g: { title: 'G · 協定', blurb: '客戶端跟 Server「怎麼說話」：信封、命令、版本。' },
  h: { title: 'H · 遊戲 Server', blurb: '房間、節奏（tick）、權威狀態——局怎麼跑起來。' },
  i: { title: 'I · 資料', blurb: '什麼進記憶體、什麼進 DB／Redis，別每幀都寫庫。' },
  j: { title: 'J · 上線前必看', blurb: '日誌、指標、優雅關閉、壓測與基本安全。' },
  k: { title: 'K · 收束與自評', blurb: '檢查自己學會了什麼，以及下一步可以擴什麼。' },
};

export function getVolumeMeta(volume: string) {
  return volumeMeta[volume] ?? { title: volume.toUpperCase(), blurb: '' };
}

function walkMarkdown(dir: string, acc: string[] = []): string[] {
  if (!fs.existsSync(dir)) return acc;
  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    const full = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      if (entry.name === '_templates' || entry.name === 'paths') continue;
      walkMarkdown(full, acc);
    } else if (entry.name.endsWith('.md')) {
      acc.push(full);
    }
  }
  return acc;
}

function toSlug(filePath: string): string {
  const rel = path.relative(contentRoot, filePath).replace(/\\/g, '/');
  return rel.replace(/\.md$/, '');
}

export function getAllLessons(): Lesson[] {
  const files = walkMarkdown(contentRoot);
  const lessons: Lesson[] = [];

  for (const filePath of files) {
    let raw = fs.readFileSync(filePath, 'utf8');
    if (raw.charCodeAt(0) === 0xfeff) raw = raw.slice(1);
    const { data, content } = matter(raw);
    const lessonId = String(data.lessonId ?? data.id ?? toSlug(filePath));
    lessons.push({
      lessonId,
      title: String(data.title ?? lessonId),
      description: String(data.description ?? ''),
      volume: String(data.volume ?? 'misc'),
      order: Number(data.order ?? 0),
      level: String(data.level ?? 'l2'),
      status: String(data.status ?? 'draft'),
      path_required: Boolean(data.path_required ?? false),
      tags: Array.isArray(data.tags) ? data.tags.map(String) : [],
      example: String(data.example ?? ''),
      prev: data.prev == null || data.prev === '' ? null : String(data.prev),
      next: data.next == null || data.next === '' ? null : String(data.next),
      body: content,
      slug: toSlug(filePath),
      filePath,
    });
  }

  return lessons.sort((a, b) => {
    if (a.volume !== b.volume) return a.volume.localeCompare(b.volume);
    return a.order - b.order;
  });
}

export function getLessonByLessonId(lessonId: string): Lesson | undefined {
  return getAllLessons().find((l) => l.lessonId === lessonId);
}

export function getLessonBySlug(slug: string): Lesson | undefined {
  return getAllLessons().find((l) => l.slug === slug);
}

export function getVolumes(): { id: string; title: string; blurb: string; lessons: Lesson[] }[] {
  const all = getAllLessons();
  const map = new Map<string, Lesson[]>();
  for (const l of all) {
    const list = map.get(l.volume) ?? [];
    list.push(l);
    map.set(l.volume, list);
  }
  const order = ['p0', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k'];
  const ids = [...new Set([...order, ...map.keys()])].filter((id) => map.has(id));
  return ids.map((id) => ({
    id,
    ...getVolumeMeta(id),
    lessons: map.get(id) ?? [],
  }));
}

export function getMainPathLessons(): Lesson[] {
  return getAllLessons().filter((l) => l.path_required);
}
