<script lang="ts">
  import FilterablePieList from '$components/FilterablePieList.svelte';
  import type { PageData } from './$types';

  export let data: PageData;
</script>

<h1>Your Rankings</h1>

<FilterablePieList
  pies={[...data.pies, ...data.customPies]}
  pieRankings={[...data.rankings, ...data.customRankings]}
  imgprssr={data.imgprssr}
/>

<style>
  h1 {
    text-align: center;
  }
</style>
