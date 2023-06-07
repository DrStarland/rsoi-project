FROM golang:1.20
WORKDIR /pain/
COPY ./go.mod .
# COPY ./go.sum .
RUN ls
COPY ./vendor* .
RUN go mod tidy -e && go mod vendor
COPY . .
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o app ./cmd/tasks

FROM alpine:3.18.0
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /pain/app ./
EXPOSE 8060
CMD ["./app"]