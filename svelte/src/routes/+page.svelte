<script lang="ts">
	import { ppCategory } from '$components/utilities/formatCategory';
	import type { PageData } from './$types';

	export let data: PageData;

	const { topPie } = data;
</script>
<svelte:head>
	<title>Mince Pie Rank :: Homepage</title>
</svelte:head>

<div>
	<p>Mince Pie Rank is currently closed until the next Christmas period!</p>
</div>

{#if topPie}
	<h2>Current Top Pie</h2>
	<div>
		<img src={`${topPie.image_file}?width=250&filter=gaussian`} alt={topPie.displayname} />
		<div>
			{#each topPie.labels as label}
				<a href={`/years/${topPie.year}/categories/${label}`}>{ppCategory(label)}</a>
			{/each}
		</div>
		<a href={`/years/${topPie.year}/makers/${topPie.makerid}/${topPie.id}`}>{topPie.displayname}</a>

	</div>
{/if}

<h2>Brands</h2>
<ul>
	{#each data.makers as maker}
		<li>
			<a href={`/years/2022/makers/${maker.id}`}>
				<img src={`${maker.logo}?width=250&filter=gaussian`} alt={`${maker.name} Logo`}>
				<h4>{maker.name}</h4>
			</a>
		</li>
	{/each}
</ul>
