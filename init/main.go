package main

import (
	"log"
	"net/http"
	"todostudy/network"

	"github.com/urfave/negroni"
)

// cmd 호출

func main() {
	mux := network.MakeHandler()
	n := negroni.Classic()
	n.UseHandler(mux)
	err := http.ListenAndServe(":8080", n)
	if err != nil {
		log.Fatal(err)
	}
}
