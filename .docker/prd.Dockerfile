FROM golang:1.16-alpine as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    TZ=Asia/Seoul

WORKDIR /build

COPY ./go.mod go.sum main.go ./

COPY ./src ./src

RUN go mod download

RUN go build -o main .

WORKDIR /dist

RUN cp /build/main .

FROM scratch

COPY --from=builder /dist/main .

COPY .env .

EXPOSE 9000

ENTRYPOINT ["/main"]
