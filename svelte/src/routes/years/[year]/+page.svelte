<script lang="ts">
	import CategoryLinkList from '$components/CategoryLinkList.svelte';
	import Breadcrumb from "$components/Breadcrumb.svelte";

	import type { PageData } from './$types';
	import BrandLinkList from '$components/BrandLinkList.svelte';

	export let data: PageData;
</script>

<Breadcrumb />
<CategoryLinkList year={data.year} categories={data.categories} />
<BrandLinkList makers={data.makers} year={data.year}/>