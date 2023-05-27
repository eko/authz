FROM node:20.2-alpine as builder

WORKDIR /app

COPY ./public ./public/
COPY ./src ./src/
COPY ./package.json .
COPY ./package-lock.json .
COPY ./tsconfig.json .

RUN npm install --force && \
    npm run build

FROM nginxinc/nginx-unprivileged:1.23-alpine

COPY --from=builder /app/build /usr/share/nginx/html
COPY ./.docker/vhost.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
