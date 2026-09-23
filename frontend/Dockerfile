FROM node:25-alpine AS build
WORKDIR /app

RUN corepack enable && corepack prepare pnpm@11.21.0 --activate

COPY frontend/package.json frontend/pnpm-lock.yaml frontend/pnpm-workspace.yaml ./

RUN pnpm config set fetch-retry-maxtimeout 120000 && \
    pnpm config set fetch-retries 5 && \
    pnpm install --frozen-lockfile

COPY frontend/ .
ARG PUBLIC_API_BASE=/api
ENV PUBLIC_API_BASE=$PUBLIC_API_BASE
RUN pnpm run build

FROM caddy:2-alpine
COPY --from=build /app/build /srv
COPY Caddyfile /etc/caddy/Caddyfile