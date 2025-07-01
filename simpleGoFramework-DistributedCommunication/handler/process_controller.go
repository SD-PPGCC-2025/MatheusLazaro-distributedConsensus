package handler

import (
    "encoding/json"
    "net/http"
    "strconv"
    "time"
    "fmt"
)

// Códigos ANSI para colorir a saída no terminal
const (
    ColorReset   = "\033[0m"
    CyanBright   = "\033[96m"
    Blue         = "\033[34m"
    Green        = "\033[32m"
    Yellow       = "\033[33m"
    Red          = "\033[31m"
    White        = "\033[97m"
)

// ProcessController lida com requisições de processamento assíncrono
type ProcessController struct {
    service ProcessService // Serviço que realiza o processamento real
}

// NewProcessController cria um novo controller com injeção de dependência do serviço
func NewProcessController(service ProcessService) *ProcessController {
    return &ProcessController{service: service}
}

// RegisterRoutes registra a rota HTTP /process para esse controller
func (pc *ProcessController) RegisterRoutes(mux *http.ServeMux) {
    mux.HandleFunc("/process", pc.handleRequest)
}

// handleRequest processa requisições POST em /process
// Espera um parâmetro delay na query string e opcionalmente um body com message
func (pc *ProcessController) handleRequest(w http.ResponseWriter, r *http.Request) {
    // Exibição de cabeçalho visual no terminal
    fmt.Println()
    fmt.Printf("%s╔══════════════════════════════════════════════════════════════╗%s\n", CyanBright, ColorReset)
    fmt.Printf("%s║            NOVO CICLO DE REQUISIÇÃO RECEBIDO                ║%s\n", CyanBright, ColorReset)
    fmt.Printf("%s╚══════════════════════════════════════════════════════════════╝%s\n", CyanBright, ColorReset)

    fmt.Printf("%s[Servidor] Aguardando pedido...%s\n", Blue, ColorReset)

    // Só aceita requisições POST
    if r.Method != http.MethodPost {
        http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
        return
    }

    // Lê e valida o parâmetro "delay" da query string
    delayStr := r.URL.Query().Get("delay")
    delay, err := strconv.Atoi(delayStr)
    if err != nil || delay <= 0 {
        http.Error(w, "Parâmetro 'delay' inválido", http.StatusBadRequest)
        return
    }

    // Lê o corpo da requisição com uma mensagem opcional
    var body struct {
        Message string `json:"message"`
    }
    _ = json.NewDecoder(r.Body).Decode(&body)

    // Gera um ID único para o tratador
    handlerID := time.Now().UnixNano()

    fmt.Printf("%s[Servidor] Criando tratador #%d para delay=%d segundos%s\n", Blue, handlerID, delay, ColorReset)

    // Executa o tratador em uma goroutine (concorrente)
    /**
        Thread: não usa diretamente. É gerenciado pelo runtime do Go, mapeia goroutines em threads nativas do SO.
        Goroutine: "função leve" que roda de forma concorrente com outras.
        Concorrência: o código é executado em paralelo (múltiplos núcleos)
        Paralelismo real: via múltiplas threads nativas

        goroutine não bloqueia o servidor, ou seja, o servidor pode continuar recebendo outras requisições.
        Quando a goroutine termina, ela imprime a resposta no terminal.
    */
    go pc.service.HandleAsync(handlerID, delay, body.Message)


    // Responde imediatamente ao cliente, mesmo com o tratamento ocorrendo em background
    fmt.Printf("%s[Servidor] Respondendo requisição para tratador #%d%s\n", Blue, handlerID, ColorReset)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":   "Tratador criado",
        "handler":  handlerID,
        "info":     "O tratamento ocorre em background.",
        "message":  body.Message,
    })
}
