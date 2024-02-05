package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var tkn = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhcGkgcHJvamVjdCIsInN1YiI6IjEwMSIsImV4cCI6MTY5OTQ2OTc0OSwiaWF0IjoxNjk5NDY2NzQ5fQ.o7maJ6wHU6YZS16TcRubmKzjxU7Izee0ge4TSOFehU7PKc1p90tQUVQQ4JWCJ_gBzTY2TrBVMqlh-hzeT5PZTPUJ5FSlTUABiIoysKD_GRu3P2pFLuGFfYzxkkJIXPp9AnhfEoojo0lo7Sn5khELU3E9oaHz4JTP2Co_QvBXSJwZGMDaIqQEVM_QeIBqpXxhaeccaED3kCLcIAqyp3dllQ_pJIoXap_xiqrgFWMx_o1poybfwnQ1MVa7Z0oPcHMzwVPy6TgluadNxBJ7DaFtC-dteiGdsLdqJQJY-76sJZxgNJtyJLWPr9b_JlmmdO8JOdbMiska85Vega2EAomquA`

func main() {

	publicPEM, err := os.ReadFile("pubkey.pem")
	if err != nil {
		// If there's an error reading the file, print an error message and stop execution
		log.Fatalln("not able to read pem file")
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	if err != nil {
		// If there's an error parsing the public key, log the error and stop execution
		log.Fatalln(err)
	}

	var c jwt.RegisteredClaims

	k := func(*jwt.Token) (interface{}, error) {
		return publicKey, nil
	}

	// Parsing the JWT token with the claims
	token, err := jwt.ParseWithClaims(tkn, &c, k)
	if err != nil {
		// If error while parsing the token, print the error and exit
		log.Fatal(err)
	}
	if !token.Valid {
		// If the token is not valid, log the error and exit
		log.Fatal(err)
	}
	fmt.Println("token is valid, we wil allow the request to go through")
	fmt.Printf("%+v\n", c)

}
