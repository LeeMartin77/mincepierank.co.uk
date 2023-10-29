<script lang="ts">
  import FilterablePieList from '$components/FilterablePieList.svelte';
  import LinkButton from '$components/generic/LinkButton.svelte';
  import type { PageData } from './$types';

  export let data: PageData;
</script>

<h1>Your Pies</h1>
<LinkButton href="/create">Create new</LinkButton>

<FilterablePieList
  pies={data.customPies}
  pieRankings={data.customRankings}
  imgprssr={data.imgprssr}
/>

<style>
  h1 {
    text-align: center;
  }
</style>
