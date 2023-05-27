<script lang="ts">
	import BrandLinkList from '$components/BrandLinkList.svelte';
	// @ts-ignore
	import PieLinkCard from '$components/PieLinkCard.svelte';
	import type { PageData } from './$types';

	export let data: PageData;
</script>
<svelte:head>
	<title>Mince Pie Rank :: Homepage</title>
</svelte:head>

<div>
	<p style="text-align: center;">Mince Pie Rank is currently closed until the next Christmas period!</p>
</div>

{#if data.topPie && data.topPieRanking}
	<h2>Current Top Pie</h2>
	<PieLinkCard pie={data.topPie} pieListRanking={data.topPieRanking} />
{/if}

<BrandLinkList makers={data.makers} year={2022}/>

<style>
	h2 {
		text-align: center;
	}
</style>