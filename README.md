# HTTP/2 Tile Demo

HTTP/2 speed test demo. Tries to load a bunch of images via HTTP/2 or
HTTP/1.

## How to Run Locally:

1. `./gen_keys.sh` to generate cert/key for the HTTPS + HTTP/2 server. Requires
   `openssl`. Set the subject name to `localhost`.
2. `go build`
3. `./http2_tile_demo` to run
4. Visit `https://localhost:8080` for HTTP/2
5. Visit `http://localhost:8081` for HTTP/1
