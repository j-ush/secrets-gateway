// Main entry point for the secrets-gateway Go implementation
// Designed to run as a non-root, read-only container for maximum hardening
// Implements the Five-Layer Defense-in-Depth framework
package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"secrets-gateway/internal/audit"
	"secrets-gateway/internal/auth"
	"secrets-gateway/internal/proxy"
)

// Default configuration constants
const (
	defaultBackendURL     = "http://secrets-backend.internal.svc.cluster.local"
	defaultPort           = "8080"
	defaultRefreshInterval = 5 * time.Minute
)

func main() {
	logger := audit.NewLogger()

	// Get backend secrets manager URL (provider-agnostic)
	backendAddr := os.Getenv("SECRETS_BACKEND_URL")
	if backendAddr == "" {
		backendAddr = defaultBackendURL
	}

	target, err := url.Parse(backendAddr)
	if err != nil {
		log.Fatalf("Failed to parse backend URL: %v", err)
	}

	// L1: Initialize IRSA authentication (Gateway authenticates using cloud-native identity)
	authClient, err := auth.NewIRSAClient()
	if err != nil {
		log.Fatalf("Failed to initialize IRSA client: %v", err)
	}

	// Perform initial authentication
	logger.Info("Performing initial authentication...")
	if _, err := authClient.GetToken(nil, ""); err != nil {
		logger.Warn("Initial token fetch failed (will retry)", "error", err)
	}

	// L2: Background token refresh - mitigates staleness for non-renewable credentials
	go func() {
		ticker := time.NewTicker(defaultRefreshInterval)
		defer ticker.Stop()
		for range ticker.C {
			logger.Info("Refreshing authentication token...")
			if _, err := authClient.GetToken(nil, ""); err != nil {
				logger.Error("Token refresh failed", "error", err)
			}
		}
	}()

	// Create proxy handler with target backend
	proxyHandler := proxy.NewHandlerWithTarget(authClient, logger, target)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	logger.Info("Secrets Gateway starting", "port", port, "backend", backendAddr)
	if err := http.ListenAndServe(":"+port, proxyHandler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

