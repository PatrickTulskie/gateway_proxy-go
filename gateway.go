package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

// Get env var or default
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getListenAddress() string {
	port := getEnv("PORT", "1338")
	return ":" + port
}

func logSetup() {
	service_a_url := os.Getenv("SERVICE_A_URL")
	service_b_url := os.Getenv("SERVICE_B_URL")
	default_service_url := os.Getenv("DEFAULT_SERVICE_URL")

	log.Printf("Server will run on: %s\n", getListenAddress())
	log.Printf("Service A url: %s\n", service_a_url)
	log.Printf("Service B url: %s\n", service_b_url)
	log.Printf("Default url: %s\n", default_service_url)
}

func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse(target)

	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	proxy.ServeHTTP(res, req)
}

func getServiceUrl(request *http.Request) string {
	switch request.URL.Path {
	case "/secure/serva", "/serva":
		return os.Getenv("SERVICE_A_URL")
	case "/secure/servb", "/servb":
		return os.Getenv("SERVICE_B_URL")
	}

	return os.Getenv("DEFAULT_SERVICE_URL")
}

func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {

	serviceUrl := getServiceUrl(req)

	// Figure out if our header was set correctly
	secureEndpoint := strings.HasPrefix(req.URL.Path, "/secure")

	if secureEndpoint && req.Header.Get("AuthToken") != os.Getenv("SECURE_SERVICE_TOKEN") {
		res.WriteHeader(400)
		res.Write([]byte("Unauthorized"))
		return
	}

	serveReverseProxy(serviceUrl, res, req)
}

func main() {
	logSetup()

	http.HandleFunc("/", handleRequestAndRedirect)
	if err := http.ListenAndServe(getListenAddress(), nil); err != nil {
		panic(err)
	}
}
