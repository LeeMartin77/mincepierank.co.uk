<script lang="ts">
  import Breadcrumb from '$components/Breadcrumb.svelte';
  import { onMount } from 'svelte';
  import type { LayoutData } from './$types';
  import Menu from '$components/Menu.svelte';
  export let data: LayoutData;

  let noticeVisible = false;
  let noticestring = `${data.readonly ? 'READONLY' : ''}:${data.notice ? data.notice : ''}`;

  let cookieNotice = false;

  const noticekey = 'noticestring';

  const cookiekey = 'cookiekey';

  onMount(() => {
    if (data.session?.user?.email) {
      // need to look into some form of event for this
      fetch(`/api/claim-anonymous`, { method: 'POST' });
    }

    noticeVisible = data.readonly || !!data.notice;

    let lsstr = localStorage.getItem(noticekey);
    if (lsstr && lsstr === noticestring) {
      noticeVisible = false;
    }

    cookieNotice = !window.location.href.endsWith('/about/cookies');
    if (localStorage.getItem(cookiekey) === 'accepted') {
      cookieNotice = false;
    }
  });
</script>

{#if noticeVisible}
  <div class="notices">
    <button
      on:click={() => {
        localStorage.setItem(noticekey, noticestring);
        noticeVisible = false;
      }}>X</button
    >
    <div>
      {#if data.readonly}
        <p>Mince Pie Rank is read-only until the next Christmas period</p>
      {/if}
      {#if data.notice}
        <p>{data.notice}</p>
      {/if}
    </div>
  </div>
{/if}
{#if cookieNotice}
  <div class="cookie-notice">
    <div>
      <p>We use cookies for authentication and associating mince pie rankings to a user</p>
      <button
        on:click={() => {
          localStorage.setItem(cookiekey, 'accepted');
          cookieNotice = false;
        }}>Accept</button
      >
      <button
        on:click={() => {
          cookieNotice = false;
        }}>Dismiss</button
      >
      <a href="/about/cookies">Read Full Policy</a>
    </div>
  </div>
{/if}
<div class="page-container">
  <Menu activeYear={data.activeYear} customPiesEnabled={data.customPiesEnabled} />
  <div class="content-container">
    <Breadcrumb />
    <div class="content">
      <slot />
    </div>
    <a style="padding-bottom: 6em" href="/about">About</a>
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
    overflow-y: scroll;
    max-height: calc(100vh - 2em);
  }
  .notices {
    text-align: center;
    border: 1px dashed grey;
    border-radius: 1em;
    position: fixed;
    width: calc(100vw - 4em);
    background-color: white;
    top: 1em;
    left: 1em;
    display: flex;
    padding: 1em;
    align-items: center;
    gap: 1em;
    z-index: 99999;
  }

  .notices button {
    height: 2em;
    width: 2em;
    background: none;
    border: 1px solid black;
    border-radius: 0.5em;
  }

  .cookie-notice {
    text-align: center;
    border: 1px dashed grey;
    border-radius: 1em;
    position: fixed;
    width: calc(100vw - 4em);
    background-color: white;
    top: 1em;
    left: 1em;
    display: flex;
    padding: 1em;
    align-items: center;
    gap: 1em;
    z-index: 9999999999999;
  }

  .cookie-notice button {
    display: block;
    width: 100%;
    cursor: pointer;
    border-radius: 1em;
    box-shadow: 0em 0em 0.3em rgba(0, 0, 0, 0.2);
    background: white;
    border: none;
    padding: 1em;
    text-align: center;
    font-weight: bold;
    color: #111;
    margin-bottom: 1em;
  }

  .cookie-notice a {
    display: block;
    border-radius: 1em;
    padding: 0.5em;
    text-align: center;
    font-weight: bold;
    color: #111;
  }
  .cookie-notice a:visited {
    color: #111;
  }

  .notices button:hover {
    cursor: pointer;
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
