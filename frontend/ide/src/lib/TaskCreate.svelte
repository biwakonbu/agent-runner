<script lang="ts">
  import { createEventDispatcher, onMount } from "svelte";
  import { Button, Input } from "../design-system";
  import { Logger } from "../services/logger";
  // @ts-ignore - Wails自動生成ファイル
  import { CreateTask, GetAvailablePools } from "../../wailsjs/go/main/App";

  const log = Logger.withComponent("TaskCreate");
  const dispatch = createEventDispatcher<{
    created: void;
  }>();

  // Pool の型定義
  interface Pool {
    id: string;
    name: string;
    description?: string;
  }

  // デフォルト Pool 一覧
  const defaultPools: Pool[] = [
    { id: "default", name: "Default" },
    { id: "codegen", name: "Codegen" },
    { id: "test", name: "Test" },
  ];

  let title = "";
  let poolId = "default";
  let isSubmitting = false;
  let error = "";
  let pools: Pool[] = [];
  let loadingPools = true;

  // Pool 一覧を取得
  onMount(async () => {
    try {
      const result = await GetAvailablePools();
      if (result && result.length > 0) {
        pools = result;
        log.info("loaded available pools", { count: pools.length });
      } else {
        pools = defaultPools;
        log.warn("no pools returned, using defaults");
      }
    } catch (e) {
      log.error("failed to load pools", { error: e });
      pools = defaultPools;
    } finally {
      loadingPools = false;
    }
  });

  async function handleSubmit() {
    if (!title.trim()) {
      error = "タイトルを入力してください";
      log.warn("submission blocked: title is empty");
      return;
    }

    error = "";
    isSubmitting = true;
    log.info("creating task", { title: title.trim(), poolId });

    try {
      await CreateTask(title.trim(), poolId);
      log.info("task created successfully", { title: title.trim() });
      title = "";
      dispatch("created");
    } catch (e) {
      log.error("task creation failed", { error: e });
      error = "タスクの作成に失敗しました";
    } finally {
      isSubmitting = false;
    }
  }

  function handleKeydown(event: KeyboardEvent) {
    if (event.key === "Enter" && !event.shiftKey) {
      event.preventDefault();
      handleSubmit();
    }
  }
</script>

<form class="task-create-form" on:submit|preventDefault={handleSubmit}>
  <!-- タイトル入力 -->
  <Input
    label="タイトル"
    placeholder="タスクのタイトルを入力"
    bind:value={title}
    disabled={isSubmitting}
    error={!!error}
    errorMessage={error}
    on:keydown={handleKeydown}
  />

  <!-- Pool選択 -->
  <div class="form-group">
    <label for="task-pool" class="form-label">Pool</label>
    <select
      id="task-pool"
      class="form-select"
      bind:value={poolId}
      disabled={isSubmitting || loadingPools}
    >
      {#if loadingPools}
        <option value="">読み込み中...</option>
      {:else}
        {#each pools as pool}
          <option value={pool.id}>{pool.name}</option>
        {/each}
      {/if}
    </select>
  </div>

  <!-- 送信ボタン -->
  <div class="form-actions">
    <Button
      type="submit"
      variant="primary"
      disabled={!title.trim() || loadingPools}
      loading={isSubmitting}
      loadingLabel="作成中..."
    >
      タスクを作成
    </Button>
  </div>
</form>

<style>
  .task-create-form {
    display: flex;
    flex-direction: column;
    gap: var(--mv-spacing-md);
    padding: var(--mv-spacing-md);
  }

  /* フォームグループ（セレクト用） */
  .form-group {
    display: flex;
    flex-direction: column;
    gap: var(--mv-spacing-xxs);
  }

  .form-label {
    font-size: var(--mv-font-size-sm);
    font-weight: var(--mv-font-weight-medium);
    color: var(--mv-color-text-secondary);
  }

  /* セレクト */
  .form-select {
    padding: var(--mv-spacing-sm);
    font-size: var(--mv-font-size-md);
    font-family: var(--mv-font-sans);
    color: var(--mv-color-text-primary);
    background: var(--mv-color-surface-secondary);
    border: var(--mv-border-width-thin) solid var(--mv-color-border-default);
    border-radius: var(--mv-radius-sm);
    transition:
      border-color var(--mv-transition-hover),
      box-shadow var(--mv-transition-hover);
    cursor: pointer;
    appearance: none;
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath fill='%23888888' d='M6 8L1 3h10z'/%3E%3C/svg%3E");
    background-repeat: no-repeat;
    background-position: right var(--mv-spacing-sm) center;
    padding-right: var(--mv-spacing-xl);
  }

  .form-select:focus {
    outline: none;
    border-color: var(--mv-color-border-focus);
    box-shadow: var(--mv-shadow-focus);
  }

  .form-select:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .form-select option {
    background: var(--mv-color-surface-primary);
    color: var(--mv-color-text-primary);
  }

  /* アクション */
  .form-actions {
    display: flex;
    justify-content: flex-end;
    padding-top: var(--mv-spacing-xs);
  }
</style>
