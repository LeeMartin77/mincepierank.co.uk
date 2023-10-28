<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  const dispatcher = createEventDispatcher();

  type ImageSelectOption = {
    imageLink: string;
    text: string;
    value: string;
  };

  export let options: ImageSelectOption[];
  export let value: string | undefined = undefined;
  export let disabled = false;
  let valueOption: ImageSelectOption | undefined;
  export let imgprssr: string;
  export let placeholder = 'Select an Item';

  let open = false;

  const handleSelected = (opt: ImageSelectOption) => {
    value = opt.value;
    valueOption = opt;
    dispatcher('change', opt.value);
    open = false;
  };
</script>

{#if !open}
  <button
    class={`select-button`}
    {disabled}
    on:click={() => {
      if (!disabled) open = true;
    }}
  >
    {#if valueOption}
      <img
        width={'50px'}
        alt={valueOption.text}
        src={imgprssr + valueOption.imageLink + '?width=50&filter=gaussian'}
      />
      <b>{valueOption.text}</b>
    {:else}
      <span>{placeholder}</span>
    {/if}
    <svg xmlns="http://www.w3.org/2000/svg" height="24" viewBox="0 -960 960 960" width="24"
      ><path d="M480-360 280-560h400L480-360Z" /></svg
    >
  </button>
{:else}
  <div class="select-list-root">
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div
      class="cover"
      on:click={() => {
        open = false;
      }}
    />
    <div class="select-list">
      {#each options as option}
        <button on:click={() => handleSelected(option)}>
          <img
            width={'50px'}
            alt={option.text}
            src={imgprssr + option.imageLink + '?width=50&filter=gaussian'}
          />
          <b>{option.text}</b>
        </button>
      {/each}
    </div>
  </div>
{/if}

<style>
  .select-button {
    background: white;
    border: none;
    padding: 1em;
    display: flex;
    align-items: center;
    gap: 1em;
    width: calc(100% - 2em);
    cursor: pointer;
    height: 50px;
    overflow: hidden;
    box-shadow: 0.1em 0.1em 0.4em rgba(0, 0, 0, 0.2);
  }

  .select-button:disabled {
    border-bottom: 1px dashed black;
    cursor: default;
    box-shadow: none;
  }

  .select-button:disabled:hover {
    background: white;
  }
  .select-button svg {
    margin-left: auto;
    margin-right: 0;
  }
  .select-button:hover {
    background: #ddd;
  }

  .select-list-root {
    border: none;
    border-bottom: 1px none white;
    width: 100%;
    height: 50px;
    position: relative;
    overflow: visible;
  }

  .cover {
    position: fixed;
    z-index: 999999998;
    width: 100vw;
    height: 100vh;
    top: 0;
    left: 0;
  }

  .select-list {
    position: absolute;
    display: flex;
    flex-direction: column;
    overflow-y: scroll;
    max-height: 15em;
    width: 100%;
    top: -1em;
    gap: 1px;
    background-color: #000;
    box-shadow: 0.3em 0.2em 0.6em rgba(0, 0, 0, 0.5);
    z-index: 999999999;
  }
  .select-list button {
    width: 100%;
    height: 50px;
    background: none;
    border: none;
    padding: 0.5em 0.5em;
    display: flex;
    flex-direction: row;
    align-items: center;
    overflow: hidden;
    gap: 1em;
    background-color: #eee;
    cursor: pointer;
  }
  .select-list button:hover {
    background-color: #ddd;
  }
</style>
