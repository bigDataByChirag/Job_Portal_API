package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// iss (issuer): Issuer of the JWT
// sub (subject): Subject of the JWT (the user)
// aud (audience): Recipient for which the JWT is intended
// exp (expiration time): Time after which the JWT expires
// nbf (not before time): Time before which the JWT must not be accepted for processing
// iat (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
// jti (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed (allows a token to be used only once)

//openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
// openssl rsa -in private.pem -pubout -out pubkey.pem
/*
{
  "sub": "1234567890",
  "name": "John Doe",
  "admin": true,
  "iat": 1516239022
}
*/

func main() {
	privatePem, err := os.ReadFile("private.pem")
	if err != nil {
		log.Fatalln(err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePem)
	if err != nil {
		log.Fatalln(err)
	}

	c := jwt.RegisteredClaims{
		Issuer:    "api project",
		Subject:   "101",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(50 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, c)

	encodedToken,err := token.SignedString(privateKey)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(encodedToken)
}
