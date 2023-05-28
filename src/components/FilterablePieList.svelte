<script lang="ts">
  import type { MakerPie, PieRankingSummary } from '$lib/storage';
  import PieLinkCard from '$components/PieLinkCard.svelte';
  import { ppCategory } from './utilities/formatCategory';
  import { mapPiesAndRankings } from './utilities/mapPiesAndRankings';
  export let pies: MakerPie[];
  export let pieRankings: PieRankingSummary[];
  export let fixedCategory: string | undefined = undefined;
  $: mpr = mapPiesAndRankings(pies, pieRankings);
</script>

{#if fixedCategory}
  <h3>{ppCategory(fixedCategory)}</h3>
{/if}

<div class="pie-list">
  {#each [...mpr.rankingOrder, ...Array.from(mpr.unrankedPies)] as pieId}
    <div class="pie-link-card-container">
      <PieLinkCard pie={mpr.mappedPies[pieId]} pieListRanking={mpr.mappedRankings[pieId]} />
    </div>
  {/each}
</div>

<style>
  h3 {
    text-align: center;
  }
  .pie-list {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    width: 100%;
    gap: 1em;
  }
  .pie-link-card-container {
    max-width: 320px;
  }
</style>
