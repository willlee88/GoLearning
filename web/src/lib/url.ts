/** Join site base path (e.g. `/` or `/GoLearning/`) with a path. */
export function url(path = ''): string {
  const base = import.meta.env.BASE_URL || '/';
  const clean = String(path).replace(/^\//, '');
  if (!clean) return base;
  return `${base}${clean}`;
}
