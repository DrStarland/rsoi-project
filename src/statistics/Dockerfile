# FROM golang:1.17

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .

# RUN mkdir -p logs temp

# RUN go build -o ./app

# ENTRYPOINT [ "./app" ]

FROM golang:1.20
WORKDIR /pain/
COPY ./go.mod .
COPY ./go.sum .
COPY ./vendor* .
RUN go mod tidy && go mod vendor
COPY . .

RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o app ./
# --no-cache 
FROM alpine:3.18.0 
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=0 /pain/app ./
EXPOSE 8080
RUN mkdir -p logs temp
CMD ["./app"]
