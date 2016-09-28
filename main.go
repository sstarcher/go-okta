package main

import (
	"log"
	"os"
)

var organization string

func main() {
	client := NewClient(os.Getenv("OKTA_ORGANIZATION"))
	auth, err := client.authenticate(os.Getenv("OKTA_USERNAME"), os.Getenv("OKTA_PASSWORD"))
	if err != nil {
		log.Fatal(err)
	}

	session, err := client.session(auth.SessionToken)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(session)
}
