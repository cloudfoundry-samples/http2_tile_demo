# HTTP/2 Tile Demo

HTTP/2 speed test demo. Tries to load a bunch of images via HTTP/2 or
HTTP/1.

## How to Run Locally:

1. `./gen_keys.sh` to generate cert/key for the HTTPS + HTTP/2 server. Requires
   `openssl`. Set the subject name to `localhost`.
2. `go build`
3. `PORT=8080 PROTO=h2 ./http2_tile_demo` to serve HTTP/2 over TLS
3. `PORT=8080 PROTO=h2c ./http2_tile_demo` to serve HTTP/2 without TLS
3. `PORT=8080 PROTO=http1 ./http2_tile_demo` to serve HTTP/1.1 without TLS
