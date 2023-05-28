<script lang="ts">
  import BrandLinkList from '$components/BrandLinkList.svelte';
  import type { PageData } from './$types';

  export let data: PageData;
</script>

<BrandLinkList makers={data.makers} year={data.year} />
