<script lang="ts">
  import FilterablePieList from '$components/FilterablePieList.svelte';
  import type { PageData } from './$types';

  export let data: PageData;
</script>

<h1>All Pies</h1>

<FilterablePieList pies={data.pies} pieRankings={data.pieRankings} imgprssr={data.imgprssr} />

<style>
  h1 {
    text-align: center;
  }
</style>
