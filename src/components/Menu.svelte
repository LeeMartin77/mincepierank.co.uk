<script lang="ts">
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-nocheck
  // oof
  import { signIn } from '@auth/sveltekit/client';
  import { page } from '$app/stores';
  export let activeYear: number;
  export let customPiesEnabled = false;

  let screenWidth: number;

  $: menuOpen = false;

  function clickOutside(node: any) {
    const handleClick = (event: any) => {
      if (node && !node.contains(event.target) && !event.defaultPrevented) {
        node.dispatchEvent(new CustomEvent('clickOutside', node));
      }
    };

    document.addEventListener('click', handleClick, true);

    return {
      destroy() {
        document.removeEventListener('click', handleClick, true);
      }
    };
  }

  const toggleMenu = () => {
    menuOpen = !menuOpen;
  };
</script>

<svelte:window bind:innerWidth={screenWidth} />
{#if screenWidth > 640 || menuOpen}
  <div class="menu-container" use:clickOutside on:clickOutside={() => (menuOpen = false)}>
    <div class="menu">
      <a on:click={toggleMenu} href="/">Home</a>
      <a on:click={toggleMenu} href="/quickrank">Quickrank</a>
      <a on:click={toggleMenu} href="/years/{activeYear}/all-pies">All Pies</a>
      <a on:click={toggleMenu} href="/years/{activeYear}/brands">Brands</a>
      <a on:click={toggleMenu} href="/years/{activeYear}/categories">Categories</a>
      <a on:click={toggleMenu} href="/profile/rankings/{activeYear}">My Rankings</a>
      {#if !!$page.data.session?.user?.email && customPiesEnabled}
        <a on:click={toggleMenu} href="/profile/pies/{activeYear}">My Custom Pies</a>
      {/if}
    </div>
    <div class="menu-footer">
      {#if !!$page.data.session?.user?.email}
        <a on:click={toggleMenu} href="/profile"
          >{#if !!$page.data.session.user.name}
            Hi {$page.data.session.user.name}
          {:else}
            Welcome back!
          {/if}</a
        >
      {:else}
        <button on:click={() => signIn()}>Sign In</button>
      {/if}
      {#if screenWidth <= 640}
        <button on:click={toggleMenu}>Close Menu</button>
      {/if}
    </div>
  </div>
{:else}
  <button class="menu-button" on:click={toggleMenu}>Menu</button>
{/if}

<style>
  .menu-button {
    position: fixed;
    bottom: 1em;
    left: calc(50% - 5em);
    width: 8em;
    border-radius: 2em;
    margin: 1em;
    z-index: 99999999;
    box-shadow: 0.2em 0.1em 0.5em rgba(0, 0, 0, 0.1);
    background: rgba(255, 255, 255, 0.2);
    backdrop-filter: blur(5px);
    -webkit-backdrop-filter: blur(5px);
    border: 1px solid rgba(255, 255, 255, 0.3);
    cursor: pointer;
  }
  .menu-container {
    display: flex;
    flex-direction: column;
    align-items: left;
    top: 0;
    left: 0;
    max-width: 320px;
    height: 100dvh;
    color: #111;
    background: rgba(255, 255, 255, 0.2);
    box-shadow: 0 4px 30px rgba(0, 0, 0, 0.1);
    backdrop-filter: blur(5px);
    -webkit-backdrop-filter: blur(5px);
    border-right: 1px solid rgba(255, 255, 255, 0.3);
    z-index: 99999999999999999;
  }
  .menu {
    display: flex;
    flex-direction: column;
    overflow-y: auto;
  }

  .menu-footer {
    margin-top: auto;
    margin-bottom: 0;
    border-top: 1px solid rgba(255, 255, 255, 0.3);
    display: flex;
    flex-direction: column;
  }

  a {
    padding: 1em;
    text-decoration: none;
    font-weight: bold;
  }

  a:hover {
    background-color: rgba(0, 0, 0, 0.5);
  }

  button {
    padding: 1em;
    text-decoration: none;
    font-weight: bold;
    width: 100%;
    color: #111;
    background-color: rgba(255, 255, 255, 0.3);
    backdrop-filter: blur(5px);
    -webkit-backdrop-filter: blur(5px);
    border: none;
  }

  button:hover:not(.menu-button) {
    background-color: rgba(255, 255, 255, 0.8);
    cursor: pointer;
  }

  @media only screen and (max-width: 640px) {
    .menu-container {
      position: fixed;
      width: 100%;
    }
  }
</style>
