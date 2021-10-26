# frontend-build-stage
FROM node:lts AS frontend-build-stage

WORKDIR /app

COPY . .

RUN npm install -g pnpm
RUN pnpm install

ENV API_SERVER_URL=http://localhost:8888
EXPOSE 5000

ENTRYPOINT pnpm run dev