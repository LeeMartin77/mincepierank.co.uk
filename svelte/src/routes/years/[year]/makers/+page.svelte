<script lang="ts">
  import type { PageData } from './$types';

  export let data: PageData;
</script>
<ul>
	{#each data.makers as maker}
		<li><a href={`/years/2022/makers/${maker.id}`}>
			{maker.name}
		</a></li>
	{/each}
</ul>