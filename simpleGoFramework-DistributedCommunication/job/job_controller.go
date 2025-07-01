package job

import (
    "encoding/json"
    "net/http"
    "strings"
    // "time"

    // "github.com/matheusLazaroCC-UFG/simpleGoFramework/framework"
)

// JobController expõe endpoints para criar e consultar jobs
type JobController struct {
    service JobService
}

// NewJobController injeta o JobService
func NewJobController(service JobService) *JobController {
    return &JobController{service: service}
}

// RegisterRoutes implementa a interface Controller
func (jc *JobController) RegisterRoutes(mux *http.ServeMux) {
    // POST /jobs  => cria job
    // GET /jobs   => lista jobs
    // GET /jobs/{id} => retorna status e resultado
    mux.HandleFunc("/jobs", jc.handleJobs)
    mux.HandleFunc("/jobs/", jc.handleJobByID)
}

// handleJobs lida com POST e GET em /jobs
func (jc *JobController) handleJobs(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodPost:
        jc.createJob(w, r)
    case http.MethodGet:
        jc.listJobs(w, r)
    default:
        http.Error(w, "método não suportado", http.StatusMethodNotAllowed)
    }
}

// createJob cria uma nova tarefa (job) e retorna o ID
func (jc *JobController) createJob(w http.ResponseWriter, r *http.Request) {
    var body struct {
        Number int64 `json:"number"`
    }
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, "corpo inválido", http.StatusBadRequest)
        return
    }
    if body.Number < 0 {
        http.Error(w, "número deve ser >= 0", http.StatusBadRequest)
        return
    }

    job, err := jc.service.CreateJob(body.Number)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(job)
}

// listJobs retorna todos os jobs
func (jc *JobController) listJobs(w http.ResponseWriter, r *http.Request) {
    jobs := jc.service.ListJobs()
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(jobs)
}

// handleJobByID lida com GET /jobs/{id}
func (jc *JobController) handleJobByID(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "método não suportado", http.StatusMethodNotAllowed)
        return
    }

    // Extrai {id} de /jobs/{id}
    pathID := strings.TrimPrefix(r.URL.Path, "/jobs/")
    job, err := jc.service.GetJob(pathID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(job)
}
