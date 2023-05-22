<script lang="ts">
	import type { Maker } from "$lib/storage";

  export let maker: Maker;
</script>

<a href={`/years/2022/makers/${maker.id}`}>
  <img src={`${maker.logo}?width=250&filter=gaussian`} alt={`${maker.name} Logo`}>
  <h4>{maker.name}</h4>
</a>