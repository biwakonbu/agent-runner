import type { Meta, StoryObj } from "@storybook/svelte";
import TaskCreatePreview from "./TaskCreatePreview.svelte";

const meta = {
  title: "IDE/TaskCreate",
  component: TaskCreatePreview,
  tags: ["autodocs"],
  argTypes: {
    pools: {
      control: "object",
      description: "利用可能なPool一覧",
    },
    loadingPools: {
      control: "boolean",
      description: "Pool読み込み中フラグ",
    },
    isSubmitting: {
      control: "boolean",
      description: "送信中フラグ",
    },
    error: {
      control: "text",
      description: "エラーメッセージ",
    },
    initialTitle: {
      control: "text",
      description: "初期タイトル値",
    },
    initialPoolId: {
      control: "text",
      description: "初期Pool ID",
    },
  },
  parameters: {
    layout: "centered",
    docs: {
      description: {
        component:
          "新規タスク作成フォーム。タイトル入力とPool選択ができます。",
      },
    },
  },
} satisfies Meta<TaskCreatePreview>;

export default meta;
type Story = StoryObj<typeof meta>;

// デフォルトPool一覧
const defaultPools = [
  { id: "default", name: "Default" },
  { id: "codegen", name: "Codegen" },
  { id: "test", name: "Test" },
];

// カスタムPool一覧
const customPools = [
  { id: "frontend", name: "Frontend", description: "フロントエンド開発" },
  { id: "backend", name: "Backend", description: "バックエンド開発" },
  { id: "infra", name: "Infrastructure", description: "インフラ関連" },
  { id: "docs", name: "Documentation", description: "ドキュメント作成" },
];

// デフォルト状態
export const Default: Story = {
  args: {
    pools: defaultPools,
    loadingPools: false,
    isSubmitting: false,
    error: "",
    initialTitle: "",
    initialPoolId: "default",
  },
};

// タイトル入力済み
export const WithTitle: Story = {
  args: {
    pools: defaultPools,
    loadingPools: false,
    isSubmitting: false,
    error: "",
    initialTitle: "ユーザー認証機能の実装",
    initialPoolId: "codegen",
  },
  parameters: {
    docs: {
      description: {
        story: "タイトルが入力された状態。",
      },
    },
  },
};

// カスタムPool一覧
export const WithCustomPools: Story = {
  args: {
    pools: customPools,
    loadingPools: false,
    isSubmitting: false,
    error: "",
    initialTitle: "",
    initialPoolId: "frontend",
  },
  parameters: {
    docs: {
      description: {
        story: "カスタムPool一覧が設定された状態。",
      },
    },
  },
};

// Pool読み込み中
export const LoadingPools: Story = {
  args: {
    pools: [],
    loadingPools: true,
    isSubmitting: false,
    error: "",
    initialTitle: "",
    initialPoolId: "",
  },
  parameters: {
    docs: {
      description: {
        story: "Pool一覧を読み込み中。セレクトは無効化されています。",
      },
    },
  },
};

// 送信中
export const Submitting: Story = {
  args: {
    pools: defaultPools,
    loadingPools: false,
    isSubmitting: true,
    error: "",
    initialTitle: "新機能の実装",
    initialPoolId: "codegen",
  },
  parameters: {
    docs: {
      description: {
        story: "タスク作成を送信中。ボタンにローディング表示。",
      },
    },
  },
};

// エラー表示
export const WithError: Story = {
  args: {
    pools: defaultPools,
    loadingPools: false,
    isSubmitting: false,
    error: "タイトルを入力してください",
    initialTitle: "",
    initialPoolId: "default",
  },
  parameters: {
    docs: {
      description: {
        story: "バリデーションエラーが表示された状態。",
      },
    },
  },
};

// 長いタイトル
export const LongTitle: Story = {
  args: {
    pools: defaultPools,
    loadingPools: false,
    isSubmitting: false,
    error: "",
    initialTitle:
      "非常に長いタスクタイトルです。これはフォーム入力フィールドの挙動をテストするために使用しています。長いテキストがどのように表示されるか確認できます。",
    initialPoolId: "codegen",
  },
  parameters: {
    docs: {
      description: {
        story: "長いタイトルが入力された状態。",
      },
    },
  },
};

// 多数のPool
export const ManyPools: Story = {
  args: {
    pools: [
      { id: "pool-1", name: "Pool 1" },
      { id: "pool-2", name: "Pool 2" },
      { id: "pool-3", name: "Pool 3" },
      { id: "pool-4", name: "Pool 4" },
      { id: "pool-5", name: "Pool 5" },
      { id: "pool-6", name: "Pool 6" },
      { id: "pool-7", name: "Pool 7" },
      { id: "pool-8", name: "Pool 8" },
    ],
    loadingPools: false,
    isSubmitting: false,
    error: "",
    initialTitle: "",
    initialPoolId: "pool-1",
  },
  parameters: {
    docs: {
      description: {
        story: "多数のPoolが存在する場合。",
      },
    },
  },
};

// Poolなし
export const NoPools: Story = {
  args: {
    pools: [],
    loadingPools: false,
    isSubmitting: false,
    error: "",
    initialTitle: "タスク名",
    initialPoolId: "",
  },
  parameters: {
    docs: {
      description: {
        story: "Poolが存在しない状態。",
      },
    },
  },
};
