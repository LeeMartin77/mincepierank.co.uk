<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  export let name: string;
  export let displayName: string;
  export let value: number | undefined;

  const dispatch = createEventDispatcher();
  const change = () => dispatch('change');
</script>

<label for={name}>{displayName}: {value ?? 0}/5</label>
<div class="ranking-wrapper">
  {#each [1, 2, 3, 4, 5] as ranking}
    <button
      aria-label={`Rank ${displayName} as ${ranking}/5`}
      class={`base-ranking-button-class ${
        value && value >= ranking ? 'ranking-positive' : 'ranking-negative'
      }`}
      on:click={() => {
        value = ranking;
        change();
      }}
    >
      {ranking}
    </button>
  {/each}
</div>

<style>
  .base-ranking-button-class {
    border: none;
    background-color: transparent;
    padding: 1em;
    cursor: pointer;
    border-radius: 0.5em;
  }
  .ranking-positive {
    background-color: rgb(211, 155, 0);
  }
  .ranking-negative {
    background-color: rgba(0, 0, 0, 0.1);
  }

  .ranking-wrapper {
    display: flex;
    flex-direction: row;
    width: 100%;
    gap: 0.5em;
    justify-content: center;
  }
</style>
