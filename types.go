package main

var signingKey = []byte("ThisIsMySecretKey")

type JwtToken struct {
	Token string `json:"token"`
}

type Exception struct {
	Message string
}
