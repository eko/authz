FROM node:26-alpine3.23 as builder

WORKDIR /app

COPY ./package.json ./package-lock.json ./
RUN npm ci

COPY ./index.html ./tsconfig.json ./vite.config.ts ./
COPY ./public ./public/
COPY ./src ./src/

RUN npm run build

FROM nginxinc/nginx-unprivileged:1.31-alpine

COPY --from=builder /app/build /usr/share/nginx/html
COPY ./.docker/vhost.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
