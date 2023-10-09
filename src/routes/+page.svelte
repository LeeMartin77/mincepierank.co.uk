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
  <PieLinkCard pie={data.topPie} pieListRanking={data.topPieRanking} imgprssr={data.imgprssr} />
{/if}

<a class="quickrank-link" href="/quickrank">Rank a pie now!</a>

<BrandLinkList makers={data.makers} year={data.activeYear} imgprssr={data.imgprssr} />

<style>
  h2 {
    text-align: center;
  }
  .quickrank-link {
    width: calc(100% - 2em);
    padding: 1em;
    display: block;
    text-align: center;
    border: 1px solid #111;
    border-radius: 1em;
    margin-top: 2em;
    margin-bottom: 2em;
    text-decoration: none;
    font-weight: bold;
  }
</style>
