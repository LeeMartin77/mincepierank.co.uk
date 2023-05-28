<script lang="ts">
  import CategoryLinkList from '$components/CategoryLinkList.svelte';
  import type { PageData } from './$types';

  export let data: PageData;
</script>

<CategoryLinkList year={data.year} categories={data.categories} />
