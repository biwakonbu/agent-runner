<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  // @ts-ignore
  import { SelectWorkspace } from '../../wailsjs/go/main/App';

  const dispatch = createEventDispatcher();

  async function select() {
    try {
      const id = await SelectWorkspace();
      if (id) {
        dispatch('selected', id);
      }
    } catch (e) {
      console.error(e);
    }
  }
</script>

<div class="container">
  <h1>Multiverse IDE</h1>
  <p>Select a project to get started.</p>
  <button on:click={select}>Open Workspace</button>
</div>

<style>
  .container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100vh;
    color: white;
  }
  button {
    padding: 10px 20px;
    font-size: 1.2em;
    cursor: pointer;
  }
</style>
