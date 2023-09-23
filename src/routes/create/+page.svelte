<script lang="ts">
  import { signIn } from '@auth/sveltekit/client';
  import { page } from '$app/stores';
  //import type { PageData } from './$types';

  //export let data: PageData;
</script>

{#if $page.data.session && $page.data.session.user && $page.data.session.user.email}
  <form method="POST" action="?/upload" enctype="multipart/form-data">
    <div>
      <label for="maker">maker</label>
      <input name="maker" />
    </div>
    <div>
      <label for="location">location</label>
      <textarea name="location" />
    </div>
    <div>
      <label for="displayname">displayname</label>
      <input name="displayname" />
    </div>
    <div>
      <label for="fresh">fresh</label>
      <input type="checkbox" name="fresh" />
    </div>
    <div>
      <label for="labels">labels</label>
      <input name="labels" />
    </div>
    <div>
      <label for="web_link">web_link</label>
      <input name="web_link" />
    </div>
    <div>
      <label for="pack_count">pack_count</label>
      <input type="number" name="pack_count" />
    </div>
    <div>
      <label for="pack_price_in_pence">pack_price_in_pence</label>
      <input type="number" name="pack_price_in_pence" />
    </div>
    <input type="file" name="image" />
    <button>Create</button>
  </form>
{:else}
  <span class="notSignedInText">Sign in with google to create a pie!</span>
  <button on:click={() => signIn()}>Sign In</button>
{/if}
