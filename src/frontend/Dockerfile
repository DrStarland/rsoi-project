FROM node:20-alpine3.18

WORKDIR /app
RUN npm install -g serve

COPY package*.json ./
RUN npm install --silent


COPY . .


RUN npm run build


EXPOSE 3000

CMD ["serve", "-s", "build"]

# #STAGE 1
# FROM node:20-alpine3.18 AS build
# WORKDIR /usr/src/app
# COPY package*.json ./
# RUN npm install --force
# COPY . .
# RUN npm run build

# #STAGE 2
# FROM nginx:1.25.1-alpine
# COPY nginx.conf /etc/nginx/nginx.conf
# COPY --from=build /usr/src/app/build /usr/share/nginx/html
