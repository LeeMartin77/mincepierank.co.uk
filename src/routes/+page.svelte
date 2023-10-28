<script lang="ts">
  import BrandLinkList from '$components/BrandLinkList.svelte';
  import PieLinkCard from '$components/PieLinkCard.svelte';
  import type { PageData } from './$types';

  export let data: PageData;
</script>

<svelte:head>
  <title>Mince Pie Rank :: Homepage</title>
</svelte:head>

{#if data.topPie && data.topPieRanking}
  <h2>Current Top Pie</h2>
  <PieLinkCard
    pie={data.topPie}
    pieListRanking={data.topPieRanking}
    imgprssr={data.imgprssr}
    raised={true}
  />
{/if}

{#if !data.readonly}
  <a class="quickrank-link" href="/quickrank">Click to rank your pie!</a>
  {#if data.latestRanking && data.latestPie}
    <h3>Latest Ranking</h3>
    <PieLinkCard
      pie={data.latestPie}
      pieListRanking={{ ...data.latestRanking, count: undefined }}
      when={data.latestRanking.last_updated}
      imgprssr={data.imgprssr}
      raised={true}
    />
  {/if}
{/if}

<BrandLinkList makers={data.makers} year={data.activeYear} imgprssr={data.imgprssr} />

<style>
  h2,
  h3 {
    text-align: center;
  }
  .quickrank-link {
    width: calc(100% - 2em);
    padding: 1em;
    display: block;
    text-align: center;
    box-shadow: 0em 0em 0.3em rgba(0, 0, 0, 0.2);
    color: #fff;
    background-color: rgba(0, 0, 0, 0.8);
    border-radius: 1em;
    margin-top: 2em;
    margin-bottom: 2em;
    text-decoration: underline;
    font-weight: bold;
  }

  .quickrank-link:visited {
    color: #fff;
  }
</style>
