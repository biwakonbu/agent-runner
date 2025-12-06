import { writable, derived } from 'svelte/store';

// Wails バインディングは wails generate module で更新後に自動生成される
// 現時点ではモック実装を提供

export interface ChatMessage {
    id: string;
    sessionId: string;
    role: 'user' | 'assistant' | 'system';
    content: string;
    timestamp: string;
    generatedTasks?: string[];
}

export interface ChatSession {
    id: string;
    workspaceId: string;
    createdAt: string;
    updatedAt: string;
}

export interface ChatResponse {
    message: ChatMessage;
    generatedTasks: Array<{
        id: string;
        title: string;
        description: string;
        status: string;
        phaseName: string;
        dependencies: string[];
    }>;
    understanding: string;
    error?: string;
}

interface ChatState {
    currentSessionId: string | null;
    messages: ChatMessage[];
    isLoading: boolean;
    error: string | null;
}

function createChatStore() {
    const initialState: ChatState = {
        currentSessionId: null,
        messages: [],
        isLoading: false,
        error: null
    };

    const { subscribe, update, set } = writable<ChatState>(initialState);

    // Wails API のダイナミックインポート（ビルド時に存在しない場合に対応）
    let wailsApp: any = null;

    const loadWailsBindings = async () => {
        try {
            // Wails バインディングを動的にインポート
            const module = await import('../../wailsjs/go/main/App');
            
            // window.go が存在することを確認（Storybook 等では存在しない）
            if (!('go' in window)) {
                console.warn('[Chat] window.go not found, using mock mode');
                return false;
            }

            wailsApp = module;
            return true;
        } catch {
            console.warn('[Chat] Wails bindings not available, using mock mode');
            return false;
        }
    };

    return {
        subscribe,

        // セッション作成
        createSession: async (): Promise<string | null> => {
            update(s => ({ ...s, isLoading: true, error: null }));

            try {
                await loadWailsBindings();

                if (wailsApp?.CreateChatSession) {
                    const session = await wailsApp.CreateChatSession();
                    if (session) {
                        update(s => ({
                            ...s,
                            currentSessionId: session.id,
                            messages: [],
                            isLoading: false
                        }));
                        return session.id;
                    }
                }

                // Mock fallback
                const mockSessionId = crypto.randomUUID();
                const systemMessage: ChatMessage = {
                    id: crypto.randomUUID(),
                    sessionId: mockSessionId,
                    role: 'system',
                    content: 'チャットセッションが開始されました。開発したい機能や解決したい課題を教えてください。',
                    timestamp: new Date().toISOString()
                };

                update(s => ({
                    ...s,
                    currentSessionId: mockSessionId,
                    messages: [systemMessage],
                    isLoading: false
                }));

                return mockSessionId;
            } catch (e) {
                const error = e instanceof Error ? e.message : 'セッション作成に失敗しました';
                update(s => ({ ...s, isLoading: false, error }));
                return null;
            }
        },

        // メッセージ送信
        sendMessage: async (text: string): Promise<ChatResponse | null> => {
            if (!text.trim()) return null;

            let currentState: ChatState;
            const unsubscribe = subscribe(s => { currentState = s; });
            unsubscribe();

            if (!currentState!.currentSessionId) {
                console.error('[Chat] No active session');
                return null;
            }

            const sessionId = currentState!.currentSessionId;

            // Optimistic update: ユーザーメッセージを即座に表示
            const userMessage: ChatMessage = {
                id: crypto.randomUUID(),
                sessionId,
                role: 'user',
                content: text,
                timestamp: new Date().toISOString()
            };

            update(s => ({
                ...s,
                messages: [...s.messages, userMessage],
                isLoading: true,
                error: null
            }));

            try {
                await loadWailsBindings();

                if (wailsApp?.SendChatMessage) {
                    const response = await wailsApp.SendChatMessage(sessionId, text);

                    if (response.error) {
                        throw new Error(response.error);
                    }

                    // アシスタント応答を追加
                    update(s => ({
                        ...s,
                        messages: [...s.messages, response.message],
                        isLoading: false
                    }));

                    return response;
                }

                // Mock fallback
                await new Promise(resolve => setTimeout(resolve, 1000));

                const mockResponse: ChatResponse = {
                    message: {
                        id: crypto.randomUUID(),
                        sessionId,
                        role: 'assistant',
                        content: `Mock: ユーザーの要求を理解しました\n\n以下の 2 個のタスクを作成しました：\n\n### 概念設計\n- **Mock概念設計タスク**: モック用の概念設計タスクです\n\n### 実装\n- **Mock実装タスク**: モック用の実装タスクです`,
                        timestamp: new Date().toISOString(),
                        generatedTasks: ['mock-task-1', 'mock-task-2']
                    },
                    generatedTasks: [
                        {
                            id: 'mock-task-1',
                            title: 'Mock概念設計タスク',
                            description: 'モック用の概念設計タスクです',
                            status: 'PENDING',
                            phaseName: '概念設計',
                            dependencies: []
                        },
                        {
                            id: 'mock-task-2',
                            title: 'Mock実装タスク',
                            description: 'モック用の実装タスクです',
                            status: 'PENDING',
                            phaseName: '実装',
                            dependencies: ['mock-task-1']
                        }
                    ],
                    understanding: 'Mock: ユーザーの要求を理解しました'
                };

                update(s => ({
                    ...s,
                    messages: [...s.messages, mockResponse.message],
                    isLoading: false
                }));

                return mockResponse;
            } catch (e) {
                const error = e instanceof Error ? e.message : 'メッセージ送信に失敗しました';
                update(s => ({ ...s, isLoading: false, error }));
                return null;
            }
        },

        // 履歴読み込み
        loadHistory: async (sessionId: string): Promise<void> => {
            update(s => ({ ...s, isLoading: true, error: null }));

            try {
                await loadWailsBindings();

                if (wailsApp?.GetChatHistory) {
                    const messages = await wailsApp.GetChatHistory(sessionId);
                    update(s => ({
                        ...s,
                        currentSessionId: sessionId,
                        messages: messages || [],
                        isLoading: false
                    }));
                    return;
                }

                // Mock fallback
                update(s => ({
                    ...s,
                    currentSessionId: sessionId,
                    messages: [],
                    isLoading: false
                }));
            } catch (e) {
                const error = e instanceof Error ? e.message : '履歴の読み込みに失敗しました';
                update(s => ({ ...s, isLoading: false, error }));
            }
        },

        // 状態リセット
        reset: () => {
            set(initialState);
        }
    };
}

export const chatStore = createChatStore();

// 派生ストア: ローディング状態
export const isChatLoading = derived(chatStore, $chat => $chat.isLoading);

// 派生ストア: エラー状態
export const chatError = derived(chatStore, $chat => $chat.error);

// 派生ストア: 現在のセッションID
export const currentSessionId = derived(chatStore, $chat => $chat.currentSessionId);

// 派生ストア: メッセージ一覧
export const chatMessages = derived(chatStore, $chat => $chat.messages);
