package observable

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Estrutura de evento de luz
type LightEvent struct {
	Time   time.Time
	Color  [3]int
	Action string
	Name   string
}

// Lista de eventos (parseando strings uma única vez)
var LightEvents = []LightEvent{
	{Time: mustParseTime("2025-10-12T03:52:33.063Z"), Color: [3]int{255, 255, 255}, Action: "start", Name: "Luz Branca"},
	{Time: mustParseTime("2025-10-12T03:53:51.668Z"), Color: [3]int{255, 0, 0}, Action: "start", Name: "Luz Vermelha"},
	{Time: mustParseTime("2025-10-12T02:10:08.070Z"), Color: [3]int{0, 0, 255}, Action: "start", Name: "Luz Azul"},
}

func mustParseTime(timeStr string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, timeStr)
	if err != nil {
		log.Fatalf("[Observable] Erro ao parsear tempo '%s': %v", timeStr, err)
	}
	return t.UTC()
}

func playLightEffect(device, action string, color [3]int) {
	fmt.Printf("💡 EFEITO: Dispositivo: %s | ação=%s | cor=%v\n", device, action, color)
}

type PDTRequest struct {
	Timestamp string `json:"timestamp"`
}

func handleProgramDateTime(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req PDTRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	currentTime, err := time.Parse(time.RFC3339Nano, req.Timestamp)
	if err != nil {
		http.Error(w, "Timestamp inválido", http.StatusBadRequest)
		return
	}
	currentTime = currentTime.UTC()

	log.Printf("[Observable] PDT recebido: %s", currentTime.Format(time.RFC3339Nano))

	checkLightEvents(currentTime, "sepe_device_id")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func checkLightEvents(currentTime time.Time, device string) {
	if currentTime.IsZero() {
		return
	}

	for _, event := range LightEvents {
		if event.Time == currentTime {
			fmt.Println("AQUIIII")
			playLightEffect(device, event.Action, event.Color)

			// O tempo atual é MAIOR ou IGUAL (After ou Equal) ao tempo do evento?

			// Verifica se a hora foi atingida E (o evento nunca foi disparado ou o último foi o anterior)
			fmt.Printf("[PDT SYNC] %s atingido em %s\n",
				event.Name, event.Time.Format(time.RFC3339))

			// Dispara o efeito de luz
			playLightEffect(device, event.Action, event.Color)

			// Atualiza o índice (ponteiro é necessário para modificar a variável externa)

			// Se estiver buscando (seek) em VOD, você continuaria no loop para pegar eventos passados.
			// Mas para execução sequencial, um 'break' aqui pode ser apropriado.
		}
	}
}

func StartObservable() {
	mux := http.NewServeMux()
	mux.HandleFunc("/program-date-time", handleProgramDateTime)

	log.Println("[Observable] Escutando em http://localhost:8888/program-date-time")
	if err := http.ListenAndServe(":8888", mux); err != nil {
		log.Fatalf("Erro Observable: %v", err)
	}
}
