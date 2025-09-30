# Используем официальный Node.js образ
FROM node:22.18-alpine AS base

# установка пакетов
FROM base AS deps
RUN apk add --no-cache libc6-compat
WORKDIR /app

# Копируем package.json и package-lock.json
COPY package.json package-lock.json* ./
RUN npm ci --only=production --omit=dev --prefer-offline --no-audit && npm cache clean --force

# Стадия сборки приложения
FROM base AS builder
WORKDIR /app
RUN apk add --no-cache libc6-compat
COPY --from=deps /app/node_modules ./node_modules
COPY . .

RUN npx orval --config ovral.config.js

# Собираем приложение
RUN npm run build

# Стадия запуска
FROM base AS runner
WORKDIR /app

ENV NODE_ENV=production

# Создаем пользователя без привилегий
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 8080
ENV HOST=0.0.0.0
ENV PORT=8080

CMD ["node", "server.js"]