package main

import (
	"go-playlist-controller/observable"
	"go-playlist-controller/server"
	"go-playlist-controller/watcher"
	"log"
)

func main() {
	go server.StartServer()
	go observable.StartObservable()
	go watcher.StartWatching()

	log.Println("🎬 Sistema iniciado (Server 8080, Observable 8888, Watcher ativo)")
	select {}
}
