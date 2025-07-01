package framework

import (
    "log"
    "net/http"
)

// Config define as configurações básicas do servidor
type Config struct {
    Port string // Porta onde o servidor irá escutar (ex: "5000")
}

// App representa a aplicação web construída com o mini-framework
type App struct {
    mux    *http.ServeMux // Multiplexador de rotas HTTP
    config *Config         // Configurações da aplicação
}

// NewApp cria uma nova instância da aplicação, inicializando o mux e as configurações
func NewApp(config *Config) *App {
    return &App{
        mux:    http.NewServeMux(),
        config: config,
    }
}

// RegisterController registra os endpoints de um controller na aplicação
func (a *App) RegisterController(c Controller) {
    c.RegisterRoutes(a.mux)
}

// Start inicia o servidor HTTP com as rotas registradas
func (a *App) Start() {
    log.Printf("Servidor iniciado na porta %s...\n", a.config.Port)
    log.Fatal(http.ListenAndServe(":"+a.config.Port, a.mux)) // Inicia o servidor e escuta a porta
}
