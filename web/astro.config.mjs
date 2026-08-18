// @ts-check
import { defineConfig } from 'astro/config';
import tailwindcss from '@tailwindcss/vite';

// GitHub Pages project site: https://willlee88.github.io/GoLearning/
// Local / Docker root deploy: leave GITHUB_PAGES unset → base `/`
const githubPages = process.env.GITHUB_PAGES === 'true';

// https://astro.build/config
export default defineConfig({
  site: 'https://willlee88.github.io',
  base: githubPages ? '/GoLearning/' : '/',
  trailingSlash: 'always',
  vite: {
    plugins: [tailwindcss()],
  },
});
