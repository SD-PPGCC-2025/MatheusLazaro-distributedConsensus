package job

import (
    "fmt"
    "math/big"
    "sync"
    "time"
)

// JobStatus representa o estado de uma tarefa
type JobStatus string

const (
    StatusPending JobStatus = "PENDING"
    StatusDone    JobStatus = "DONE"
)

// Job representa uma tarefa (job) em background
type Job struct {
    ID       string    `json:"id"`
    Number   int64     `json:"number"`
    Status   JobStatus `json:"status"`
    Result   string    `json:"result,omitempty"` // fatorial em string
    Created  time.Time `json:"created_at"`
    Finished time.Time `json:"finished_at,omitempty"`
}

// JobService define as operações que podemos fazer com jobs
type JobService interface {
    CreateJob(number int64) (Job, error)
    ListJobs() []Job
    GetJob(id string) (Job, error)
}

// jobServiceImpl implementação concreta do JobService
type jobServiceImpl struct {
    mu   sync.Mutex
    jobs map[string]Job
}

// NewJobService cria uma instância do serviço
func NewJobService() JobService {
    return &jobServiceImpl{
        jobs: make(map[string]Job),
    }
}

// CreateJob cria e inicia o processamento de um job (cálculo do fatorial)
func (s *jobServiceImpl) CreateJob(number int64) (Job, error) {
    // Gerar um ID simples (pode ser UUID ou outro)
    jobID := fmt.Sprintf("job-%d", time.Now().UnixNano())

    // Monta objeto Job (status PENDING)
    newJob := Job{
        ID:      jobID,
        Number:  number,
        Status:  StatusPending,
        Created: time.Now(),
    }

    s.mu.Lock()
    // Armazena job em memória (PENDING)
    s.jobs[jobID] = newJob
    s.mu.Unlock()

    // Dispara goroutine para processar (calcular fatorial)
    go s.processJob(jobID, number)

    return newJob, nil
}

// processJob efetua o cálculo pesado (fatorial) e atualiza status
func (s *jobServiceImpl) processJob(jobID string, n int64) {
    // Calcula o fatorial (operação potencialmente pesada)
    fact := factorial(n)

    // Ao terminar, atualiza o job com resultado e status DONE
    s.mu.Lock()
    j := s.jobs[jobID]
    j.Result = fact.String()
    j.Status = StatusDone
    j.Finished = time.Now()
    s.jobs[jobID] = j
    s.mu.Unlock()
}

// ListJobs retorna todos os jobs em memória
func (s *jobServiceImpl) ListJobs() []Job {
    s.mu.Lock()
    defer s.mu.Unlock()

    jobs := make([]Job, 0, len(s.jobs))
    for _, v := range s.jobs {
        jobs = append(jobs, v)
    }
    return jobs
}

// GetJob retorna um job específico pelo ID
func (s *jobServiceImpl) GetJob(id string) (Job, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    job, ok := s.jobs[id]
    if !ok {
        return Job{}, fmt.Errorf("job %s não encontrado", id)
    }
    return job, nil
}

// factorial faz o cálculo do fatorial de n usando big.Int
func factorial(n int64) *big.Int {
    result := big.NewInt(1)
    for i := int64(2); i <= n; i++ {
        result.Mul(result, big.NewInt(i))
    }
    return result
}
