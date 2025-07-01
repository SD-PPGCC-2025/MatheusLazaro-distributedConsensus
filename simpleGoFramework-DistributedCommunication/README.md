# simpleGoFramework

```go
go mod tidy
go run main.go
```

### Requisição 1
```sh
curl -X POST "http://localhost:5000/process?delay=3" \
  -H "Content-Type: application/json" \
  -d '{"message": "processando dados da tarefa X"}'
```


Framework (pasta framework): gera a estrutura de servidor com App e Controller.

job_service.go: demonstra multithreading (goroutines) no processamento do fatorial, atualizando o estado do job.

job_controller.go: API REST que cria e consulta as tarefas.

Cada requisição POST “solicita” uma tarefa, e o servidor responde rapidamente (status PENDING), enquanto o processamento roda em paralelo.