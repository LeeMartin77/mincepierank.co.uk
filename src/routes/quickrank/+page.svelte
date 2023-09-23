<script lang="ts">
  import type { PageData } from './$types';
  import PieRankingInterface from '$components/PieRankingInterface.svelte';
  import type { MakerPie } from '$lib/storage';
  import { imgprssrPrefix } from '$lib/imgprssr';

  export let data: PageData;

  let selectedmaker = '';
  let selectedpie = '';

  let pies: MakerPie[] = [];

  let loadingPies = false;

  const handleMakerChange = (slmkr: string) => {
    loadingPies = true;
    pies = [];
    fetch(`/api/pies?year=${data.activeYear}&makerid=${slmkr}`)
      .then((res) => {
        if (res.status === 200) {
          return res.json();
        }
      })
      .then((ps: MakerPie[]) => {
        pies = ps;
        selectedpie = '';
      })
      .finally(() => {
        loadingPies = false;
      });
  };

  $: pieindex = pies ? pies.findIndex((x) => x.id === selectedpie) : -1;
</script>

<div>
  <div>
    <label for="maker">Maker</label>
    <select
      name="maker"
      bind:value={selectedmaker}
      on:change={() => handleMakerChange(selectedmaker)}
    >
      <option value="">Select Maker</option>
      {#each data.makers as maker}
        <option value={maker.id}>{maker.name}</option>
      {/each}
    </select>
  </div>
  <div>
    <label for="pie">Pie</label>
    <select
      disabled={loadingPies || !selectedmaker || pies.length === 0}
      name="pie"
      bind:value={selectedpie}
      on:change={() => {
        selectedpie = selectedpie;
      }}
    >
      <option value="">Select Pie</option>
      {#if pies}
        {#each pies as pie}
          <option value={pie.id}>{pie.displayname}</option>
        {/each}
      {/if}
    </select>
  </div>
  {#if pies && pieindex > -1}
    <div>
      <h3>{pies[pieindex].displayname}</h3>

      <img
        src={imgprssrPrefix(
          pies[pieindex].image_file + '?filter=gaussian&width=250',
          data.imgprssr
        )}
        alt={pies[pieindex].displayname}
      />
    </div>
  {/if}
  {#if pies && selectedmaker && selectedpie}
    <div>
      <PieRankingInterface
        readonly={false}
        year={data.activeYear}
        makerid={selectedmaker}
        pieid={selectedpie}
      />
    </div>
  {/if}
</div>
