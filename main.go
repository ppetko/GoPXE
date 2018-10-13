package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ppetko/gopxe/conf"
	"github.com/ppetko/gopxe/routes"
)

// This is the main package
// Output is webserver om port
func main() {
	conf.Setup()
	port := flag.Lookup("port").Value.(flag.Getter).Get().(string)

	routes := routes.New()
	log.Printf("Serving on port: %s", port)
	if err := http.ListenAndServe(":"+port, routes); err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}

}
