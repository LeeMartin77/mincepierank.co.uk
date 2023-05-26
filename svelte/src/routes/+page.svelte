<script lang="ts">
	import BrandLinkCard from '$components/BrandLinkCard.svelte';
	import PieLinkCard from '$components/PieLinkCard.svelte';
	import type { PageData } from './$types';

	export let data: PageData;
</script>
<svelte:head>
	<title>Mince Pie Rank :: Homepage</title>
</svelte:head>

<div>
	<p>Mince Pie Rank is currently closed until the next Christmas period!</p>
</div>

{#if data.topPie && data.topPieRanking}
	<h2>Current Top Pie</h2>
	<PieLinkCard pie={data.topPie} pieListRanking={data.topPieRanking} />
{/if}

<h2>Brands</h2>
{#each data.makers as maker}
	<BrandLinkCard year={2022} maker={maker} />
{/each}
