<script lang="ts">
  import { ppCategory } from "./utilities/formatCategory";
  export let year: number;
  export let categories: string[];
</script>

{#each categories as category}
  <a href={`/years/${year}/categories/${category}`}>{ppCategory(category)}</a>
{/each}
