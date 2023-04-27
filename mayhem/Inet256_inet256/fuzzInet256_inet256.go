package fuzzInet256_inet256

import (
    fuzz "github.com/AdaLogics/go-fuzz-headers"

    "github.com/inet256/inet256/pkg/inet256"
    
)

func mayhemit(bytes []byte) int {

    fuzzConsumer := fuzz.NewConsumer(bytes)
    fuzzBytes, _ := fuzzConsumer.GetBytes()

    inet256.ParseAddrBase64(fuzzBytes)
    return 0
}

func Fuzz(data []byte) int {
    _ = mayhemit(data)
    return 0
}