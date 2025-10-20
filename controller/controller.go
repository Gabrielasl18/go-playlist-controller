package controller

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	go_m3u8 "github.com/globocom/go-m3u8"
)

const playlistPath = "../hls/playlist_1080p.m3u8"

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

var (
	currentEffect string
	effectMutex   sync.RWMutex
)

func playLightEffect(device, action string, color [3]int, name string) string {
	msg := fmt.Sprintf("💡 EFEITO: Dispositivo=%s | ação=%s | cor=%v | nome=%s\n",
		device, action, color, name)
	fmt.Print(msg)
	return msg
}

func processMediaSequence(mediaSeq int, device string) string {
	var result string
	for _, event := range LightEvents {
		if event.Seg == mediaSeq {
			result += playLightEffect(device, event.Action, event.Color, event.Name)
		}
	}
	if result == "" {
		result = fmt.Sprintf("Nenhum efeito no segmento %d\n", mediaSeq)
	}
	// Atualiza efeito atual protegido por mutex
	effectMutex.Lock()
	currentEffect = result
	effectMutex.Unlock()
	return result
}

func readAndCorrectPlaylist(filename string) (io.ReadCloser, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir arquivo: %w", err)
	}
	defer file.Close()

	contentBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler todo o conteúdo: %w", err)
	}

	correctedContent := strings.ReplaceAll(string(contentBytes), "+0000", "Z")
	return io.NopCloser(bytes.NewBufferString(correctedContent)), nil
}

func extractMediaSequence() (int, error) {
	playlistReader, err := readAndCorrectPlaylist(playlistPath)
	if err != nil {
		return 0, err
	}
	defer playlistReader.Close()

	playlist, err := go_m3u8.ParsePlaylist(playlistReader)
	if err != nil {
		return 0, err
	}
	return playlist.MediaSequence, nil
}

func StartCycle() {
	fmt.Println("---------------------------------------------------------")
	fmt.Println("[Controller] Monitorando playlist a cada 1s...")
	fmt.Println("---------------------------------------------------------")
	for {
		mediaSeq, err := extractMediaSequence()
		if err != nil {
			fmt.Printf("[Controller] Erro ao extrair media sequence: %v\n", err)
		} else {
			processMediaSequence(mediaSeq, "sepe_device_id")
		}
		time.Sleep(1 * time.Second)
	}
}

func LightHandler(w http.ResponseWriter, r *http.Request) {
	effectMutex.RLock()
	defer effectMutex.RUnlock()
	w.Write([]byte(currentEffect))
}
