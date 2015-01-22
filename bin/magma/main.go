package main

import (
	"log"
	"net/http"

	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json"
	_ "github.com/matm/scytale"
)

const (
	port    = "localhost:8080"
	rootUrl = "/api"
)

func main() {
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(NewMagma(), "")
	http.Handle(rootUrl, s)

	log.Printf("Listening on %v ...", port)
	log.Println("Send requests to", rootUrl)
	log.Fatal(http.ListenAndServe(port, nil))
}
