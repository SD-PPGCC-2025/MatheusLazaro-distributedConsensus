package main

import (
    "github.com/matheusLazaroCC-UFG/simpleGoFramework/framework"
    "github.com/matheusLazaroCC-UFG/simpleGoFramework/job"
	"github.com/matheusLazaroCC-UFG/simpleGoFramework/handler"
    "log"
)

func main() {
    // Configurações básicas
    config := &framework.Config{
        Port: "5000",
    }

    // Cria a aplicação principal
    app := framework.NewApp(config)

    // Instancia o serviço de "jobs" (tarefas em background)
    jobService := job.NewJobService()
    // Cria o controller de "jobs" e injeta o serviço
    jobController := job.NewJobController(jobService)
    // Registra o controller no mini-framework
    app.RegisterController(jobController)

	// Instancia o serviço de "handler" (tarefas em background)
    processService := handler.NewProcessService()
    // Cria o controller de "jobs" e injeta o serviço
    processController := handler.NewProcessController(processService)
    // Registra o controller no mini-framework
    app.RegisterController(processController)

    log.Println("Servidor rodando em http://localhost:5000")

    // Inicia o servidor
    app.Start()
}
