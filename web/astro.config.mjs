// @ts-check
import { defineConfig } from 'astro/config';
import tailwindcss from '@tailwindcss/vite';

// - default / Docker / portable zip (serve from folder root): base `/`
// - GITHUB_PAGES=true: base `/GoLearning/`
const githubPages = process.env.GITHUB_PAGES === 'true';

// https://astro.build/config
export default defineConfig({
  site: githubPages ? 'https://willlee88.github.io' : undefined,
  base: githubPages ? '/GoLearning/' : '/',
  trailingSlash: 'always',
  vite: {
    plugins: [tailwindcss()],
  },
});
