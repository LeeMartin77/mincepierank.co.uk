<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  export let name: string;
  export let displayName: string;
  export let value: number | undefined;

  const dispatch = createEventDispatcher();
  const change = () => dispatch('change');
</script>

<div class="ranking-wrapper">
  <label for={name}>{displayName}</label>
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
</div>

<style>
  .ranking-wrapper {
    display: flex;
    align-items: center;
  }
  .ranking-wrapper > label {
    margin: 0;
    font-size: 1em;
  }
  .base-ranking-button-class {
    border: none;
    background-color: transparent;
    padding: 0.5em;
    cursor: pointer;
    border-radius: 0.5em;
    display: block;
    width: 2.5em;
    height: 2.5em;
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
