<script lang="ts">
  // @ts-ignore
  import { RunTask } from '../../wailsjs/go/main/App';

  export let task: any | null = null;

  async function run() {
    if (!task) return;
    try {
      await RunTask(task.id);
    } catch (e) {
      console.error(e);
    }
  }
</script>

<div class="task-detail">
  {#if task}
    <header>
      <h1>{task.title}</h1>
      <div class="actions">
        <button on:click={run} disabled={task.status === 'RUNNING'}>
          {task.status === 'RUNNING' ? 'Running...' : 'Run'}
        </button>
      </div>
    </header>
    <div class="info">
      <p><strong>ID:</strong> {task.id}</p>
      <p><strong>Status:</strong> {task.status}</p>
      <p><strong>Pool:</strong> {task.poolId}</p>
      <p><strong>Created:</strong> {new Date(task.createdAt).toLocaleString()}</p>
    </div>
    <!-- Attempts list would go here -->
  {:else}
    <div class="empty">Select a task to view details</div>
  {/if}
</div>

<style>
  .task-detail {
    flex: 1;
    padding: 20px;
    overflow-y: auto;
  }
  header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid #444;
    padding-bottom: 10px;
    margin-bottom: 20px;
  }
  .empty {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    color: #888;
  }
  button {
    padding: 8px 16px;
    background: #2196f3;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  button:disabled {
    background: #555;
    cursor: not-allowed;
  }
</style>
