# Builder
FROM golang:1.14.2-alpine3.11 as builder

RUN apk update && apk upgrade && \
    apk --update add git make build-base

WORKDIR /app

COPY . .

RUN go generate ./...
RUN go build -o goBinary .

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && apk --no-cache add ca-certificates && \
    apk --update --no-cache add tzdata

ENV TZ=Asia/Jakarta

WORKDIR /app 

EXPOSE 9090

COPY --from=builder /app/goBinary /app

CMD /app/goBinary