<script lang="ts">
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-nocheck
  // oof
  import { signIn } from '@auth/sveltekit/client';
  import { page } from '$app/stores';
  export let activeYear: number;

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
    <a on:click={toggleMenu} href="/">Home</a>
    <a on:click={toggleMenu} href="/quickrank">Quickrank</a>
    <a on:click={toggleMenu} href="/years/{activeYear}/all-pies">All Pies</a>
    <a on:click={toggleMenu} href="/years/{activeYear}/brands">Brands</a>
    <a on:click={toggleMenu} href="/profile/rankings/{activeYear}">My Rankings</a>
    {#if !!$page.data.session?.user?.email}
      <a on:click={toggleMenu} href="/profile/pies/{activeYear}">My Custom Pies</a>
    {/if}

    {#if !!$page.data.session?.user?.email}
      <a on:click={toggleMenu} class="menu-footer" href="/profile"
        >{#if !!$page.data.session.user.name}
          Hi {$page.data.session.user.name}
        {:else}
          Welcome back!
        {/if}</a
      >
    {:else}
      <button class="menu-footer" on:click={() => signIn()}>Sign In</button>
    {/if}
    {#if screenWidth <= 640}
      <button
        style="width:100%; position:relative; margin-top: 1em; left: 0;"
        class="menu-button"
        on:click={toggleMenu}>Close Menu</button
      >
    {/if}
  </div>
{:else}
  <button class="menu-button" on:click={toggleMenu}>Menu</button>
{/if}

<style>
  .menu-button {
    position: fixed;
    bottom: 1em;
    left: calc(50% - 4em);
    width: 8em;
  }
  .menu-container {
    display: flex;
    flex-direction: column;
    align-items: left;
    top: 0;
    left: 0;
    max-width: 320px;
    height: 100vh;
    color: #111;
    background-color: #fff;
    border-right: 1px solid #111;
  }

  .menu-footer {
    margin-top: auto;
    margin-bottom: 1em;
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
    background-color: #ccc;
    border: none;
  }

  button:hover {
    background-color: #bbb;
    cursor: pointer;
  }

  @media (prefers-color-scheme: dark) {
    .menu-container {
      color: #fff;
      background-color: #111;
      border-right: 1px solid #fff;
    }

    a {
      color: #fff;
    }
    a:visited {
      color: #ddd;
    }

    button {
      color: #fff;
      background-color: #333;
    }

    button:hover {
      background-color: #444;
    }
  }

  @media only screen and (max-width: 640px) {
    .menu-container {
      position: fixed;
      width: 100%;
    }
  }
</style>
