<script lang="ts">
  import BrandLinkList from '$components/BrandLinkList.svelte';
  import type { PageData } from './$types';

  export let data: PageData;
</script>

<svelte:head>
  <title>{data.year} Pies by Brand :: Mince Pie Rank</title>
  <meta name="description" content={`All brands we loaded data for in ${data.year}`} />
  <meta name="keywords" content={`Mince Pies, UK, Ranking, ${data.year}`} />
</svelte:head>

<BrandLinkList makers={data.makers} year={data.year} imgprssr={data.imgprssr} />
