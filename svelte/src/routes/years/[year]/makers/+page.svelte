<script lang="ts">
	import BrandLinkCard from '$components/BrandLinkCard.svelte';
	import type { PageData } from './$types';

	export let data: PageData;
</script>

{#each data.makers as maker}
	<BrandLinkCard year={data.year} maker={maker} />
{/each}