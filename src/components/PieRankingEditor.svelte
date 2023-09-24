<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  const dispatcher = createEventDispatcher();

  type HotRanking = {
    pastry?: number;
    filling?: number;
    topping?: number;
    looks?: number;
    value?: number;
    notes?: string;
  };
  export let hotRanking: HotRanking;

  const handleRankingChange = (hotRanking: HotRanking) => {
    dispatcher('change', hotRanking);
  };
</script>

<div>
  <h3>My Ranking</h3>
  <label for="pastry">Pastry</label>
  <input
    name="pastry"
    type="number"
    min="1"
    max="5"
    bind:value={hotRanking.pastry}
    on:change={() => handleRankingChange(hotRanking)}
  />
  <label for="filling">Filling</label>
  <input
    name="filling"
    type="number"
    min="1"
    max="5"
    bind:value={hotRanking.filling}
    on:change={() => handleRankingChange(hotRanking)}
  />
  <label for="topping">Topping</label>
  <input
    name="topping"
    type="number"
    min="1"
    max="5"
    bind:value={hotRanking.topping}
    on:change={() => handleRankingChange(hotRanking)}
  />
  <label for="looks">Looks</label>
  <input
    name="looks"
    type="number"
    min="1"
    max="5"
    bind:value={hotRanking.looks}
    on:change={() => handleRankingChange(hotRanking)}
  />
  <label for="value">Value</label>
  <input
    name="value"
    type="number"
    min="1"
    max="5"
    bind:value={hotRanking.value}
    on:change={() => handleRankingChange(hotRanking)}
  />
</div>

<style>
  label {
    font-size: 1em;
  }
</style>
