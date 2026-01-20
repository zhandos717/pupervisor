package api

import (
	"io/fs"
	"net/http"

	"pupervisor/internal/handlers"
	"pupervisor/internal/middleware"
	"pupervisor/internal/service"

	"github.com/gorilla/mux"
)

type Router struct {
	*mux.Router
}

func NewRouter(pm *service.ProcessManager, templatesFS, staticFS fs.FS) (*Router, error) {
	r := mux.NewRouter()

	tmplHandler, err := handlers.NewTemplateHandler(templatesFS)
	if err != nil {
		return nil, err
	}

	procHandler := handlers.NewProcessHandler(pm)

	// Health check endpoints
	r.HandleFunc("/health", handlers.HealthCheck).Methods(http.MethodGet)
	r.HandleFunc("/ready", handlers.ReadyCheck).Methods(http.MethodGet)

	// Web UI routes
	r.HandleFunc("/", tmplHandler.ServeTemplate("dashboard")).Methods(http.MethodGet)
	r.HandleFunc("/processes", tmplHandler.ServeTemplate("processes")).Methods(http.MethodGet)
	r.HandleFunc("/logs", tmplHandler.ServeTemplate("logs")).Methods(http.MethodGet)
	r.HandleFunc("/settings", tmplHandler.ServeTemplate("settings")).Methods(http.MethodGet)

	// Serve static files
	staticHandler := http.FileServer(http.FS(staticFS))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticHandler))

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/processes", procHandler.GetProcesses).Methods(http.MethodGet)
	api.HandleFunc("/processes/{name}/start", procHandler.StartProcess).Methods(http.MethodPost)
	api.HandleFunc("/processes/{name}/stop", procHandler.StopProcess).Methods(http.MethodPost)
	api.HandleFunc("/processes/{name}/restart", procHandler.RestartProcess).Methods(http.MethodPost)
	api.HandleFunc("/logs", procHandler.GetLogs).Methods(http.MethodGet)
	api.HandleFunc("/logs/worker", procHandler.GetWorkerLogs).Methods(http.MethodGet)
	api.HandleFunc("/logs/system", procHandler.GetSystemLogs).Methods(http.MethodGet)
	api.HandleFunc("/logs/worker/{workerName}", procHandler.GetWorkerSpecificLogs).Methods(http.MethodGet)

	// Apply middleware
	r.Use(middleware.Recovery)
	r.Use(middleware.Logging)

	return &Router{Router: r}, nil
}
