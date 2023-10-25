<script lang="ts">
  import { signIn, signOut } from '@auth/sveltekit/client';
  import { page } from '$app/stores';
  import type { PageData } from './$types';
  import LinkButton from '$components/generic/LinkButton.svelte';

  export let data: PageData;
</script>

<svelte:head>
  <title>My Profile :: Mince Pie Rank</title>
</svelte:head>

<div style="display: flex; flex-direction: column; align-items: center; gap: 1em;">
  {#if $page.data.session && $page.data.session.user}
    <span class="signedInText">
      <small>Signed in as</small><br />
      <strong>{$page.data.session.user.email || $page.data.session.user.name}</strong>
    </span>
    <button
      style="padding: 1em; font-weight: bold; border-radius: 1em;"
      on:click={() => signOut()}
      class="button">Sign out</button
    >

    {#if data.customPiesEnabled}
      <LinkButton href="/profile/pies">All Custom pies</LinkButton>
    {/if}
  {:else}
    <span class="notSignedInText">You are not signed in</span>
    <button on:click={() => signIn()}>Sign In</button>
  {/if}

  <LinkButton href="/profile/rankings">All Rankings</LinkButton>
</div>

<style>
</style>
