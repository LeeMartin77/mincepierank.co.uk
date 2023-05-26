<script lang="ts">
	import FilterablePieList from '$components/FilterablePieList.svelte';
	import type { PageData } from './$types';

	export let data: PageData;
</script>

<FilterablePieList 
	pies={data.pies}
	pieRankings={data.pieRankings}
	fixedCategory={data.category}/>