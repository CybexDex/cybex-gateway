package app

import (
	"context"
	"net/http"
	"strings"

	u "coding.net/bobxuyang/cy-gateway-BN/utils"
	jwt "github.com/dgrijalva/jwt-go"
)

//Token ...
type Token struct {
	UserID uint
	jwt.StandardClaims
}

//JwtAuthentication ...
var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/api/account/new", "/api/account/login"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path                                     //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			u.Respond(w, u.Message(false, "Missing auth token"), http.StatusForbidden)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			u.Respond(w, u.Message(false, "Invalid/Malformed auth token"), http.StatusForbidden)
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in
		tk := &Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			// TODO: need to remove hard code
			// TODO: need to remove hard code
			return []byte("token_password"), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			u.Respond(w, u.Message(false, "Malformed authentication token"), http.StatusForbidden)
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			u.Respond(w, u.Message(false, "Token is not valid"), http.StatusForbidden)
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		// TODO: validation check
		//u.Infof("User [%s] login", tk.UserID)
		ctx := context.WithValue(r.Context(), "UserID", tk.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}
