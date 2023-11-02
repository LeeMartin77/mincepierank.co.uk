<script lang="ts">
  import { signIn } from '@auth/sveltekit/client';
  import { page } from '$app/stores';
  import type { PageData } from './$types';

  export let data: PageData;
  export let form;

  let showExtra = false;

  const essentialdata: { [key: string]: any | undefined } = {
    displayname: undefined,
    maker: undefined,
    image: undefined
  };

  const optionaldata: { [key: string]: any | undefined } = {
    location: '',
    labels: '',
    web_link: '',
    pack_count: 1,
    pack_price_in_pence: 0
  };
</script>

<svelte:head>
  <title>Upload a Pie :: Mince Pie Rank</title>
</svelte:head>

{#if form?.error}
  <div>
    <p>{form.error}</p>
  </div>
{/if}
{#if $page.data.session && $page.data.session.user && $page.data.session.user.email}
  <form method="POST" action="?/upload" enctype="multipart/form-data">
    <input type="hidden" name="fresh" />
    <div>
      <label for="displayname">What's It Called?</label>
      <input name="displayname" bind:value={essentialdata.displayname} />
    </div>
    <div>
      <label for="maker">Who made it?</label>
      <input name="maker" bind:value={essentialdata.maker} />
    </div>
    {#if essentialdata.image && essentialdata.image[0].size > data.maxFileSize}
      <div>
        <span style="margin-left: auto; margin-right: 0;"> Image is too large! </span>
      </div>
    {/if}
    <div>
      <label for="image">Got a picture?</label>
      <input type="file" name="image" bind:files={essentialdata.image} />
    </div>
    {#if !showExtra}
      <input type="hidden" name="location" bind:value={optionaldata.location} />
      <input type="hidden" name="labels" bind:value={optionaldata.labels} />
      <input type="hidden" name="web_link" bind:value={optionaldata.web_link} />
      <input type="hidden" name="pack_count" bind:value={optionaldata.pack_count} />
      <input
        type="hidden"
        name="pack_price_in_pence"
        bind:value={optionaldata.pack_price_in_pence}
      />
      <button
        type="button"
        on:click={() => {
          showExtra = true;
        }}
        style="display: flex; align-items: center; justify-content:center;"
        ><span>Enter extra data?</span>
        <svg xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 -960 960 960" width="24"
          ><path d="M480-360 280-560h400L480-360Z" /></svg
        ></button
      >
    {/if}
    {#if showExtra}
      <button
        type="button"
        on:click={() => {
          showExtra = false;
        }}
        style="display: flex; align-items: center; justify-content:center;"
        ><span>Collapse</span>
        <svg xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 -960 960 960" width="24"
          ><path d="M480-360 280-560h400L480-360Z" /></svg
        ></button
      >
      <div>
        <label for="location">Where did you find it?</label>
        <textarea name="location" bind:value={optionaldata.location} />
      </div>
      <div>
        <label for="labels">How would you categorise it?</label>
        <input
          name="labels"
          placeholder="premium,rum,puff pastry"
          bind:value={optionaldata.labels}
        />
      </div>
      <div>
        <label for="web_link">Is there a website?</label>
        <input name="web_link" bind:value={optionaldata.web_link} />
      </div>
      <div>
        <label for="pack_count">How many in a pack?</label>
        <input type="number" name="pack_count" bind:value={optionaldata.pack_count} />
      </div>
      <div>
        <label for="pack_price_in_pence">How much is a pack in pence?</label>
        <input
          type="number"
          name="pack_price_in_pence"
          bind:value={optionaldata.pack_price_in_pence}
        />
      </div>
    {/if}
    <button
      type="submit"
      disabled={Object.values(essentialdata).some((x) => x === undefined) ||
        !essentialdata.image ||
        essentialdata.image[0].size > data.maxFileSize}>Create!</button
    >
  </form>
{:else}
  <span class="notSignedInText">Sign in with google to create a pie!</span>
  <button on:click={() => signIn()}>Sign In</button>
{/if}

<style>
  form {
    display: flex;
    flex-direction: column;
    gap: 1em;
  }
  form div {
    display: flex;
    align-items: center;
  }
  form label {
    margin: 0;
  }
  form input,
  form textarea {
    box-shadow: 0em 0em 0.3em rgba(0, 0, 0, 0.2);
  }
  button {
    display: block;
    box-shadow: 0.2em 0.2em 0.6em rgba(0, 0, 0, 0.4);
    border-radius: 1em;
    padding: 1em;
    text-align: center;
    font-weight: bold;
    color: #111;
    border: none;
    background: white;
  }

  button[type='submit'] {
    background-color: rgba(0, 0, 0, 0.8);
    color: white;
  }
  button[type='submit']:disabled {
    background-color: rgba(0, 0, 0, 0.4);
    color: white;
    box-shadow: 0.1em 0.1em 0.3em rgba(0, 0, 0, 0.4);
  }
</style>
