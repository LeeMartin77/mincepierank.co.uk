<script lang="ts">
	import type { PageData } from './$types';

	export let data: PageData;
</script>

{#each data.pies as pie}
	<a href={`/years/${pie.year}/makers/${pie.makerid}/${pie.id}`}>{pie.displayname}</a>
{/each}
