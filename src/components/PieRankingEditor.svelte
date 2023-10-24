<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import RankingInput from './subcomponents/RankingInput.svelte';
  import Card from './generic/Card.svelte';
  import LoadingIndicator from './generic/LoadingIndicator.svelte';
  const dispatcher = createEventDispatcher();

  type HotRanking = {
    pastry?: number;
    filling?: number;
    topping?: number;
    looks?: number;
    value?: number;
  };
  export let hotRanking: HotRanking;
  export let submitting: boolean;
  export let submitted = false;

  const debounce = (func: any, timeout = 500) => {
    let timer: any;
    return (...args: any[]) => {
      clearTimeout(timer);
      timer = setTimeout(() => {
        // eslint-disable-next-line
        // @ts-ignore
        func.apply(this, args);
      }, timeout);
    };
  };

  const handleRankingChange = (hotRanking: HotRanking) => {
    debounce(dispatcher('change', hotRanking));
  };

  const required: (keyof HotRanking)[] = ['pastry', 'filling', 'topping', 'looks', 'value'];

  $: showHint =
    !required.every((r) => hotRanking[r] === undefined) &&
    !required.every((r) => hotRanking[r] !== undefined);
</script>

<Card>
  <div class="wrapper">
    <div class="loading-wrapper">
      {#if submitting}
        <LoadingIndicator width={18} />
      {/if}
      {#if showHint}
        <span>Select all scores to save!</span>
      {:else if submitted}
        <span>Submitted!</span>
      {/if}
    </div>
    <h3 style="text-align: center;">My Ranking</h3>
    <div class="ranking-column">
      <RankingInput
        displayName={'Pastry'}
        name={'pastry'}
        bind:value={hotRanking.pastry}
        on:change={() => handleRankingChange(hotRanking)}
      />
      <RankingInput
        displayName={'Filling'}
        name={'filling'}
        bind:value={hotRanking.filling}
        on:change={() => handleRankingChange(hotRanking)}
      />
      <RankingInput
        displayName={'Topping'}
        name={'topping'}
        bind:value={hotRanking.topping}
        on:change={() => handleRankingChange(hotRanking)}
      />
      <RankingInput
        displayName={'Looks'}
        name={'looks'}
        bind:value={hotRanking.looks}
        on:change={() => handleRankingChange(hotRanking)}
      />
      <RankingInput
        displayName={'Value'}
        name={'value'}
        bind:value={hotRanking.value}
        on:change={() => handleRankingChange(hotRanking)}
      />
    </div>
  </div>
</Card>

<style>
  .wrapper {
    position: relative;
  }
  .loading-wrapper {
    position: absolute;
  }
  .loading-wrapper > span {
    font-size: 0.7em;
  }
  .ranking-column {
    display: flex;
    flex-direction: column;
    gap: 1em;
  }
</style>
