<script>
  import { signIn, signOut } from '@auth/sveltekit/client';
  import { page } from '$app/stores';
</script>

<p>
  {#if $page.data.session && $page.data.session.user}
    <span class="signedInText">
      <small>Signed in as</small><br />
      <strong>{$page.data.session.user.email || $page.data.session.user.name}</strong>
    </span>
    <button on:click={() => signOut()} class="button">Sign out</button>
  {:else}
    <span class="notSignedInText">You are not signed in</span>
    <button on:click={() => signIn()}>Sign In</button>
  {/if}
</p>
