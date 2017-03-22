package main

import (
	"flag"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"math/rand"
	"os"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var (
	listen   = flag.String("listen", fmt.Sprintf(":%s", os.Getenv("PORT")), "TCP address to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")
)

var src = rand.NewSource(time.Now().UnixNano())
var r13k = RandStringBytesMaskImprSrc(13 * 1024)
var r1m = RandStringBytesMaskImprSrc(1024 * 1024)

func main() {
	flag.Parse()

	h := requestHandler
	if *compress {
		h = fasthttp.CompressHandler(h)
	}

	if err := fasthttp.ListenAndServe(*listen, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/1m":
		fmt.Fprintf(ctx, r1m)
	case "/13k":
		fmt.Fprintf(ctx, r13k)
	case "/1b":
		fmt.Fprintf(ctx, "A")
	case "/":
		fmt.Fprintf(ctx, "A")
	default:
		ctx.Error("not found", fasthttp.StatusNotFound)
	}
}
