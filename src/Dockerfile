#STAGE 1
FROM node:20-alpine3.18 AS build
WORKDIR /usr/src/app
COPY ./costs-n-tasks/package.json ./costs-n-tasks/package-lock.json ./
RUN npm install --force
COPY ./costs-n-tasks/ .
RUN npm run build

#STAGE 2
FROM nginx:1.25.1-alpine
COPY nginx.conf /etc/nginx/nginx.conf
COPY --from=build /usr/src/app/dist/costs-n-tasks /usr/share/nginx/html