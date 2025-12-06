/**
 * WBS（Work Breakdown Structure）ビュー用ストア
 *
 * タスクをマイルストーン・フェーズ別にツリー構造化し、
 * 折りたたみ状態と進捗率を管理
 */

import { writable, derived, get } from 'svelte/store';
import { tasks } from './taskStore';
import type { Task, PhaseName, TaskStatus } from '../types';

// WBS ツリーノードの型
export interface WBSNode {
  id: string;
  type: 'milestone' | 'phase' | 'task';
  label: string;
  milestone?: string;
  phaseName?: PhaseName;
  task?: Task;
  children: WBSNode[];
  level: number;
  progress: {
    completed: number;
    total: number;
    percentage: number;
  };
}

// フェーズの表示順序
const phaseOrder: PhaseName[] = ['概念設計', '実装設計', '実装', '検証', ''];

// フェーズのラベル
const phaseLabels: Record<PhaseName, string> = {
  概念設計: 'Concept Design',
  実装設計: 'Architecture Design',
  実装: 'Implementation',
  検証: 'Verification',
  '': 'Other',
};

const defaultMilestoneLabel = 'General';

// 折りたたみ状態ストア
function createExpandedStore() {
  const { subscribe, update, set } = writable<Set<string>>(new Set());

  return {
    subscribe,

    // ノードを展開/折りたたみ切り替え
    toggle: (nodeId: string) => {
      update((expanded) => {
        const newSet = new Set(expanded);
        if (newSet.has(nodeId)) {
          newSet.delete(nodeId);
        } else {
          newSet.add(nodeId);
        }
        return newSet;
      });
    },

    // ノードを展開
    expand: (nodeId: string) => {
      update((expanded) => {
        const newSet = new Set(expanded);
        newSet.add(nodeId);
        return newSet;
      });
    },

    // ノードを折りたたむ
    collapse: (nodeId: string) => {
      update((expanded) => {
        const newSet = new Set(expanded);
        newSet.delete(nodeId);
        return newSet;
      });
    },

    // 全て展開
    expandAll: () => {
      const currentTasks = get(tasks);
      const milestoneKeys = new Set<string>();
      currentTasks.forEach((t) => milestoneKeys.add(t.milestone?.trim() || 'default'));

      const next = new Set<string>();
      for (const m of milestoneKeys) {
        const milestoneId = m === '' ? 'default' : m;
        next.add(`milestone-${milestoneId}`);
        for (const p of phaseOrder) {
          next.add(`phase-${milestoneId}-${p}`);
        }
      }
      set(next);
    },

    // 全て折りたたむ
    collapseAll: () => {
      set(new Set());
    },

    // 初期状態にリセット（全フェーズを展開）
    reset: () => {
      const currentTasks = get(tasks);
      const milestoneKeys = new Set<string>();
      currentTasks.forEach((t) => milestoneKeys.add(t.milestone?.trim() || 'default'));

      const next = new Set<string>();
      for (const m of milestoneKeys) {
        const milestoneId = m === '' ? 'default' : m;
        next.add(`milestone-${milestoneId}`);
        for (const p of phaseOrder) {
          next.add(`phase-${milestoneId}-${p}`);
        }
      }
      set(next);
    },
  };
}

export const expandedNodes = createExpandedStore();

// ビューモード（Graph or WBS）
export type ViewMode = 'graph' | 'wbs';

function createViewModeStore() {
  const { subscribe, set } = writable<ViewMode>('graph');

  return {
    subscribe,
    setGraph: () => set('graph'),
    setWBS: () => set('wbs'),
    toggle: () => {
      let current: ViewMode;
      subscribe((v) => (current = v))();
      set(current! === 'graph' ? 'wbs' : 'graph');
    },
  };
}

export const viewMode = createViewModeStore();

// タスクの完了判定
function isTaskCompleted(status: TaskStatus): boolean {
  return status === 'SUCCEEDED' || status === 'COMPLETED' || status === 'CANCELED';
}

// 進捗を計算
function calculateProgress(tasks: Task[]): {
  completed: number;
  total: number;
  percentage: number;
} {
  const total = tasks.length;
  const completed = tasks.filter((t) => isTaskCompleted(t.status)).length;
  const percentage = total > 0 ? Math.round((completed / total) * 100) : 0;
  return { completed, total, percentage };
}

// WBS ツリー構造を生成
export const wbsTree = derived(tasks, ($tasks): WBSNode[] => {
  // マイルストーン -> フェーズ -> タスク の三層で構築
  const milestones = new Map<string, Map<PhaseName, Task[]>>();

  for (const task of $tasks) {
    const milestoneKey = task.milestone?.trim() || '';
    const phase = (task.phaseName || '') as PhaseName;
    const phaseKey = phaseOrder.includes(phase) ? phase : '';

    if (!milestones.has(milestoneKey)) {
      const phaseMap = new Map<PhaseName, Task[]>();
      for (const p of phaseOrder) {
        phaseMap.set(p, []);
      }
      milestones.set(milestoneKey, phaseMap);
    }

    const phaseMap = milestones.get(milestoneKey)!;
    phaseMap.get(phaseKey)!.push(task);
  }

  // ツリー構造を生成
  const tree: WBSNode[] = [];

  for (const [milestoneKey, phaseMap] of milestones) {
    // ミドル層（フェーズ）ノードを構築
    const phaseNodes: WBSNode[] = [];
    let milestoneTasks: Task[] = [];

    const milestoneId = milestoneKey !== '' ? milestoneKey : 'default';

    for (const phase of phaseOrder) {
      const phaseTasks = phaseMap.get(phase) || [];
      if (phaseTasks.length === 0) {
        continue;
      }

      milestoneTasks = milestoneTasks.concat(phaseTasks);

      const phaseProgress = calculateProgress(phaseTasks);
      phaseNodes.push({
        id: `phase-${milestoneId}-${phase}`,
        type: 'phase',
        label: phaseLabels[phase],
        phaseName: phase,
        milestone: milestoneKey,
        children: phaseTasks.map((task) => ({
          id: task.id,
          type: 'task' as const,
          label: task.title,
          task,
          children: [],
          level: 2,
          progress: {
            completed: isTaskCompleted(task.status) ? 1 : 0,
            total: 1,
            percentage: isTaskCompleted(task.status) ? 100 : 0,
          },
        })),
        level: 1,
        progress: phaseProgress,
      });
    }

    // タスクが無いマイルストーンはスキップ
    if (phaseNodes.length === 0) {
      continue;
    }

    const milestoneProgress = calculateProgress(milestoneTasks);
    const milestoneLabel = milestoneKey !== '' ? milestoneKey : defaultMilestoneLabel;

    tree.push({
      id: `milestone-${milestoneId}`,
      type: 'milestone',
      label: milestoneLabel,
      milestone: milestoneKey,
      phaseName: '',
      children: phaseNodes,
      level: 0,
      progress: milestoneProgress,
    });
  }

  return tree;
});

// 全体の進捗率
export const overallProgress = derived(tasks, ($tasks) => {
  return calculateProgress($tasks);
});

// フラット化したWBSノードリスト（展開状態を考慮）
export const flattenedWBSNodes = derived(
  [wbsTree, expandedNodes],
  ([$tree, $expanded]): WBSNode[] => {
    const result: WBSNode[] = [];

    function flatten(nodes: WBSNode[], parentExpanded: boolean) {
      for (const node of nodes) {
        if (parentExpanded) {
          result.push(node);
        }

        if (node.children.length > 0) {
          const isExpanded = $expanded.has(node.id);
          flatten(node.children, parentExpanded && isExpanded);
        }
      }
    }

    flatten($tree, true);
    return result;
  }
);
