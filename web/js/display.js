// Display module for handling UI updates
export function showSkeletonLoading(processStatusDiv) {
    processStatusDiv.innerHTML = '';
    for (let i = 0; i < 3; i++) {
        const skeletonCard = document.createElement('div');
        skeletonCard.className = 'process-card skeleton';
        skeletonCard.innerHTML = `
            <div class="skeleton skeleton-text" style="width: 60%; height: 1.2rem; margin-bottom: 1rem;"></div>
            <div class="skeleton skeleton-text" style="width: 80%; height: 0.8rem; margin-bottom: 0.5rem;"></div>
            <div class="skeleton skeleton-text" style="width: 70%; height: 0.8rem; margin-bottom: 1rem;"></div>
            <div style="display: flex; gap: 0.5rem;">
                <div class="skeleton" style="flex: 1; height: 2rem;"></div>
                <div class="skeleton" style="flex: 1; height: 2rem;"></div>
                <div class="skeleton" style="flex: 1; height: 2rem;"></div>
            </div>
        `;
        processStatusDiv.appendChild(skeletonCard);
    }
}

export function displayProcesses(processStatusDiv, processes) {
    processStatusDiv.innerHTML = '';
    
    if (processes && processes.length > 0) {
        processes.forEach(process => {
            const processCard = document.createElement('div');
            processCard.className = `process-card ${process.status.toLowerCase()}`;
            
            const statusClass = process.status === 'running' ? 'status-running' :
                              process.status === 'stopped' ? 'status-stopped' : 'status-paused';
            
            processCard.innerHTML = `
                <div class="process-info">
                    <h3>
                        <span class="status-indicator ${statusClass}"></span>
                        ${process.name}
                    </h3>
                    <div class="process-details">
                        <span>Status: ${process.status}</span>
                        <span>PID: ${process.pid || 'N/A'}</span>
                    </div>
                    <div class="process-details">
                        <span>Uptime: ${process.uptime || 'N/A'}</span>
                        <span>Memory: ${process.memory || 'N/A'}</span>
                    </div>
                    <div class="process-actions">
                        <button class="btn btn-success" data-action="start" data-process="${process.name}">Start</button>
                        <button class="btn btn-danger" data-action="stop" data-process="${process.name}">Stop</button>
                        <button class="btn btn-warning" data-action="restart" data-process="${process.name}">Restart</button>
                    </div>
                </div>
            `;
            
            processStatusDiv.appendChild(processCard);
        });
    } else {
        processStatusDiv.innerHTML = '<p>No processes found.</p>';
    }
}

export function displayWorkerLogs(workerLogsContainer, logs, workerFilter, logLevelFilter) {
    workerLogsContainer.innerHTML = '';
    
    if (logs && logs.length > 0) {
        // Apply worker filter
        const selectedWorker = workerFilter.value;
        let filteredLogs = logs;
        if (selectedWorker !== 'all') {
            filteredLogs = logs.filter(log => log.worker === selectedWorker);
        }
        
        // Apply log level filter
        const filterLevel = logLevelFilter.value;
        if (filterLevel !== 'all') {
            filteredLogs = filteredLogs.filter(log => log.level === filterLevel);
        }
        
        filteredLogs.forEach(log => {
            const logEntry = document.createElement('div');
            logEntry.className = `log-entry ${log.level.toLowerCase()}`;
            logEntry.innerHTML = `
                <span class="timestamp">[${formatTimestamp(log.timestamp)}]</span>
                <span class="message">${log.worker ? `[${log.worker}] ` : ''}${log.message}</span>
            `;
            workerLogsContainer.appendChild(logEntry);
        });
        
        // Scroll to bottom
        workerLogsContainer.scrollTop = workerLogsContainer.scrollHeight;
    } else {
        workerLogsContainer.innerHTML = '<div class="log-entry"><span class="timestamp">[--:--:--]</span><span class="message">No worker logs available.</span></div>';
    }
}

export function displaySystemLogs(systemLogsContainer, logs) {
    systemLogsContainer.innerHTML = '';
    
    if (logs && logs.length > 0) {
        logs.forEach(log => {
            const logEntry = document.createElement('div');
            logEntry.className = `log-entry ${log.level.toLowerCase()}`;
            logEntry.innerHTML = `
                <span class="timestamp">[${formatTimestamp(log.timestamp)}]</span>
                <span class="message">${log.message}</span>
            `;
            systemLogsContainer.appendChild(logEntry);
        });
        
        // Scroll to bottom
        systemLogsContainer.scrollTop = systemLogsContainer.scrollHeight;
    } else {
        systemLogsContainer.innerHTML = '<div class="log-entry"><span class="timestamp">[--:--:--]</span><span class="message">No system logs available.</span></div>';
    }
}

export function displaySpecificWorkerLogs(specificWorkerLogsContainer, logs, append = false, followLogs = false) {
    if (!append) {
        specificWorkerLogsContainer.innerHTML = '';
    }
    
    if (logs && logs.length > 0) {
        logs.forEach(log => {
            const logEntry = document.createElement('div');
            logEntry.className = `log-entry ${log.level.toLowerCase()}`;
            logEntry.innerHTML = `
                <span class="timestamp">[${formatTimestamp(log.timestamp)}]</span>
                <span class="message">${log.worker ? `[${log.worker}] ` : ''}${log.message}</span>
            `;
            specificWorkerLogsContainer.appendChild(logEntry);
        });
        
        // Scroll to bottom if following logs
        if (followLogs) {
            specificWorkerLogsContainer.scrollTop = specificWorkerLogsContainer.scrollHeight;
        }
    } else {
        if (!append) {
            specificWorkerLogsContainer.innerHTML = '<div class="log-entry"><span class="timestamp">[--:--:--]</span><span class="message">No logs available for this worker.</span></div>';
        }
    }
}

function formatTimestamp(timestamp) {
    try {
        const date = new Date(timestamp);
        return date.toLocaleTimeString();
    } catch (e) {
        return timestamp;
    }
}

export function addMessageToWorkerLogs(workerLogsContainer, message, level, worker = 'system') {
    const logEntry = document.createElement('div');
    logEntry.className = `log-entry ${level}`;
    logEntry.innerHTML = `
        <span class="timestamp">[${formatTimestamp(new Date().toISOString())}]</span>
        <span class="message">[${worker}] ${message}</span>
    `;
    
    // Add to worker logs container
    workerLogsContainer.appendChild(logEntry);
    
    // Scroll to bottom
    workerLogsContainer.scrollTop = workerLogsContainer.scrollHeight;
}