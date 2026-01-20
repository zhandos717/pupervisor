// Main application module
import { fetchProcesses, fetchWorkerLogs, fetchSystemLogs, fetchSpecificWorkerLogs } from './api.js';
import { showSkeletonLoading, displayProcesses, displayWorkerLogs, displaySystemLogs, displaySpecificWorkerLogs } from './display.js';
import { startProcess, stopProcess, restartProcess } from './process-actions.js';

document.addEventListener('DOMContentLoaded', function() {
    const refreshBtn = document.getElementById('refresh-btn');
    const clearLogsBtn = document.getElementById('clear-logs-btn');
    const logLevelFilter = document.getElementById('log-level-filter');
    const workerFilter = document.getElementById('worker-filter');
    const specificWorkerSelector = document.getElementById('specific-worker-selector');
    const followLogsBtn = document.getElementById('follow-logs-btn');
    const processStatusDiv = document.getElementById('process-status');
    const workerLogsContainer = document.getElementById('worker-logs-container');
    const systemLogsContainer = document.getElementById('system-logs-container');
    const specificWorkerLogsContainer = document.getElementById('specific-worker-logs-container');
    
    // Track follow logs state
    let followLogs = false;
    let currentWorker = '';
    let followLogsInterval = null;
    
    // Load initial data
    loadProcessStatus();
    loadWorkerLogs();
    loadSystemLogs();
    
    // Set up refresh button
    refreshBtn.addEventListener('click', function() {
        loadProcessStatus();
        loadWorkerLogs();
        loadSystemLogs();
        if (currentWorker) {
            loadSpecificWorkerLogs(currentWorker);
        }
    });
    
    // Set up clear logs button
    clearLogsBtn.addEventListener('click', function() {
        workerLogsContainer.innerHTML = '';
    });
    
    // Set up log level filter
    logLevelFilter.addEventListener('change', function() {
        loadWorkerLogs();
    });
    
    // Set up worker filter
    workerFilter.addEventListener('change', function() {
        loadWorkerLogs();
    });
    
    // Set up specific worker selector
    specificWorkerSelector.addEventListener('change', function() {
        currentWorker = specificWorkerSelector.value;
        if (currentWorker) {
            loadSpecificWorkerLogs(currentWorker);
            // Start following logs if the button is active
            if (followLogs) {
                startFollowingLogs();
            }
        } else {
            displaySpecificWorkerLogs(specificWorkerLogsContainer, [{timestamp: '--:--:--', message: 'Select a worker to view logs.', level: 'info'}], false, followLogs);
            stopFollowingLogs();
        }
    });
    
    // Set up follow logs button
    followLogsBtn.addEventListener('click', function() {
        followLogs = !followLogs;
        followLogsBtn.classList.toggle('btn-primary', followLogs);
        followLogsBtn.classList.toggle('btn-secondary', !followLogs);
        followLogsBtn.textContent = followLogs ? 'Stop Following' : 'Follow Logs';
        
        if (followLogs && currentWorker) {
            startFollowingLogs();
        } else {
            stopFollowingLogs();
        }
    });
    
    // Add event delegation for process action buttons
    processStatusDiv.addEventListener('click', function(event) {
        if (event.target.hasAttribute('data-action')) {
            const action = event.target.getAttribute('data-action');
            const processName = event.target.getAttribute('data-process');
            
            switch(action) {
                case 'start':
                    startProcess(processName, event.target);
                    loadProcessStatus(); // Refresh status after action
                    break;
                case 'stop':
                    stopProcess(processName, event.target);
                    loadProcessStatus(); // Refresh status after action
                    break;
                case 'restart':
                    restartProcess(processName, event.target);
                    loadProcessStatus(); // Refresh status after action
                    break;
            }
        }
    });
    
    // Refresh data every 30 seconds
    setInterval(function() {
        loadProcessStatus();
        loadWorkerLogs();
        loadSystemLogs();
        if (currentWorker && !followLogs) { // Only update if not in follow mode
            loadSpecificWorkerLogs(currentWorker);
        }
    }, 3000);
    
    function startFollowingLogs() {
        stopFollowingLogs(); // Clear any existing interval
        followLogsInterval = setInterval(function() {
            if (currentWorker) {
                loadSpecificWorkerLogs(currentWorker, true); // Append new logs instead of replacing
            }
        }, 2000); // Update every 2 seconds
    }
    
    function stopFollowingLogs() {
        if (followLogsInterval) {
            clearInterval(followLogsInterval);
            followLogsInterval = null;
        }
    }
    
    async function loadProcessStatus() {
        showSkeletonLoading(processStatusDiv);
        try {
            const processes = await fetchProcesses();
            displayProcesses(processStatusDiv, processes);
        } finally {
            hideSkeletonLoading();
        }
    }
    
    async function loadWorkerLogs() {
        const logs = await fetchWorkerLogs();
        displayWorkerLogs(workerLogsContainer, logs, workerFilter, logLevelFilter);
    }
    
    async function loadSystemLogs() {
        const logs = await fetchSystemLogs();
        displaySystemLogs(systemLogsContainer, logs);
    }
    
    async function loadSpecificWorkerLogs(workerName, append = false) {
        if (!workerName) return;
        
        const logs = await fetchSpecificWorkerLogs(workerName);
        displaySpecificWorkerLogs(specificWorkerLogsContainer, logs, append, followLogs);
    }
    
    function hideSkeletonLoading() {
        // Skeleton loading is automatically replaced when real data is loaded
    }
});