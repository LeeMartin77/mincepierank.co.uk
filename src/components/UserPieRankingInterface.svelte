<script lang="ts">
  import type { UserPieRanking } from '$lib/storage';
  import { onMount, createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher();

  export let readonly = false;

  export let year: number;
  export let pieid: string;

  $: {
    if (mounted && year && pieid) {
      loadPie();
    }
  }

  let initialLoad = true;
  let submitting = false;
  let mounted = false;
  let userRanking: UserPieRanking;

  type HotRanking = Partial<Omit<UserPieRanking, 'userid'>> &
    Omit<UserPieRanking, 'userid' | 'pastry' | 'filling' | 'topping' | 'looks' | 'value'>;

  let hotRanking: HotRanking = {
    year,
    pieid
  };

  const loadPie = () => {
    initialLoad = true;
    fetch(`/api/user-ranking?year=${year}&pieid=${pieid}`)
      .then((res) => {
        if (res.status === 200) {
          return res.json();
        }
        return Promise.resolve({
          year,
          pieid
        });
      })
      .then((rnking: UserPieRanking) => {
        userRanking = rnking;
        hotRanking = rnking;
      })
      .finally(() => {
        initialLoad = false;
      });
  };

  onMount(() => {
    mounted = true;
  });

  const handleRankingChange = (rnkn: HotRanking) => {
    if (!rnkn.filling || !rnkn.pastry || !rnkn.topping || !rnkn.value || !rnkn.looks) {
      return;
    }
    submitting = true;
    fetch(`/api/user-ranking`, { method: 'POST', body: JSON.stringify(rnkn) })
      .then((res) => {
        if (res.status === 200) {
          return res.json();
        }
      })
      .then((rnking: UserPieRanking) => {
        userRanking = rnking;
        hotRanking = rnking;
        dispatch('newRanking', rnking);
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
