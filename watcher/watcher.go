package watcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	go_m3u8 "github.com/globocom/go-m3u8"
)

// Caminho do arquivo da playlist no disco local
const playlistPath = "playlist_360p.m3u8"

// Endereço do endpoint do Observable
const observableURL = "http://localhost:8888/program-date-time"

// readAndCorrectPlaylist lê o arquivo, substitui "+0000" por "Z" e retorna o conteúdo corrigido como io.ReadCloser.
func readAndCorrectPlaylist(filename string) (io.ReadCloser, error) {
	log.Printf("[Reader] Abrindo arquivo: %s", filename)
	// 1. Abre o arquivo
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir arquivo: %w", err)
	}
	defer file.Close() // Fecha o arquivo de disco

	// 2. Lê todo o conteúdo do arquivo para um slice de bytes
	contentBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler todo o conteúdo do arquivo: %w", err)
	}

	// 3. Converte para string e faz a substituição do formato
	contentString := string(contentBytes)
	correctedContent := strings.ReplaceAll(contentString, "+0000", "Z")

	fmt.Println(correctedContent)

	// 4. Cria um novo Reader (io.ReadCloser) a partir da string corrigida na memória
	log.Println("[Reader] Formato PDT corrigido de +0000 para Z na memória.")
	return io.NopCloser(bytes.NewBufferString(correctedContent)), nil
}

// extractAndSendPDT faz a leitura, correção, parsing e envio dos PDTs.
func extractAndSendPDT() error {
	playlistReader, err := readAndCorrectPlaylist(playlistPath)
	if err != nil {
		log.Fatalf("erro ao preparar playlist: %v", err)
	}
	defer playlistReader.Close()

	log.Println("[Parser] Fazendo parsing do conteúdo corrigido...")
	playlist, err := go_m3u8.ParsePlaylist(playlistReader)
	if err != nil {
		log.Fatalf("erro ao parsear playlist: %v", err)
	}

	segments := playlist.Segments()

	log.Printf("[Parser] Encontrados %d segmentos na playlist.\n", len(segments))

	foundPDT := false

	// Itera e imprime os dados de cada segmento
	for _, seg := range segments {
		if seg == nil {
			continue
		}
		sendPDTToObservable(seg.HLSElement.Details["ProgramDateTime"])
		foundPDT = true
	}

	if !foundPDT {
		log.Println("[Watcher] Aviso: Nenhuma tag #EXT-X-PROGRAM-DATE-TIME encontrada nos segmentos (Possível playlist VOD sem PDT ou erro de formatação/parsing).")
	}

	return nil // Retorna nil se o processo de leitura e envio for bem-sucedido
}

// sendPDTToObservable envia o PDT extraído via POST para o endpoint
func sendPDTToObservable(timestamp string) error {
	body := map[string]string{"timestamp": timestamp}
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("erro ao serializar JSON: %w", err)
	}

	resp, err := http.Post(observableURL, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("erro ao enviar POST para %s: %w", observableURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Tenta ler a resposta para incluir na mensagem de erro
		errorBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("observable retornou status %d. Resposta: %s", resp.StatusCode, string(errorBody))
	}

	log.Printf("[Watcher] Enviado PDT: %s", timestamp)
	return nil
}

// StartWatching inicia o loop de monitoramento da playlist
func StartWatching() {
	log.Println("---------------------------------------------------------")
	log.Println("[Watcher] Monitorando arquivo local do HLS a cada 1s...")
	log.Println("---------------------------------------------------------")
	for {
		// Chama a função de extração e processamento do arquivo local
		if err := extractAndSendPDT(); err != nil {
			// Loga erros que saem da cadeia de leitura/parsing/envio
			log.Printf("[Watcher] Erro no processamento do arquivo: %v", err)
		}

		time.Sleep(1 * time.Second)
	}
}
