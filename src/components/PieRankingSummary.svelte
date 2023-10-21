<script lang="ts">
  import type { PieRankingSummary } from '$lib/storage';
  import PieGraph from './subcomponents/PieGraph.svelte';

  const bgcolor =
    'filter: invert(96%) sepia(0%) saturate(508%) hue-rotate(145deg) brightness(91%) contrast(90%);';

  const rankingprops: ('pastry' | 'filling' | 'topping' | 'looks' | 'value')[] = [
    'pastry',
    'filling',
    'topping',
    'looks',
    'value'
  ];

  export let rankingSummary:
    | (Omit<PieRankingSummary, 'makerid' | 'average' | 'count'> & {
        average?: number;
        count?: number;
      })
    | undefined;
</script>

{#if rankingSummary}
  <div>
    {#if rankingSummary.count !== undefined}
      <h3>Ranking Summary</h3>
      <h4>{rankingSummary.count} Rankings | {rankingSummary.average} Average</h4>
    {/if}
    {#each rankingprops as property}
      <h5>
        {property.charAt(0).toUpperCase()}{property.slice(1)}: {rankingSummary[property].toFixed(
          2
        )}/5.00
      </h5>
      <PieGraph score={rankingSummary[property]} />
    {/each}
  </div>
{/if}
