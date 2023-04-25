FROM golang:1.19 as builder

COPY . /inet256
WORKDIR /inet256/mayhem

RUN go install github.com/dvyukov/go-fuzz/go-fuzz@latest github.com/dvyukov/go-fuzz/go-fuzz-build@latest
RUN go get github.com/dvyukov/go-fuzz/go-fuzz-dep
RUN go get github.com/AdaLogics/go-fuzz-headers
RUN apt update && apt install -y clang

RUN cd Inet256_client && go-fuzz-build -libfuzzer -o fuzzInet256_client.a && \
    clang -fsanitize=fuzzer fuzzInet256_client.a -o fuzzInet256_client.libfuzzer

RUN cd Inet256_inet256 && go-fuzz-build -libfuzzer -o fuzzInet256_inet256.a && \
    clang -fsanitize=fuzzer fuzzInet256_inet256.a -o fuzzInet256_inet256.libfuzzer

FROM debian:bookworm-slim
COPY --from=builder /inet256/mayhem/Inet256_client/fuzzInet256_client.libfuzzer /
COPY --from=builder /inet256/mayhem/Inet256_inet256/fuzzInet256_inet256.libfuzzer /