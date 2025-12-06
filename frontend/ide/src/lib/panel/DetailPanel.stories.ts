import type { Meta, StoryObj } from "@storybook/svelte";
import DetailPanelPreview from "./DetailPanelPreview.svelte";
import type { Task, Attempt } from "../../types";

const meta = {
  title: "IDE/DetailPanel",
  component: DetailPanelPreview,
  tags: ["autodocs"],
  argTypes: {
    task: {
      control: "object",
      description: "選択されたタスク",
    },
    attempts: {
      control: "object",
      description: "実行履歴",
    },
    isRunning: {
      control: "boolean",
      description: "タスク実行中フラグ",
    },
    loadingAttempts: {
      control: "boolean",
      description: "実行履歴読み込み中フラグ",
    },
  },
  parameters: {
    layout: "centered",
    docs: {
      description: {
        component:
          "タスクの詳細情報を表示するサイドパネル。タスク名、ステータス、メタ情報、実行履歴を表示します。",
      },
    },
  },
  decorators: [
    () => ({
      Component: DetailPanelPreview,
      props: {},
      // Wrapper でパネル表示用の高さを確保
    }),
  ],
} satisfies Meta<DetailPanelPreview>;

export default meta;
type Story = StoryObj<typeof meta>;

// サンプルタスク
const baseTask: Task = {
  id: "task-12345",
  title: "ユーザー認証機能の実装",
  status: "PENDING",
  poolId: "codegen",
  createdAt: "2024-12-05T10:30:00Z",
  updatedAt: "2024-12-05T10:30:00Z",
  dependencies: [],
};

// サンプル実行履歴
const sampleAttempts: Attempt[] = [
  {
    id: "attempt-3",
    taskId: "task-12345",
    status: "SUCCEEDED",
    startedAt: "2024-12-05T11:00:00Z",
    finishedAt: "2024-12-05T11:05:00Z",
  },
  {
    id: "attempt-2",
    taskId: "task-12345",
    status: "FAILED",
    startedAt: "2024-12-05T10:45:00Z",
    finishedAt: "2024-12-05T10:48:00Z",
    errorSummary: "テストが失敗しました: 認証トークンの検証エラー",
  },
  {
    id: "attempt-1",
    taskId: "task-12345",
    status: "FAILED",
    startedAt: "2024-12-05T10:35:00Z",
    finishedAt: "2024-12-05T10:38:00Z",
    errorSummary: "コンパイルエラー: 未定義の変数 'userId'",
  },
];

// 空状態（タスク未選択）
export const Empty: Story = {
  args: {
    task: null,
    attempts: [],
    isRunning: false,
    loadingAttempts: false,
  },
  parameters: {
    docs: {
      description: {
        story: "タスクが選択されていない状態。",
      },
    },
  },
};

// タスク選択状態（PENDING）
export const SelectedTaskPending: Story = {
  args: {
    task: { ...baseTask, status: "PENDING" },
    attempts: [],
    isRunning: false,
    loadingAttempts: false,
  },
  parameters: {
    docs: {
      description: {
        story: "待機中のタスクを選択した状態。",
      },
    },
  },
};

// タスク選択状態（READY）
export const SelectedTaskReady: Story = {
  args: {
    task: { ...baseTask, status: "READY" },
    attempts: [],
    isRunning: false,
    loadingAttempts: false,
  },
};

// 実行中タスク
export const RunningTask: Story = {
  args: {
    task: {
      ...baseTask,
      status: "RUNNING",
      startedAt: "2024-12-05T10:35:00Z",
    },
    attempts: [
      {
        id: "attempt-1",
        taskId: "task-12345",
        status: "RUNNING",
        startedAt: "2024-12-05T10:35:00Z",
      },
    ],
    isRunning: true,
    loadingAttempts: false,
  },
  parameters: {
    docs: {
      description: {
        story: "実行中のタスク。実行ボタンは無効化されています。",
      },
    },
  },
};

// 成功タスク
export const SucceededTask: Story = {
  args: {
    task: {
      ...baseTask,
      status: "SUCCEEDED",
      startedAt: "2024-12-05T10:35:00Z",
      doneAt: "2024-12-05T10:40:00Z",
    },
    attempts: [
      {
        id: "attempt-1",
        taskId: "task-12345",
        status: "SUCCEEDED",
        startedAt: "2024-12-05T10:35:00Z",
        finishedAt: "2024-12-05T10:40:00Z",
      },
    ],
    isRunning: false,
    loadingAttempts: false,
  },
  parameters: {
    docs: {
      description: {
        story: "成功したタスク。完了日時が表示されます。",
      },
    },
  },
};

// 失敗タスク（エラー表示）
export const FailedTask: Story = {
  args: {
    task: {
      ...baseTask,
      status: "FAILED",
      startedAt: "2024-12-05T10:35:00Z",
      doneAt: "2024-12-05T10:38:00Z",
    },
    attempts: [
      {
        id: "attempt-1",
        taskId: "task-12345",
        status: "FAILED",
        startedAt: "2024-12-05T10:35:00Z",
        finishedAt: "2024-12-05T10:38:00Z",
        errorSummary:
          "エラー: テストが失敗しました\n\n詳細:\n- ユーザー認証のテストケース 3/10 が失敗\n- タイムアウト: 30秒を超過",
      },
    ],
    isRunning: false,
    loadingAttempts: false,
  },
  parameters: {
    docs: {
      description: {
        story: "失敗したタスク。エラーサマリが表示されます。",
      },
    },
  },
};

// キャンセル済みタスク
export const CanceledTask: Story = {
  args: {
    task: {
      ...baseTask,
      status: "CANCELED",
      startedAt: "2024-12-05T10:35:00Z",
      doneAt: "2024-12-05T10:36:00Z",
    },
    attempts: [
      {
        id: "attempt-1",
        taskId: "task-12345",
        status: "CANCELED",
        startedAt: "2024-12-05T10:35:00Z",
        finishedAt: "2024-12-05T10:36:00Z",
      },
    ],
    isRunning: false,
    loadingAttempts: false,
  },
};

// ブロック中タスク
export const BlockedTask: Story = {
  args: {
    task: {
      ...baseTask,
      status: "BLOCKED",
    },
    attempts: [],
    isRunning: false,
    loadingAttempts: false,
  },
  parameters: {
    docs: {
      description: {
        story: "依存タスクが未完了のためブロックされているタスク。",
      },
    },
  },
};

// 複数の実行履歴あり
export const WithMultipleAttempts: Story = {
  args: {
    task: {
      ...baseTask,
      status: "SUCCEEDED",
      startedAt: "2024-12-05T10:35:00Z",
      doneAt: "2024-12-05T11:05:00Z",
    },
    attempts: sampleAttempts,
    isRunning: false,
    loadingAttempts: false,
  },
  parameters: {
    docs: {
      description: {
        story: "複数回の実行履歴があるタスク。",
      },
    },
  },
};

// 実行履歴読み込み中
export const LoadingAttempts: Story = {
  args: {
    task: { ...baseTask, status: "SUCCEEDED" },
    attempts: [],
    isRunning: false,
    loadingAttempts: true,
  },
  parameters: {
    docs: {
      description: {
        story: "実行履歴を読み込み中。",
      },
    },
  },
};

// 長いタイトル
export const LongTitle: Story = {
  args: {
    task: {
      ...baseTask,
      title:
        "非常に長いタスクタイトルです。これは省略表示やワードラップのテストのために使用しています。タスク管理システムでは様々な長さのタイトルを扱う必要があります。",
    },
    attempts: [],
    isRunning: false,
    loadingAttempts: false,
  },
  parameters: {
    docs: {
      description: {
        story: "長いタイトルは折り返して表示されます。",
      },
    },
  },
};

// 長いエラーメッセージ
export const LongErrorMessage: Story = {
  args: {
    task: {
      ...baseTask,
      status: "FAILED",
      startedAt: "2024-12-05T10:35:00Z",
      doneAt: "2024-12-05T10:38:00Z",
    },
    attempts: [
      {
        id: "attempt-1",
        taskId: "task-12345",
        status: "FAILED",
        startedAt: "2024-12-05T10:35:00Z",
        finishedAt: "2024-12-05T10:38:00Z",
        errorSummary: `Error: Test suite failed to run

    TypeError: Cannot read properties of undefined (reading 'map')

      at Object.<anonymous> (/workspace/project/src/services/auth.ts:45:12)
      at Module._compile (node:internal/modules/cjs/loader:1105:14)
      at Object.Module._extensions..js (node:internal/modules/cjs/loader:1159:10)
      at Module.load (node:internal/modules/cjs/loader:981:32)
      at Function.Module._load (node:internal/modules/cjs/loader:822:12)

    Test Suites: 1 failed, 1 total
    Tests:       3 failed, 7 passed, 10 total
    Snapshots:   0 total
    Time:        4.523 s`,
      },
    ],
    isRunning: false,
    loadingAttempts: false,
  },
  parameters: {
    docs: {
      description: {
        story: "長いエラーメッセージはスクロール可能なエリアに表示されます。",
      },
    },
  },
};
