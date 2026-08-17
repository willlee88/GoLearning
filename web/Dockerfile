# syntax=docker/dockerfile:1
# Build context: monorepo root (needs content/ + web/)
FROM node:22-bookworm-slim AS build

WORKDIR /app

# Lessons read ../content from web/ at build time
COPY content /app/content
COPY web/package.json web/package-lock.json /app/web/
WORKDIR /app/web
RUN npm ci

COPY web/ /app/web/
# Ensure content is still sibling
RUN npm run build

FROM nginx:1.27-alpine AS runtime

COPY web/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/web/dist /usr/share/nginx/html

EXPOSE 80
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://127.0.0.1/ >/dev/null || exit 1
