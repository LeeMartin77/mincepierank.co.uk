FROM --platform=$BUILDPLATFORM node:18.13.0 as builder
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM --platform=$BUILDPLATFORM node:18.13.0 as dep-builder
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci --omit=dev

FROM node:18.13.0-bullseye-slim AS runner
WORKDIR /app

ENV NODE_ENV production

COPY package.json package-lock.json ./

COPY --from=dep-builder /app/node_modules ./node_modules
COPY --from=builder /app/build ./build

CMD ["node", "build"]