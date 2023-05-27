<script lang="ts">
	import BrandLinkList from '$components/BrandLinkList.svelte';
	import Breadcrumb from '$components/Breadcrumb.svelte';
	import type { PageData } from './$types';

	export let data: PageData;
</script>

<Breadcrumb />

<BrandLinkList makers={data.makers} year={data.year} />
