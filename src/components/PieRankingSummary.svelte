<script lang="ts">
  import type { PieRankingSummary } from '$lib/storage';

  export let rankingSummary: Omit<PieRankingSummary, 'makerid'> | undefined;
</script>

{#if rankingSummary}
  <div>
    <h3>Ranking Summary</h3>
    <h4>{rankingSummary.count} Rankings</h4>
    <dl>
      <dt>Pastry</dt>
      <dd>{rankingSummary.pastry.toFixed(2)}</dd>
      <dt>Filling</dt>
      <dd>{rankingSummary.filling.toFixed(2)}</dd>
      <dt>Topping</dt>
      <dd>{rankingSummary.topping.toFixed(2)}</dd>
      <dt>Looks</dt>
      <dd>{rankingSummary.looks.toFixed(2)}</dd>
      <dt>value</dt>
      <dd>{rankingSummary.value.toFixed(2)}</dd>
    </dl>
  </div>
{/if}
