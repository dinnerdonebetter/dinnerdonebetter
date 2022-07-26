package main

import (
	"log"
	"os"

	passwordvalidator "github.com/wagslane/go-password-validator"
)

func main() {
	log.Println(passwordvalidator.Validate(os.Args[1], 60))
}
