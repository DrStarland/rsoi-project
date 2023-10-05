FROM golang:1.20
WORKDIR /pain/
COPY ./go.mod .
COPY ./vendor .
RUN go mod tidy && go mod vendor
COPY . .
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o app ./cmd/users

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /pain/app ./
EXPOSE 8040
CMD ["./app"]