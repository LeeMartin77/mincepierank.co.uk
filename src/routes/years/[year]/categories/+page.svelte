<script lang="ts">
  import Breadcrumb from '$components/Breadcrumb.svelte';
  import CategoryLinkList from '$components/CategoryLinkList.svelte';
  import type { PageData } from './$types';

  export let data: PageData;
</script>

<Breadcrumb />

<CategoryLinkList year={data.year} categories={data.categories} />
