<script lang="ts">
	import { ppCategory } from "$components/utilities/formatCategory";
	import type { PageData } from "./$types";

  export let data: PageData;
</script>

<h2>Categories index</h2>

{#each data.categories as category}
  <a href={`/years/${data.year}/categories/${category}`}>{ppCategory(category)}</a>
{/each}