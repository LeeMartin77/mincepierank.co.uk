<script>
  import { signIn, signOut } from '@auth/sveltekit/client';
  import { page } from '$app/stores';
</script>

{#if $page.data.session && $page.data.session.user}
  <a class="year-link" href="/profile/rankings/2022">2022</a>
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
