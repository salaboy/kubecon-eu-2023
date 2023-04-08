package main

import (
	"net/url"
	"strings"

	"github.com/http-wasm/http-wasm-guest-tinygo/handler"
	"github.com/http-wasm/http-wasm-guest-tinygo/handler/api"
)

func main() {
	handler.HandleRequestFn = handleRequest
}

// handleRequest implements a simple HTTP router.
func handleRequest(req api.Request, resp api.Response) (next bool, reqCtx uint32) {
	if uri := req.GetURI(); strings.HasPrefix(uri, "/v1.0/invoke/write-app") {
		u, err := url.ParseRequestURI(req.GetURI())
		if err != nil {
			panic(err)
		}
		q := u.Query()
		q.Set("message", strings.ReplaceAll(q.Get("message"), "a", "*"))
		u.RawQuery = q.Encode()
		req.SetURI(u.String())
		next = true // proceed to the next handler on the host.
		return
	}

	next = true // proceed to the next handler on the host.
	return

}
