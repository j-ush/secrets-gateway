// Package auth provides IRSA and token caching logic
package auth

import (
	"context"
	"sync"
	"time"
)

// IRSAClient handles AWS IRSA authentication and token caching
type IRSAClient struct {
	tokenCache map[string]*CachedToken
	mu         sync.RWMutex
}

// CachedToken represents a cached authentication token
type CachedToken struct {
	Token     string
	ExpiresAt time.Time
}

// NewIRSAClient creates a new IRSA client
func NewIRSAClient() (*IRSAClient, error) {
	return &IRSAClient{
		tokenCache: make(map[string]*CachedToken),
	}, nil
}

// GetToken retrieves a token for the given role, using cache when available
func (c *IRSAClient) GetToken(ctx context.Context, roleARN string) (string, error) {
	c.mu.RLock()
	if cached, ok := c.tokenCache[roleARN]; ok && time.Now().Before(cached.ExpiresAt) {
		c.mu.RUnlock()
		return cached.Token, nil
	}
	c.mu.RUnlock()

	// Fetch new token (implementation depends on AWS SDK)
	token, expiresAt, err := c.fetchNewToken(ctx, roleARN)
	if err != nil {
		return "", err
	}

	c.mu.Lock()
	c.tokenCache[roleARN] = &CachedToken{
		Token:     token,
		ExpiresAt: expiresAt,
	}
	c.mu.Unlock()

	return token, nil
}

func (c *IRSAClient) fetchNewToken(ctx context.Context, roleARN string) (string, time.Time, error) {
	// TODO: Implement AWS STS AssumeRoleWithWebIdentity
	return "", time.Time{}, nil
}

