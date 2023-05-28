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
  <PieLinkCard pie={data.topPie} pieListRanking={data.topPieRanking} />
{/if}

<BrandLinkList makers={data.makers} year={2022} />

<style>
  h2 {
    text-align: center;
  }
</style>
