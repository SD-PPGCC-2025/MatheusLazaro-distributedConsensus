package framework

import (
    "net/http"
)

// Controller define a interface para todos os controllers da aplicação
// Cada controller deve implementar a função RegisterRoutes para registrar suas rotas
type Controller interface {
    RegisterRoutes(mux *http.ServeMux)
}
