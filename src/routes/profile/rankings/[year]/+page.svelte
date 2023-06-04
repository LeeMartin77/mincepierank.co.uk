<script lang="ts">
  import FilterablePieList from '$components/FilterablePieList.svelte';
  import type { PageData } from './$types';

  export let data: PageData;
</script>

<h1>Your Rankings</h1>

<FilterablePieList pies={data.pies} pieRankings={data.rankings} />

<style>
  h1 {
    text-align: center;
  }
</style>
