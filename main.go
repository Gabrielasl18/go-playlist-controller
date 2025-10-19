package main

import (
	"go-playlist-controller/controller"
	"log"
)

func main() {
	go controller.StartCycle()

	log.Println("🎬 Sistema iniciado (Server 8080, Controller unificado ativo)")
	select {}
}
