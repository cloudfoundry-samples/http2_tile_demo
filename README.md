# HTTP/2 Tile Demo

HTTP/2 speed test demo. Tries to load a bunch of images via HTTP/2 or
HTTP/1.

## Running Locally:

1. `./gen_keys.sh` to generate cert/key for the HTTPS + HTTP/2 server. Requires
   `openssl`. Set the subject name to `localhost`.
1. `go build`
1. `PORT=8080 PROTO=h2 ./http2_tile_demo` to serve HTTP/2 over TLS
1. `PORT=8080 PROTO=h2c ./http2_tile_demo` to serve HTTP/2 without TLS
1. `PORT=8080 PROTO=http1 ./http2_tile_demo` to serve HTTP/1.1 without TLS

## Running On Cloud Foundry:
1. `cf push` with the provided `manifest.yml`. This will deploy a HTTP/1 and
   HTTP/2 version of the demo.
1. Enable end-to-end HTTP/2 routing for the HTTP/2 demo. This is currently being
   implemented for Cloud Foundry, so check CF docs for detailed steps.
