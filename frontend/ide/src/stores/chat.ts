/**
 * チャットデータ管理ストア
 */

import { writable, get } from 'svelte/store';
import { GetChatHistory, SendChatMessage, CreateChatSession } from '../../wailsjs/go/main/App';
import { EventsOn } from '../../wailsjs/wailsjs/runtime/runtime';
import type { ChatMessage } from '../types';

export interface ChatResponse {
    message: ChatMessage;
    generatedTasks: Array<{
        id: string;
        title: string;
        status: string;
        // ... 他のフィールド
    }>;
    understanding: string;
    conflicts?: Array<{
        file: string;
        tasks: string[];
        warning: string;
    }>;
    error?: string;
}

// Chat Log Interface
export interface ChatLogEntry {
    step: string;
    message: string;
    timestamp: string;
}

// チャットメッセージストア
function createMessagesStore() {
  const { subscribe, set, update } = writable<ChatMessage[]>([]);

  return {
    subscribe,
    setMessages: (messages: ChatMessage[]) => set(messages),
    addMessage: (message: ChatMessage) => update((msgs) => [...msgs, message]),
    clear: () => set([]),
  };
}

// チャットログストア
function createChatLogStore() {
  const { subscribe, update } = writable<ChatLogEntry[]>([]);

  return {
    subscribe,
    addLog: (entry: ChatLogEntry) => {
      update((logs) => [...logs, entry]);
    },
    clear: () => update(() => []),
  };
}

export const chatMessages = createMessagesStore();
export const chatLog = createChatLogStore();
export const currentSessionId = writable<string | null>(null);
export const isChatLoading = writable<boolean>(false);
export const chatError = writable<string | null>(null);

// ストア
const chatStore = {
    // セッション作成
    createSession: async () => {
        try {
            const session = await CreateChatSession();
            if (!session?.id) return;

            // #region agent log
            fetch('http://127.0.0.1:7242/ingest/e0c5926c-4256-4f95-83f1-ee92ab435f0c', {
              method: 'POST',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify({
                sessionId: 'debug-session',
                runId: 'pre-fix',
                hypothesisId: 'D',
                location: 'chat.ts:createSession',
                message: 'chat session created',
                data: { sessionId: session.id },
                timestamp: Date.now(),
              }),
            }).catch(() => {});
            // #endregion agent log

            currentSessionId.set(session.id);
            // セッション切替時に既存ログとメッセージをクリア
            chatMessages.clear();
            chatLog.clear();
            chatError.set(null);

            try {
                const history = await GetChatHistory(session.id);
                chatMessages.setMessages(history);
            } catch (e) {
                console.error('Failed to load chat history:', e);
            }
        } catch (e) {
            console.error('Failed to create session:', e);
            // #region agent log
            fetch('http://127.0.0.1:7242/ingest/e0c5926c-4256-4f95-83f1-ee92ab435f0c', {
              method: 'POST',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify({
                sessionId: 'debug-session',
                runId: 'pre-fix',
                hypothesisId: 'D',
                location: 'chat.ts:createSession',
                message: 'failed to create session',
                data: { error: e instanceof Error ? e.message : String(e) },
                timestamp: Date.now(),
              }),
            }).catch(() => {});
            // #endregion agent log
        }
    },

    // メッセージ送信
    sendMessage: async (content: string): Promise<ChatResponse | null> => {
        isChatLoading.set(true);
        chatError.set(null);
        let sessionId: string | null = get(currentSessionId);

        if (!sessionId) {
            console.warn('No active session. Attempting to recreate...');
            await chatStore.createSession();
            sessionId = get(currentSessionId);
            if (!sessionId) {
                console.error('No active session inside sendMessage');
                isChatLoading.set(false);
                chatError.set('No active chat session');
                return null;
            }
        }

        try {
            // #region agent log
            fetch('http://127.0.0.1:7242/ingest/e0c5926c-4256-4f95-83f1-ee92ab435f0c', {
              method: 'POST',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify({
                sessionId: 'debug-session',
                runId: 'pre-fix',
                hypothesisId: 'E',
                location: 'chat.ts:sendMessage',
                message: 'sending chat message',
                data: { sessionId, length: content.length },
                timestamp: Date.now(),
              }),
            }).catch(() => {});
            // #endregion agent log

            const response = await SendChatMessage(sessionId, content);
            
            if (response.error) {
                console.error('Chat error:', response.error);
                chatError.set(response.error);
            } else {
                 // ユーザーメッセージとアシスタントメッセージは backend 側で保存済みなので
                 // 履歴を再取得するか、レスポンスから追加する
                 // ここではレスポンスから追加
                 // user message is implicitly added by optimistic update usually, but here simple:
                 
                 // Wait, ChatHandler saves user message already.
                 // We should reload history or append both manually?
                 // Let's reload history to be safe and consistent
                 const history = await GetChatHistory(sessionId!);
                 chatMessages.setMessages(history);
                 chatError.set(null);
            }

            // #region agent log
            fetch('http://127.0.0.1:7242/ingest/e0c5926c-4256-4f95-83f1-ee92ab435f0c', {
              method: 'POST',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify({
                sessionId: 'debug-session',
                runId: 'pre-fix',
                hypothesisId: 'E',
                location: 'chat.ts:sendMessage',
                message: 'chat response received',
                data: {
                  sessionId,
                  hasError: Boolean(response.error),
                  generatedTasks: response.generatedTasks?.length ?? 0,
                },
                timestamp: Date.now(),
              }),
            }).catch(() => {});
            // #endregion agent log

            return response as ChatResponse;

        } catch (e) {
            console.error('Failed to send message:', e);
            chatError.set(e instanceof Error ? e.message : 'Failed to send message');
            // #region agent log
            fetch('http://127.0.0.1:7242/ingest/e0c5926c-4256-4f95-83f1-ee92ab435f0c', {
              method: 'POST',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify({
                sessionId: 'debug-session',
                runId: 'pre-fix',
                hypothesisId: 'E',
                location: 'chat.ts:sendMessage',
                message: 'chat send failed',
                data: { error: e instanceof Error ? e.message : String(e) },
                timestamp: Date.now(),
              }),
            }).catch(() => {});
            // #endregion agent log
            return null;
        } finally {
            isChatLoading.set(false);
        }
    }
};

export { chatStore };

// Wailsイベントリスナーの初期化
export function initChatEvents() {
    EventsOn('chat:progress', (event: { step: string; message: string; timestamp: string }) => {
        console.log('Chat Progress:', event);
        chatLog.addLog({
            step: event.step,
            message: event.message,
            timestamp: event.timestamp
        });

        // ローディング状態の制御
        if (event.step === 'Completed' || event.step === 'Failed') {
            isChatLoading.set(false);
        } else {
            isChatLoading.set(true);
        }
    });
}
