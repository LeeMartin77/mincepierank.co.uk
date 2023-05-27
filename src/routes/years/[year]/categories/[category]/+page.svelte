<script lang="ts">
	import Breadcrumb from '$components/Breadcrumb.svelte';
	import FilterablePieList from '$components/FilterablePieList.svelte';
	import type { PageData } from './$types';

	export let data: PageData;
</script>

<Breadcrumb />

<FilterablePieList pies={data.pies} pieRankings={data.pieRankings} fixedCategory={data.category} />
