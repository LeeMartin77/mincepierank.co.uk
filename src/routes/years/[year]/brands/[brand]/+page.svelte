<script lang="ts">
  import FilterablePieList from '$components/FilterablePieList.svelte';
  import type { PageData } from './$types';

  export let data: PageData;
</script>

<svelte:head>
  <title
    >{data.makerData.name}{data.makerData.name.endsWith('s') ? "'" : "'s"} Pies :: Mince Pie Rank</title
  >
  <meta
    name="description"
    content={`${data.makerData.name}${
      data.makerData.name.endsWith('s') ? "'" : "'s"
    } pies for the year of ${data.year}`}
  />
  <meta name="keywords" content={`Mince Pies, UK, Ranking, ${data.year}, ${data.makerData.name}`} />
</svelte:head>

<h1>{data.makerData.name}</h1>

<FilterablePieList pies={data.pies} pieRankings={data.pieRankings} imgprssr={data.imgprssr} />

<style>
  h1 {
    text-align: center;
  }
</style>
