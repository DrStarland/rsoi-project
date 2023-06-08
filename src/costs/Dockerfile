FROM golang:1.20
WORKDIR /pain/
COPY ./go.mod .
COPY ./vendor .
RUN go mod tidy && go mod vendor
COPY . .
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o app ./cmd/costs

FROM alpine:3.18.0
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /pain/app ./
EXPOSE 8060
CMD ["./app"]