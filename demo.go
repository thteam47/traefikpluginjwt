// Package plugindemo a demo plugin.

package plugindemo



import (

	"strings"

	"bytes"

	"context"

	"fmt"

	"net/http"

	"text/template"

	"encoding/base64"

	"crypto/hmac"

	"crypto/sha256"

)



// Config the plugin configuration.

type Config struct {

	AuthHeader string `json:"authHeader,omitempty"`

}



// CreateConfig creates the default plugin configuration.

func CreateConfig() *Config {

	return &Config{}

}



// Demo a Demo plugin.

type Demo struct {

	next     http.Handler

	authHeader  string

	name     string

	template *template.Template

}



// New created a new Demo plugin.

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	if len(config.AuthHeader) == 0 {

		config.AuthHeader = "Authorization"

	}



	return &Demo{

		authHeader:  config.AuthHeader,

		next:     next,

		name:     name,

		template: template.New("demo").Delims("[[", "]]"),

	}, nil

}



func (a *Demo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	

	headerToken := req.Header.Get(a.authHeader)



	if len(headerToken) == 0 {

		http.Error(rw, "Request error 1", http.StatusUnauthorized)

		return

	}

	

	token, preprocessError  := preprocessJWT(headerToken, "Bearer")

	if preprocessError != nil {

		http.Error(rw, "Request error 2", http.StatusBadRequest)

		return

	}

	

	_, verificationError := verifyJWT(token, "thteam")

	if verificationError != nil {

		http.Error(rw, "Not allowed", http.StatusUnauthorized)

		return

	}





		payload, decodeErr := decodeBase64(token.payload)

		if decodeErr != nil {

			http.Error(rw, "Request error", http.StatusBadRequest)

			return

		}



		req.Header.Add("X-Auth-User", payload)

		req.Header.Del("Authorization")

		a.next.ServeHTTP(rw, req)

	

}

type Token struct {

	header string

	payload string

	verification string

}



// verifyJWT Verifies jwt token with secret

func verifyJWT(token Token, secret string) (bool, error) {

	mac := hmac.New(sha256.New, []byte(secret))

	message := token.header + "." + token.payload

	mac.Write([]byte(message))

	expectedMAC := mac.Sum(nil)

	

	decodedVerification, errDecode := base64.RawURLEncoding.DecodeString(token.verification)

	if errDecode != nil {

		return false, errDecode

	}



	if hmac.Equal(decodedVerification, expectedMAC) {

		return true, nil

	}

	return false, nil

	// TODO Add time check to jwt verification

}



// preprocessJWT Takes the request header string, strips prefix and whitespaces and returns a Token

func preprocessJWT(reqHeader string, prefix string) (Token, error) {

	// fmt.Println("==> [processHeader] SplitAfter")

	// structuredHeader := strings.SplitAfter(reqHeader, "Bearer ")[1]

	cleanedString := strings.TrimPrefix(reqHeader, prefix)

	cleanedString = strings.TrimSpace(cleanedString)

	// fmt.Println("<== [processHeader] SplitAfter", cleanedString)



	var token Token



	tokenSplit := strings.Split(cleanedString, ".")



	if len(tokenSplit) != 3 {

		return token, fmt.Errorf("Invalid token")

	}



	token.header = tokenSplit[0]

	token.payload = tokenSplit[1]

	token.verification = tokenSplit[2]



	return token, nil

}



// decodeBase64 Decode base64 to string

func decodeBase64(baseString string) (string, error) {

	byte, decodeErr := base64.RawURLEncoding.DecodeString(baseString)

	if decodeErr != nil {

		return baseString, fmt.Errorf("Error decoding")

	}

	return string(byte), nil

}