<script lang="ts">
  import Breadcrumb from '$components/Breadcrumb.svelte';
  import type { LayoutData } from './$types';
  export let data: LayoutData;
</script>

{#if data.readonly || data.notice}
  <div class="notices">
    {#if data.readonly}
      <p>Mince Pie Rank is read-only until the next Christmas period</p>
    {/if}
    {#if data.notice}
      <p>{data.notice}</p>
    {/if}
  </div>
{/if}

<Breadcrumb />

<div class="content-container">
  <div class="content">
    <slot />
  </div>
</div>

<style>
  .content-container {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  .notices {
    text-align: center;
    border: 1px dashed grey;
    border-radius: 1em;
  }

  @media only screen and (max-width: 640px) {
    .content {
      max-width: 460px;
    }
  }
  @media only screen and (min-width: 640px) and (max-width: 960px) {
    .content {
      min-width: 640px;
      max-width: 960px;
    }
  }

  @media only screen and (min-width: 960px) {
    .content {
      min-width: 960px;
      max-width: 1020px;
    }
  }
</style>
