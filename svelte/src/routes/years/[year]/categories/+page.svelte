<script lang="ts">
	import CategoryLinkList from "$components/CategoryLinkList.svelte";
	import FilterablePieList from "$components/FilterablePieList.svelte";
	import { ppCategory } from "$components/utilities/formatCategory";
	import type { PageData } from "./$types";

  export let data: PageData;
</script>

<h2>Categories</h2>

<CategoryLinkList year={data.year} categories={data.categories} />

<FilterablePieList 
	pies={data.pies}
	pieRankings={data.pieRankings}/>