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
  <div style="display: flex; flex-direction: column; align-items:center;">
    {#if rankingSummary.count !== undefined}
      <h3>Average</h3>
      <PieGraph score={rankingSummary.average} centerAlign={true} />
      <h5>{rankingSummary.count} Rankings</h5>
    {/if}
    <table>
      <tbody>
        {#each rankingprops as property}
          <tr>
            <td><b>{property.charAt(0).toUpperCase()}{property.slice(1)}</b></td>
            <td><PieGraph score={rankingSummary[property]} /></td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
{/if}

<style>
  td {
    padding: 0.5em 0.75em;
  }

  td:first-of-type {
    text-align: right;
  }
</style>
