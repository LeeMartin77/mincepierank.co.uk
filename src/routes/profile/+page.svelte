<script lang="ts">
  import { signIn, signOut } from '@auth/sveltekit/client';
  import { page } from '$app/stores';
  import type { PageData } from './$types';

  export let data: PageData;
</script>

{#if $page.data.session && $page.data.session.user}
  {#each data.years as year}
    <a class="year-link" href="/profile/rankings/{year}">Rankings for {year}</a>
    <a class="year-link" href="/profile/pies/{year}">Custom pies for {year}</a>
  {/each}
  <span class="signedInText">
    <small>Signed in as</small><br />
    <strong>{$page.data.session.user.email || $page.data.session.user.name}</strong>
  </span>
  <button on:click={() => signOut()} class="button">Sign out</button>
{:else}
  <span class="notSignedInText">You are not signed in</span>
  <button on:click={() => signIn()}>Sign In</button>
{/if}

<style>
  .year-link {
    display: block;
    padding: 1em;
    color: initial;
    text-align: center;
    text-decoration: none;
    font-weight: bold;
    font-size: 2em;
    padding: 0.5em;
    border: 1px solid black;
    border-radius: 1em;
  }
  .year-link:visited {
    color: initial;
  }
</style>
