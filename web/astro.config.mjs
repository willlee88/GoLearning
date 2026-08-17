// @ts-check
import { defineConfig } from 'astro/config';

import tailwindcss from '@tailwindcss/vite';

// https://astro.build/config
export default defineConfig({
  // Static deploy under nginx: directories are .../page/index.html
  trailingSlash: 'always',
  vite: {
    plugins: [tailwindcss()]
  }
});