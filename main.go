package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func handleErr(message string, err error) {
	if err != nil {
		fmt.Printf("[Error] %s: %s\n", message, err)
		os.Exit(1)
	}
}

var (
	//go:embed assets/cf_logo.png
	imgContents []byte
)

func serveImgFactory() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("X-Http-Version", r.Proto)
		_, err := io.Copy(w, bytes.NewBuffer(imgContents))
		handleErr("Writing image", err)
	}
}

type TemplateArgs struct {
	HttpVersion string
	Images      template.HTML
	Script      template.JS
}

func generateImage() string {
	fingerprint := rand.Int63n(math.MaxInt64)
	return fmt.Sprintf( /* language=html */ `<img src="/images/test_%[1]d.png" alt="Test Image %[1]d" height="20" />`, fingerprint)
}

func generateImages(quantity int) string {
	var sb strings.Builder
	for i := 0; i < quantity; i++ {
		sb.WriteString(generateImage())
	}
	return sb.String()
}

var (
	//go:embed assets/load_time.js
	loadTimeScript string

	//go:embed assets/index.gohtml
	indexHTML string
)

func serveHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	templateArgs := TemplateArgs{
		HttpVersion: r.Proto,
		Images:      template.HTML(generateImages(1000)),
		Script:      template.JS(loadTimeScript),
	}
	err := template.Must(template.New("").Parse(indexHTML)).Execute(w, templateArgs)
	handleErr("Writing HTML", err)
}

func serveHttp2(address string) {
	http2Server := &http.Server{Addr: address}
	fmt.Printf("Listening on https://%s", address)
	err := http2Server.ListenAndServeTLS("server.crt", "server.key")
	handleErr("Serving H2", err)
}

func serveH2c(address string) {
	http2Server := &http.Server{Addr: address, Handler: h2c.NewHandler(http.DefaultServeMux, &http2.Server{})}
	fmt.Printf("Listening on http://%s", address)
	err := http2Server.ListenAndServe()
	handleErr("Serving H2C", err)
}

func serveHttp1(address string) {
	http2Server := &http.Server{Addr: address}
	fmt.Printf("Listening on http://%s", address)
	err := http2Server.ListenAndServe()
	handleErr("Serving HTTP1", err)
}

func main() {
	http.HandleFunc("/", serveHTML)
	serveImg := serveImgFactory()
	http.HandleFunc("/images/", serveImg)

	proto := strings.ToLower(os.Getenv("PROTO"))
	port := os.Getenv("PORT")

	address := fmt.Sprintf("0.0.0.0:%s", port)

	switch proto {
	case "h2":
		serveHttp2(address)
	case "h2c":
		serveH2c(address)
	case "http1":
		serveHttp1(address)
	default:
		fmt.Println("No protocol set. Specify PROTO environment variable. Valid values: 'h2', 'h2c', 'http1'.")
	}
}
