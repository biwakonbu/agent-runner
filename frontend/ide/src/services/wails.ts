/**
 * Wails バインディングサービス層
 *
 * Wails が生成する any 型のバインディングを Zod で検証し、
 * 型安全な値としてアプリケーションに提供する
 */

import {
  TaskSchema,
  TaskArraySchema,
  WorkspaceSchema,
  type Task,
  type Workspace,
} from '../schemas';

// Wails 生成バインディングのインポート
// 実行時に存在しない場合のフォールバックを提供
let wailsApp: {
  SelectWorkspace: () => Promise<string>;
  GetWorkspace: (id: string) => Promise<unknown>;
  ListTasks: () => Promise<unknown>;
  CreateTask: (title: string, poolId: string) => Promise<unknown>;
  RunTask: (taskId: string) => Promise<void>;
} | null = null;

// Wails バインディングを動的にロード
async function getWailsApp() {
  if (wailsApp) return wailsApp;

  try {
    // Wails 環境でのみ利用可能
    const module = await import('../../wailsjs/go/main/App');
    wailsApp = module;
    return wailsApp;
  } catch {
    console.warn('Wails bindings not available. Using mock mode.');
    return null;
  }
}

// パース結果の型
type ParseResult<T> =
  | { success: true; data: T }
  | { success: false; error: Error };

/**
 * Workspace を選択するダイアログを開く
 */
export async function selectWorkspace(): Promise<string | null> {
  const app = await getWailsApp();
  if (!app) return null;

  const result = await app.SelectWorkspace();
  return result || null;
}

/**
 * Workspace を ID で取得
 */
export async function getWorkspace(
  id: string
): Promise<ParseResult<Workspace>> {
  const app = await getWailsApp();
  if (!app) {
    return { success: false, error: new Error('Wails not available') };
  }

  try {
    const raw = await app.GetWorkspace(id);
    const parsed = WorkspaceSchema.safeParse(raw);

    if (!parsed.success) {
      console.error('Workspace validation failed:', parsed.error);
      return { success: false, error: new Error(parsed.error.message) };
    }

    return { success: true, data: parsed.data };
  } catch (err) {
    return {
      success: false,
      error: err instanceof Error ? err : new Error(String(err)),
    };
  }
}

/**
 * タスク一覧を取得
 */
export async function listTasks(): Promise<ParseResult<Task[]>> {
  const app = await getWailsApp();
  if (!app) {
    return { success: false, error: new Error('Wails not available') };
  }

  try {
    const raw = await app.ListTasks();
    const parsed = TaskArraySchema.safeParse(raw);

    if (!parsed.success) {
      console.error('Task list validation failed:', parsed.error);
      return { success: false, error: new Error(parsed.error.message) };
    }

    return { success: true, data: parsed.data };
  } catch (err) {
    return {
      success: false,
      error: err instanceof Error ? err : new Error(String(err)),
    };
  }
}

/**
 * タスクを作成
 */
export async function createTask(
  title: string,
  poolId: string
): Promise<ParseResult<Task>> {
  const app = await getWailsApp();
  if (!app) {
    return { success: false, error: new Error('Wails not available') };
  }

  try {
    const raw = await app.CreateTask(title, poolId);
    const parsed = TaskSchema.safeParse(raw);

    if (!parsed.success) {
      console.error('Created task validation failed:', parsed.error);
      return { success: false, error: new Error(parsed.error.message) };
    }

    return { success: true, data: parsed.data };
  } catch (err) {
    return {
      success: false,
      error: err instanceof Error ? err : new Error(String(err)),
    };
  }
}

/**
 * タスクを実行
 */
export async function runTask(taskId: string): Promise<ParseResult<void>> {
  const app = await getWailsApp();
  if (!app) {
    return { success: false, error: new Error('Wails not available') };
  }

  try {
    await app.RunTask(taskId);
    return { success: true, data: undefined };
  } catch (err) {
    return {
      success: false,
      error: err instanceof Error ? err : new Error(String(err)),
    };
  }
}
