<script lang="ts">
  import CategoryLinkList from '$components/CategoryLinkList.svelte';

  import type { PageData } from './$types';
  import BrandLinkList from '$components/BrandLinkList.svelte';

  export let data: PageData;
</script>

<CategoryLinkList year={data.year} categories={data.categories} />
<BrandLinkList makers={data.makers} year={data.year} />
