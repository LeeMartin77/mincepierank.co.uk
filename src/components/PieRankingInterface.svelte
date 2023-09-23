<script lang="ts">
  import { ppCategory } from '$components/utilities/formatCategory';
  import { imgprssrPrefix } from '$lib/imgprssr';
  import type { MakerPieRanking } from '$lib/storage';
  import { onMount } from 'svelte';

  export let readonly: boolean = false;

  export let pie: {
    year: number;
    makerid: string;
    id: string;
  };

  let initialLoad = true;
  let submitting = false;
  let userRanking: MakerPieRanking;

  type HotRanking = Partial<Omit<MakerPieRanking, 'userid'>> &
    Omit<MakerPieRanking, 'userid' | 'pastry' | 'filling' | 'topping' | 'looks' | 'value'>;

  let hotRanking: HotRanking = {
    year: pie.year,
    makerid: pie.makerid,
    pieid: pie.id
  };

  onMount(() => {
    fetch(`/api/ranking?year=${pie.year}&makerid=${pie.makerid}&pieid=${pie.id}`)
      .then((res) => {
        if (res.status === 200) {
          return res.json();
        }
      })
      .then((rnking: MakerPieRanking) => {
        userRanking = rnking;
        hotRanking = rnking;
      })
      .finally(() => {
        initialLoad = false;
      });
  });

  const handleRankingChange = (rnkn: HotRanking) => {
    if (
      rnkn.filling === undefined ||
      rnkn.pastry === undefined ||
      rnkn.topping === undefined ||
      rnkn.value === undefined ||
      rnkn.looks === undefined
    ) {
      return;
    }
    submitting = true;
    fetch(`/api/ranking`, { method: 'POST', body: JSON.stringify(rnkn) })
      .then((res) => {
        if (res.status === 200) {
          return res.json();
        }
      })
      .then((rnking: MakerPieRanking) => {
        userRanking = rnking;
        hotRanking = rnking;
      })
      .finally(() => {
        submitting = false;
      });
  };
</script>

{#if initialLoad}
  <div>Loading</div>
{:else if userRanking && readonly}
  <div>
    <h3>My Ranking</h3>
    <dl>
      <dt>Pastry</dt>
      <dd>{userRanking.pastry}</dd>
      <dt>Filling</dt>
      <dd>{userRanking.filling}</dd>
      <dt>Topping</dt>
      <dd>{userRanking.topping}</dd>
      <dt>Looks</dt>
      <dd>{userRanking.looks}</dd>
      <dt>value</dt>
      <dd>{userRanking.value}</dd>
    </dl>
  </div>
{:else if !readonly}
  <div>
    <h3>My Ranking</h3>
    <dl>
      <dt>Pastry</dt>
      <dd>
        <input
          type="number"
          min="0"
          max="5"
          bind:value={hotRanking.pastry}
          on:change={() => handleRankingChange(hotRanking)}
        />
      </dd>
      <dt>Filling</dt>
      <dd>
        <input
          type="number"
          min="0"
          max="5"
          bind:value={hotRanking.filling}
          on:change={() => handleRankingChange(hotRanking)}
        />
      </dd>
      <dt>Topping</dt>
      <dd>
        <input
          type="number"
          min="0"
          max="5"
          bind:value={hotRanking.topping}
          on:change={() => handleRankingChange(hotRanking)}
        />
      </dd>
      <dt>Looks</dt>
      <dd>
        <input
          type="number"
          min="0"
          max="5"
          bind:value={hotRanking.looks}
          on:change={() => handleRankingChange(hotRanking)}
        />
      </dd>
      <dt>value</dt>
      <dd>
        <input
          type="number"
          min="0"
          max="5"
          bind:value={hotRanking.value}
          on:change={() => handleRankingChange(hotRanking)}
        />
      </dd>
    </dl>
  </div>
{/if}
