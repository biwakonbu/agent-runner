/**
 * E2E テスト用 Wails Runtime モック
 *
 * EventsOn などの Wails ランタイム関数をモック
 */
console.log('[Mock] Wails runtime mock loaded');

// イベントリスナー管理
const eventListeners = new Map();

export function EventsOn(eventName, callback) {
    console.log("[Mock] EventsOn registered:", eventName);
    if (!eventListeners.has(eventName)) {
        eventListeners.set(eventName, []);
    }
    eventListeners.get(eventName).push(callback);
}

export function EventsOff(eventName) {
    console.log("[Mock] EventsOff:", eventName);
    eventListeners.delete(eventName);
}

export function EventsEmit(eventName, data) {
    console.log("[Mock] EventsEmit:", eventName, data);
    const listeners = eventListeners.get(eventName) || [];
    listeners.forEach(callback => callback(data));
}

// テスト用にイベントを発火するヘルパー
if (typeof window !== 'undefined') {
    window.__mockEmitEvent = EventsEmit;
}
