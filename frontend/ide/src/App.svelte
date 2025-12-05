<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import WorkspaceSelector from './lib/WorkspaceSelector.svelte';
  import TaskList from './lib/TaskList.svelte';
  import TaskDetail from './lib/TaskDetail.svelte';
  import TaskCreate from './lib/TaskCreate.svelte';
  // @ts-ignore
  import { ListTasks, GetWorkspace } from '../wailsjs/go/main/App';

  let workspaceId: string | null = null;
  let tasks: any[] = [];
  let selectedTask: any | null = null;
  let interval: any;

  async function loadTasks() {
    if (!workspaceId) return;
    try {
      tasks = await ListTasks();
      // Update selected task if it exists
      if (selectedTask) {
        const updated = tasks.find(t => t.id === selectedTask.id);
        if (updated) selectedTask = updated;
      }
    } catch (e) {
      console.error(e);
    }
  }

  function onWorkspaceSelected(event: CustomEvent<string>) {
    workspaceId = event.detail;
    loadTasks();
    interval = setInterval(loadTasks, 2000);
  }

  function onTaskSelect(event: CustomEvent<any>) {
    selectedTask = event.detail;
  }

  function onTaskCreated() {
    loadTasks();
  }

  onDestroy(() => {
    if (interval) clearInterval(interval);
  });
</script>

<main>
  {#if !workspaceId}
    <WorkspaceSelector on:selected={onWorkspaceSelected} />
  {:else}
    <div class="layout">
      <div class="sidebar">
        <TaskList {tasks} on:select={onTaskSelect} />
        <TaskCreate on:created={onTaskCreated} />
      </div>
      <TaskDetail task={selectedTask} />
    </div>
  {/if}
</main>

<style>
  main {
    height: 100vh;
    display: flex;
    flex-direction: column;
    background: #1e1e1e;
    color: #eee;
    font-family: sans-serif;
  }
  .layout {
    display: flex;
    flex: 1;
    overflow: hidden;
  }
  .sidebar {
    display: flex;
    flex-direction: column;
    border-right: 1px solid #444;
  }
</style>
