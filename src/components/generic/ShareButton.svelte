<script lang="ts">
  export let text: string;
  export let url: string;
  export let title: string | undefined = undefined;
  export let style: string | undefined = undefined;

  let showThanks = false;

  const handleShare = () => {
    navigator
      .share({
        text,
        url,
        title
      })
      .then(() => {
        // thanks for sharing!
        showThanks = true;
      });
  };
</script>

<button class="share-button" {style} on:click={() => handleShare()}><slot /></button>

{#if showThanks}
  <div class="thanks-notice">
    <div>
      <p>Thanks for sharing!</p>
      <button
        on:click={() => {
          showThanks = false;
        }}>Dismiss</button
      >
    </div>
  </div>
{/if}

<style>
  .share-button {
    display: block;
    box-shadow: 0em 0em 0.3em rgba(0, 0, 0, 0.2);
    background: rgba(255, 255, 255, 0.2);
    backdrop-filter: blur(5px);
    -webkit-backdrop-filter: blur(5px);
    border: 1px solid rgba(255, 255, 255, 0.3);
    border-radius: 1em;
    padding: 1em;
    text-align: center;
    font-weight: bold;
    color: #111;
    cursor: pointer;
  }

  .thanks-notice {
    text-align: center;
    border-radius: 1em;
    position: fixed;
    width: calc(100vw - 4em);
    box-shadow: 0.2em 0.1em 0.5em rgba(0, 0, 0, 0.1);
    background: rgba(255, 255, 255, 0.2);
    backdrop-filter: blur(5px);
    -webkit-backdrop-filter: blur(5px);
    border: 1px solid rgba(255, 255, 255, 0.3);
    bottom: 1em;
    left: 1em;
    display: flex;
    flex-direction: column;
    padding: 1em;
    align-items: center;
    gap: 1em;
    z-index: 9999999999999;
  }

  .thanks-notice button {
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
</style>
