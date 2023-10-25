<script lang="ts">
  import CategoryLinkList from '$components/CategoryLinkList.svelte';
  import type { PageData } from './$types';

  export let data: PageData;
</script>

<svelte:head>
  <title>{data.year} Pies by Category :: Mince Pie Rank</title>
</svelte:head>

<CategoryLinkList year={data.year} categories={data.categories} />
