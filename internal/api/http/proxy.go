package proxy

import (
	"gateway/internal/config"
	"io"
	"net/http"
	"strings"

	"github.com/Blockary/platform-core/http/server"
)

func RegisterEndpoint() {
	server.Router.HandleFunc("/", forwardRequest)
}

var serviceConfig, _ = config.LoadConfig("config.json")

func forwardRequest(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	servicePath := pathParts[1]
	serviceHost := serviceConfig.Services[servicePath]
	proxyRequest(w, r, serviceHost)
}

func proxyRequest(w http.ResponseWriter, r *http.Request, target string) {
	req, err := http.NewRequest(r.Method, target+r.RequestURI, r.Body)
	if err != nil {
		http.Error(w, "Failed to create request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for key, value := range r.Header {
		req.Header[key] = value
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to reach service", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	for key, value := range resp.Header {
		w.Header()[key] = value
	}

	io.Copy(w, resp.Body)
}
