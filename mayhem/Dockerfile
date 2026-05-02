FROM golang:1.25 AS builder

COPY . /go/src/github.com/inet256/inet256
WORKDIR /go/src/github.com/inet256/inet256

RUN go install github.com/dvyukov/go-fuzz/go-fuzz@latest github.com/dvyukov/go-fuzz/go-fuzz-build@latest
RUN GOFLAGS=-mod=mod go get github.com/dvyukov/go-fuzz/go-fuzz-dep
RUN GOFLAGS=-mod=mod go get github.com/AdaLogics/go-fuzz-headers
RUN apt update && apt install -y clang

RUN go-fuzz-build -libfuzzer -o mayhem/Inet256_client/fuzzInet256_client.a ./mayhem/Inet256_client && \
    clang -fsanitize=fuzzer mayhem/Inet256_client/fuzzInet256_client.a -o mayhem/Inet256_client/fuzzInet256_client.libfuzzer

RUN go-fuzz-build -libfuzzer -o mayhem/Inet256_inet256/fuzzInet256_inet256.a ./mayhem/Inet256_inet256 && \
    clang -fsanitize=fuzzer mayhem/Inet256_inet256/fuzzInet256_inet256.a -o mayhem/Inet256_inet256/fuzzInet256_inet256.libfuzzer

FROM debian:bookworm-slim
COPY --from=builder /go/src/github.com/inet256/inet256/mayhem/Inet256_client/fuzzInet256_client.libfuzzer /
COPY --from=builder /go/src/github.com/inet256/inet256/mayhem/Inet256_inet256/fuzzInet256_inet256.libfuzzer /
