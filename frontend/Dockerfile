FROM node:22.15-alpine3.21 as builder

WORKDIR /app

COPY ./public ./public/
COPY ./src ./src/
COPY ./package.json .
COPY ./package-lock.json .
COPY ./tsconfig.json .

RUN npm install --force && \
    npm run build

FROM nginxinc/nginx-unprivileged:1.25-alpine

COPY --from=builder /app/build /usr/share/nginx/html
COPY ./.docker/vhost.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
