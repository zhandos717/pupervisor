package service

import (
	"errors"
	"sync"
	"time"

	"pupervisor/internal/models"
)

var (
	ErrProcessNotFound = errors.New("process not found")
)

type ProcessService struct {
	mu           sync.RWMutex
	processes    []models.Process
	logs         []models.LogEntry
	workerLogs   []models.LogEntry
	systemLogs   []models.LogEntry
	maxLogSize   int
}

func NewProcessService() *ProcessService {
	ps := &ProcessService{
		maxLogSize: 1000,
		processes: []models.Process{
			{Name: "web-server", Status: "running", Pid: 1234, Uptime: "2h 15m", Memory: "45MB", CPU: "2.5%"},
			{Name: "api-service", Status: "stopped", Pid: 0, Uptime: "N/A", Memory: "N/A", CPU: "0%"},
			{Name: "database", Status: "running", Pid: 5678, Uptime: "24h 30m", Memory: "120MB", CPU: "5.2%"},
			{Name: "worker-queue", Status: "paused", Pid: 9012, Uptime: "5h 42m", Memory: "23MB", CPU: "1.1%"},
		},
		logs: []models.LogEntry{
			{Timestamp: time.Now().Add(-30 * time.Minute).Format(time.RFC3339), Message: "Server started successfully", Level: "info"},
			{Timestamp: time.Now().Add(-25 * time.Minute).Format(time.RFC3339), Message: "Process web-server started", Level: "info"},
			{Timestamp: time.Now().Add(-20 * time.Minute).Format(time.RFC3339), Message: "Process api-service stopped", Level: "warning"},
			{Timestamp: time.Now().Add(-15 * time.Minute).Format(time.RFC3339), Message: "Memory usage high for database process", Level: "warning"},
		},
		workerLogs: []models.LogEntry{
			{Timestamp: time.Now().Add(-30 * time.Minute).Format(time.RFC3339), Message: "Worker initialized", Level: "info", Worker: "worker-1"},
			{Timestamp: time.Now().Add(-25 * time.Minute).Format(time.RFC3339), Message: "Processing job #123", Level: "info", Worker: "worker-1"},
			{Timestamp: time.Now().Add(-20 * time.Minute).Format(time.RFC3339), Message: "Job #123 completed successfully", Level: "info", Worker: "worker-1"},
			{Timestamp: time.Now().Add(-15 * time.Minute).Format(time.RFC3339), Message: "Error processing job #124: timeout", Level: "error", Worker: "worker-1"},
			{Timestamp: time.Now().Add(-10 * time.Minute).Format(time.RFC3339), Message: "Retrying job #124", Level: "warning", Worker: "worker-1"},
		},
		systemLogs: []models.LogEntry{
			{Timestamp: time.Now().Add(-30 * time.Minute).Format(time.RFC3339), Message: "System started", Level: "info"},
			{Timestamp: time.Now().Add(-25 * time.Minute).Format(time.RFC3339), Message: "Health check passed", Level: "info"},
			{Timestamp: time.Now().Add(-20 * time.Minute).Format(time.RFC3339), Message: "Memory usage at 65%", Level: "info"},
			{Timestamp: time.Now().Add(-15 * time.Minute).Format(time.RFC3339), Message: "CPU usage high", Level: "warning"},
		},
	}
	return ps
}

func (ps *ProcessService) GetProcesses() []models.Process {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	result := make([]models.Process, len(ps.processes))
	copy(result, ps.processes)
	return result
}

func (ps *ProcessService) FindProcess(name string) (models.Process, bool) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, p := range ps.processes {
		if p.Name == name {
			return p, true
		}
	}
	return models.Process{}, false
}

func (ps *ProcessService) StartProcess(name string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for i, p := range ps.processes {
		if p.Name == name {
			ps.processes[i].Status = "running"
			ps.processes[i].Pid = int(time.Now().UnixNano() % 65535)
			ps.processes[i].Uptime = "Just started"
			ps.processes[i].Memory = "0MB"
			ps.processes[i].CPU = "0%"

			ps.addLogUnsafe(models.LogEntry{
				Timestamp: time.Now().Format(time.RFC3339),
				Message:   "Process " + name + " started",
				Level:     "info",
			})
			return nil
		}
	}
	return ErrProcessNotFound
}

func (ps *ProcessService) StopProcess(name string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for i, p := range ps.processes {
		if p.Name == name {
			ps.processes[i].Status = "stopped"
			ps.processes[i].Pid = 0
			ps.processes[i].Uptime = "N/A"
			ps.processes[i].Memory = "N/A"
			ps.processes[i].CPU = "0%"

			ps.addLogUnsafe(models.LogEntry{
				Timestamp: time.Now().Format(time.RFC3339),
				Message:   "Process " + name + " stopped",
				Level:     "info",
			})
			return nil
		}
	}
	return ErrProcessNotFound
}

func (ps *ProcessService) RestartProcess(name string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for i, p := range ps.processes {
		if p.Name == name {
			ps.processes[i].Status = "running"
			ps.processes[i].Pid = int(time.Now().UnixNano() % 65535)
			ps.processes[i].Uptime = "Just restarted"
			ps.processes[i].Memory = "0MB"
			ps.processes[i].CPU = "0%"

			ps.addLogUnsafe(models.LogEntry{
				Timestamp: time.Now().Format(time.RFC3339),
				Message:   "Process " + name + " restarted",
				Level:     "info",
			})
			return nil
		}
	}
	return ErrProcessNotFound
}

func (ps *ProcessService) addLogUnsafe(entry models.LogEntry) {
	ps.logs = append(ps.logs, entry)
	if len(ps.logs) > ps.maxLogSize {
		ps.logs = ps.logs[len(ps.logs)-ps.maxLogSize:]
	}
}

func (ps *ProcessService) GetLogs(limit int) []models.LogEntry {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	return ps.getLastN(ps.logs, limit)
}

func (ps *ProcessService) GetWorkerLogs(limit int) []models.LogEntry {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	return ps.getLastN(ps.workerLogs, limit)
}

func (ps *ProcessService) GetSystemLogs(limit int) []models.LogEntry {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	return ps.getLastN(ps.systemLogs, limit)
}

func (ps *ProcessService) GetWorkerSpecificLogs(workerName string, limit int) []models.LogEntry {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	var filtered []models.LogEntry
	for _, log := range ps.workerLogs {
		if log.Worker == workerName {
			filtered = append(filtered, log)
		}
	}
	return ps.getLastN(filtered, limit)
}

func (ps *ProcessService) getLastN(logs []models.LogEntry, n int) []models.LogEntry {
	if n <= 0 || len(logs) == 0 {
		return []models.LogEntry{}
	}

	start := 0
	if len(logs) > n {
		start = len(logs) - n
	}

	result := make([]models.LogEntry, len(logs[start:]))
	copy(result, logs[start:])
	return result
}

func (ps *ProcessService) GetStats() (activeProcesses, totalProcesses int) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, p := range ps.processes {
		if p.Status == "running" {
			activeProcesses++
		}
	}
	return activeProcesses, len(ps.processes)
}
