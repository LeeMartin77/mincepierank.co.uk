<script lang="ts">
  const bgcolor =
    'filter: invert(96%) sepia(0%) saturate(508%) hue-rotate(145deg) brightness(91%) contrast(90%);';

  export let score: number = 0;
  export let topRank: number = 5;

  export let centerAlign = false;

  $: wrapperStyle = centerAlign
    ? 'display: flex; flex-direction: column; align-items: center;'
    : '';

  $: ranks = Array.apply(null, Array(topRank)).map((_, index) => index + 1);
</script>

<div style={wrapperStyle}>
  <div style="display: flex;">
    {#each ranks as rank}
      {#if score >= rank}
        <div><img alt="ranking icon" src="/rankicon.svg" width="24" /></div>
      {:else if score - (rank - 1) > 0}
        <div style="display: flex; width: 24px">
          <div style={`width: ${(score - (rank - 1)) * 100}%; overflow: hidden;`}>
            <img alt="ranking icon" src="/rankicon.svg" width="24" />
          </div>
          <div
            style={`width: calc(24px - ${
              (score - (rank - 1)) * 100
            }%); overflow: hidden; ${bgcolor}; display: flex; justify-content: flex-end; align-items: baseline;`}
          >
            <img alt="ranking icon" src="/rankicon.svg" width="24" />
          </div>
        </div>
      {:else}
        <div><img alt="ranking icon" src="/rankicon.svg" width="24" style={bgcolor} /></div>
      {/if}
    {/each}
  </div>
  <span>{score.toFixed(2)}/5.00</span>
</div>
