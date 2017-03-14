package main

import (
	"flag"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var (
	listen   = flag.String("listen", ":8080", "TCP address to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")
	dir      = flag.String("dir", "/opt", "Directory to serve static files from (path /timefil)")
	file     = flag.String("file", "index.html", "File to serve")
	delay    = flag.Int("delay", 0, "Delay to respond the request (path /delay)")
	rmin     = flag.Int("rmin", 1, "Minimum time (millisecond) to respond the request (path /range)")
	rmax     = flag.Int("rmax", 10, "Maximum time (millisecond) to respond the request  (path /range)")
	response = flag.String("response", "A", "Default response  (path /range and /delay)")
)

var src = rand.NewSource(time.Now().UnixNano())
var r13k = RandStringBytesMaskImprSrc(13 * 1024)
var r1m = RandStringBytesMaskImprSrc(1024 * 1024)
var timeList = [26]int{5, 5, 5, 5, 5, 10, 10, 10, 10, 10, 5, 5, 5, 5, 5, 10, 10, 10, 10, 10, 50, 50, 100, 500, 1000}

func main() {
	flag.Parse()

	rw := fasthttp.PathRewriteFunc(func(ctx *fasthttp.RequestCtx) []byte {
		return []byte("/" + *file)
	})

	fs := &fasthttp.FS{
		Root:               *dir,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: true,
		Compress:           *compress,
		AcceptByteRange:    false,
		PathRewrite:        rw,
	}

	fsHandler := fs.NewRequestHandler()

	h := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/1m", "/1m.ghtml":
			fmt.Fprintf(ctx, r1m)
		case "/13k", "/13k.ghtml":
			fmt.Fprintf(ctx, r13k)
		case "/1b", "/1b.ghtml":
			fmt.Fprintf(ctx, "A")
		case "/delay", "/delay.ghtml":
			doneCh := make(chan struct{})
			go func() {
				time.Sleep(time.Millisecond * time.Duration(*delay))
				close(doneCh)
			}()

			select {
			case <-doneCh:
				fmt.Fprintf(ctx, *response)
			}
		case "/range", "/range.ghtml":
			doneCh := make(chan struct{})
			go func() {
				rand.Seed(time.Now().Unix())
				workDuration := time.Millisecond * time.Duration(rand.Intn(*rmax-*rmin+1)+*rmin)
				time.Sleep(workDuration)
				close(doneCh)
			}()

			select {
			case <-doneCh:
				fmt.Fprintf(ctx, *response)
			}
		case "/timefile", "/timefile.ghtml":
			doneCh := make(chan struct{})
			go func() {
				rand.Seed(time.Now().Unix())
				workDuration := time.Millisecond * time.Duration(timeList[rand.Intn(26)])
				time.Sleep(workDuration)
				close(doneCh)
			}()

			select {
			case <-doneCh:
				fsHandler(ctx)
			}
		case "/checkgaleb":
		  // healthcheck for multicontent
			fmt.Fprintf(ctx, "WORKING")
		case "/":
			fmt.Fprintf(ctx, "/1b /13k /1m /delay /range /timefile\n")
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}

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
