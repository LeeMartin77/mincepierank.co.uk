<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import RankingInput from './subcomponents/RankingInput.svelte';
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
  <RankingInput
    value={hotRanking.pastry}
    displayName={'Pastry'}
    name={'pastry'}
    on:change={() => handleRankingChange(hotRanking)}
  />
  <RankingInput
    value={hotRanking.filling}
    displayName={'Filling'}
    name={'filling'}
    on:change={() => handleRankingChange(hotRanking)}
  />
  <RankingInput
    value={hotRanking.topping}
    displayName={'Topping'}
    name={'topping'}
    on:change={() => handleRankingChange(hotRanking)}
  />
  <RankingInput
    value={hotRanking.looks}
    displayName={'Looks'}
    name={'looks'}
    on:change={() => handleRankingChange(hotRanking)}
  />
  <RankingInput
    value={hotRanking.value}
    displayName={'Value'}
    name={'value'}
    on:change={() => handleRankingChange(hotRanking)}
  />
</div>

<style>
  label {
    font-size: 1em;
  }
</style>
