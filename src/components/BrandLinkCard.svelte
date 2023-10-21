<script lang="ts">
  import { imgprssrPrefix } from '$lib/imgprssr';
  import type { Maker } from '$lib/storage';
  import Card from './generic/Card.svelte';
  import LinkButton from './generic/LinkButton.svelte';

  export let maker: Maker;
  export let year: number;
  export let imgprssr: string;
</script>

<Card>
  <a class="brand-link-image" href={`/years/${year}/brands/${maker.id}`}>
    <div>
      <img
        src={`${imgprssrPrefix(maker.logo, imgprssr)}?width=200&filter=gaussian`}
        alt={`${maker.name} Logo`}
      />
    </div>
  </a>
  <LinkButton href={`/years/${year}/brands/${maker.id}`}>{maker.name}</LinkButton>
</Card>

<style>
  .brand-link-image {
    display: flex;
    height: 250px;
    align-items: center;
  }
</style>
