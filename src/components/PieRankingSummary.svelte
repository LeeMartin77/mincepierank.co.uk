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
      <dd>{rankingSummary.pastry}</dd>
      <dt>Filling</dt>
      <dd>{rankingSummary.filling}</dd>
      <dt>Topping</dt>
      <dd>{rankingSummary.topping}</dd>
      <dt>Looks</dt>
      <dd>{rankingSummary.looks}</dd>
      <dt>value</dt>
      <dd>{rankingSummary.value}</dd>
    </dl>
  </div>
{/if}
