<script lang="ts">
  import type { MakerPieRanking } from '$lib/storage';
  import { onMount } from 'svelte';

  export let readonly = false;

  export let year: number;
  export let makerid: string;
  export let pieid: string;

  $: {
    if (year && makerid && pieid) {
      loadPie();
    }
  }

  let initialLoad = true;
  let submitting = false;
  let userRanking: MakerPieRanking;

  type HotRanking = Partial<Omit<MakerPieRanking, 'userid'>> &
    Omit<MakerPieRanking, 'userid' | 'pastry' | 'filling' | 'topping' | 'looks' | 'value'>;

  let hotRanking: HotRanking = {
    year,
    makerid,
    pieid
  };

  const loadPie = () => {
    initialLoad = true;
    fetch(`/api/ranking?year=${year}&makerid=${makerid}&pieid=${pieid}`)
      .then((res) => {
        if (res.status === 200) {
          return res.json();
        }
        return Promise.resolve({
          year,
          makerid,
          pieid
        });
      })
      .then((rnking: MakerPieRanking) => {
        userRanking = rnking;
        hotRanking = rnking;
      })
      .finally(() => {
        initialLoad = false;
      });
  };

  onMount(loadPie);

  const handleRankingChange = (rnkn: HotRanking) => {
    if (!rnkn.filling || !rnkn.pastry || !rnkn.topping || !rnkn.value || !rnkn.looks) {
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
          min="1"
          max="5"
          bind:value={hotRanking.pastry}
          on:change={() => handleRankingChange(hotRanking)}
        />
      </dd>
      <dt>Filling</dt>
      <dd>
        <input
          type="number"
          min="1"
          max="5"
          bind:value={hotRanking.filling}
          on:change={() => handleRankingChange(hotRanking)}
        />
      </dd>
      <dt>Topping</dt>
      <dd>
        <input
          type="number"
          min="1"
          max="5"
          bind:value={hotRanking.topping}
          on:change={() => handleRankingChange(hotRanking)}
        />
      </dd>
      <dt>Looks</dt>
      <dd>
        <input
          type="number"
          min="1"
          max="5"
          bind:value={hotRanking.looks}
          on:change={() => handleRankingChange(hotRanking)}
        />
      </dd>
      <dt>value</dt>
      <dd>
        <input
          type="number"
          min="1"
          max="5"
          bind:value={hotRanking.value}
          on:change={() => handleRankingChange(hotRanking)}
        />
      </dd>
    </dl>
  </div>
{/if}
