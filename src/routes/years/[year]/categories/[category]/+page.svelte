<script lang="ts">
  import FilterablePieList from '$components/FilterablePieList.svelte';
  import type { PageData } from './$types';

  export let data: PageData;
</script>

<svelte:head>
  <title>{data.category[0].toUpperCase()}{data.category.slice(1)} Pies :: Mince Pie Rank</title>
</svelte:head>

<FilterablePieList
  pies={data.pies}
  pieRankings={data.pieRankings}
  fixedCategory={data.category}
  imgprssr={data.imgprssr}
/>
