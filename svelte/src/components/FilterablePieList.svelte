<script lang="ts">
	import type { MakerPie, PieRankingSummary } from "$lib/storage";
	import PieLinkCard from "./PieLinkCard.svelte";
	import { ppCategory } from "./utilities/formatCategory";
	import { mapPiesAndRankings, type MappedPiesAndRankings } from "./utilities/mapPiesAndRankings";
  export let pies: MakerPie[];
  export let pieRankings: PieRankingSummary[];
  export let fixedCategory: string | undefined = undefined;
  $: mappedPiesAndRankings = mapPiesAndRankings(pies, pieRankings)
</script>

{#if fixedCategory}
  <h3>{ppCategory(fixedCategory)}</h3>
{/if}

{#each pies as pie}
  <PieLinkCard pie={pie} />
{/each}