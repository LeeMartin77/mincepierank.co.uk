<script lang="ts">
  export let raised = false;
</script>

{#if raised}
  <div class="card raised">
    <slot />
  </div>
{:else}
  <div class="card">
    <slot />
  </div>
{/if}

<style>
  .card {
    padding: 1em;
    border-radius: 1em;
    display: flex;
    flex-direction: column;
    box-shadow: 0.2em 0.1em 0.5em rgba(0, 0, 0, 0.1);
    background: rgba(255, 255, 255, 0.2);
    backdrop-filter: blur(5px);
    -webkit-backdrop-filter: blur(5px);
    border: 1px solid rgba(255, 255, 255, 0.3);
  }

  .raised {
    box-shadow: 0.3em 0.2em 0.6em rgba(0, 0, 0, 0.2);
  }
</style>
