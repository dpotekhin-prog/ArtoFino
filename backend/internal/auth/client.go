package auth

import (
	"ArtoFino/backend/internal/config"
	"context"
	"fmt"

	"github.com/coreos/go-oidc"
)

type KeycloakConfig struct {
	AuthServerURL string
	Realm         string
	ClientID      string
	ClientSecret  string
}

type KeycloakClient struct {
	Verifier *oidc.IDTokenVerifier
}

func NewKeycloakClient(cfg config.KeycloakConfig) (*KeycloakClient, error) {
	ctx := context.Background()

	issuer := fmt.Sprintf("%s/realms/%s", cfg.AuthServerURL, cfg.Realm)

	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("keycloak provider init error: %w", err)
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: cfg.ClientID,
	})

	return &KeycloakClient{
		Verifier: verifier,
	}, nil
}
