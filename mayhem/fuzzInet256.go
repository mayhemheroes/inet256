package fuzzInet256

import (
    "strconv"
    fuzz "github.com/AdaLogics/go-fuzz-headers"

    "github.com/inet256/inet256/client/go_client/inet256client"
    "github.com/inet256/inet256/pkg/bitstr"
    "github.com/inet256/inet256/pkg/inet256"
    
)

func mayhemit(bytes []byte) int {

    var num int
    if len(bytes) > 10 {
        num, _ = strconv.Atoi(string(bytes[0]))
        bytes = bytes[1:]

        switch num {

        case 0:
            var testStruct bitstr.BytesLSB

            testStruct.Bytes = bytes
            testStruct.Begin = 0
            testStruct.End = len(bytes)

            testStruct.Len()

            return 0

        case 1:
            fuzzConsumer := fuzz.NewConsumer(bytes)
            var testAddr []byte
            err := fuzzConsumer.CreateSlice(&testAddr)
            if err != nil {
                return 0
            }

            inet256.AddrFromBytes(testAddr)

            return 0

        case 2:
            fuzzConsumer := fuzz.NewConsumer(bytes)
            var fuzzBytes []byte
            err := fuzzConsumer.CreateSlice(&fuzzBytes)
            if err != nil {
                return 0
            }

            inet256.ParseAddrBase64(fuzzBytes)

            return 0

        case 3:
            fuzzConsumer := fuzz.NewConsumer(bytes)
            var x []byte
            var prefix []byte
            var nbits int

            err := fuzzConsumer.CreateSlice(&x)
            if err != nil {
                return 0
            }

            err = fuzzConsumer.CreateSlice(&prefix)
            if err != nil {
                return 0
            }

            err = fuzzConsumer.CreateSlice(&nbits)
            if err != nil {
                return 0
            }

            inet256.HasPrefix(x, prefix, nbits)
            return 0
    
        default:
            fuzzConsumer := fuzz.NewConsumer(bytes)
            var fuzzString string
            err := fuzzConsumer.CreateSlice(&fuzzString)
            if err != nil {
                return 0
            }

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