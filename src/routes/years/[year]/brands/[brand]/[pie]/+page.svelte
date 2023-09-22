<script lang="ts">
  import { ppCategory } from '$components/utilities/formatCategory';
  import { imgprssrPrefix } from '$lib/imgprssr';
  import type { MakerPieRanking } from '$lib/storage';
  import { onMount } from 'svelte';
  import type { PageData } from './$types';

  export let data: PageData;

  let initialLoad = true;
  let submitting = false;
  let userRanking: MakerPieRanking;

  type HotRanking = Partial<Omit<MakerPieRanking, 'userid'>> &
    Omit<MakerPieRanking, 'userid' | 'pastry' | 'filling' | 'topping' | 'looks' | 'value'>;

  let hotRanking: HotRanking = {
    year: data.pie.year,
    makerid: data.pie.makerid,
    pieid: data.pie.id
  };

  onMount(() => {
    fetch(`/api/ranking?year=${data.pie.year}&makerid=${data.pie.makerid}&pieid=${data.pie.id}`)
      .then((res) => {
        if (res.status === 200) {
          return res.json();
        }
      })
      .then((rnking: MakerPieRanking) => {
        userRanking = rnking;
        hotRanking = rnking;
      })
      .finally(() => {
        initialLoad = false;
      });
  });

  const handleRankingChange = (rnkn: HotRanking) => {
    if (
      rnkn.filling === undefined ||
      rnkn.pastry === undefined ||
      rnkn.topping === undefined ||
      rnkn.value === undefined ||
      rnkn.looks === undefined
    ) {
      return;
    }
    submitting = true;
    fetch(`/api/ranking`, { method: 'POST', body: JSON.stringify(rnkn) })
      .then((res) => {
        if (res.status === 200) {
          return res.json();
        }
      })
      .then((rnking: MakerPieRanking) => {
        userRanking = rnking;
        hotRanking = rnking;
      })
      .finally(() => {
        submitting = false;
      });
  };
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

{#if initialLoad}
  <div>Loading</div>
{:else if userRanking && (data.readonly || data.activeYear !== data.pie.year)}
  <div>
    <h3>My Ranking</h3>
    <dl>
      <dt>Pastry</dt>
      <dd>{userRanking.pastry}</dd>
      <dt>Filling</dt>
      <dd>{userRanking.filling}</dd>
      <dt>Topping</dt>
      <dd>{userRanking.topping}</dd>
      <dt>Looks</dt>
      <dd>{userRanking.looks}</dd>
      <dt>value</dt>
      <dd>{userRanking.value}</dd>
    </dl>
  </div>
{:else if !data.readonly && data.activeYear === data.pie.year}
  <div>
    <h3>My Ranking</h3>
    <dl>
      <dt>Pastry</dt>
      <dd>
        <input
          type="number"
          min="0"
          max="5"
          bind:value={hotRanking.pastry}
          on:change={() => handleRankingChange(hotRanking)}
        />
      </dd>
      <dt>Filling</dt>
      <dd>
        <input
          type="number"
          min="0"
          max="5"
          bind:value={hotRanking.filling}
          on:change={() => handleRankingChange(hotRanking)}
        />
      </dd>
      <dt>Topping</dt>
      <dd>
        <input
          type="number"
          min="0"
          max="5"
          bind:value={hotRanking.topping}
          on:change={() => handleRankingChange(hotRanking)}
        />
      </dd>
      <dt>Looks</dt>
      <dd>
        <input
          type="number"
          min="0"
          max="5"
          bind:value={hotRanking.looks}
          on:change={() => handleRankingChange(hotRanking)}
        />
      </dd>
      <dt>value</dt>
      <dd>
        <input
          type="number"
          min="0"
          max="5"
          bind:value={hotRanking.value}
          on:change={() => handleRankingChange(hotRanking)}
        />
      </dd>
    </dl>
  </div>
{/if}
