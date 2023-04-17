FROM golang:1.19 as builder

COPY . /inet256
WORKDIR /inet256/mayhem

RUN go install github.com/dvyukov/go-fuzz/go-fuzz@latest github.com/dvyukov/go-fuzz/go-fuzz-build@latest
RUN go get github.com/dvyukov/go-fuzz/go-fuzz-dep
RUN go get github.com/AdaLogics/go-fuzz-headers
RUN apt update && apt install -y clang

RUN go-fuzz-build -libfuzzer -o fuzzInet256.a && \
    clang -fsanitize=fuzzer fuzzInet256.a -o fuzzInet256.libfuzzer

FROM debian:bookworm-slim
COPY --from=builder /inet256/mayhem/fuzzInet256.libfuzzer /