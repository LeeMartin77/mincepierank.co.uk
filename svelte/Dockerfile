FROM --platform=$BUILDPLATFORM node:18.13.0 as builder
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM node:18.13.0-bullseye-slim AS runner
WORKDIR /app

ENV NODE_ENV production

COPY package.json package-lock.json ./

COPY --from=builder /app/build ./build

CMD ["node", "build"]