<script lang="ts">
  import type { UserPieRanking } from '$lib/storage';
  import { onMount, createEventDispatcher } from 'svelte';
  import PieRankingEditor from './PieRankingEditor.svelte';
  import PieRankingSummary from './PieRankingSummary.svelte';
  import Card from './generic/Card.svelte';
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
{:else if userRanking && userRanking.pastry && readonly}
  <Card>
    <h3>My Ranking</h3>
    <PieRankingSummary rankingSummary={userRanking} />
  </Card>
{:else if !readonly}
  <PieRankingEditor
    {hotRanking}
    {submitting}
    on:change={(e) => {
      handleRankingChange(e.detail);
    }}
  />
{/if}
