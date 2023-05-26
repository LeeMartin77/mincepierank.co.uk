<script lang="ts">
	import type { MakerPie, PieRankingSummary } from "$lib/storage";
	import PieLinkCard from "./PieLinkCard.svelte";
	import { ppCategory } from "./utilities/formatCategory";
	import { mapPiesAndRankings, type MappedPiesAndRankings } from "./utilities/mapPiesAndRankings";
  export let pies: MakerPie[];
  export let pieRankings: PieRankingSummary[];
  export let fixedCategory: string | undefined = undefined;
  $: mpr = mapPiesAndRankings(pies, pieRankings)
</script>

{#if fixedCategory}
  <h3>{ppCategory(fixedCategory)}</h3>
{/if}

<div class="pie-list">
{#each mpr.rankingOrder as pieId}
  <PieLinkCard 
    pie={mpr.mappedPies[pieId]} 
    pieListRanking={mpr.mappedRankings[pieId]} />
{/each}
{#each Array.from(mpr.unrankedPies) as pieId}
  <PieLinkCard 
    pie={mpr.mappedPies[pieId]} 
    pieListRanking={mpr.mappedRankings[pieId]} />
{/each}
</div>

<style>
  .pie-list {
    display: flex;
    flex-direction: column;
    gap: 1em;
  }
</style>