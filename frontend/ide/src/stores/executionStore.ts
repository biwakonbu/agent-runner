/**
 * 実行状態管理ストア
 *
 * ExecutionOrchestrator の状態をフロントエンドで管理
 * Wails バインディングが生成された後、実際の API に接続
 */

import { writable } from 'svelte/store';
import { EventsOn } from '../../wailsjs/runtime/runtime';
import { StartExecution, PauseExecution, ResumeExecution, StopExecution, GetExecutionState } from '../../wailsjs/go/main/App';
import { toasts } from './toastStore';

export type ExecutionState = 'IDLE' | 'RUNNING' | 'PAUSED';

export const executionState = writable<ExecutionState>('IDLE');

// Wails イベントリスナー初期化（Wails バインディング生成後に有効化）
export function initExecutionEvents(): void {
    EventsOn('execution:stateChange', (event: { newState: ExecutionState }) => {
        executionState.set(event.newState);
    });
}

// 実行開始
export async function startExecution(): Promise<void> {
    try {
        await StartExecution();
        toasts.add('Execution started', 'success');
    } catch (e: any) {
        console.error('Failed to start execution:', e);
        const errorMsg = String(e);
        if (errorMsg.includes('CLI session not found')) {
            toasts.add('CLI session not found. Please login via terminal.', 'error', 10000, {
                label: 'Retry',
                onClick: () => startExecution(),
            });
        } else {
            toasts.add(`Failed to start: ${e}`, 'error');
        }
    }
}

// 実行一時停止
export async function pauseExecution(): Promise<void> {
    try {
        await PauseExecution();
        toasts.add('Execution paused', 'info');
    } catch (e) {
        console.error('Failed to pause execution:', e);
        toasts.add(`Failed to pause: ${e}`, 'error');
    }
}

// 実行再開
export async function resumeExecution(): Promise<void> {
    try {
        await ResumeExecution();
        toasts.add('Execution resumed', 'success');
    } catch (e) {
        console.error('Failed to resume execution:', e);
        toasts.add(`Failed to resume: ${e}`, 'error');
    }
}

// 実行停止
export async function stopExecution(): Promise<void> {
    try {
        await StopExecution();
        toasts.add('Execution stopped', 'info');
    } catch (e) {
        console.error('Failed to stop execution:', e);
        toasts.add(`Failed to stop: ${e}`, 'error');
    }
}

// サーバーから実行状態を同期
export async function syncExecutionState(): Promise<void> {
    try {
        const state = await GetExecutionState();
        executionState.set(state as ExecutionState);
    } catch (e) {
        console.error('Failed to sync execution state:', e);
    }
}
