package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Caminho local para o arquivo da playlist
const playlistFile = "playlist_360p.m3u8"

// handlerPlaylist serve o conteúdo atual da playlist
func handlerPlaylist(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(playlistFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro lendo playlist: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

	log.Printf("[Server] Enviando playlist (%s) para %s", playlistFile, r.RemoteAddr)
}

// StartServer inicia o servidor HTTP na porta 8080
func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/playlist360p.m3u8", handlerPlaylist)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	log.Println("[Server] Servindo playlist em http://localhost:8080/playlist360p.m3u8")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
