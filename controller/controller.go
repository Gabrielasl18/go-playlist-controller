package controller

import (
	"bytes" // Mantido caso queira reintroduzir comunicação, mas não usado
	"fmt"
	"io"
	"log" // Mantido caso queira reintroduzir comunicação, mas não usado
	"os"
	"strings"
	"time"

	go_m3u8 "github.com/globocom/go-m3u8"
)

// --- Constantes e Estruturas ---

const playlistPath = "./hls/playlist_360p.m3u8"

// observableURL não é mais necessário

type LightEvent struct {
	Seg    int
	Color  [3]int
	Action string
	Name   string
}

var LightEvents = []LightEvent{
	{Seg: 2, Color: [3]int{255, 255, 255}, Action: "start", Name: "Luz Branca"},
	{Seg: 12, Color: [3]int{255, 0, 0}, Action: "start", Name: "Luz Vermelha"},
	{Seg: 20, Color: [3]int{0, 0, 255}, Action: "start", Name: "Luz Azul"},
}

// --- Lógica do Observable (Consumidor) ---

func playLightEffect(device, action string, color [3]int) {
	fmt.Printf("💡 EFEITO: Dispositivo: %s | ação=%s | cor=%v\n", device, action, color)
}

// Função que era o 'checkLightEvents' no observable
func processMediaSequence(mediaSeq int, device string) {
	for _, event := range LightEvents {
		if event.Seg == mediaSeq {
			playLightEffect(device, event.Action, event.Color)
			fmt.Printf("[PDT SYNC] atingido em %d\n", mediaSeq)
			// Chamada repetida no código original, mantida para compatibilidade funcional
			playLightEffect(device, event.Action, event.Color)
		}
	}
}

// --- Lógica do Watcher (Produtor) ---

func readAndCorrectPlaylist(filename string) (io.ReadCloser, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir arquivo: %w", err)
	}
	defer file.Close()

	contentBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler todo o conteúdo do arquivo: %w", err)
	}

	// Correção de formatação
	contentString := string(contentBytes)
	correctedContent := strings.ReplaceAll(contentString, "+0000", "Z")

	return io.NopCloser(bytes.NewBufferString(correctedContent)), nil
}

// Função que era o 'extractAndSendPDT' no watcher, agora chama diretamente o processamento
func extractAndProcessPDT() error {
	playlistReader, err := readAndCorrectPlaylist(playlistPath)
	if err != nil {
		return fmt.Errorf("erro ao preparar playlist: %w", err)
	}
	defer playlistReader.Close()

	playlist, err := go_m3u8.ParsePlaylist(playlistReader)
	if err != nil {
		return fmt.Errorf("erro ao parsear playlist: %w", err)
	}

	mediaSeq := playlist.MediaSequence
	processMediaSequence(mediaSeq, "sepe_device_id") // ID do dispositivo hardcoded no original

	log.Printf("[Sync] Processado MediaSequence: %d", mediaSeq)
	return nil
}

// Função de Início do Ciclo (Substitui StartWatching)
func StartCycle() {
	log.Println("---------------------------------------------------------")
	log.Println("[Controller] Monitorando arquivo local do HLS a cada 1s...")
	log.Println("---------------------------------------------------------")
	for {
		if err := extractAndProcessPDT(); err != nil {
			log.Printf("[Controller] Erro no processamento do arquivo: %v", err)
		}

		time.Sleep(1 * time.Second)
	}
}
