FROM golang:1.21-alpine AS builder

RUN apk add --no-cache upx binutils

WORKDIR /build

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags '-w -s -extldflags "-static"' \
    -trimpath \
    -o tiny-test . && \
    strip --strip-all tiny-test && \
    upx --best --lzma tiny-test

FROM scratch

COPY --from=builder /build/tiny-test /tiny-test

EXPOSE 8080

ENTRYPOINT ["/tiny-test"]

