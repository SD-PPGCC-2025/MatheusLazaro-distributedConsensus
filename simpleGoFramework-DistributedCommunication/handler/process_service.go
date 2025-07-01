package handler

import (
    "fmt"
    "time"
)

// ProcessService define a interface para simular o tratamento de requisições com delay
type ProcessService interface {
    Handle(delay int) string
    HandleAsync(handlerID int64, delay int, msg string)
}

// processServiceImpl é a implementação da interface ProcessService
type processServiceImpl struct{}

// NewProcessService cria uma nova instância do serviço de processamento
func NewProcessService() ProcessService {
    return &processServiceImpl{}
}

// Handle simula o processamento com um delay artificial
func (s *processServiceImpl) Handle(delay int) string {
    time.Sleep(time.Duration(delay) * time.Second) // Espera N segundos
    return fmt.Sprintf("Processamento concluído em %d segundos", delay)
}

// HandleAsync gerencia o processamento assíncrono
func (s *processServiceImpl) HandleAsync(handlerID int64, delay int, msg string) {
    fmt.Printf("\n───────────────────────────── Tratador #%d ─────────────────────────────\n", handlerID)
    fmt.Printf("%s[Tratador #%d] Tratando pedido...%s\n", Green, handlerID, ColorReset)

    if msg != "" {
        fmt.Printf("%s[Tratador #%d] Mensagem: %s%s\n", Yellow, handlerID, msg, ColorReset)
    }

    // Chama o serviço de processamento com delay
    result := s.Handle(delay)

    fmt.Printf("%s[Tratador #%d] Respondeu: \"%s\"%s\n", White, handlerID, result, ColorReset)
    fmt.Printf("%s[Tratador #%d] Tratador termina.%s\n", Red, handlerID, ColorReset)
}
