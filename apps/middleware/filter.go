package main

import (
	"net/url"
	"strings"

	"github.com/enescakir/emoji"
	"github.com/http-wasm/http-wasm-guest-tinygo/handler"
	"github.com/http-wasm/http-wasm-guest-tinygo/handler/api"
)

func main() {
	handler.HandleRequestFn = handleRequest
}

// handleRequest implements a simple HTTP router.
func handleRequest(req api.Request, resp api.Response) (next bool, reqCtx uint32) {
	if uri := req.GetURI(); !strings.HasPrefix(uri, "/dapr/") {
		u, err := url.ParseRequestURI(req.GetURI())
		if err != nil {
			panic(err)
		}
		q := u.Query()
		// Our platform doesn't support cats yet.
		if strings.Contains(q.Get("message"), "_cat") ||
			strings.Contains(q.Get("message"), ":cat:") ||
			strings.Contains(q.Get("message"), ":cat2:") {
			q.Set("message", strings.ReplaceAll(q.Get("message"), "cat", "***"))
		} else {
			q.Set("message", emoji.Parse(q.Get("message")))
		}
		u.RawQuery = q.Encode()
		req.SetURI(u.String())
	}

	next = true // proceed to the next handler on the host.
	return

}
