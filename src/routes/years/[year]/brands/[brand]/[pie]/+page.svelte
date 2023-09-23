<script lang="ts">
  import { ppCategory } from '$components/utilities/formatCategory';
  import { imgprssrPrefix } from '$lib/imgprssr';
  import type { PageData } from './$types';
  import PieRankingInterface from '$components/PieRankingInterface.svelte';
  import PieRankingSummary from '$components/PieRankingSummary.svelte';

  export let data: PageData;

  $: rankingSummary = data.ranking;
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

<PieRankingSummary {rankingSummary} />

<PieRankingInterface
  readonly={data.readonly || data.activeYear !== data.pie.year}
  year={data.pie.year}
  makerid={data.pie.makerid}
  pieid={data.pie.id}
  on:newRanking={({ detail }) => {
    fetch(
      `/api/ranking/summary?year=${detail.year}&makerid=${detail.makerid}&pieid=${detail.pieid}`
    )
      .then((res) => {
        if (res.status === 200) {
          return res.json();
        }
      })
      .then((rnking) => {
        rankingSummary = rnking;
      });
  }}
/>
