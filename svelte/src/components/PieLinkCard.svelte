<script lang="ts">
	import type { MakerPie } from "$lib/storage";
	import { ppCategory } from "./utilities/formatCategory";
  export let pie: MakerPie;

</script>

<div class="pie-card">
  <div>
    <img width="30%" src={`${pie.image_file}?width=250&height=250&filter=gaussian`} alt={pie.displayname} />
    <div>
      
    </div>
  </div>
  <div class="category-links">
    {#each pie.labels as label}
      <a class="category-link" href={`/years/${pie.year}/categories/${label}`}>{ppCategory(label)}</a>
    {/each}
  </div>
  <a class="pie-link" href={`/years/${pie.year}/makers/${pie.makerid}/${pie.id}`}>{pie.displayname}</a>
</div>

<style>
.pie-card {
  padding: 1em;
  border: 1px solid black;
  border-radius: 1em;
}
.pie-link {
  display: block;
  border: 1px solid black;
  border-radius: 1em;
  padding: 1em;
  text-decoration: none;
  text-align: center;
  font-weight: bold;
  color: initial;
}
.pie-link:visited {
  color: initial;
}

.category-links {
  display: flex;
  justify-content: center;
  gap: 0.5em;
  margin-top: 0.5em;
  margin-bottom: 0.5em;
}

.category-link {
  display: block;
  border: 1px solid black;
  font-size: 0.8em;
  border-radius: 2em;
  padding: 0.5em;
  text-decoration: none;
  text-align: center;
  color: initial;
}
.category-link:visited {
  color: initial;
}
</style>