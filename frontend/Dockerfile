FROM node:19.4-alpine as builder

WORKDIR /app

COPY ./public ./public/
COPY ./src ./src/
COPY ./package.json .
COPY ./package-lock.json .
COPY ./tsconfig.json .

RUN npm install && \
    npm run build

FROM nginx:1.23-alpine

COPY --from=builder /app/build /usr/share/nginx/html
COPY ./.docker/vhost.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
