/**
 * Zod スキーマ定義
 *
 * Wails が生成する any 型のバインディングを
 * ランタイムで検証し、型安全な値に変換する
 */

import { z } from 'zod';

// TaskStatus スキーマ
export const TaskStatusSchema = z.enum([
  'PENDING',
  'READY',
  'RUNNING',
  'SUCCEEDED',
  'FAILED',
  'CANCELED',
  'BLOCKED',
]);

export type TaskStatus = z.infer<typeof TaskStatusSchema>;

// Task スキーマ
export const TaskSchema = z.object({
  id: z.string(),
  title: z.string(),
  status: TaskStatusSchema,
  poolId: z.string(),
  createdAt: z.string().datetime({ offset: true }).or(z.string()),
  updatedAt: z.string().datetime({ offset: true }).or(z.string()),
  startedAt: z.string().datetime({ offset: true }).or(z.string()).optional(),
  doneAt: z.string().datetime({ offset: true }).or(z.string()).optional(),
});

export type Task = z.infer<typeof TaskSchema>;

// Task 配列スキーマ
export const TaskArraySchema = z.array(TaskSchema);

// Workspace スキーマ
export const WorkspaceSchema = z.object({
  version: z.string(),
  projectRoot: z.string(),
  displayName: z.string(),
  createdAt: z.string().datetime({ offset: true }).or(z.string()),
  lastOpenedAt: z.string().datetime({ offset: true }).or(z.string()),
});

export type Workspace = z.infer<typeof WorkspaceSchema>;

// グリッド配置用のタスクノード
export const TaskNodeSchema = z.object({
  task: TaskSchema,
  col: z.number().int().nonnegative(),
  row: z.number().int().nonnegative(),
});

export type TaskNode = z.infer<typeof TaskNodeSchema>;

// ステータスからCSS変数名サフィックスへの変換
export function statusToCssClass(status: TaskStatus): string {
  return status.toLowerCase();
}

// ステータスの表示名
export const statusLabels: Record<TaskStatus, string> = {
  PENDING: '待機中',
  READY: '準備完了',
  RUNNING: '実行中',
  SUCCEEDED: '成功',
  FAILED: '失敗',
  CANCELED: 'キャンセル',
  BLOCKED: 'ブロック',
};
