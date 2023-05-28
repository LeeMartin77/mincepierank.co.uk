<script lang="ts">
  import { ppCategory } from '$components/utilities/formatCategory';
  import type { PageData } from './$types';

  export let data: PageData;
  $: pie = data.pie;
  $: ranking = data.ranking;
</script>

<h1>{pie.displayname}</h1>
<img src={pie.image_file + '?filter=gaussian&width=250'} alt={pie.displayname} />
<div>
  <h3>Details</h3>
  <dl>
    <dt>Categories</dt>
    <dd>{pie.labels.map(ppCategory).join(', ')}</dd>
    <dt>Pack Count</dt>
    <dd>{pie.pack_count}</dd>
    <dt>Price</dt>
    <dd>{pie.pack_price_in_pence}</dd>
    <dt>Website</dt>
    <dd><a href={pie.web_link}>Link to Website</a></dd>
  </dl>
</div>

{#if ranking}
  <div>
    <h3>Ranking Summary</h3>
    <dl>
      <dt>Pastry</dt>
      <dd>{ranking.pastry}</dd>
      <dt>Filling</dt>
      <dd>{ranking.filling}</dd>
      <dt>Topping</dt>
      <dd>{ranking.topping}</dd>
      <dt>Looks</dt>
      <dd>{ranking.looks}</dd>
      <dt>value</dt>
      <dd>{ranking.value}</dd>
    </dl>
  </div>
{/if}
