<script lang="ts">
  import { ppCategory } from '$components/utilities/formatCategory';
  import { imgprssrPrefix } from '$lib/imgprssr';
  import type { PageData } from './$types';
  import PieRankingInterface from '$components/PieRankingInterface.svelte';
  import PieRankingSummary from '$components/PieRankingSummary.svelte';
  import { formatCurrency } from '$components/utilities/formatCurrency';
  import { onMount } from 'svelte';

  export let data: PageData;

  $: rankingSummary = data.ranking;
  const reloadRanking = () => {
    fetch(
      `/api/ranking/summary?year=${data.pie.year}&makerid=${data.pie.makerid}&pieid=${data.pie.id}`
    )
      .then((res) => {
        if (res.status === 200) {
          return res.json();
        }
      })
      .then((rnking) => {
        rankingSummary = rnking;
      });
  };

  onMount(() => {
    reloadRanking();
  });
</script>

<svelte:head>
  <title
    >{data.pie.displayname}{data.maker ? ` from ${data.maker.name}` : ''}{data.ranking
      ? ` :: ${data.ranking.average.toFixed(1)}/5`
      : ''} :: Mince Pie Rank</title
  >
</svelte:head>
<div class="page-wrapper">
  <div style="display:flex; flex-direction: column; gap: 0.5em;">
    <h1 style="margin-bottom: 0em;">{data.pie.displayname}</h1>
    {#if data.maker}
      <h3>from {data.maker.name}</h3>
    {/if}
    <a href={data.pie.web_link}>Link to Website</a>
    <img
      src={imgprssrPrefix(data.pie.image_file + '?filter=gaussian&width=250', data.imgprssr)}
      alt={data.pie.displayname}
    />
  </div>
  <div>
    <div class="category-links">
      {#each data.pie.labels as label}
        <a class="category-link" href={`/years/${data.pie.year}/categories/${label}`}
          >{ppCategory(label)}</a
        >
      {/each}
    </div>
    <h3>Pieformation</h3>
    <table>
      <tbody>
        <tr>
          <td><b>Pack Count</b></td>
          <td>{data.pie.pack_count}</td>
        </tr>
        <tr>
          <td><b>Price</b></td>
          <td>{formatCurrency(data.pie.pack_price_in_pence)}</td>
        </tr>
        <tr>
          <td><b>Cost per pie</b></td>
          <td>{formatCurrency(data.pie.pack_price_in_pence / data.pie.pack_count)}</td>
        </tr>
      </tbody>
    </table>
  </div>

  <div>
    <PieRankingSummary {rankingSummary} />
  </div>
  <div>
    <PieRankingInterface
      readonly={data.readonly || data.activeYear !== data.pie.year}
      year={data.pie.year}
      makerid={data.pie.makerid}
      pieid={data.pie.id}
      on:newRanking={() => {
        setTimeout(reloadRanking, 500);
      }}
    />
  </div>
</div>

<style>
  .page-wrapper {
    display: flex;
    flex-direction: row;
    justify-content: center;
    flex-wrap: wrap;
    gap: 2em;
  }
  .page-wrapper > div {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  td {
    padding: 0.25em 0.75em;
  }

  table tr td {
    border-bottom: 1px solid black;
  }

  table tr:last-of-type td {
    border: none;
  }
  .category-links {
    display: flex;
    flex-direction: row;
    justify-content: center;
    gap: 0.5em;
    margin-top: 0.5em;
    margin-bottom: 0.5em;
  }

  .category-link {
    display: block;
    box-shadow: 0em 0em 0.3em rgba(0, 0, 0, 0.1);
    background: rgba(255, 255, 255, 0.2);
    backdrop-filter: blur(5px);
    -webkit-backdrop-filter: blur(5px);
    border: 1px solid rgba(255, 255, 255, 0.3);
    font-size: 0.8em;
    border-radius: 2em;
    padding: 0.5em;
    text-decoration: none;
    text-align: center;
    color: initial;
  }
  .category-link:visited {
    color: initial;
  }
</style>
