# syntax=docker/dockerfile:1

ARG NODE_VERSION=22-alpine

FROM node:${NODE_VERSION} AS deps
WORKDIR /app
ENV NEXT_TELEMETRY_DISABLED=1
COPY package.json ./
RUN npm install --package-lock=false --no-audit --no-fund

FROM deps AS build
ARG IM_API_INTERNAL_BASE_URL=http://go-api:9000
ENV NEXT_TELEMETRY_DISABLED=1
ENV IM_API_INTERNAL_BASE_URL=${IM_API_INTERNAL_BASE_URL}
COPY . .
RUN mkdir -p public && npm run build

FROM node:${NODE_VERSION} AS runner
WORKDIR /app
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV HOSTNAME=0.0.0.0
ENV PORT=3000
ARG IM_API_INTERNAL_BASE_URL=http://go-api:9000
ENV IM_API_INTERNAL_BASE_URL=${IM_API_INTERNAL_BASE_URL}

COPY package.json ./
RUN npm install --omit=dev --package-lock=false --no-audit --no-fund && npm cache clean --force
COPY --from=build /app/.next ./.next
COPY --from=build /app/next.config.mjs ./next.config.mjs
COPY --from=build /app/public ./public

EXPOSE 3000
CMD ["npm", "run", "start"]
