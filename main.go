package main

import (
	json "encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/dgrijalva/jwt-go"
)

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("NotImplemented"))
})

var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("The API is up"))
})

var JwtAuthoritay = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	var signingKey = []byte("ThisIsMySecretKey")

	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)

	// Create a map to store our claims
	claims := token.Claims.(jwt.MapClaims)

	// Set token claims
	claims["admin"] = true
	claims["name"] = "New User"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Sign the token with our secret
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		// Return to stop processing request
		return
	}

	// Set json content type on response header
	w.Header().Set("Content-Type", "application/json")

	// Write json encoded token to response
	json.NewEncoder(w).Encode(JwtToken{Token: tokenString})

	// Finally, write the token to the browser window
	// 	w.Write([]byte(tokenString))
})

// var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
// 	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
// 		return mySigningKey, nil
// 	},
// 	SigningMethod: jwt.SigningMethodHS256,
// })

func ValidationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")

		if authorizationHeader == "" {
			http.Error(w, "Authorization header must be present", http.StatusUnauthorized)
			return
		}

		bearer := strings.Split(authorizationHeader, " ")

		if len(bearer) == 2 {
			token, err := jwt.Parse(bearer[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					// 	http.Error(w, Error("Authorization header must be present", http.StatusUnauthorized))
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return signingKey, nil
			})
			if err != nil {
				json.NewEncoder(w).Encode(Exception{Message: err.Error()})
				return
			}

			if token.Valid {
				// context.Set(r, "decoded", token.Claims)
				next(w, r)

			} else if ve, ok := err.(*jwt.ValidationError); ok {
				// Check verification errors
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					fmt.Println("That's not even a token")
					json.NewEncoder(w).Encode(Exception{Message: "Token is malformed"})
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					// Token is either expired or not active yet
					fmt.Println("Timing is everything")
					json.NewEncoder(w).Encode(Exception{Message: "Token is expired or not valid yet"})
				} else {
					fmt.Println("Couldn't handle this token:", err)
					json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
				}
			} else {
				fmt.Println("Couldn't handle this token:", err)
				json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
			}
		}
	})
}

func main() {
	// Here we are instantiating the gorilla/mux router
	r := mux.NewRouter()

	// On the default page we will simply serve our static index page.
	r.Handle("/", http.FileServer(http.Dir("./views/")))
	// We will setup our server so we can serve static assest like images, css from the /static/{file} route
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.Handle("/status", StatusHandler).Methods("GET")
	r.Handle("/token", JwtAuthoritay).Methods("GET")
	r.Handle("/admin", ValidationMiddleware(StatusHandler)).Methods("GET")

	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r))
}
