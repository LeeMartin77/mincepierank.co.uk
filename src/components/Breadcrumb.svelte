<script lang="ts">
  import { page } from '$app/stores';
  import { ppCategory } from './utilities/formatCategory';
  $: params = $page.params;
  $: route = ($page.route.id || '/').slice(1).split('/');
  $: breadcrumbs = route.reduce<{ url: string; name: string }[]>(
    (acc, curr) => {
      if (curr === '') {
        return acc;
      }
      const last = acc[acc.length - 1];
      if (curr.indexOf('[') !== -1) {
        // is a param
        const paramName = curr.substring(1, curr.length - 1);
        if (
          params[paramName].match(/^[{]?[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}[}]?$/)
        ) {
          return acc;
        }
        acc.push({
          url: last.url + params[paramName] + '/',
          name: ppCategory(params[paramName])
        });
      } else {
        acc.push({
          url: last.url + curr + '/',
          name: ppCategory(curr)
        });
      }
      return acc;
    },
    [{ url: '/', name: 'Home' }]
  );
</script>

{#if breadcrumbs.length > 1}
  <div class="breadcrumb-container">
    {#each breadcrumbs as breadcrumb, i}
      {#if i > 0}
        <div>/</div>
      {/if}
      <a href={breadcrumb.url}>{breadcrumb.name}</a>
    {/each}
  </div>
{/if}

<style>
  .breadcrumb-container {
    display: flex;
    flex-wrap: wrap;
    gap: 1em;
    background-color: white;
    box-shadow: 0em 0.2em 0.2em 0 rgba(0, 0, 0, 0.2);
    background: rgba(255, 255, 255, 0.2);
    backdrop-filter: blur(5px);
    -webkit-backdrop-filter: blur(5px);
    border: 1px solid rgba(255, 255, 255, 0.3);
    margin: 1em;
    padding: 1em;
  }
  .breadcrumb-container > a {
    color: #111;
  }
  .breadcrumb-container > a:visited {
    color: #111;
  }
</style>
