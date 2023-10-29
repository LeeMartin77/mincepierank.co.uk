<script lang="ts">
  import { ppCategory } from './utilities/formatCategory';
  export let year: number;
  export let categories: string[];
</script>

<h2 style="text-align: center">Categories</h2>
<div class="category-links">
  {#each categories as category}
    <a class="category-link" href={`/years/${year}/categories/${category}`}
      >{ppCategory(category)}</a
    >
  {/each}
  {#if categories.length === 0}
    <span>No Categories</span>
  {/if}
</div>

<style>
  .category-links {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    gap: 0.5em;
    margin-top: 0.5em;
    margin-bottom: 0.5em;
  }

  .category-link {
    display: block;
    box-shadow: 0em 0em 0.3em rgba(0, 0, 0, 0.1);
    background: rgba(255, 255, 255, 0.2);
    backdrop-filter: blur(5px);
    -webkit-backdrop-filter: blur(5px);
    border: 1px solid rgba(255, 255, 255, 0.3);
    font-size: 1em;
    font-weight: bold;
    border-radius: 0.5em;
    padding: 0.5em 1em;
    text-align: center;
    color: #111;
  }
  .category-link:visited {
    color: #111;
  }
</style>
