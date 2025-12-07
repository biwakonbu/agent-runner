/**
 * E2E テスト用 Wails Models モック
 *
 * orchestrator などの Go モデル型をモック
 */
console.log('[Mock] Wails models mock loaded');

export const orchestrator = {
    BacklogType: {
        FAILURE: 'FAILURE',
        QUESTION: 'QUESTION',
        BLOCKER: 'BLOCKER'
    }
};
