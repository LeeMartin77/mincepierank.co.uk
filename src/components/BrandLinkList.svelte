<script lang="ts">
	import type { Maker } from "$lib/storage";
  import BrandLinkCard from "./BrandLinkCard.svelte";

  export let makers: Maker[];
  export let year: number;
</script>

<h2>Brands</h2>
<div class="brand-list">
{#each makers as maker}
  <BrandLinkCard maker={maker} year={year} />
{/each}
</div>

<style>
  h2 {
    text-align: center;
  }
  .brand-list {
    display: flex;
    flex-direction: column;
    gap: 1em;
  }
</style>