<script lang="ts">
	import type { Maker } from '$lib/storage';

	export let maker: Maker;
	export let year: number;
</script>

<div class="brand-link-card">
	<a class="brand-link-image" href={`/years/${year}/brands/${maker.id}`}>
		<div>
			<img src={`${maker.logo}?width=200&filter=gaussian`} alt={`${maker.name} Logo`} />
		</div>
	</a>

	<a class="brand-link" href={`/years/${year}/brands/${maker.id}`}>
		{maker.name}
	</a>
</div>

<style>
	.brand-link-card {
		display: flex;
		gap: 0.5em;
		flex-direction: column;
		align-items: center;
		padding: 1em 1.5em;
		border: 1px solid black;
		border-radius: 1em;
	}
	.brand-link-image {
		display: flex;
		height: 250px;
		align-items: center;
	}
	.brand-link {
		color: initial;
		text-decoration: none;
		font-weight: bold;
		padding: 0.5em;
		border: 1px solid black;
		width: 100%;
		text-align: center;
		border-radius: 0.5em;
	}
	.brand-link:visited {
		color: initial;
	}
</style>
