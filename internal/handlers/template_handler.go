package handlers

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"sync"

	"pupervisor/internal/models"
	"pupervisor/internal/service"
)

type PageData struct {
	Title                  string
	PageTitle              string
	CurrentPage            string
	ActiveProcesses        int
	TotalProcesses         int
	TotalWorkers           int
	CpuUsage               int
	MemoryUsage            int
	ActiveProcessesPercent int
	WorkersOnlinePercent   int
	Processes              []models.Process
	Logs                   []models.LogEntry
	WorkerLogs             []models.LogEntry
	SystemLogs             []models.LogEntry
	Workers                []string
	SystemName             string
	RefreshInterval        string
	LogRetention           string
	EmailNotifications     bool
	PushNotifications      bool
	CriticalAlerts         bool
	ProcessEvents          bool
}

type TemplateHandler struct {
	templates *template.Template
	svc       *service.ProcessService
	mu        sync.RWMutex
}

func NewTemplateHandler(templatesFS fs.FS, svc *service.ProcessService) (*TemplateHandler, error) {
	tmpl, err := template.ParseFS(templatesFS, "*.html")
	if err != nil {
		return nil, err
	}

	return &TemplateHandler{
		templates: tmpl,
		svc:       svc,
	}, nil
}

func (th *TemplateHandler) buildPageData(currentPage, pageTitle string) PageData {
	processes := th.svc.GetProcesses()
	activeProcesses, totalProcesses := th.svc.GetStats()
	logs := th.svc.GetLogs(10)
	workerLogs := th.svc.GetWorkerLogs(10)
	systemLogs := th.svc.GetSystemLogs(10)

	activePercent := 0
	if totalProcesses > 0 {
		activePercent = (activeProcesses * 100) / totalProcesses
	}

	return PageData{
		Title:                  "Pupervisor - " + pageTitle,
		PageTitle:              pageTitle,
		CurrentPage:            currentPage,
		ActiveProcesses:        activeProcesses,
		TotalProcesses:         totalProcesses,
		TotalWorkers:           24,
		CpuUsage:               42,
		MemoryUsage:            68,
		ActiveProcessesPercent: activePercent,
		WorkersOnlinePercent:   85,
		SystemName:             "Pupervisor System",
		RefreshInterval:        "10s",
		LogRetention:           "7d",
		EmailNotifications:     false,
		PushNotifications:      true,
		CriticalAlerts:         true,
		ProcessEvents:          true,
		Workers:                []string{"worker-1", "worker-2", "worker-3"},
		Processes:              processes,
		Logs:                   logs,
		WorkerLogs:             workerLogs,
		SystemLogs:             systemLogs,
	}
}

func (th *TemplateHandler) ServeTemplate(templateName, currentPage, pageTitle string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := th.buildPageData(currentPage, pageTitle)

		th.mu.RLock()
		tmpl := th.templates
		th.mu.RUnlock()

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if err := tmpl.ExecuteTemplate(w, templateName+".html", data); err != nil {
			log.Printf("Error executing template %s: %v", templateName, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
