package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"pupervisor/internal/service"

	"github.com/gorilla/mux"
)

type ProcessHandler struct {
	svc *service.ProcessService
}

func NewProcessHandler(svc *service.ProcessService) *ProcessHandler {
	return &ProcessHandler{svc: svc}
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func (h *ProcessHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func (h *ProcessHandler) writeError(w http.ResponseWriter, status int, err error, message string) {
	h.writeJSON(w, status, ErrorResponse{
		Error:   err.Error(),
		Message: message,
	})
}

func (h *ProcessHandler) GetProcesses(w http.ResponseWriter, r *http.Request) {
	processes := h.svc.GetProcesses()
	h.writeJSON(w, http.StatusOK, processes)
}

func (h *ProcessHandler) StartProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if err := h.svc.StartProcess(name); err != nil {
		if errors.Is(err, service.ErrProcessNotFound) {
			h.writeError(w, http.StatusNotFound, err, "Process not found: "+name)
			return
		}
		h.writeError(w, http.StatusInternalServerError, err, "Failed to start process")
		return
	}

	h.writeJSON(w, http.StatusOK, SuccessResponse{
		Status:  "started",
		Message: "Process " + name + " started successfully",
	})
}

func (h *ProcessHandler) StopProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if err := h.svc.StopProcess(name); err != nil {
		if errors.Is(err, service.ErrProcessNotFound) {
			h.writeError(w, http.StatusNotFound, err, "Process not found: "+name)
			return
		}
		h.writeError(w, http.StatusInternalServerError, err, "Failed to stop process")
		return
	}

	h.writeJSON(w, http.StatusOK, SuccessResponse{
		Status:  "stopped",
		Message: "Process " + name + " stopped successfully",
	})
}

func (h *ProcessHandler) RestartProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if err := h.svc.RestartProcess(name); err != nil {
		if errors.Is(err, service.ErrProcessNotFound) {
			h.writeError(w, http.StatusNotFound, err, "Process not found: "+name)
			return
		}
		h.writeError(w, http.StatusInternalServerError, err, "Failed to restart process")
		return
	}

	h.writeJSON(w, http.StatusOK, SuccessResponse{
		Status:  "restarted",
		Message: "Process " + name + " restarted successfully",
	})
}

func (h *ProcessHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	logs := h.svc.GetLogs(50)
	h.writeJSON(w, http.StatusOK, logs)
}

func (h *ProcessHandler) GetWorkerLogs(w http.ResponseWriter, r *http.Request) {
	logs := h.svc.GetWorkerLogs(50)
	h.writeJSON(w, http.StatusOK, logs)
}

func (h *ProcessHandler) GetSystemLogs(w http.ResponseWriter, r *http.Request) {
	logs := h.svc.GetSystemLogs(50)
	h.writeJSON(w, http.StatusOK, logs)
}

func (h *ProcessHandler) GetWorkerSpecificLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	workerName := vars["workerName"]

	logs := h.svc.GetWorkerSpecificLogs(workerName, 50)
	h.writeJSON(w, http.StatusOK, logs)
}
