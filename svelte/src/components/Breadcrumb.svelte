<script lang='ts'>
  import { page } from '$app/stores';
	import { ppCategory } from './utilities/formatCategory';
  $: params = $page.params
  $: route = ($page.route.id || "/").slice(1).split('/');
  $: breadcrumbs = route.reduce<{url: string, name: string}[]>((acc, curr) => {
      const last = acc[acc.length - 1];
      if (curr.indexOf('[') !== -1) {
        // is a param
        const paramName = curr.substring(1, curr.length - 1)
        acc.push({
          url: last.url + params[paramName] + '/',
          name: ppCategory(params[paramName])
        })
      } else {
        acc.push({
          url: last.url + curr + '/',
          name: ppCategory(curr)
        })
      }
      return acc;
    
  }, [{ url: '/', name: 'Home'}])
</script>

<div class="breadcrumb-container">
  {#each breadcrumbs as breadcrumb}
    <a href={breadcrumb.url}>{breadcrumb.name}</a>
  {/each}
</div>

<style>
  .breadcrumb-container {
    display: flex;
    gap: 1em;
  }
</style>