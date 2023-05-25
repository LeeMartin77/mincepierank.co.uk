<script lang="ts">
	import type { MakerPie } from "$lib/storage";
	import PieLinkCard from "./PieLinkCard.svelte";
	import { ppCategory } from "./utilities/formatCategory";
  export let pies: MakerPie[];
  export let fixedCategory: string | undefined = undefined;

</script>

{#if fixedCategory}
  <h3>{ppCategory(fixedCategory)}</h3>
{/if}

{#each pies as pie}
  <PieLinkCard pie={pie} />
{/each}