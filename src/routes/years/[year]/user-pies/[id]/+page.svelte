<script lang="ts">
  import { ppCategory } from '$components/utilities/formatCategory';
  import { imgprssrPrefix } from '$lib/imgprssr';
  import type { PageData } from './$types';
  import UserPieRankingInterface from '$components/UserPieRankingInterface.svelte';

  export let data: PageData;
</script>

<h1>{data.pie.displayname}</h1>
<img
  src={imgprssrPrefix(data.pie.image_file + '?filter=gaussian&width=250', data.imgprssr)}
  alt={data.pie.displayname}
/>
<div>
  <h3>Details</h3>
  <dl>
    <dt>Categories</dt>
    <dd>{data.pie.labels.map(ppCategory).join(', ')}</dd>
    <dt>Pack Count</dt>
    <dd>{data.pie.pack_count}</dd>
    <dt>Price</dt>
    <dd>{data.pie.pack_price_in_pence}</dd>
    <dt>Website</dt>
    <dd><a href={data.pie.web_link}>Link to Website</a></dd>
  </dl>
</div>

{#if data.ranking}
  <div>
    <h3>Ranking Summary</h3>
    <dl>
      <dt>Pastry</dt>
      <dd>{data.ranking.pastry}</dd>
      <dt>Filling</dt>
      <dd>{data.ranking.filling}</dd>
      <dt>Topping</dt>
      <dd>{data.ranking.topping}</dd>
      <dt>Looks</dt>
      <dd>{data.ranking.looks}</dd>
      <dt>value</dt>
      <dd>{data.ranking.value}</dd>
    </dl>
  </div>
{/if}

<UserPieRankingInterface
  readonly={data.readonly || data.activeYear !== data.pie.year}
  year={data.pie.year}
  pieid={data.pie.id}
/>
