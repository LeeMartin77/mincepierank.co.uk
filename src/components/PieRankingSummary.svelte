<script lang="ts">
  import type { PieRankingSummary } from '$lib/storage';

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
      <div style="display: flex;">
        {#each [1, 2, 3, 4, 5] as rank}
          {#if rankingSummary[property] >= rank}
            <div><img alt="ranking icon" src="/rankicon.svg" width="24" /></div>
          {:else if rankingSummary[property] - (rank - 1) > 0}
            <div style="display: flex; width: 24px">
              <div
                style={`width: ${
                  (rankingSummary[property] - (rank - 1)) * 100
                }%; overflow: hidden;`}
              >
                <img alt="ranking icon" src="/rankicon.svg" width="24" />
              </div>
              <div
                style={`width: calc(24px - ${
                  (rankingSummary[property] - (rank - 1)) * 100
                }%); overflow: hidden; ${bgcolor}; display: flex; justify-content: flex-end; align-items: baseline;`}
              >
                <img alt="ranking icon" src="/rankicon.svg" width="24" />
              </div>
            </div>
          {:else}
            <div><img alt="ranking icon" src="/rankicon.svg" width="24" style={bgcolor} /></div>
          {/if}
        {/each}
      </div>
    {/each}
  </div>
{/if}
