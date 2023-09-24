<script lang="ts">
  import { signIn, signOut } from '@auth/sveltekit/client';
  import { page } from '$app/stores';
  import type { PageData } from './$types';

  //export let data: PageData;
</script>

{#if $page.data.session && $page.data.session.user}
  <span class="signedInText">
    <small>Signed in as</small><br />
    <strong>{$page.data.session.user.email || $page.data.session.user.name}</strong>
  </span>
  <button on:click={() => signOut()} class="button">Sign out</button>
  <a class="year-link" href="/profile/rankings">All Rankings</a>
  <a class="year-link" href="/profile/pies">All Custom pies</a>
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

  @media (prefers-color-scheme: dark) {
    .year-link {
      color: #fff;
      border: 1px solid #fff;
    }
    .year-link:visited {
      color: #ddd;
    }
  }
</style>
