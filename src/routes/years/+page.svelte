<script lang="ts">
  import LinkButton from '$components/generic/LinkButton.svelte';
  import type { PageData } from './$types';
  export let data: PageData;
  const years = data.years.sort().reverse();
</script>

<svelte:head>
  <title>Available Years :: Mince Pie Rank</title>
  <meta name="description" content="All the years Mince Pie Rank has data for" />
  <meta name="keywords" content="Mince Pies, UK, Ranking" />
</svelte:head>

<h1 style="min-width: 240px;">Years</h1>

<div style="width: 100%; display: flex; flex-direction: column; gap: 2em;">
  {#each years as year}
    <LinkButton href={`/years/${year}`} style={'width: calc(100% - 2em); display: block;'}
      >{year}</LinkButton
    >
  {/each}
</div>

<style>
  h1 {
    text-align: center;
  }
</style>
