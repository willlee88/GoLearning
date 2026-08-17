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
  p0: { title: 'P0 · Python → Go 心智遷移', blurb: '入口必讀：用 Python 經驗當橋，建立 Go 模型。' },
  a: { title: 'A · 語言與資料模型', blurb: '型別、slice、map、interface…' },
  b: { title: 'B · 錯誤與 API', blurb: 'error 哲學與可維護 API。' },
  c: { title: 'C · 併發', blurb: 'goroutine、channel、race。' },
  d: { title: 'D · 工程化', blurb: 'modules、表驅動、fuzz、bench。' },
  e: { title: 'E · 標準庫', blurb: 'time、io、flag、embed 地圖。' },
  f: { title: 'F · 網路', blurb: 'TCP / HTTP / WebSocket 與房間 Lab。' },
  g: { title: 'G · 協定', blurb: 'JSON 信封、命令/事件、版本化。' },
  h: { title: 'H · 遊戲 Server', blurb: 'Room FSM、Tick、權威 input、Arena Lab。' },
  i: { title: 'I · 資料', blurb: '持久化邊界、Repo、Redis、worker。' },
  j: { title: 'J · 生產化', blurb: 'slog、metrics、shutdown、壓測。' },
  k: { title: 'K · Capstone', blurb: '自評與延伸任務。' },
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
