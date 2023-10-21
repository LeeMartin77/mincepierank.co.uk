<script lang="ts">
  import Breadcrumb from '$components/Breadcrumb.svelte';
  import { onMount } from 'svelte';
  import type { LayoutData } from './$types';
  import Menu from '$components/Menu.svelte';
  export let data: LayoutData;

  onMount(() => {
    if (data.session?.user?.email) {
      // need to look into some form of event for this
      fetch(`/api/claim-anonymous`, { method: 'POST' });
    }
  });
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
<div class="page-container">
  <Menu activeYear={data.activeYear} />
  <div class="content-container">
    <Breadcrumb />
    <div class="content">
      <slot />
    </div>
  </div>
</div>

<style>
  .page-container {
    display: flex;
    flex-direction: row;
    top: 0;
    left: 0;
  }
  .content-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    margin-left: auto;
    margin-right: auto;
    padding: 1em;
  }
  .notices {
    text-align: center;
    border: 1px dashed grey;
    border-radius: 1em;
  }

  .content {
    padding-bottom: 6em;
  }

  @media only screen and (max-width: 640px) {
    .content {
      max-width: 460px;
    }
  }
  @media only screen and (min-width: 640px) and (max-width: 960px) {
    .content {
      min-width: 460px;
      max-width: 960px;
    }
  }

  @media only screen and (min-width: 1020px) {
    .content {
      min-width: 840px;
      max-width: 1020px;
    }
  }
</style>
