<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  // @ts-ignore
  import { CreateTask } from '../../wailsjs/go/main/App';

  const dispatch = createEventDispatcher();

  let title = '';
  let poolId = 'default';

  async function create() {
    if (!title) return;
    try {
      await CreateTask(title, poolId);
      dispatch('created');
      title = '';
    } catch (e) {
      console.error(e);
    }
  }
</script>

<div class="task-create">
  <h3>New Task</h3>
  <input type="text" bind:value={title} placeholder="Task Title" />
  <select bind:value={poolId}>
    <option value="default">Default Pool</option>
    <option value="codegen">Codegen</option>
    <option value="test">Test</option>
  </select>
  <button on:click={create}>Create</button>
</div>

<style>
  .task-create {
    padding: 10px;
    border-top: 1px solid #444;
    display: flex;
    gap: 10px;
  }
  input {
    flex: 1;
    padding: 5px;
  }
  button {
    padding: 5px 10px;
  }
</style>
