package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"
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

func serveImg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("X-Http-Version", r.Proto)
	_, err := io.Copy(w, bytes.NewBuffer(imgContents))
	handleErr("Writing image", err)
}

type TemplateArgs struct {
	HttpVersion       string
	ImageFingerprints []int
	Script            template.JS
}

func generateImageFingerprints(n int) []int {
	result := make([]int, n)
	for i := range result {
		result[i] = rand.Intn(math.MaxInt64)
	}
	return result
}

var (
	//go:embed assets/load_time.js
	loadTimeScript string

	//go:embed assets/index.gohtml
	indexHTML string
)

func serveHTML(w http.ResponseWriter, r *http.Request) {
	n, err := intQueryParam(r, "n", 1000)
	if err != nil {
		http.Error(w, "failed to parse n (int)", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	templateArgs := TemplateArgs{
		HttpVersion:       r.Proto,
		ImageFingerprints: generateImageFingerprints(n),
		Script:            template.JS(loadTimeScript),
	}
	var buf bytes.Buffer
	err = template.Must(template.New("").Parse(indexHTML)).Execute(&buf, templateArgs)
	if err != nil {
		http.Error(w, "failed to generate page", http.StatusInternalServerError)
		return
	}
	if pusher, ok := w.(http.Pusher); ok {
		ps := getURLs(&buf, n)
		log.Printf("pushing %d frames", len(ps))
		for _, p := range ps {
			if err := pusher.Push(p, nil); err != nil {
				log.Printf("Failed to push: %v", err)
				break
			}
		}
	} else {
		log.Printf("w does not implement http.Pusher")
	}

	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, &buf)
}

func intQueryParam(r *http.Request, name string, defaultValue int) (int, error) {
	result := defaultValue
	if nStr := r.URL.Query().Get(name); nStr != "" {
		qn, err := strconv.Atoi(nStr)
		if err != nil {
			return defaultValue, err
		}
		result = qn
	}
	return result, nil
}

func getURLs(r io.Reader, max int) []string {
	paths := make([]string, 0, max)
	root, err := html.Parse(r)
	if err != nil {
		panic("failed to parse html")
	}
	var visit func(node *html.Node)
	visit = func(node *html.Node) {
		if node == nil || len(paths) >= max {
			return
		}
		switch node.Type {
		case html.ElementNode:
			for _, attr := range node.Attr {
				if attr.Key == "src" && strings.HasPrefix(attr.Val, "/") {
					paths = append(paths, attr.Val)
				}
			}
			for n := node.FirstChild; n != nil; n = n.NextSibling {
				visit(n)
			}
			visit(node.NextSibling)
		case html.DoctypeNode:
			visit(node.NextSibling)
		case html.DocumentNode:
			visit(node.FirstChild)
		}
	}
	visit(root)
	return paths
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
