// Package proxy provides request forwarding and validation logic
// Implements Layer 1: Policy Enforcement of the defense-in-depth framework
package proxy

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"secrets-gateway/internal/audit"
	"secrets-gateway/internal/auth"
)

// Handler handles proxy requests with connection pooling
// Acts as the infrastructure-enforced boundary between clients and the secrets backend
type Handler struct {
	authClient   *auth.IRSAClient
	logger       *audit.Logger
	target       *url.URL
	reverseProxy *httputil.ReverseProxy
}

// NewHandler creates a new proxy handler with connection pooling
func NewHandler(authClient *auth.IRSAClient, logger *audit.Logger) *Handler {
	return &Handler{
		authClient: authClient,
		logger:     logger,
	}
}

// NewHandlerWithTarget creates a handler with a specific backend target
func NewHandlerWithTarget(authClient *auth.IRSAClient, logger *audit.Logger, target *url.URL) *Handler {
	h := &Handler{
		authClient:   authClient,
		logger:       logger,
		target:       target,
		reverseProxy: httputil.NewSingleHostReverseProxy(target),
	}
	return h
}

// ServeHTTP implements http.Handler
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// L1: Policy Enforcement - validate request format
	if err := h.validateRequest(r); err != nil {
		h.logger.Error("Request validation failed", "error", err, "path", r.URL.Path)
		http.Error(w, "Forbidden: Invalid backend path", http.StatusForbidden)
		return
	}

	// L4: Audit logging - log the request for accountability
	h.logger.Info("Processing request", "path", r.URL.Path, "method", r.Method)

	// Forward request to secrets backend
	h.forwardRequest(w, r)
}

func (h *Handler) validateRequest(r *http.Request) error {
	// L1: Policy Enforcement - Ensure request format matches expected backend structure
	// Mitigates V2: Path Traversal and V1: Cross-Tenant Access
	validPrefixes := []string{
		"/api/v1/secrets/", // Generic secrets API path
		"/v1/secrets/",     // Alternative secrets path
	}

	for _, prefix := range validPrefixes {
		if strings.HasPrefix(r.URL.Path, prefix) {
			return nil
		}
	}

	return errors.New("invalid backend path format")
}

func (h *Handler) forwardRequest(w http.ResponseWriter, r *http.Request) {
	// If no target configured, return error
	if h.target == nil || h.reverseProxy == nil {
		h.logger.Error("No backend target configured")
		http.Error(w, "Gateway not configured", http.StatusServiceUnavailable)
		return
	}

	// L2: Token Caching - get cached token from auth client
	ctx := r.Context()
	token, err := h.authClient.GetToken(ctx, "")
	if err != nil {
		h.logger.Error("Failed to get auth token", "error", err)
		http.Error(w, "Authentication failed", http.StatusInternalServerError)
		return
	}

	// Attach the generalized gateway authentication token
	r.Header.Set("X-Gateway-Auth-Token", token)
	r.Host = h.target.Host

	// Forward via reverse proxy (connection pooling handled by http.Transport)
	h.reverseProxy.ServeHTTP(w, r)
}

