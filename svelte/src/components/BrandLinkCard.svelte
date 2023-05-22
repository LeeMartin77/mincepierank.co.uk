<script lang="ts">
	import type { Maker } from "$lib/storage";

  export let maker: Maker;
  export let year: number;
</script>

<a href={`/years/${year}/makers/${maker.id}`}>
  <img src={`${maker.logo}?width=250&filter=gaussian`} alt={`${maker.name} Logo`}>
  <h4>{maker.name}</h4>
</a>