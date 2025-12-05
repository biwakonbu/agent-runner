<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  export let tasks: any[] = [];

  const dispatch = createEventDispatcher();

  function select(task: any) {
    dispatch('select', task);
  }
</script>

<div class="task-list">
  <h2>Tasks</h2>
  <ul>
    {#each tasks as task}
      <li class:running={task.status === 'RUNNING'}>
        <button type="button" on:click={() => select(task)}>
          <span class="status {task.status}">{task.status}</span>
          <span class="title">{task.title}</span>
          <span class="pool">{task.poolId}</span>
        </button>
      </li>
    {/each}
  </ul>
</div>

<style>
  .task-list {
    width: 300px;
    border-right: 1px solid #444;
    overflow-y: auto;
    background: #1e1e1e;
  }
  ul {
    list-style: none;
    padding: 0;
    margin: 0;
  }
  li {
    border-bottom: 1px solid #333;
  }
  li button {
    width: 100%;
    padding: 10px;
    background: transparent;
    border: none;
    cursor: pointer;
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    color: inherit;
    font: inherit;
    text-align: left;
  }
  li button:hover {
    background: #2a2a2a;
  }
  li button:focus {
    outline: 2px solid #4caf50;
    outline-offset: -2px;
  }
  .status {
    font-size: 0.8em;
    font-weight: bold;
    margin-bottom: 4px;
  }
  .status.RUNNING { color: #4caf50; }
  .status.FAILED { color: #f44336; }
  .status.PENDING { color: #ff9800; }
  .title {
    font-weight: bold;
  }
  .pool {
    font-size: 0.8em;
    color: #888;
  }
</style>
