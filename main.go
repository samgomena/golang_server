package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	
	"github.com/dgrijalva/jwt-go"
)

var NotImplemented = http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("NotImplemented"))
})

var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("API is up"))
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
    tokenString, _ := token.SignedString(signingKey)

    // Finally, write the token to the browser window
    w.Write([]byte(tokenString))
})

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
  ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
    return mySigningKey, nil
  },
  SigningMethod: jwt.SigningMethodHS256,
})

func main() {
	// Here we are instantiating the gorilla/mux router
	r := mux.NewRouter()

	// On the default page we will simply serve our static index page.
	r.Handle("/", http.FileServer(http.Dir("./views/")))
	// We will setup our server so we can serve static assest like images, css from the /static/{file} route
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	
    r.Handle("/status", StatusHandler).Methods("GET")
    r.Handle("/token", JwtAuthoritay).Methods("GET")
    
    r.Handle("/admin", jwtMiddleware.Handler(StatusHandler)).Methods("GET")
    // r.Handle("/products", NotImplemented).Methods("GET")
    // r.Handle("/products/{slug}/feedback", NotImplemented).Methods("POST")

	// Our application will run on port 3000. Here we declare the port and pass in our router.
	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r))
}
