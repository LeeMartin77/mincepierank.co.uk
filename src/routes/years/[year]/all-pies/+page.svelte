<script lang="ts">
  import FilterablePieList from '$components/FilterablePieList.svelte';
  import type { PageData } from './$types';

  export let data: PageData;
</script>

<svelte:head>
  <title>All {data.year} Pies :: Mince Pie Rank</title>
  <meta name="description" content={`All the pies from major brands for ${data.year}`} />
  <meta name="keywords" content={`Mince Pies, UK, Ranking, ${data.year}`} />
</svelte:head>

<h1>All Pies</h1>

<FilterablePieList pies={data.pies} pieRankings={data.pieRankings} imgprssr={data.imgprssr} />

<style>
  h1 {
    text-align: center;
  }
</style>
