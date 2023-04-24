# FROM golang:1.19 AS builder

# ENV GO111MODULE=on

# WORKDIR /app

# COPY . .

# RUN go mod tidy
# RUN go build -o app ./cmd/gateway

# EXPOSE 8080

# ENTRYPOINT ["./app"]

# # Этап, на котором выполняется сборка приложения
# FROM golang:1.16-alpine as builder
# WORKDIR /build
# COPY go.mod . # go.sum
# RUN go mod download
# COPY . .
# RUN go build -o /main main.go
# # Финальный этап, копируем собранное приложение
# FROM alpine:3
# COPY --from=builder main /bin/main
# ENTRYPOINT ["/bin/main"]

# # 1 шаг
# FROM golang:1.19 AS build_stage
# WORKDIR /app
# COPY . .
# RUN go mod tidy
# RUN go build -o /main ./cmd/gateway

# # 2 шаг
# FROM alpine AS run_stage
# WORKDIR /app
# COPY --from=build_stage /main /main
# RUN ls
# # RUN chmod +x /bin/main

# EXPOSE 8080

# ENTRYPOINT ["/main"]

FROM golang:1.19
WORKDIR /pain/
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o app ./cmd/gateway

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /pain/app ./
EXPOSE 8080
CMD ["./app"]