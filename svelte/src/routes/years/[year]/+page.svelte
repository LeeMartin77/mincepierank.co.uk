<script lang="ts">
  import type { PageData } from './$types';

  export let data: PageData;
</script>

<h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://kit.svelte.dev">kit.svelte.dev</a> to read the documentation</p>

{#if data.topPie}
<div>
  <img src={data.topPie.image_file} alt={data.topPie.displayname} />
  <h2>{data.topPie.displayname}</h2>
</div>
{/if}

<ul>
	{#each data.makers as maker}
		<li><a href={`/${maker.id}`}>
			{maker.name}
		</a></li>
	{/each}
</ul>