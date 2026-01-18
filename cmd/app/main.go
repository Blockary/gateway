package main

import (
	proxy "gateway/internal/api/http"

	"github.com/Blockary/platform-core/http/server"
)

func main() {
	proxy.RegisterEndpoint()
	server.StartServer()
}
