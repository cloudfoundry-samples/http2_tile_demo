package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

func handleErr(message string, err error) {
	if err != nil {
		fmt.Printf("[Error] %s: %s\n", message, err)
		os.Exit(1)
	}
}

func serveImgFactory() func(http.ResponseWriter, *http.Request) {
	img, err := os.Open("cf_logo.png")
	handleErr("Opening image", err)
	defer img.Close()
	imgContents, err := io.ReadAll(img)
	handleErr("Copy image", err)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("X-Http-Version", r.Proto)
		_, err = io.Copy(w, bytes.NewBuffer(imgContents))
		handleErr("Writing image", err)
	}
}

type TemplateArgs map[string]interface{}

func generateImage() string {
	fingerprint := rand.Int63n(math.MaxInt64)
	return fmt.Sprintf(`<img src="/images/test_%d.png" height="20" />`, fingerprint)
}

func generateImages(quantity int) string {
	var sb strings.Builder
	for i := 0; i < quantity; i++ {
		sb.WriteString(generateImage())
	}
	return sb.String()
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	loadTimeScript := template.JS(`
window.addEventListener("load", function() {
	window.setTimeout(function() {
		document.getElementById("time").innerHTML = (window.performance.timing.loadEventEnd - window.performance.timing.navigationStart);
	}, 0);
});`)
	templateArgs := TemplateArgs{
		"HttpVersion": r.Proto,
		"Images":      template.HTML(generateImages(1000)),
		"Script":      loadTimeScript,
	}
	templateString := `<!DOCTYPE html>
<html>
<title>HTTP/2 Image Tile Demo</title>
<body>
<h1>HTTP/2 Image Tile Demo</h1>
<div>HTTP Version: {{.HttpVersion}}</div>
<div>Load Time: <span id="time"></span></div>
{{.Images}}
<script>
{{.Script}}
</script>
</body>
</html>`
	renderableTemplate := template.Must(template.New("").Parse(templateString))
	err := renderableTemplate.Execute(w, templateArgs)
	handleErr("Writing HTML", err)
}

func serveHttp2(doneChan chan string) {
	http2Server := &http.Server{Addr: ":8080"}
	fmt.Println("Listening on https://localhost:8080")
	err := http2Server.ListenAndServeTLS("server.crt", "server.key")
	handleErr("Serving HTTP2", err)
	doneChan <- "Done HTTP/2"
}

func serveHttp1(doneChan chan string) {
	http2Server := &http.Server{Addr: ":8081"}
	fmt.Println("Listening on http://localhost:8081")
	err := http2Server.ListenAndServe()
	handleErr("Serving HTTP1", err)
	doneChan <- "Done HTTP/1"
}

func main() {
	http.HandleFunc("/", serveHTML)
	serveImg := serveImgFactory()
	http.HandleFunc("/images/", serveImg)

	doneChan := make(chan string)

	go serveHttp1(doneChan)
	go serveHttp2(doneChan)

	fmt.Println(<-doneChan)
}
