FROM node:20-slim

RUN corepack enable && corepack prepare pnpm@latest --activate

WORKDIR /app

COPY package.json pnpm-lock.yaml frontend/package.json pnpm-workspace.yaml ./

RUN pnpm install

COPY . .

RUN pnpm install

CMD ["pnpm", "run", "-r", "dev"]