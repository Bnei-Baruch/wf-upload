package api

import (
	"context"
	"github.com/Bnei-Baruch/wf-upload/common"
	"net/http"
	"strings"
)

type Roles struct {
	Roles []string `json:"roles"`
}

type IDTokenClaims struct {
	Acr               string           `json:"acr"`
	AllowedOrigins    []string         `json:"allowed-origins"`
	Aud               interface{}      `json:"aud"`
	AuthTime          int              `json:"auth_time"`
	Azp               string           `json:"azp"`
	Email             string           `json:"email"`
	Exp               int              `json:"exp"`
	FamilyName        string           `json:"family_name"`
	GivenName         string           `json:"given_name"`
	Iat               int              `json:"iat"`
	Iss               string           `json:"iss"`
	Jti               string           `json:"jti"`
	Name              string           `json:"name"`
	Nbf               int              `json:"nbf"`
	Nonce             string           `json:"nonce"`
	PreferredUsername string           `json:"preferred_username"`
	RealmAccess       Roles            `json:"realm_access"`
	ResourceAccess    map[string]Roles `json:"resource_access"`
	SessionState      string           `json:"session_state"`
	Sub               string           `json:"sub"`
	Typ               string           `json:"typ"`
}

func (a *App) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if common.SKIP_AUTH {
			next.ServeHTTP(w, r)
			return
		}

		// Right now we don't need to check token on download
		if r.Method == "GET" {
			next.ServeHTTP(w, r)
			return
		}

		auth := parseToken(r)

		if auth == "" {
			respondWithError(w, http.StatusBadRequest, "no `Authorization` header set")
			return
		}

		// Authorization header provided, let's verify.
		token, err := a.tokenVerifier.Verify(context.TODO(), auth)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		// parse claims
		var claims IDTokenClaims
		if err := token.Claims(&claims); err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Check permission
		if !checkPermission(claims.RealmAccess.Roles) {
			respondWithError(w, http.StatusForbidden, "Access denied")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func checkPermission(roles []string) bool {
	if roles != nil {
		for _, r := range roles {
			if r == "wp_plugin_upload" {
				return true
			}
		}
	}
	return false
}

func parseToken(r *http.Request) string {
	var token = ""
	authHeader := strings.Split(strings.TrimSpace(r.Header.Get("Authorization")), " ")
	if len(authHeader) == 2 && strings.ToLower(authHeader[0]) == "bearer" && len(authHeader[1]) > 0 {
		token = authHeader[1]
	}
	return token
}
