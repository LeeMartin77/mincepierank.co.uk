<script lang="ts">
	import BrandLinkCard from '$components/BrandLinkCard.svelte';
	import Breadcrumb from '$components/Breadcrumb.svelte';
	import CategoryLinkList from '$components/CategoryLinkList.svelte';
	import type { PageData } from './$types';

	export let data: PageData;
</script>

<Breadcrumb />

<CategoryLinkList year={data.year} categories={data.categories} />

{#each data.makers as maker}
	<BrandLinkCard year={data.year} maker={maker} />
{/each}