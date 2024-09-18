// Package security provides utilities for setting up network transport credentials.
package security

import (
	"fmt"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// LoadClientCredentials sets up and returns the appropriate transport credentials for the client.
func LoadClientCredentials(certPath string) (credentials.TransportCredentials, error) {
	if certPath != "" {
		creds, err := credentials.NewClientTLSFromFile(certPath, "")
		if err != nil {
			return nil, fmt.Errorf("failed to load client TLS credentials: %w", err)
		}
		return creds, nil
	}

	// If no certificate is provided, use insecure credentials (for development or non-production)
	return insecure.NewCredentials(), nil
}

// LoadServerCredentials sets up and returns the appropriate transport credentials for the server.
func LoadServerCredentials(certPath, keyPath string) (credentials.TransportCredentials, error) {
	if certPath != "" {
		creds, err := credentials.NewServerTLSFromFile(certPath, keyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load server TLS credentials: %w", err)
		}
		return creds, nil
	}

	// If no certificate is provided, use insecure credentials (for development or non-production)
	return insecure.NewCredentials(), nil
}
