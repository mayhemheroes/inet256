package fuzzInet256

import (
    "strconv"
    fuzz "github.com/AdaLogics/go-fuzz-headers"

    "github.com/inet256/inet256/client/go_client/inet256client"
    "github.com/inet256/inet256/pkg/bitstr"
)

func mayhemit(bytes []byte) int {

    var num int
    if len(bytes) > 10 {
        num, _ = strconv.Atoi(string(bytes[0]))
        bytes = bytes[1:]

        switch num {

        case 0:
            fuzzConsumer := fuzz.NewConsumer(bytes)
            var testStruct bitstr.BytesLSB
            var testInt int

            testStruct.Bytes = bytes
            testStruct.Begin = 0
            testStruct.End = len(bytes)

            err := fuzzConsumer.CreateSlice(&testInt)
            if err != nil {
                return 0
            }

            testStruct.Len()

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