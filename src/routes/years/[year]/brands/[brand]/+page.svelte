<script lang="ts">
  import FilterablePieList from '$components/FilterablePieList.svelte';
  import type { PageData } from './$types';

  export let data: PageData;
</script>

<h1>{data.makerData.name}</h1>

<FilterablePieList pies={data.pies} pieRankings={data.pieRankings} />

<style>
  h1 {
    text-align: center;
  }
</style>
