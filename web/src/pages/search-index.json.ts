import type { APIRoute } from 'astro';
import { getAllLessons } from '../lib/lessons';

export const GET: APIRoute = () => {
  const items = getAllLessons().map((l) => ({
    lessonId: l.lessonId,
    title: l.title,
    description: l.description,
    volume: l.volume,
    status: l.status,
    level: l.level,
    tags: l.tags,
    slug: l.slug,
    // lightweight body excerpt for search
    text: l.body.replace(/```[\s\S]*?```/g, ' ').replace(/[#>*_`|\[\]()]/g, ' ').slice(0, 4000),
  }));

  return new Response(JSON.stringify({ generatedAt: new Date().toISOString(), items }), {
    headers: {
      'Content-Type': 'application/json; charset=utf-8',
      'Cache-Control': 'no-cache',
    },
  });
};
