// Package pluginJWT a JWT plugin.
package pluginJWT

import (
	"context"
	"net/http"
)

// Config the plugin configuration.
type Config struct {
	Secret string `json:"secret,omitempty"`
	Authorization string `json:"authorization,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// JWT a JWT plugin.
type JWT struct {
	next     http.Handler
	name     string
	secret	 string
	authorization  string
}

// New created a new JWT plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Secret) == 0 {
		config.Secret = "SECRET"
	}
	if len(config.Authorization) == 0 {
		config.Authorization = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InRoYWlwdjlAdmlldHRlbGN5YmVyLmNvbSIsIm5hbWUiOiJQaGFtIFZhbiBUaGFpIiwicm9sZSI6InVzZXIifQ.tbh6XFFJIm7TNym0WNTMoLP3KBQdx1RpLOP-fgdopxQ"
	}

	return &JWT{
		next:     next,
		name:     name,
		secret: config.Secret,
		authorization: config.Authorization,
	}, nil
}

func (j *JWT) ServeHTTP(res http.ResponseWriter, req *http.Request) {
		req.Header.Add("X-Auth-User", j.authorization)
		j.next.ServeHTTP(res, req)
	// headerToken := req.Header.Get(j.authorization)

	// if len(headerToken) == 0 {
	// 	http.Error(res, "Request error", http.StatusUnauthorized)
	// 	return
	// }
	
	// token, preprocessError  := preprocessJWT(headerToken, "Bearer")
	// if preprocessError != nil {
	// 	http.Error(res, "Request error", http.StatusBadRequest)
	// 	return
	// }
	
	// verified, verificationError := verifyJWT(token, j.secret)
	// if verificationError != nil {
	// 	http.Error(res, "Not allowed", http.StatusUnauthorized)
	// 	return
	// }

	// if (verified) {
	// 	// If true decode payload
	// 	payload, decodeErr := decodeBase64(token.payload)
	// 	if decodeErr != nil {
	// 		http.Error(res, "Request error", http.StatusBadRequest)
	// 		return
	// 	}

	// 	// TODO Check for outside of ASCII range characters
		
	// 	// Inject header as proxypayload or configured name
	// 	req.Header.Add("X-Auth-User", payload)
	// 	j.next.ServeHTTP(res, req)
	// } else {
	// 	http.Error(res, "Not allowed", http.StatusUnauthorized)
	// }
}
