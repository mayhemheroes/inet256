package fuzzInet256_client

import (
    fuzz "github.com/AdaLogics/go-fuzz-headers"

    "github.com/inet256/inet256/client/go_client/inet256client"
    "github.com/inet256/inet256/pkg/bitstr"
    "github.com/inet256/inet256/pkg/inet256"
    
)

func mayhemit(bytes []byte) int {

    if len(bytes) > 10 {
        num := int(bytes[0])
        bytes = bytes[1:]
        fuzzConsumer := fuzz.NewConsumer(bytes)

        switch num {

        case 0:
            var testStruct bitstr.BytesLSB

            testStruct.Bytes = bytes
            testStruct.Begin = 0
            testStruct.End = len(bytes)

            testStruct.Len()

            return 0

        case 1:
            inet256.AddrFromBytes(bytes)
            return 0

        case 2:
            x, _ := fuzzConsumer.GetBytes()
            prefix, _ := fuzzConsumer.GetBytes()
            nbits, _ := fuzzConsumer.GetInt()

            inet256.HasPrefix(x, prefix, nbits)
            return 0
    
        default:
            fuzzString, _ := fuzzConsumer.GetString()

            inet256client.NewClient(fuzzString)
            return 0
        }
    }
    return 0
}

func Fuzz(data []byte) int {
    _ = mayhemit(data)
    return 0
}