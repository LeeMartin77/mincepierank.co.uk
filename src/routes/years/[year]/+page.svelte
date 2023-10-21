<script lang="ts">
  import CategoryLinkList from '$components/CategoryLinkList.svelte';

  import type { PageData } from './$types';
  import BrandLinkList from '$components/BrandLinkList.svelte';
  import LinkButton from '$components/generic/LinkButton.svelte';

  export let data: PageData;
</script>

<LinkButton href={`/years/${data.year}/all-pies`}>View All Pies</LinkButton>

<CategoryLinkList year={data.year} categories={data.categories} />
<BrandLinkList makers={data.makers} year={data.year} imgprssr={data.imgprssr} />

<style>
  a {
    display: block;
    border: 1px solid black;
    font-size: 1em;
    font-weight: bold;
    border-radius: 0.5em;
    padding: 0.5em 1em;
    text-decoration: none;
    text-align: center;
    color: #111;
  }
  a:visited {
    color: #111;
  }
</style>
