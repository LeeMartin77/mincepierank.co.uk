<script lang="ts">
  import type { MakerPieRanking } from '$lib/storage';
  import { onMount, createEventDispatcher } from 'svelte';
  import PieRankingEditor from './PieRankingEditor.svelte';
  const dispatch = createEventDispatcher();

  export let readonly = false;

  export let year: number;
  export let makerid: string;
  export let pieid: string;

  $: {
    if (mounted && year && makerid && pieid) {
      loadPie();
    }
  }

  let initialLoad = true;
  let submitting = false;
  let mounted = false;
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

  onMount(() => {
    mounted = true;
  });

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
        dispatch('newRanking', rnking);
      })
      .finally(() => {
        submitting = false;
      });
  };
</script>

{#if initialLoad}
  <div>Loading</div>
{:else if userRanking && userRanking.pastry && readonly}
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
  <PieRankingEditor
    {hotRanking}
    on:change={(e) => {
      handleRankingChange(e.detail);
    }}
  />
{/if}
