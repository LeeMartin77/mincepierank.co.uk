<script lang="ts">
  import Card from '$components/generic/Card.svelte';
  import { imgprssrPrefix } from '$lib/imgprssr';
  import type { PageData } from './$types';
  import UserPieRankingInterface from '$components/UserPieRankingInterface.svelte';
  import PieRankingSummary from '$components/PieRankingSummary.svelte';
  import { onMount } from 'svelte';
  import { formatCurrency } from '$components/utilities/formatCurrency';

  export let data: PageData;

  $: rankingSummary = data.ranking;

  const reloadRanking = () => {
    return fetch(`/api/user-ranking/summary?year=${data.pie.year}&pieid=${data.pie.id}`).then(
      (res) => {
        if (res.status === 200) {
          return res.json();
        }
      }
    );
  };

  onMount(async () => {
    rankingSummary = await reloadRanking();
  });
</script>

{#if data.isMine}
  <Card>
    <div class="my-pie-buttons">
      <form method="GET" action={`/create`}>
        <input type="hidden" name="id" value={data.pie.id} />
        <input type="hidden" name="year" value={data.pie.year} />
        <button>Edit</button>
      </form>
      <form method="POST" action="?/delete">
        <input type="hidden" name="id" value={data.pie.id} />
        <input type="hidden" name="year" value={data.pie.year} />
        <button>Delete</button>
      </form>
    </div>
  </Card>
{/if}
<div class="page-wrapper">
  <div>
    <h1 style="margin-bottom: 0em;">{data.pie.displayname}</h1>
    {#if data.pie.maker}
      <h3>from {data.pie.maker}</h3>
    {/if}
    {#if data.pie.web_link}
      <a href={data.pie.web_link}>{data.pie.web_link}</a>
    {/if}
    <img
      src={imgprssrPrefix(data.pie.image_file + '?filter=gaussian&width=250', data.imgprssr)}
      alt={data.pie.displayname}
    />
  </div>
  {#if data.pie.location || data.pie.labels || data.pie.pack_count || data.pie.pack_price_in_pence}
    <div>
      <h3>Pieformation</h3>
      <table>
        <tbody>
          {#if data.pie.location}
            <tr>
              <td><b>Where</b></td>
              <td>{data.pie.location}</td>
            </tr>
          {/if}
          {#if data.pie.labels}
            <tr>
              <td><b>Categories</b></td>
              <td>{data.pie.labels.join(', ')}</td>
            </tr>
          {/if}
          {#if data.pie.pack_count}
            <tr>
              <td><b>Pack Count</b></td>
              <td>{data.pie.pack_count}</td>
            </tr>
          {/if}
          {#if data.pie.pack_price_in_pence}
            <tr>
              <td><b>Price</b></td>
              <td>{formatCurrency(data.pie.pack_price_in_pence)}</td>
            </tr>
          {/if}
          {#if data.pie.pack_price_in_pence && data.pie.pack_count}
            <tr>
              <td><b>Cost per pie</b></td>
              <td>{formatCurrency(data.pie.pack_price_in_pence / data.pie.pack_count)}</td>
            </tr>
          {/if}
        </tbody>
      </table>
    </div>
  {/if}
  <div>
    <PieRankingSummary {rankingSummary} />
  </div>
  <div>
    <UserPieRankingInterface
      readonly={data.readonly || data.activeYear !== data.pie.year}
      year={data.pie.year}
      pieid={data.pie.id}
      on:newRanking={async () => {
        setTimeout(async () => {
          rankingSummary = await reloadRanking();
        }, 500);
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

  .my-pie-buttons {
    display: flex;
    gap: 1em;
    justify-content: center;
  }

  .my-pie-buttons button {
    display: block;
    width: 10em;
    box-shadow: 0em 0em 0.3em rgba(0, 0, 0, 0.2);
    background: rgba(255, 255, 255, 0.2);
    backdrop-filter: blur(5px);
    -webkit-backdrop-filter: blur(5px);
    border: 1px solid rgba(255, 255, 255, 0.3);
    border-radius: 1em;
    padding: 1em;
    text-align: center;
    font-weight: bold;
    color: #111;
  }
</style>
