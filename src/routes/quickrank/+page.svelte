<script lang="ts">
  import type { PageData } from './$types';
  import PieRankingInterface from '$components/PieRankingInterface.svelte';
  import type { MakerPie, PieRankingSummary } from '$lib/storage';
  import LinkButton from '$components/generic/LinkButton.svelte';
  import PieLinkCard from '$components/PieLinkCard.svelte';
  import LoadingIndicator from '$components/generic/LoadingIndicator.svelte';

  export let data: PageData;

  let selectedmaker = '';
  let selectedpie = '';

  let pies: MakerPie[] = [];

  let ranking: PieRankingSummary | undefined = undefined;

  let loadingPies = false;
  let loadingPieRanking = false;

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

  const handlePieChange = (slmpie: string, setloading = true) => {
    loadingPieRanking = setloading;
    if (setloading) {
      ranking = undefined;
    }
    fetch(`/api/ranking/summary?year=${data.activeYear}&makerid=${selectedmaker}&pieid=${slmpie}`)
      .then((res) => {
        if (res.status === 200) {
          return res.json();
        }
      })
      .then((ps: PieRankingSummary) => {
        ranking = ps;
      })
      .finally(() => {
        loadingPieRanking = false;
      });
  };

  $: pieindex = pies ? pies.findIndex((x) => x.id === selectedpie) : -1;
</script>

<svelte:head>
  <title>Quickrank :: Mince Pie Rank</title>
</svelte:head>

<div style="display: flex; flex-direction: column; gap: 2em;">
  {#if data.customPiesEnabled}
    <div style="margin-bottom: 1em;">
      <p>Can't find what you're looking for?</p>
      <LinkButton href="/create">Create a custom pie!</LinkButton>
    </div>
  {/if}
  <div style="display: flex; flex-direction: column; align-items: center;">
    <label for="maker">Brand</label>
    <select
      name="maker"
      bind:value={selectedmaker}
      on:change={() => handleMakerChange(selectedmaker)}
    >
      <option value="">Select Brand</option>
      {#each data.makers as maker}
        <option value={maker.id}>{maker.name}</option>
      {/each}
    </select>
  </div>
  <div style="margin-bottom: 1em; display: flex; flex-direction: column; align-items: center;">
    <label for="pie">Pie</label>
    <select
      disabled={loadingPies || !selectedmaker || pies.length === 0}
      name="pie"
      bind:value={selectedpie}
      on:change={() => {
        selectedpie = selectedpie;
        handlePieChange(selectedpie);
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
    <div style="position: relative;">
      <div style="position:absolute; top:20%; right: 25%; z-index: 99;">
        {#if loadingPieRanking}
          <LoadingIndicator />
        {/if}
      </div>
      <PieLinkCard pie={pies[pieindex]} imgprssr={data.imgprssr} pieListRanking={ranking} />
    </div>
  {/if}
  {#if pies && selectedmaker && selectedpie}
    <div>
      <PieRankingInterface
        readonly={false}
        year={data.activeYear}
        makerid={selectedmaker}
        pieid={selectedpie}
        on:newRanking={(newrnking) => {
          loadingPieRanking = true;
          setTimeout(() => {
            handlePieChange(newrnking.detail.pieid, false);
          }, 500);
        }}
      />
    </div>
  {/if}
</div>

<style>
  select {
    background-color: rgba(0, 0, 0, 0.04);
    box-shadow: 0 0.2em 0.4em rgba(0, 0, 0, 0.1);
    width: 16em;
  }
  select:disabled {
    background-color: white;
    box-shadow: none;
  }
</style>
