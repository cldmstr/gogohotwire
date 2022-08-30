package main

import (
	"log"

	"github.com/cldmstr/gogohotwire/internal/app"
)

func main() {
	a := app.New()

	err := a.Run()
	if err != nil {
		log.Fatalln("%+v", err)
	}
}
