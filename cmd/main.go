package main

import (
	"fmt"
	"go-playlist-controller/controller"
	"log"
	"net/http"
)

func main() {
	port := 8080

	// Serve todo o diretório pai
	http.Handle("/", http.FileServer(http.Dir("../")))

	// Endpoint para pegar efeito atual das luzes
	http.HandleFunc("/light", controller.LightHandler)

	go controller.StartCycle()

	fmt.Printf("Servidor iniciado!\n")
	fmt.Printf("Acesse no navegador: http://localhost:%d/\n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
