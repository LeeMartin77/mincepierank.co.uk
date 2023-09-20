<script lang="ts">
  import type { Maker } from '$lib/storage';
  import BrandLinkCard from './BrandLinkCard.svelte';

  export let makers: Maker[];
  export let year: number;
  export let imgprssr: string;
</script>

<h2>Brands</h2>
<div class="brand-list">
  {#each makers as maker}
    <div class="brand-link-container">
      <BrandLinkCard {maker} {year} {imgprssr} />
    </div>
  {/each}
  {#if makers.length === 0}
    <span>No Makers</span>
  {/if}
</div>

<style>
  h2 {
    text-align: center;
  }
  .brand-list {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    justify-content: center;
    gap: 1em;
  }
  .brand-link-container {
    max-width: 320px;
  }
</style>
