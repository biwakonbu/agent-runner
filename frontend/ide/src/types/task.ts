/**
 * タスク関連の型定義
 *
 * Zod スキーマからの再エクスポート
 * 既存コードとの互換性を維持
 */

export {
  type TaskStatus,
  type PhaseName,
  type Task,
  type TaskNode,
  statusToCssClass,
  statusLabels,
  type AttemptStatus,
  type Attempt,
  attemptStatusLabels,
  type PoolSummary,
} from '../schemas';
