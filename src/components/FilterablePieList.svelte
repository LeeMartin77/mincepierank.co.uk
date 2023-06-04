<script lang="ts">
  import type { MakerPie, PieRankingSummary } from '$lib/storage';
  import PieLinkCard from '$components/PieLinkCard.svelte';
  import { ppCategory } from './utilities/formatCategory';
  import { mapPiesAndRankings } from './utilities/mapPiesAndRankings';
  export let pies: MakerPie[];
  export let pieRankings: PieRankingSummary[];
  export let fixedCategory: string | undefined = undefined;
  export let hideUnranked: boolean = false;
  $: availableCategoryIds = Array.from(
    new Set(
      pies.reduce<string[]>((acc, curr) => {
        return [...acc, ...curr.labels.filter((lb) => lb !== fixedCategory)];
      }, [])
    )
  );
  $: filteredCategoryIds = [] as string[];
  $: mpr = mapPiesAndRankings(pies, pieRankings, new Set(filteredCategoryIds));
  $: orderedFilteredPieIds = [
    ...mpr.rankingOrder,
    ...(hideUnranked ? [] : Array.from(mpr.unrankedPies))
  ];
  const toggleFilter = (categoryId: string) => {
    const index = filteredCategoryIds.indexOf(categoryId);
    if (index !== -1) {
      filteredCategoryIds.splice(index, 1);
      filteredCategoryIds = filteredCategoryIds;
    } else {
      filteredCategoryIds = [...filteredCategoryIds, categoryId];
    }
  };
</script>

{#if fixedCategory}
  <h3>{ppCategory(fixedCategory)}</h3>
{/if}

<div class="filters">
  {#each availableCategoryIds as categoryId}
    <button
      class={filteredCategoryIds.includes(categoryId) ? 'selected' : ''}
      on:click={() => toggleFilter(categoryId)}>{ppCategory(categoryId)}</button
    >
  {/each}
</div>

<div class="pie-list">
  {#if orderedFilteredPieIds.length === 0}
    <span>No pies visibile with selected filters - try removing a filter</span>
  {/if}
  {#each orderedFilteredPieIds as pieId}
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

  .filters {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    gap: 0.5em;
    margin: 1em;
  }
  .filters button {
    cursor: pointer;
    border: 1px dotted grey;
    border-radius: 0.5em;
    background-color: transparent;
    padding: 0.5em;
  }
  .filters button:hover {
    background-color: rgba(0, 0, 0, 0.1);
  }
  .filters button.selected {
    border: 1px solid black;
  }
</style>
