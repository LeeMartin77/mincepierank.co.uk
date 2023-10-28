<script lang="ts">
  import type { PageData } from './$types';
  import PieRankingInterface from '$components/PieRankingInterface.svelte';
  import type { MakerPie, PieRankingSummary } from '$lib/storage';
  import LinkButton from '$components/generic/LinkButton.svelte';
  import PieLinkCard from '$components/PieLinkCard.svelte';
  import LoadingIndicator from '$components/generic/LoadingIndicator.svelte';
  import ImageSelect from '$components/generic/ImageSelect.svelte';
  import { imgprssrPrefix } from '$lib/imgprssr';

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
    <ImageSelect
      placeholder="Select Brand"
      options={data.makers.map((x) => ({
        imageLink: imgprssrPrefix(x.logo + '?width=50&filter=gaussian', data.imgprssr),
        text: x.name,
        value: x.id
      }))}
      on:change={({ detail }) => {
        selectedmaker = detail;
        handleMakerChange(detail);
      }}
    />
  </div>
  <div style="margin-bottom: 1em; display: flex; flex-direction: column; align-items: center;">
    <label for="pie">Pie</label>
    <ImageSelect
      placeholder="Select Pie"
      disabled={loadingPies || !selectedmaker || pies.length === 0}
      options={pies.map((x) => ({
        imageLink: imgprssrPrefix(x.image_file + '?width=50&filter=gaussian', data.imgprssr),
        text: x.displayname,
        value: x.id
      }))}
      on:change={({ detail }) => {
        selectedpie = detail;
        handlePieChange(detail);
      }}
    />
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
</style>
