<script lang="ts">
	import type { MakerPie } from "$lib/storage";
	import { ppCategory } from "./utilities/formatCategory";
  export let pie: MakerPie;

</script>

<div>
  <img src={`${pie.image_file}?width=250&filter=gaussian`} alt={pie.displayname} />
  <div>
    {#each pie.labels as label}
      <a href={`/years/${pie.year}/categories/${label}`}>{ppCategory(label)}</a>
    {/each}
  </div>
  <a href={`/years/${pie.year}/makers/${pie.makerid}/${pie.id}`}>{pie.displayname}</a>
</div>