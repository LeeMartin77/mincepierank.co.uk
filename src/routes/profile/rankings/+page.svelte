<script lang="ts">
  import LinkButton from '$components/generic/LinkButton.svelte';
  import type { PageData } from './$types';

  export let data: PageData;
</script>

<h1 style="min-width: 240px;">Years</h1>

<div style="width: 100%; display: flex; flex-direction: column; gap: 2em;">
  {#each data.years as year}
    <LinkButton
      href={`/profile/rankings/${year}`}
      style={'width: calc(100% - 2em); display: block;'}>{year}</LinkButton
    >
  {/each}
</div>

<style>
  h1 {
    text-align: center;
  }
</style>
