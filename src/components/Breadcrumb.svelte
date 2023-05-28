<script lang="ts">
  import { page } from '$app/stores';
  import { ppCategory } from './utilities/formatCategory';
  $: params = $page.params;
  $: route = ($page.route.id || '/').slice(1).split('/');
  $: breadcrumbs = route.reduce<{ url: string; name: string }[]>(
    (acc, curr) => {
      if (curr === "") {
        return acc;
      }
      const last = acc[acc.length - 1];
      if (curr.indexOf('[') !== -1) {
        // is a param
        const paramName = curr.substring(1, curr.length - 1);
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
    border: 1px solid grey;
    margin: 1em;
    padding: 1em;
    border-radius: 1em;
  }
  .breadcrumb-container > a {
    color: initial;
  }
  .breadcrumb-container > a:visited {
    color: initial;
  }
</style>
