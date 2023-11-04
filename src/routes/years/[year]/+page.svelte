<script lang="ts">
  import CategoryLinkList from '$components/CategoryLinkList.svelte';

  import type { PageData } from './$types';
  import BrandLinkList from '$components/BrandLinkList.svelte';
  import LinkButton from '$components/generic/LinkButton.svelte';

  export let data: PageData;
</script>

<svelte:head>
  <title>{data.year} Pies :: Mince Pie Rank</title>
  <meta
    name="description"
    content={`Mince pie categories and brands we have data for ${data.year}`}
  />
  <meta name="keywords" content={`Mince Pies, UK, Ranking, ${data.year}`} />
</svelte:head>

<LinkButton href={`/years/${data.year}/all-pies`}>View All Pies</LinkButton>

<CategoryLinkList year={data.year} categories={data.categories} />
<BrandLinkList makers={data.makers} year={data.year} imgprssr={data.imgprssr} />

<style>
</style>
