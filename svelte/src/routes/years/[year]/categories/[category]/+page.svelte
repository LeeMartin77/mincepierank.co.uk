<script lang="ts">
	import PieLinkCard from '$components/PieLinkCard.svelte';
	import type { PageData } from './$types';

	export let data: PageData;
</script>

{#each data.pies as pie}
	<PieLinkCard pie={pie} />
{/each}
