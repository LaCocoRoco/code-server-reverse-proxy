package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Configurable values from environment variables
	domainSuffix := os.Getenv("DOMAIN_SUFFIX")
	targetBaseURL := os.Getenv("TARGET_BASE_URL")
	serverPort := os.Getenv("SERVER_PORT")

	if domainSuffix == "" || targetBaseURL == "" || serverPort == "" {
		log.Fatal("Environment variables DOMAIN_SUFFIX, TARGET_BASE_URL, or SERVER_PORT are not set")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if the request is to /proxy/{port}
		re := regexp.MustCompile(`^/proxy/(\d+)(/.*)?`)
		matches := re.FindStringSubmatch(r.URL.Path)
		if len(matches) >= 2 {
			port := matches[1]
			path := "/"
			if len(matches) == 3 {
				path = matches[2]
			}

			// Construct the new URL with the port as a subdomain
			newURL := "https://" + port + "." + domainSuffix + path + "?" + r.URL.RawQuery

			// Redirect to the new URL
			http.Redirect(w, r, newURL, http.StatusMovedPermanently)
			return
		}

		// Check if the request is to a subdomain
		hostParts := regexp.MustCompile(`^(\d+)\.` + regexp.QuoteMeta(domainSuffix)).FindStringSubmatch(r.Host)
		if len(hostParts) == 2 {
			port := hostParts[1]

			// Construct the target URL
			targetURL := targetBaseURL + ":" + port

			// Set up reverse proxy
			target, err := url.Parse(targetURL)
			if err != nil {
				http.Error(w, "Invalid target URL", http.StatusInternalServerError)
				return
			}
			proxy := httputil.NewSingleHostReverseProxy(target)
			proxy.ServeHTTP(w, r)
			return
		}

		http.Error(w, "Invalid request", http.StatusBadRequest)
	})

	// Start server on configured port
	log.Printf("Starting server on port %s", serverPort)
	http.ListenAndServe(":"+serverPort, nil)
}
