// API module for handling API requests
export async function fetchProcesses() {
    try {
        const response = await fetch('/api/processes');
        if (response.ok) {
            return await response.json();
        } else {
            console.error('Failed to fetch process status:', response.status);
            return [];
        }
    } catch (error) {
        console.error('Error fetching process status:', error);
        return [];
    }
}

export async function fetchWorkerLogs() {
    try {
        const response = await fetch('/api/logs/worker');
        if (response.ok) {
            return await response.json();
        } else {
            console.error('Failed to fetch worker logs:', response.status);
            return [];
        }
    } catch (error) {
        console.error('Error fetching worker logs:', error);
        return [];
    }
}

export async function fetchSystemLogs() {
    try {
        const response = await fetch('/api/logs/system');
        if (response.ok) {
            return await response.json();
        } else {
            console.error('Failed to fetch system logs:', response.status);
            return [];
        }
    } catch (error) {
        console.error('Error fetching system logs:', error);
        return [];
    }
}

export async function fetchSpecificWorkerLogs(workerName) {
    if (!workerName) return [];
    
    try {
        const response = await fetch(`/api/logs/worker/${workerName}`);
        if (response.ok) {
            return await response.json();
        } else {
            console.error('Failed to fetch specific worker logs:', response.status);
            return [];
        }
    } catch (error) {
        console.error('Error fetching specific worker logs:', error);
        return [];
    }
}

export async function startProcess(processName) {
    try {
        const response = await fetch(`/api/processes/${processName}/start`, {
            method: 'POST'
        });
        return response.ok;
    } catch (error) {
        console.error('Error starting process:', error);
        return false;
    }
}

export async function stopProcess(processName) {
    try {
        const response = await fetch(`/api/processes/${processName}/stop`, {
            method: 'POST'
        });
        return response.ok;
    } catch (error) {
        console.error('Error stopping process:', error);
        return false;
    }
}

export async function restartProcess(processName) {
    try {
        const response = await fetch(`/api/processes/${processName}/restart`, {
            method: 'POST'
        });
        return response.ok;
    } catch (error) {
        console.error('Error restarting process:', error);
        return false;
    }
}