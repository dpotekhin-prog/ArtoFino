package auth

// HasRole checks both realm roles and client roles.
func HasRole(claims map[string]interface{}, role string, clientID string) bool {
	// --- Realm roles ---
	if realmAccess, ok := claims["realm_access"].(map[string]interface{}); ok {
		if roles, ok := realmAccess["roles"].([]interface{}); ok {
			for _, r := range roles {
				if rs, ok := r.(string); ok && rs == role {
					return true
				}
			}
		}
	}

	// --- Client roles ---
	if clientID != "" {
		if resourceAccess, ok := claims["resource_access"].(map[string]interface{}); ok {
			if client, ok := resourceAccess[clientID].(map[string]interface{}); ok {
				if roles, ok := client["roles"].([]interface{}); ok {
					for _, r := range roles {
						if rs, ok := r.(string); ok && rs == role {
							return true
						}
					}
				}
			}
		}
	}

	return false
}
