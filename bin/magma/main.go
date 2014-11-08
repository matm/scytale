package main

import (
	"log"
	"net/http"

	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json"
	_ "secret"
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

	/*
		password, err := secret.ReadPassword(pwdMinLen, false)
		if err != nil {
			log.Fatal(err)
		}
		ar := secret.NewZipArchive(password)
		if *output == "" {
			*output = "."
		}
		if err := ar.Extract(flag.Arg(0), *output); err != nil {
			log.Fatal(err)
		}
		return
		password, err := secret.ReadPassword(pwdMinLen, true)
		if err != nil {
			log.Fatal(err)
		}
		ar := secret.NewZipArchive(password)
		if err := ar.Create(*output, flag.Args()); err != nil {
			log.Fatal(err)
		}
	*/
}
