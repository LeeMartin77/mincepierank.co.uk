<script lang="ts">
	import PieLinkCard from '$components/PieLinkCard.svelte';
	import { ppCategory } from '$components/utilities/formatCategory';
	import type { PageData } from './$types';

	export let data: PageData;
</script>

<h1>{data.makerData.name}</h1>

{#each data.categories as category}
  <a href={`/years/${data.year}/categories/${category}`}>{ppCategory(category)}</a>
{/each}

{#each data.pies as pie}
	<PieLinkCard pie={pie} />
{/each}
