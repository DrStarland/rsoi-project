FROM golang:1.19
WORKDIR /pain/
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o app ./cmd/tickets

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /pain/app ./
EXPOSE 8070
CMD ["./app"]