// Process actions module for handling process control
import { startProcess as apiStartProcess, stopProcess as apiStopProcess, restartProcess as apiRestartProcess } from './api.js';
import { addMessageToWorkerLogs } from './display.js';

export async function startProcess(processName, btn) {
    const originalText = btn.innerHTML;
    btn.innerHTML = '<span class="loading-indicator"></span> Starting...';
    btn.disabled = true;
    
    try {
        const success = await apiStartProcess(processName);
        
        if (success) {
            addMessageToWorkerLogs(document.getElementById('worker-logs-container'), `Process ${processName} started successfully`, 'info', 'system');
        } else {
            addMessageToWorkerLogs(document.getElementById('worker-logs-container'), `Failed to start process ${processName}`, 'error', 'system');
        }
    } catch (error) {
        console.error('Error starting process:', error);
        addMessageToWorkerLogs(document.getElementById('worker-logs-container'), `Error starting process ${processName}`, 'error', 'system');
    } finally {
        btn.innerHTML = originalText;
        btn.disabled = false;
    }
}

export async function stopProcess(processName, btn) {
    const originalText = btn.innerHTML;
    btn.innerHTML = '<span class="loading-indicator"></span> Stopping...';
    btn.disabled = true;
    
    try {
        const success = await apiStopProcess(processName);
        
        if (success) {
            addMessageToWorkerLogs(document.getElementById('worker-logs-container'), `Process ${processName} stopped successfully`, 'info', 'system');
        } else {
            addMessageToWorkerLogs(document.getElementById('worker-logs-container'), `Failed to stop process ${processName}`, 'error', 'system');
        }
    } catch (error) {
        console.error('Error stopping process:', error);
        addMessageToWorkerLogs(document.getElementById('worker-logs-container'), `Error stopping process ${processName}`, 'error', 'system');
    } finally {
        btn.innerHTML = originalText;
        btn.disabled = false;
    }
}

export async function restartProcess(processName, btn) {
    const originalText = btn.innerHTML;
    btn.innerHTML = '<span class="loading-indicator"></span> Restarting...';
    btn.disabled = true;
    
    try {
        const success = await apiRestartProcess(processName);
        
        if (success) {
            addMessageToWorkerLogs(document.getElementById('worker-logs-container'), `Process ${processName} restarted successfully`, 'info', 'system');
        } else {
            addMessageToWorkerLogs(document.getElementById('worker-logs-container'), `Failed to restart process ${processName}`, 'error', 'system');
        }
    } catch (error) {
        console.error('Error restarting process:', error);
        addMessageToWorkerLogs(document.getElementById('worker-logs-container'), `Error restarting process ${processName}`, 'error', 'system');
    } finally {
        btn.innerHTML = originalText;
        btn.disabled = false;
    }
}