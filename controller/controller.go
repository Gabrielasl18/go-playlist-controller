package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	go_m3u8 "github.com/globocom/go-m3u8"
)

const (
	playlistPath  = "../hls/playlist_1080p.m3u8"
	lightJSONPath = "../controller/light_events.json"
)

type LightEvent struct {
	SegStart int    `json:"SegStart"`
	SegEnd   int    `json:"SegEnd"`
	Action   string `json:"Action"`
	Color    [3]int `json:"Color"`
	Name     string `json:"Name"`
}

var (
	LightEvents   []LightEvent
	currentEffect string
	effectMutex   sync.RWMutex
)

func LoadLightEvents() error {
	data, err := os.ReadFile(lightJSONPath)
	if err != nil {
		return fmt.Errorf("erro ao ler JSON de eventos: %w", err)
	}
	if err := json.Unmarshal(data, &LightEvents); err != nil {
		return fmt.Errorf("erro ao parsear JSON de eventos: %w", err)
	}
	fmt.Printf("[Controller] Carregados %d eventos de luz do JSON\n", len(LightEvents))
	return nil
}

func playLightEffect(action string, color [3]int, name string) string {
	msg := fmt.Sprintf("💡 EFEITO: ação=%s | cor=%v | nome=%s\n",
		action, color, name)
	fmt.Print(msg)
	return msg
}

func processMediaSequence(mediaSeq int) string {
	var result string
	for _, event := range LightEvents {
		if mediaSeq >= event.SegStart && mediaSeq <= event.SegEnd {
			result += playLightEffect(event.Action, event.Color, event.Name)
		}
	}
	if result == "" {
		result = fmt.Sprintf("Nenhum efeito no segmento %d\n", mediaSeq)
	}
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
		return nil, fmt.Errorf("erro ao ler conteúdo: %w", err)
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

	if err := LoadLightEvents(); err != nil {
		fmt.Printf("[Controller] Erro ao carregar eventos: %v\n", err)
	}

	for {
		mediaSeq, err := extractMediaSequence()
		if err != nil {
			fmt.Printf("[Controller] Erro ao extrair media sequence: %v\n", err)
		} else {
			processMediaSequence(mediaSeq)
		}
		time.Sleep(1 * time.Second)
	}
}

// LightHandler serve o efeito atual via HTTP
func LightHandler(w http.ResponseWriter, r *http.Request) {
	effectMutex.RLock()
	defer effectMutex.RUnlock()
	w.Write([]byte(currentEffect))
}
