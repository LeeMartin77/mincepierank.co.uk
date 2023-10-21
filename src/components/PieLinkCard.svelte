<script lang="ts">
  import { imgprssrPrefix } from '$lib/imgprssr';
  import type { MakerPie, UserPie } from '$lib/storage';
  import Card from './generic/Card.svelte';
  import LinkButton from './generic/LinkButton.svelte';
  import PieGraph from './subcomponents/PieGraph.svelte';
  import { ppCategory } from './utilities/formatCategory';
  import type { PieListRanking } from './utilities/mapPiesAndRankings';
  export let pie: UserPie | MakerPie;
  export let pieListRanking: PieListRanking | undefined = undefined;
  export let imgprssr: string;
  export let raised: boolean = false;
  let pielink: string;
  $: {
    if ('makerid' in pie) {
      pielink = `/years/${pie.year}/brands/${pie.makerid}/${pie.id}`;
    } else {
      pielink = `/years/${pie.year}/user-pies/${pie.id}`;
    }
  }
</script>

<Card {raised}>
  <div class="pie-card-top">
    <div style="width: 30%">
      {#if pie.image_file}
        <img
          width="100%"
          src={`${imgprssrPrefix(pie.image_file, imgprssr)}?width=250&height=250&filter=gaussian`}
          alt={pie.displayname}
        />
      {/if}
    </div>
    <div>
      <h4 style="margin-bottom: 0.25em;">Average: {(pieListRanking?.average || 0).toFixed(2)}/5</h4>
      <PieGraph score={pieListRanking?.average || 0} />
      <h5 style="margin-top: 0.25em;">{pieListRanking?.count || 0} Votes</h5>
    </div>
  </div>
  <div class="category-links">
    {#each pie.labels as label}
      <a class="category-link" href={`/years/${pie.year}/categories/${label}`}
        >{ppCategory(label)}</a
      >
    {/each}
  </div>
  <LinkButton href={pielink}>{pie.displayname}</LinkButton>
</Card>

<style>
  .pie-card-top {
    display: flex;
    gap: 2em;
    justify-content: center;
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
    box-shadow: 0em 0em 0.2em rgba(0, 0, 0, 0.2);
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
